package receiver

import (
	"bytes"
	"context"
	"strconv"
	"sync"
	"time"

	"github.com/bbdshow/bkit/errc"
	"github.com/bbdshow/qelog/api"
	"github.com/bbdshow/qelog/api/receiverpb"
	"github.com/bbdshow/qelog/pkg/model"
	"github.com/bbdshow/qelog/pkg/types"
)

func (svc *Service) JSONPacketToLogging(ctx context.Context, ip string, in *api.JSONPacket) error {
	if len(in.Data) <= 0 {
		return nil
	}
	// check module is registered
	svc.lock.RLock()
	m, ok := svc.modules[in.Module]
	svc.lock.RUnlock()
	if !ok {
		return errc.ErrNotFound.MultiMsg("module unregistered")
	}

	docs := svc.decodeJSONPacket(ip, in)

	if svc.cfg.Receiver.AlarmEnable && svc.alarm.ModuleIsEnable(in.Module) {
		go svc.alarm.IsAlarm(docs)
	}

	if svc.cfg.Receiver.MetricsEnable {
		go svc.metrics.Statistics(in.Module, ip, docs)
	}

	return svc.createLogging(ctx, m, docs)
}

func (svc *Service) PacketToLogging(ctx context.Context, ip string, in *receiverpb.Packet) error {
	if len(in.Data) <= 0 {
		return nil
	}
	svc.lock.RLock()
	m, ok := svc.modules[in.Module]
	svc.lock.RUnlock()
	if !ok {
		return errc.ErrNotFound.MultiMsg("module unregistered")
	}

	docs := svc.decodePacket(ip, in)

	if svc.cfg.Receiver.AlarmEnable && svc.alarm.ModuleIsEnable(in.Module) {
		go svc.alarm.IsAlarm(docs)
	}

	if svc.cfg.Receiver.MetricsEnable {
		go svc.metrics.Statistics(in.Module, ip, docs)
	}

	return svc.createLogging(ctx, m, docs)
}

func (svc *Service) decodePacket(ip string, in *receiverpb.Packet) []*model.Logging {
	byteItems := bytes.Split(in.Data, []byte{'\n'})
	records := make([]*model.Logging, 0, len(byteItems))

	for i, v := range byteItems {
		if v == nil || bytes.Equal(v, []byte{}) || bytes.Equal(v, []byte{'\n'}) {
			continue
		}
		r := &model.Logging{
			Module:    in.Module,
			IP:        ip,
			Full:      string(v),
			MessageID: in.Id + "_" + strconv.Itoa(i),
			TimeSec:   time.Now().Unix(),
			Size:      len(v),
		}
		dec := types.Decoder{}
		if err := types.Unmarshal(v, &dec); err == nil {
			r.Short = dec.Short()
			r.Level = dec.Level()
			r.Condition1 = dec.Condition(1)
			r.Condition2 = dec.Condition(2)
			r.Condition3 = dec.Condition(3)
			r.TraceID = dec.TraceIDHex()
			r.TimeMill = dec.TimeMill()
			r.TimeSec = r.TimeMill / 1e3
			r.Full = dec.Full()
		}
		records = append(records, r)
	}
	return records
}

func (svc *Service) decodeJSONPacket(ip string, in *api.JSONPacket) []*model.Logging {
	records := make([]*model.Logging, 0, len(in.Data))

	for i, v := range in.Data {
		if v == "" {
			continue
		}
		r := &model.Logging{
			Module:    in.Module,
			IP:        ip,
			Full:      v,
			MessageID: in.Id + "_" + strconv.Itoa(i),
			TimeSec:   time.Now().Unix(),
			Size:      len(v),
		}
		dec := types.Decoder{}
		if err := types.Unmarshal([]byte(v), &dec); err == nil {
			r.Short = dec.Short()
			r.Level = dec.Level()
			r.Condition1 = dec.Condition(1)
			r.Condition2 = dec.Condition(2)
			r.Condition3 = dec.Condition(3)
			r.TraceID = dec.TraceIDHex()
			r.TimeMill = dec.TimeMill()
			r.TimeSec = r.TimeMill / 1e3
			r.Full = dec.Full()
		}
		records = append(records, r)
	}
	return records
}

func (svc *Service) createLogging(ctx context.Context, m *module, docs []*model.Logging) error {
	aDoc, bDoc := svc.loggerDataShardingByTimestamp(m, docs)
	if ctx == nil {
		c, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		ctx = c
		defer cancel()
	}
	inserts := func(ctx context.Context, m *module, v *documents) error {
		if v == nil {
			return nil
		}
		if err := svc.ifCreateCollIndex(ctx, m, v.CollectionName); err != nil {
			return errc.ErrInternalErr.MultiErr(err)
		}

		if err := svc.d.CreateManyLogging(ctx, m.m.Database, v.CollectionName, v.Docs); err != nil {
			return errc.ErrInternalErr.MultiErr(err)
		}
		return nil
	}

	defer func() {
		freeDocuments(aDoc, bDoc)
	}()
	if err := inserts(ctx, m, aDoc); err != nil {
		return err
	}
	if err := inserts(ctx, m, bDoc); err != nil {
		return err
	}

	return nil
}

type documents struct {
	bucket         string
	CollectionName string
	Docs           []interface{}
}

var documentsPool = sync.Pool{New: func() interface{} {
	return &documents{CollectionName: "", Docs: make([]interface{}, 0, 32)}
}}

func initDocuments() *documents {
	v := documentsPool.Get().(*documents)
	v.CollectionName = ""
	v.bucket = ""
	v.Docs = v.Docs[:0]
	return v
}

func freeDocuments(docs ...*documents) {
	for _, v := range docs {
		if v != nil {
			documentsPool.Put(v)
		}
	}
}

// rely on logging time, calc log should be written to that shard
// because multi log in one packet,may be multi shard
func (svc *Service) loggerDataShardingByTimestamp(m *module, docs []*model.Logging) (d1, d2 *documents) {
	// current time shard setting, one packet max written two shard.
	currentName := ""
	d1 = initDocuments()
	for _, v := range docs {
		name := m.sc.EncodeCollName(m.m.Bucket, v.TimeSec)
		if currentName == "" {
			currentName = name
			d1.CollectionName = name
			d1.bucket = m.m.Bucket
		}
		if name != currentName {
			// two shard
			if d2 == nil {
				d2 = initDocuments()
				d2.CollectionName = name
				d2.bucket = m.m.Bucket
			}
			d2.Docs = append(d2.Docs, v)
			continue
		}
		d1.Docs = append(d1.Docs, v)
	}
	return d1, d2
}

// determine collection is exists, if not,it is created and index
// collection name unique
func (svc *Service) ifCreateCollIndex(ctx context.Context, m *module, collectionName string) error {
	svc.lock.Lock()
	defer svc.lock.Unlock()
	if _, ok := svc.collections[collectionName]; ok {
		return nil
	}
	names, err := svc.d.ListCollectionNames(ctx, m.m.Database, m.m.LoggingPrefix())
	if err != nil {
		return err
	}
	exists := false
	for _, n := range names {
		if n == collectionName {
			exists = true
		}
		svc.collections[n] = struct{}{}
	}

	if !exists {
		return svc.d.CreateLoggingIndex(m.m.Database, collectionName)
	}

	return nil
}

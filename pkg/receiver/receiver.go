package receiver

import (
	"context"
	"strconv"
	"sync"
	"time"

	"github.com/huzhongqing/qelog/pkg/config"

	"github.com/huzhongqing/qelog/pkg/receiver/metrics"

	"github.com/huzhongqing/qelog/libs/logs"
	"go.uber.org/zap"

	"github.com/huzhongqing/qelog/pb"

	"github.com/huzhongqing/qelog/pkg/receiver/alarm"

	"github.com/huzhongqing/qelog/pkg/httputil"

	"github.com/huzhongqing/qelog/pkg/common/model"
	"github.com/huzhongqing/qelog/pkg/storage"
	"github.com/huzhongqing/qelog/pkg/types"
)

type Service struct {
	store *storage.Store

	mutex       sync.RWMutex
	modules     map[string]*model.Module
	collections map[string]struct{}

	alarm   *alarm.Alarm
	metrics *metrics.Metrics
}

func NewService(store *storage.Store) *Service {
	srv := &Service{
		store:       store,
		modules:     make(map[string]*model.Module, 0),
		collections: make(map[string]struct{}, 0),
	}

	if err := srv.syncModule(); err != nil {
		panic(err)
	}

	go srv.intervalSyncModule()

	if config.GlobalConfig.AlarmEnable {
		srv.alarm = alarm.NewAlarm()

		if err := srv.syncAlarmRule(); err != nil {
			panic(err)
		}
		go srv.intervalSyncAlarmRule()
	}

	if config.GlobalConfig.MetricsEnable {
		srv.metrics = metrics.NewMetrics(srv.store)
		metrics.SetIncIntervalSec(10)
	}

	return srv
}

func (srv *Service) InsertPacket(ctx context.Context, ip string, in *pb.Packet) error {
	if len(in.Data) <= 0 {
		return nil
	}
	// 判断 module 是否有效，如果无效，则不接受写入
	srv.mutex.RLock()
	module, ok := srv.modules[in.Module]
	srv.mutex.RUnlock()
	if !ok {
		return httputil.NewError(httputil.ErrCodeNotFound, "module unregistered")
	}

	docs := srv.decodePacket(ip, in)

	if config.GlobalConfig.AlarmEnable && srv.alarm.ModuleIsEnable(in.Module) {
		// 异步执行报警逻辑
		go srv.alarm.AlarmIfHitRule(docs)
	}

	if config.GlobalConfig.MetricsEnable {
		go srv.metrics.Statistics(in.Module, ip, docs)
	}

	aDoc, bDoc := srv.loggingShardingByTimestamp(module.DBIndex, docs)

	if ctx == nil {
		ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)
	}
	inserts := func(ctx context.Context, v *InsertDocs) error {
		if v == nil {
			return nil
		}
		ok, err := srv.collectionExists(v.CollectionName)
		if err != nil {
			return httputil.ErrSystemException.MergeError(err)
		}
		if !ok {
			// 如果不存在创建索引
			if err := srv.store.UpsertCollectionIndexMany(model.LoggingIndexMany(v.CollectionName)); err != nil {
				return httputil.ErrSystemException.MergeError(err)
			}
		}

		if err := srv.store.InsertManyLogging(ctx, v.CollectionName, v.Docs); err != nil {
			return httputil.ErrSystemException.MergeError(err)
		}
		return nil
	}

	defer func() {
		if aDoc != nil {
			putInsertDocs(aDoc)
		}
		if bDoc != nil {
			putInsertDocs(bDoc)
		}
	}()
	if err := inserts(ctx, aDoc); err != nil {
		return err
	}
	if err := inserts(ctx, bDoc); err != nil {
		return err
	}

	return nil
}

func (srv *Service) decodePacket(ip string, in *pb.Packet) []*model.Logging {
	records := make([]*model.Logging, len(in.Data))
	for i, v := range in.Data {
		r := &model.Logging{
			Module:    in.Module,
			IP:        ip,
			Full:      v,
			MessageID: in.Id + "_" + strconv.Itoa(i),
			TimeSec:   time.Now().Unix(),
			Size:      len(v),
		}
		val := make(map[string]interface{})
		if err := types.Unmarshal([]byte(v), &val); err == nil {
			dec := types.Decoder{Val: val}
			r.Short = dec.Short()
			r.Level = dec.Level()
			r.Condition1 = dec.Condition(1)
			r.Condition2 = dec.Condition(2)
			r.Condition3 = dec.Condition(3)
			r.TraceID = dec.TraceID()
			r.TimeMill = dec.Time()
			r.TimeSec = r.TimeMill / 1e3
			// full 去掉已经提取出来的字段
			r.Full = dec.Full()
		}
		records[i] = r
	}
	return records
}

// 判断集合是否存在，如果不存在需要创建索引
func (srv *Service) collectionExists(collectionName string) (bool, error) {
	srv.mutex.Lock()
	defer srv.mutex.Unlock()
	if _, ok := srv.collections[collectionName]; ok {
		return true, nil
	}
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
	names, err := srv.store.ListAllCollectionNames(ctx)
	if err != nil {
		return false, err
	}
	exists := false
	for _, n := range names {
		if n == collectionName {
			exists = true
		}
		srv.collections[n] = struct{}{}
	}
	return exists, nil
}

type InsertDocs struct {
	CollectionName string
	Docs           []interface{}
}

var insertDocsPool = sync.Pool{New: func() interface{} {
	return &InsertDocs{CollectionName: "", Docs: make([]interface{}, 0, 32)}
}}

func getInsertDocs() *InsertDocs {
	return insertDocsPool.Get().(*InsertDocs)
}

func putInsertDocs(v *InsertDocs) {
	v.CollectionName = ""
	v.Docs = v.Docs[:0]
	insertDocsPool.Put(v)
}

// 因为是合并包，有少数情况下，根据时间分集合，一个包的内容会写入到不同的集合中区
func (srv *Service) loggingShardingByTimestamp(dbIndex int32, docs []*model.Logging) (a, b *InsertDocs) {
	// 当前时间分片，一组数据最多只会出现在两片上
	currentName := ""
	a = getInsertDocs()
	for _, v := range docs {
		name := model.LoggingCollectionName(dbIndex, v.TimeSec)
		if currentName == "" {
			currentName = name
			a.CollectionName = name
		}
		if name != currentName {
			// 出现了两片的情况
			if b == nil {
				b = getInsertDocs()
				b.CollectionName = name
			}
			b.Docs = append(b.Docs, v)
			continue
		}
		a.Docs = append(a.Docs, v)
	}
	return a, b
}

func (srv *Service) syncModule() error {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	docs, err := srv.store.FindAllModule(ctx)
	if err != nil {
		return err
	}
	srv.mutex.Lock()
	for _, v := range docs {
		srv.modules[v.Name] = v
	}
	srv.mutex.Unlock()
	return nil
}
func (srv *Service) intervalSyncModule() {
	tick := time.NewTicker(30 * time.Second)
	for range tick.C {
		err := srv.syncModule()
		if err != nil {
			logs.Qezap.Error("receiver.service", zap.String("syncModule", err.Error()))
		}
	}
}

func (srv *Service) syncAlarmRule() error {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	docs, err := srv.store.FindAllEnableAlarmRule(ctx)
	if err != nil {
		return err
	}
	srv.alarm.InitRuleState(docs)
	return nil
}

func (srv *Service) intervalSyncAlarmRule() {
	tick := time.NewTicker(35 * time.Second)
	for range tick.C {
		err := srv.syncAlarmRule()
		if err != nil {
			logs.Qezap.Error("receiver.service", zap.String("syncAlarmRule", err.Error()))
		}
	}
}

func (srv *Service) Sync() {
	if srv.metrics != nil {
		srv.metrics.Sync()
	}
}

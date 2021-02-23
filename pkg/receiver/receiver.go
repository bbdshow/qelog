package receiver

import (
	"bytes"
	"context"
	"strconv"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/huzhongqing/qelog/api"
	"github.com/huzhongqing/qelog/api/receiverpb"
	"github.com/huzhongqing/qelog/infra/httputil"
	"github.com/huzhongqing/qelog/infra/logs"
	"github.com/huzhongqing/qelog/pkg/common/model"
	"github.com/huzhongqing/qelog/pkg/config"
	"github.com/huzhongqing/qelog/pkg/receiver/alarm"
	"github.com/huzhongqing/qelog/pkg/receiver/metrics"
	"github.com/huzhongqing/qelog/pkg/storage"
	"github.com/huzhongqing/qelog/pkg/types"
	"go.uber.org/zap"
)

type Service struct {
	store    *storage.Store
	sharding *storage.Sharding

	mutex       sync.RWMutex
	modules     map[string]*model.Module
	collections map[string]struct{}

	alarm   *alarm.Alarm
	metrics *metrics.Metrics
}

func NewService(sharding *storage.Sharding) *Service {
	mainDB, err := sharding.MainStore()
	if err != nil {
		panic(err)
	}
	srv := &Service{
		store:       mainDB,
		sharding:    sharding,
		modules:     make(map[string]*model.Module, 0),
		collections: make(map[string]struct{}, 0),
	}

	if err := srv.syncModule(); err != nil {
		panic(err)
	}

	go srv.intervalSyncModule()

	if config.Global.AlarmEnable {
		srv.alarm = alarm.NewAlarm()

		if err := srv.syncAlarmRule(); err != nil {
			panic(err)
		}
		go srv.intervalSyncAlarmRule()
	}

	if config.Global.MetricsEnable {
		srv.metrics = metrics.NewMetrics(srv.store)
		metrics.SetIncIntervalSec(10)
	}

	return srv
}

func (srv *Service) InsertJSONPacket(ctx context.Context, ip string, in *api.JSONPacket) error {
	if len(in.Data) <= 0 {
		return nil
	}
	// 判断 module 是否有效，如果无效，则不接受写入
	srv.mutex.RLock()
	module, ok := srv.modules[in.Module]
	srv.mutex.RUnlock()
	if !ok {
		return httputil.NewError(httputil.ErrCodeNotFound, in.Module+"module unregistered")
	}

	docs := srv.decodeJSONPacket(ip, in)

	if config.Global.AlarmEnable && srv.alarm.ModuleIsEnable(in.Module) {
		// 异步执行报警逻辑
		go srv.alarm.AlarmIfHitRule(docs)
	}

	if config.Global.MetricsEnable {
		go srv.metrics.Statistics(in.Module, ip, docs)
	}

	return srv.insertLogging(ctx, module.DBIndex, docs)
}

func (srv *Service) InsertPacket(ctx context.Context, ip string, in *receiverpb.Packet) error {
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

	if config.Global.AlarmEnable && srv.alarm.ModuleIsEnable(in.Module) {
		// 异步执行报警逻辑
		go srv.alarm.AlarmIfHitRule(docs)
	}

	if config.Global.MetricsEnable {
		go srv.metrics.Statistics(in.Module, ip, docs)
	}

	return srv.insertLogging(ctx, module.DBIndex, docs)
}

func (srv *Service) insertLogging(ctx context.Context, dbIndex int32, docs []*model.Logging) error {
	aDoc, bDoc := srv.loggingShardingByTimestamp(dbIndex, docs)

	if ctx == nil {
		ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)
	}
	inserts := func(ctx context.Context, v *documents) error {
		if v == nil {
			return nil
		}

		shardingStore, err := srv.sharding.GetStore(v.DBIndex)
		if err != nil {
			return httputil.ErrArgsInvalid.MergeError(err)
		}

		ok, err := srv.collectionExists(shardingStore, v.CollectionName)
		if err != nil {
			return httputil.ErrSystemException.MergeError(err)
		}
		if !ok {
			// 如果不存在创建索引
			if err := shardingStore.UpsertCollectionIndexMany(model.LoggingIndexMany(v.CollectionName)); err != nil {
				return httputil.ErrSystemException.MergeError(err)
			}
		}

		if err := shardingStore.InsertManyLogging(ctx, v.CollectionName, v.Docs); err != nil {
			return httputil.ErrSystemException.MergeError(err)
		}
		return nil
	}

	defer func() {
		freeDocuments(aDoc, bDoc)
	}()
	if err := inserts(ctx, aDoc); err != nil {
		return err
	}
	if err := inserts(ctx, bDoc); err != nil {
		return err
	}

	return nil
}

func (srv *Service) decodePacket(ip string, in *receiverpb.Packet) []*model.Logging {
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
			// full 去掉已经提取出来的字段
			r.Full = dec.Full()
		}
		records = append(records, r)
	}
	return records
}

func (srv *Service) decodeJSONPacket(ip string, in *api.JSONPacket) []*model.Logging {
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
			// full 去掉已经提取出来的字段
			r.Full = dec.Full()
		}
		records = append(records, r)
	}
	return records
}

// 判断集合是否存在，如果不存在需要创建索引
// 因为有序号绑定，每一个集合名都是唯一的
func (srv *Service) collectionExists(store *storage.Store, collectionName string) (bool, error) {
	srv.mutex.Lock()
	defer srv.mutex.Unlock()
	if _, ok := srv.collections[collectionName]; ok {
		return true, nil
	}
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
	names, err := store.ListAllCollectionNames(ctx)
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

type documents struct {
	DBIndex        int32
	CollectionName string
	Docs           []interface{}
}

var documentsPool = sync.Pool{New: func() interface{} {
	return &documents{CollectionName: "", Docs: make([]interface{}, 0, 32)}
}}

func initDocuments() *documents {
	v := documentsPool.Get().(*documents)
	v.CollectionName = ""
	v.DBIndex = 0
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

// 因为是合并包，有少数情况下，根据时间分集合，一个包的内容会写入到不同的集合中区
func (srv *Service) loggingShardingByTimestamp(dbIndex int32, docs []*model.Logging) (a, b *documents) {
	// 当前时间分片，一组数据最多只会出现在两片上
	currentName := ""
	a = initDocuments()
	for _, v := range docs {
		name := model.LoggingCollectionName(dbIndex, v.TimeSec)
		if currentName == "" {
			currentName = name
			a.CollectionName = name
			a.DBIndex = dbIndex
		}
		if name != currentName {
			// 出现了两片的情况
			if b == nil {
				b = initDocuments()
				b.CollectionName = name
				b.DBIndex = dbIndex
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
	hooks := make([]*model.HookURL, 0)
	if _, err := srv.store.FindHookURL(ctx, bson.M{}, &hooks, nil); err != nil {
		return err
	}
	srv.alarm.InitRuleState(docs, hooks)
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

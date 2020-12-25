package receiver

import (
	"context"
	"strconv"
	"sync"
	"time"

	"github.com/huzhongqing/qelog/pkg/httputil"

	"github.com/huzhongqing/qelog/pkg/common/proto/push"

	"github.com/huzhongqing/qelog/pkg/common/model"
	"github.com/huzhongqing/qelog/pkg/storage"
	"github.com/huzhongqing/qelog/pkg/types"
)

type Service struct {
	store *storage.Store

	mutex       sync.RWMutex
	modules     map[string]*model.Module
	collections map[string]struct{}
}

func NewService(store *storage.Store) *Service {
	srv := &Service{
		store:       store,
		modules:     make(map[string]*model.Module, 0),
		collections: make(map[string]struct{}, 0),
	}

	go srv.intervalSyncModule()

	return srv
}

func (srv *Service) InsertPacket(ctx context.Context, ip string, in *push.Packet) error {
	if len(in.Data) <= 0 {
		return nil
	}
	// 判断 module 是否有效，如果无效，则不接受写入
	srv.mutex.RLock()
	mReg, ok := srv.modules[in.Module]
	srv.mutex.RUnlock()
	if !ok {
		return httputil.NewError(httputil.ErrCodeNotFound, "module unregistered")
	}

	docs := srv.decodePacket(ip, in)
	sMap := srv.loggingShardingByTimestamp(mReg.DBIndex, docs)

	if ctx == nil {
		ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)
	}
	for collectionName, docs := range sMap {
		ok, err := srv.collectionExists(collectionName)
		if err != nil {
			return err
		}
		if !ok {
			// 如果不存在创建索引
			if err := srv.store.UpsertCollectionIndexMany(model.LoggingIndexMany(collectionName)); err != nil {
				return err
			}
		}

		if err := srv.store.InsertManyLogging(ctx, collectionName, docs); err != nil {
			return err
		}
	}

	return nil
}

func (srv *Service) decodePacket(ip string, in *push.Packet) []*model.Logging {
	records := make([]*model.Logging, 0, len(in.Data))
	for i, v := range in.Data {
		r := &model.Logging{
			Module:    in.Module,
			IP:        ip,
			Full:      v,
			MessageID: in.Id + "_" + strconv.Itoa(i),
			Timestamp: time.Now().Unix(),
		}
		val := make(map[string]interface{})
		if err := types.Unmarshal([]byte(v), &val); err == nil {
			dec := types.Decoder{Val: val}
			r.Short = dec.Short()
			r.Level = dec.Level()
			r.Condition1 = dec.Condition(1)
			r.Condition2 = dec.Condition(2)
			r.Condition3 = dec.Condition(3)
			r.Full = dec.Full()
			r.Time = dec.Time()
		}
		records = append(records, r)
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

// 因为是合并包，有少数情况下，根据时间分集合，一个包的内容会写入到不同的集合中区
func (srv *Service) loggingShardingByTimestamp(dbIndex int32, docs []*model.Logging) map[string][]interface{} {
	out := make(map[string][]interface{})
	for _, v := range docs {
		name := model.LoggingCollectionName(dbIndex, v.Timestamp)
		val, ok := out[name]
		if ok {
			out[name] = append(val, v)
		} else {
			out[name] = []interface{}{v}
		}
	}
	return out
}

func (srv *Service) intervalSyncModule() {
	tick := time.NewTicker(30 * time.Second)
	for {
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		docs, err := srv.store.FindAllModule(ctx)
		if err == nil {
			srv.mutex.Lock()
			for _, v := range docs {
				srv.modules[v.Name] = v
			}
			srv.mutex.Unlock()
		}
		select {
		case <-tick.C:
		}
	}
}

package receiver

import (
	"context"
	"time"

	"github.com/huzhongqing/qelog/libs/sharding"

	"github.com/huzhongqing/qelog/pkg/common/entity"
	"github.com/huzhongqing/qelog/pkg/common/model"
	"github.com/huzhongqing/qelog/pkg/storage"
	"github.com/huzhongqing/qelog/pkg/types"
)

type Service struct {
	store    *storage.Store
	sharding *sharding.Sharding
}

func NewService(store *storage.Store) *Service {
	srv := &Service{
		store:    store,
		sharding: sharding.NewSharding(sharding.FormatMonth, "logging"),
	}

	return srv
}

func (srv *Service) InsertPacket(uk, ip string, packet entity.DataPacket) error {
	if len(packet.Data) <= 0 {
		return nil
	}
	docs := srv.decodePacket(uk, ip, packet)
	bucket := "test"

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	sMap := srv.shardingLogging(bucket, docs)

	for collectionName, docs := range sMap {
		ok, err := srv.sharding.NameExists(ctx, collectionName, srv.store.ListAllCollectionNames)
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

func (srv *Service) decodePacket(uk, ip string, packet entity.DataPacket) []*model.Logging {
	records := make([]*model.Logging, 0, len(packet.Data))
	for _, v := range packet.Data {
		r := &model.Logging{
			UniqueKey: uk,
			Module:    packet.Module,
			IP:        ip,
			Full:      v,
			TimeStamp: time.Now().Unix(),
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

func (srv *Service) shardingLogging(bucket string, docs []*model.Logging) map[string][]interface{} {
	out := make(map[string][]interface{})
	for _, v := range docs {
		name := srv.sharding.GenerateName(bucket, v.TimeStamp)
		val, ok := out[name]
		if ok {
			out[name] = append(val, v)
		} else {
			out[name] = []interface{}{v}
		}
	}
	return out
}

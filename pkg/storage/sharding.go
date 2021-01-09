package storage

import (
	"context"
	"time"

	"github.com/huzhongqing/qelog/libs/mongo"
	"github.com/huzhongqing/qelog/pkg/config"
)

// 存储分片，把不同的 dbIndex 存储到归类的 DB 实例中，以达到存储横向扩展的目的
// Note: 分片实例一旦设定。如果更改将涉及到数据迁移， 增加不影响
type Sharding struct {
	mainStore     *Store
	shardingStore map[int32]*Store
}

// Sharding Store Lib
func NewSharding(main config.MongoDB, sharding []config.MongoDBIndex) (*Sharding, error) {
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	mainDB, err := mongo.NewDatabase(ctx, main.URI, main.DataBase)
	if err != nil {
		return nil, err
	}
	s := &Sharding{
		mainStore:     New(mainDB),
		shardingStore: make(map[int32]*Store, 0),
	}
	for _, v := range sharding {
		db, err := mongo.NewDatabase(ctx, v.URI, v.DataBase)
		if err != nil {
			return nil, err
		}
		shardingStore := New(db)
		for _, i := range v.Index {
			s.shardingStore[i] = shardingStore
		}
	}

	return s, nil
}

// 配置管理的存储实例
func (s *Sharding) MainStore() (*Store, error) {
	if s.mainStore != nil {
		return s.mainStore, nil
	}
	return nil, ErrMainDBNotFound
}

// 日志记录的实例
func (s *Sharding) GetStore(index int32) (*Store, error) {
	store, ok := s.shardingStore[index]
	if !ok {
		return nil, ErrShardingDBNotFound
	}
	return store, nil
}

func (s *Sharding) Disconnect() {
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	if s.mainStore != nil {
		_ = s.mainStore.database.Client().Disconnect(ctx)
	}
	for _, store := range s.shardingStore {
		if store != nil {
			_ = store.database.Client().Disconnect(ctx)
		}
	}
}

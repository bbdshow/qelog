package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/huzhongqing/qelog/infra/mongo"
	"github.com/huzhongqing/qelog/pkg/common/model"
	"github.com/huzhongqing/qelog/pkg/config"
)

var ShardingDB *Sharding

func SetGlobalShardingDB(db *Sharding) error {
	if db == nil {
		return fmt.Errorf("db nil pointer")
	}
	ShardingDB = db
	return nil
}

// 存储分片，把不同的 dbIndex 存储到归类的 DB 实例中，以达到存储横向扩展的目的
// Note: 分片实例一旦设定。如果更改将涉及到数据迁移， 增加不影响
type Sharding struct {
	mainCfg       config.MongoMainDB
	shardingCfg   []config.MongoShardingDB
	mainStore     *Store
	shardingStore map[int]*Store
}

// Sharding Store Lib
func NewSharding(main config.MongoMainDB, sharding []config.MongoShardingDB, index int) (*Sharding, error) {
	model.SetShardingIndexSize(index)
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	mainDB, err := mongo.NewDatabase(ctx, main.URI, main.DataBase)
	if err != nil {
		return nil, err
	}
	s := &Sharding{
		mainCfg:       main,
		shardingCfg:   sharding,
		mainStore:     New(mainDB),
		shardingStore: make(map[int]*Store, 0),
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
func (s *Sharding) GetStore(index int) (*Store, error) {
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

func (s *Sharding) MainCfg() config.MongoMainDB {
	return s.mainCfg
}

func (s *Sharding) ShardingCfg() []config.MongoShardingDB {
	return s.shardingCfg
}

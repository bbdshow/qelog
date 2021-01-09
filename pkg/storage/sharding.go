package storage

import (
	"context"
	"time"

	"github.com/huzhongqing/qelog/pkg/config"
)

// 存储分片，把不同的 dbIndex 存储到归类的 DB 实例中，以达到存储横向扩展的目的
// Note: 分片实例一旦设定。如果更改将涉及到数据迁移， 增加不影响
type Sharding struct {
	mainStore     *Store
	shardingStore map[int32]*Store
}

func NewSharding(main config.MongoDB, sharding []config.ShardingDB) (*Sharding, error) {

	return nil, nil
}

// 配置管理的存储实例
func (s *Sharding) MainStore() (*Store, bool) {
	if s.mainStore != nil {
		return s.mainStore, true
	}
	return nil, false
}

// 日志记录的实例
func (s *Sharding) GetStore(index int32) (*Store, bool) {
	store, ok := s.shardingStore[index]
	return store, ok
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

package mongo

import (
	"context"
	"time"
)

// 主库配置
type MainConfig struct {
	DataBase string
	URI      string
}

// 分片数据插槽配置
type ShardSlotConfig struct {
	Index    []int
	DataBase string
	URI      string
}

// 存储分片，把不同的 index 存储到归类的 DB 实例中，以达到存储横向扩展的目的
// Note: 分片实例一旦设定。如果更改将涉及到数据迁移， 增加不影响
type Sharding struct {
	main        MainConfig
	slots       []ShardSlotConfig
	mainDB      *Database
	shardSlotDB map[int]*Database
}

// NewSharding Mongodb分区
func NewSharding(main MainConfig, slots []ShardSlotConfig) (*Sharding, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	mdb, err := NewDatabase(ctx, main.URI, main.DataBase)
	if err != nil {
		return nil, err
	}
	s := &Sharding{
		main:        main,
		slots:       slots,
		mainDB:      mdb,
		shardSlotDB: make(map[int]*Database),
	}

	for _, v := range slots {
		slot, err := NewDatabase(ctx, v.URI, v.DataBase)
		if err != nil {
			return nil, err
		}
		for _, i := range v.Index {
			s.shardSlotDB[i] = slot
		}
	}

	return s, nil
}

// MainDB 主实例库
func (s *Sharding) MainDB() (*Database, error) {
	if s.mainDB != nil {
		return s.mainDB, nil
	}
	return nil, ErrMainDBNotFound
}

// ShardSlotDB 分片数据插槽库，多用于海量数据存储
func (s *Sharding) ShardSlotDB(index int) (*Database, error) {
	db, ok := s.shardSlotDB[index]
	if !ok {
		return nil, ErrShardSlotNotFound
	}
	return db, nil
}

// MainConfig 主库配置
func (s *Sharding) MainConfig() MainConfig {
	return s.main
}

// ShardSlotsConfig 分片数据插槽配置
func (s *Sharding) ShardSlotsConfig() []ShardSlotConfig {
	slots := make([]ShardSlotConfig, len(s.slots))
	copy(slots, s.slots)
	return slots
}

// Disconnect 断开连接
func (s *Sharding) Disconnect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if s.mainDB != nil {
		_ = s.mainDB.Client().Disconnect(ctx)
	}
	for _, db := range s.shardSlotDB {
		if db != nil {
			_ = db.Client().Disconnect(ctx)
		}
	}
	return nil
}

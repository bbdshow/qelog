package model

import (
	"fmt"
	"github.com/huzhongqing/qelog/infra/mongo"
)

// 定义存储数据结构
var ShardingDB *mongo.Sharding
var MainDB *mongo.Database

func SetGlobalMainDB(db *mongo.Database) error {
	if db == nil {
		return fmt.Errorf("main db nil pointer")
	}
	MainDB = db
	return nil
}

func SetGlobalShardingDB(db *mongo.Sharding) error {
	if db == nil {
		return fmt.Errorf("sharding db nil pointer")
	}
	ShardingDB = db
	return nil
}

func AllIndex() []mongo.Index {
	all := make([]mongo.Index, 0)
	all = append(all, AlarmRuleIndexMany()...)
	all = append(all, ModuleIndexMany()...)
	all = append(all, ModuleMetricsIndexMany()...)
	all = append(all, CollStatsIndexMany()...)
	return all
}

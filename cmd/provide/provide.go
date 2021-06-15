package provide

import (
	"github.com/bbdshow/qelog/infra/logs"
	"github.com/bbdshow/qelog/infra/mongo"
	"github.com/bbdshow/qelog/pkg/common/model"
	"github.com/bbdshow/qelog/pkg/config"
	"go.uber.org/zap"
)

// 初始化一些依赖
func InitConfig() *config.Config {
	cfg := config.InitConfig(ConfigPath)
	if err := cfg.Validate(); err != nil {
		logs.Qezap.Fatal("config validate", zap.Error(err))
		return cfg
	}
	config.SetGlobalConfig(cfg)
	return cfg
}

func InitMongodb(cfg *config.Config, createIndex bool) *mongo.Sharding {
	slots := make([]mongo.ShardSlotConfig, 0)
	for _, v := range cfg.Sharding {
		slots = append(slots, mongo.ShardSlotConfig{
			Index:    v.Index,
			DataBase: v.DataBase,
			URI:      v.URI,
		})
	}
	sharding, err := mongo.NewSharding(mongo.MainConfig{
		DataBase: cfg.Main.DataBase,
		URI:      cfg.Main.URI,
	}, slots)
	if err != nil {
		logs.Qezap.Fatal("mongo connect failed ", zap.Error(err))
	}

	if err := model.SetGlobalShardingDB(sharding); err != nil {
		logs.Qezap.Fatal("SetGlobalShardingDB", zap.Error(err))
	}

	mainDB, err := sharding.MainDB()
	if err != nil {
		logs.Qezap.Fatal("mongo connect failed", zap.Error(err))
	}
	if err := model.SetGlobalMainDB(mainDB); err != nil {
		logs.Qezap.Fatal("SetGlobalMainDB", zap.Error(err))
	}
	if createIndex {
		if err := mainDB.UpsertCollectionIndexMany(model.AllIndex()); err != nil {
			logs.Qezap.Fatal("mongo create index", zap.Error(err))
		}
	}
	return sharding
}

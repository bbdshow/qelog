package tests

import (
	"github.com/bbdshow/qelog/infra/logs"
	"github.com/bbdshow/qelog/infra/mongo"
	"github.com/bbdshow/qelog/pkg/common/model"
	"github.com/bbdshow/qelog/pkg/config"
)

func InitTestDepends(cfgPath ...string) {
	InitTestConfig(cfgPath...)
	InitTestQezap(config.Global)
	InitShardingDB(config.Global)
}

func InitTestConfig(cfgPath ...string) {
	path := ""
	if len(cfgPath) > 0 {
		path = cfgPath[0]
	}
	cfg := config.InitConfig(path)
	config.SetGlobalConfig(cfg)
}

func InitTestQezap(cfg *config.Config) {
	logs.InitQezap(cfg.Logging.Addr, cfg.Logging.Module, cfg.Logging.Filename)
}

func InitShardingDB(cfg *config.Config) {
	slots := make([]mongo.ShardSlotConfig, 0)
	for _, v := range cfg.Sharding {
		slots = append(slots, mongo.ShardSlotConfig{
			Index:    v.Index,
			DataBase: v.DataBase,
			URI:      v.URI,
		})
	}
	db, err := mongo.NewSharding(mongo.MainConfig{
		DataBase: cfg.Main.DataBase,
		URI:      cfg.Main.URI,
	}, slots)
	if err != nil {
		panic(err)
	}

	_ = model.SetGlobalShardingDB(db)

	mainDB, err := db.MainDB()
	if err != nil {
		panic(err)
	}
	_ = model.SetGlobalMainDB(mainDB)

	// 创建索引
	if err := mainDB.UpsertCollectionIndexMany(model.AllIndex()); err != nil {
		panic(err)
	}
}

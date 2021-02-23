package tests

import (
	"github.com/huzhongqing/qelog/infra/logs"
	"github.com/huzhongqing/qelog/pkg/config"
	"github.com/huzhongqing/qelog/pkg/storage"
	"go.uber.org/zap"
)

func InitTestDepends() {
	InitTestConfig()
	InitTestQezap()
	InitShardingDB()
}

func InitTestConfig() {
	cfg := config.InitConfig("")
	config.SetGlobalConfig(cfg)
}

func InitTestQezap() {
	cfg := config.InitConfig("")
	logs.InitQezap(cfg.Logging.Addr, cfg.Logging.Module, cfg.Logging.Filename)
}

func InitShardingDB() {
	cfg := config.InitConfig("")
	sharding, err := storage.NewSharding(cfg.Main, cfg.Sharding, cfg.MaxShardingIndex)
	if err != nil {
		logs.Qezap.Fatal("mongo connect failed ", zap.Error(err))
	}

	if err := storage.SetGlobalShardingDB(sharding); err != nil {
		logs.Qezap.Fatal("SetGlobalShardingDB", zap.Error(err))
	}
}

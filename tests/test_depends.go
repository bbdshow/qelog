package tests

import (
	"github.com/huzhongqing/qelog/infra/logs"
	"github.com/huzhongqing/qelog/pkg/config"
	"github.com/huzhongqing/qelog/pkg/storage"
	"go.uber.org/zap"
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
	sharding, err := storage.NewSharding(cfg.Main, cfg.Sharding, cfg.ShardingIndexSize)
	if err != nil {
		logs.Qezap.Fatal("mongo connect failed ", zap.Error(err))
	}

	if err := storage.SetGlobalShardingDB(sharding); err != nil {
		logs.Qezap.Fatal("SetGlobalShardingDB", zap.Error(err))
	}
}

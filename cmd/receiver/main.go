package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/huzhongqing/qelog/libs/logs"

	"github.com/huzhongqing/qelog/pkg/common/model"
	"github.com/huzhongqing/qelog/pkg/storage"

	"github.com/huzhongqing/qelog/pkg/receiver"

	"github.com/huzhongqing/qelog/pkg/config"
)

var (
	goVersion = ""
	buildTime = ""
	gitHash   = ""

	configPath = "./configs/config.toml"
	version    = false
)

func main() {
	flag.StringVar(&configPath, "f", "./configs/config.toml", "config file default(./configs/config.toml)")
	flag.BoolVar(&version, "v", false, "show version")
	flag.Parse()

	if version {
		fmt.Printf("goVersion: %s \nbuildTime: %s \ngitHash: %s \n", goVersion, buildTime, gitHash)
		return
	}

	cfg := config.InitConfig(configPath)
	if err := cfg.Validate(); err != nil {
		panic(fmt.Sprintf("config validate %s", err.Error()))
		return
	}

	config.SetGlobalConfig(cfg)

	logs.InitQezap(nil, "")
	sharding, err := storage.NewSharding(cfg.Main, cfg.Sharding)
	if err != nil {
		logs.Qezap.Fatal("mongo connect failed", zap.Error(err))
	}

	if !cfg.Release() {
		db, err := sharding.MainStore()
		if err != nil {
			logs.Qezap.Fatal("mongo connect failed", zap.Error(err))
		}
		if err := db.Database().UpsertCollectionIndexMany(
			model.ModuleMetricsIndexMany()); err != nil {
			logs.Qezap.Fatal("mongo create index", zap.Error(err))
		}
	}

	httpSrv := receiver.NewHTTPService(sharding)

	go func() {
		logs.Qezap.Info("http server listen ", zap.String("addr", cfg.ReceiverAddr))
		if err := httpSrv.Run(cfg.ReceiverAddr); err != nil {
			logs.Qezap.Fatal("http server listen failed", zap.Error(err))
		}
	}()

	grpcSrv := receiver.NewGRPCService(sharding)
	go func() {
		logs.Qezap.Info("gRPC server listen ", zap.String("addr", cfg.ReceiverGRPCAddr))
		if err := grpcSrv.Run(cfg.ReceiverGRPCAddr); err != nil {
			logs.Qezap.Fatal("gRPC server listen failed", zap.Error(err))
		}
	}()

	signalAccept()
	_ = httpSrv.Close()
	_ = grpcSrv.Close()
	_ = sharding.Disconnect
}

func signalAccept() {
	// 不同的信号量不同的处理方式
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		si := <-ch
		switch si {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			time.Sleep(time.Second)
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}

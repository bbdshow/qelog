package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/huzhongqing/qelog/pkg/common/model"

	"github.com/huzhongqing/qelog/pkg/storage"

	"github.com/huzhongqing/qelog/libs/logs"

	"github.com/huzhongqing/qelog/pkg/manager"

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

	sharding, err := storage.NewSharding(cfg.Main, cfg.Sharding)
	if err != nil {
		log.Fatalln("mongo connect failed ", err.Error())
	}

	if !cfg.Release() {
		db, err := sharding.MainStore()
		if err != nil {
			panic(err)
		}
		if err := db.Database().UpsertCollectionIndexMany(
			model.ModuleIndexMany(),
			model.AlarmRuleIndexMany()); err != nil {
			panic("create main db index " + err.Error())
		}
	}

	logs.InitQezap(cfg.Logging.Addr, cfg.Logging.Module)

	httpSrv := manager.NewHTTPService(sharding)

	go func() {
		if err := httpSrv.Run(cfg.ManagerAddr); err != nil {
			log.Fatalln("http server listen failed ", err.Error())
		}
		log.Println("http server listen ", cfg.ManagerAddr)
	}()

	signalAccept()

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

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/huzhongqing/qelog/libs/logs"

	"github.com/huzhongqing/qelog/pkg/manager"

	"github.com/huzhongqing/qelog/libs/mongo"

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

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	database, err := mongo.NewDatabase(ctx, cfg.MainDB.URI,
		cfg.MainDB.DataBase)
	if err != nil {
		log.Fatalln("mongo connect failed ", err.Error())
	}

	config.SetGlobalConfig(cfg)
	logs.InitQezap(cfg.Logging.Addr, cfg.Logging.Module)

	httpSrv := manager.NewHTTPService(database)

	go func() {
		if err := httpSrv.Run(cfg.ManagerAddr); err != nil {
			log.Fatalln("http server listen failed ", err.Error())
		}
		log.Println("http server listen ", cfg.ManagerAddr)
	}()

	signalAccept()

	_ = database.Client().Disconnect(nil)
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

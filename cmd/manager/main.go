package main

import (
	"context"
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

func main() {
	cfg := config.InitConfig("./configs/config.toml")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	database, err := mongo.NewDatabase(ctx, cfg.MongoDB.URI,
		cfg.MongoDB.DataBase)
	if err != nil {
		log.Fatalln("mongo connect failed ", err.Error())
	}

	logs.InitQezap([]string{"127.0.0.1:31082"}, "example")

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

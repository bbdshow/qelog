package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	//rec := receiver.New(database)
	//
	//go func() {
	//	if err := rec.Run(cfg.ReceiverAddr); err != nil {
	//		log.Fatalln("http server listen failed ", err.Error())
	//	}
	//	log.Println("http server listen ", cfg.ReceiverAddr)
	//}()

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

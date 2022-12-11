package main

import (
	"fmt"
	"log"

	"github.com/bbdshow/bkit/logs"
	"github.com/bbdshow/bkit/runner"
	"github.com/bbdshow/qelog/pkg/conf"
	"github.com/bbdshow/qelog/pkg/server/grpc"
	"github.com/bbdshow/qelog/pkg/server/http"

	"go.uber.org/zap"

	"github.com/bbdshow/qelog/pkg/receiver"
)

func main() {
	if err := conf.InitConf(); err != nil {
		panic(err)
	}
	logs.InitQezap(conf.Conf.Logging)
	defer logs.Qezap.Close()
	logs.Qezap.Info("init", zap.Any("config", conf.Conf.EraseSensitive()))

	svc := receiver.NewService(conf.Conf)
	defer svc.Close()

	// http
	go func() {
		if conf.Conf.Receiver.HttpListenAddr != "" {
			httpSvc := http.NewReceiverHttpServer(conf.Conf, svc)
			if err := runner.RunServer(httpSvc,
				runner.WithListenAddr(conf.Conf.Receiver.HttpListenAddr),
			); err != nil {
				panic(fmt.Sprintf("runner exit: %v\n", err))
			}
		}
	}()

	rpcSvc := grpc.NewReceiverGRpc(conf.Conf, svc)
	if err := runner.RunServer(rpcSvc,
		runner.WithListenAddr(conf.Conf.Receiver.RpcListenAddr),
	); err != nil {
		log.Printf("runner exit: %v\n", err)
	}
}

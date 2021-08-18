package main

import (
	"github.com/bbdshow/bkit/logs"
	"github.com/bbdshow/bkit/runner"
	"github.com/bbdshow/qelog/pkg/conf"
	"github.com/bbdshow/qelog/pkg/server/grpc"
	"log"

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

	rpcSvc := grpc.NewReceiverGRpc(conf.Conf, svc)

	if err := runner.Run(rpcSvc, func() error {
		// dealloc
		return nil
	}, runner.WithListenAddr(conf.Conf.Receiver.RpcListenAddr)); err != nil {
		log.Printf("runner exit: %v\n", err)
	}
}

package main

import (
	"fmt"
	"github.com/bbdshow/qelog/cmd/provide"
	"github.com/bbdshow/qelog/infra/kit"
	"github.com/bbdshow/qelog/pkg/httpserver"
	"go.uber.org/multierr"
	"go.uber.org/zap"

	"github.com/bbdshow/qelog/infra/logs"

	"github.com/bbdshow/qelog/pkg/receiver"
)

func main() {
	provide.InitFlag()
	cfg := provide.InitConfig()
	logs.InitQezap(cfg.Logging.Addr, cfg.Logging.Module, cfg.Logging.Filename)
	db := provide.InitMongodb(cfg, false)

	httpSrv := httpserver.NewHTTPServer(cfg.Env)
	receiver.RegisterRouter(httpSrv.Engine())

	go func() {
		fmt.Println("http server listen", cfg.ReceiverAddr)
		if err := httpSrv.Run(cfg.ReceiverAddr); err != nil {
			logs.Qezap.Fatal("http server listen failed", zap.Error(err))
		}
	}()

	grpcSrv := receiver.NewGRPCService()
	go func() {
		fmt.Println("gRPC server listen", cfg.ReceiverGRPCAddr)
		if err := grpcSrv.Run(cfg.ReceiverGRPCAddr); err != nil {
			logs.Qezap.Fatal("gRPC server listen failed", zap.Error(err))
		}
	}()

	kit.SignalAccept(func() error {
		// 释放资源
		return multierr.Combine(httpSrv.Close(), grpcSrv.Close(), db.Disconnect(), logs.Qezap.Close())
	}, nil)
}

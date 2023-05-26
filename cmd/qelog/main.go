package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/bbdshow/bkit/logs"
	"github.com/bbdshow/bkit/runner"
	"github.com/bbdshow/qelog/pkg/admin"
	"github.com/bbdshow/qelog/pkg/conf"
	"github.com/bbdshow/qelog/pkg/receiver"
	"github.com/bbdshow/qelog/pkg/server/grpc"
	"github.com/bbdshow/qelog/pkg/server/http"
	"github.com/bbdshow/qelog/pkg/types"
	"go.uber.org/zap"
	"log"
	"time"
)

var (
	mode string
)

func main() {
	flag.StringVar(&mode, "mode", string(types.Single), "server mode, single | cluster_admin | cluster_receiver")

	if err := conf.InitConf(); err != nil {
		panic(err)
	}
	logs.InitQezap(conf.Conf.Logging)
	defer logs.Qezap.Close()

	time.AfterFunc(time.Second, func() {
		// Deferred execution，wait for the receiver server to be initialized
		logs.Qezap.Info("init", zap.Any("config", conf.Conf.EraseSensitive()))
	})

	ctx := context.Background()

	smode := types.GetFlagOrOSEnvServerMode(types.ServerMode(mode))
	fmt.Println("current running server mode:", smode)
	switch smode {
	case types.Single:
		go func() {
			adminServer(ctx, conf.Conf)
		}()
		time.Sleep(10 * time.Millisecond)
		receiverServer(ctx, conf.Conf)

	case types.ClusterAdmin:
		adminServer(ctx, conf.Conf)
	case types.ClusterReceiver:
		receiverServer(ctx, conf.Conf)
	default:
		logs.Qezap.Fatal("invalid server mode")
	}

}

func adminServer(ctx context.Context, conf *conf.Config) {

	pwd := types.GetOSEnvAdminPassword()
	if pwd != "" {
		conf.Admin.Password = pwd
	}

	svc := admin.NewService(conf)
	defer svc.Close()

	httpSvc := http.NewAdminHttpServer(conf, svc)
	if err := runner.RunServer(httpSvc,
		runner.WithListenAddr(conf.Admin.HttpListenAddr),
		runner.WithContext(ctx),
	); err != nil {
		log.Printf("runner exit: %v\n", err)
	}
}

func receiverServer(ctx context.Context, conf *conf.Config) {

	svc := receiver.NewService(conf)
	defer svc.Close()

	rpcSvc := grpc.NewReceiverGRpc(conf, svc)
	if err := runner.RunServer(rpcSvc,
		runner.WithListenAddr(conf.Receiver.RpcListenAddr),
		runner.WithContext(ctx),
	); err != nil {
		log.Printf("runner exit: %v\n", err)
	}
}

package main

import (
	"fmt"
	"github.com/huzhongqing/qelog/cmd/provide"
	"github.com/huzhongqing/qelog/infra/kit"
	"github.com/huzhongqing/qelog/infra/logs"
	"github.com/huzhongqing/qelog/pkg/admin"
	"github.com/huzhongqing/qelog/pkg/httpserver"
	"go.uber.org/multierr"
	"go.uber.org/zap"
)

func main() {
	provide.InitFlag()
	cfg := provide.InitConfig()
	logs.InitQezap(cfg.Logging.Addr, cfg.Logging.Module, cfg.Logging.Filename)
	db := provide.InitMongodb(cfg, true)

	httpSrv := httpserver.NewHTTPServer(cfg.Env)
	// 注册后台路由
	admin.RegisterRouter(httpSrv.Engine())

	go func() {
		fmt.Println("http server listen", cfg.AdminAddr)
		if err := httpSrv.Run(cfg.AdminAddr); err != nil {
			logs.Qezap.Fatal("http server listen failed", zap.Error(err))
		}
	}()

	logs.Qezap.Info("init", zap.Any("config", cfg.Print()), zap.String("buildInfo", provide.BuildInfo()))

	kit.SignalAccept(func() error {
		// 释放资源
		return multierr.Combine(db.Disconnect(), httpSrv.Close(), logs.Qezap.Close())
	}, nil)
}

package main

import (
	"log"

	"github.com/bbdshow/bkit/logs"
	"github.com/bbdshow/bkit/runner"
	"github.com/bbdshow/qelog/pkg/admin"
	"github.com/bbdshow/qelog/pkg/conf"
	"github.com/bbdshow/qelog/pkg/server/http"
	"go.uber.org/zap"
)

func main() {
	if err := conf.InitConf(); err != nil {
		panic(err)
	}
	logs.InitQezap(conf.Conf.Logging)
	defer logs.Qezap.Close()
	logs.Qezap.Info("init", zap.Any("config", conf.Conf.EraseSensitive()))

	svc := admin.NewService(conf.Conf)
	defer svc.Close()

	httpSvc := http.NewAdminHttpServer(conf.Conf, svc)
	if err := runner.RunServer(httpSvc,
		runner.WithListenAddr(conf.Conf.Admin.HttpListenAddr)); err != nil {
		log.Printf("runner exit: %v\n", err)
	}
}

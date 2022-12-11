package admin

import (
	"context"
	"os"
	"testing"

	"github.com/bbdshow/qelog/pkg/conf"
	"github.com/bbdshow/qelog/pkg/model"
)

var svc *Service

func TestMain(m *testing.M) {
	if err := conf.InitConf("../../configs/config.toml"); err != nil {
		panic(err)
	}
	svc = NewService(conf.Conf)
	os.Exit(m.Run())
}

func TestService_CreateModule(t *testing.T) {
	in := &model.CreateModuleReq{
		Name:     "benchmark",
		Desc:     "benchmark",
		DaySpan:  0,
		MaxMonth: 0,
	}
	if err := svc.CreateModule(context.Background(), in); err != nil {
		t.Fatal(err)
	}
}

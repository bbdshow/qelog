package receiver

import (
	"os"
	"testing"

	"github.com/bbdshow/qelog/pkg/conf"
)

var svc *Service

func TestMain(m *testing.M) {
	if err := conf.InitConf("../../configs/config.toml"); err != nil {
		panic(err)
	}
	svc = NewService(conf.Conf)
	os.Exit(m.Run())
}

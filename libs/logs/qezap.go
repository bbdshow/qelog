package logs

import (
	"github.com/huzhongqing/qelog/qezap"
	"go.uber.org/zap"
)

var Qezap *qezap.Logger

func InitQezap(addrs []string, moduleName string) {
	cfg := qezap.NewConfig(addrs, moduleName)
	if len(addrs) == 0 {
		cfg.SetEnableRemote(false)
	}
	Qezap = qezap.New(cfg, zap.DebugLevel)
}

package logs

import (
	"time"

	"github.com/huzhongqing/qelog/qezap"
	"go.uber.org/zap"
)

var Qezap *qezap.Logger

func init() {
	// 注册一个本地的 Log
	cfg := qezap.NewConfig(nil, "").SetMaxAge(6 * 30 * 24 * time.Hour)
	Qezap = qezap.New(cfg, zap.DebugLevel)
}

func InitQezap(addrs []string, moduleName string) {
	// 默认注册的
	if Qezap != nil {
		// 则追加远程模块
		if len(addrs) > 0 && moduleName != "" {
			cfg := Qezap.Config().SetEnableRemote(true).SetAddr(addrs).SetModule(moduleName)
			Qezap.AppendWriter(qezap.NewWriteRemote(cfg))
		}
		return
	}
	cfg := qezap.NewConfig(addrs, moduleName)
	if len(addrs) == 0 {
		cfg.SetEnableRemote(false)
	}

	Qezap = qezap.New(cfg, zap.DebugLevel)
}

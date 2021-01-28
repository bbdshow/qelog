package qezap

import (
	"testing"
	"time"

	"go.uber.org/zap"
)

func TestReNew(t *testing.T) {
	localCfg := NewConfig(nil, "")
	log := New(localCfg, zap.DebugLevel)
	go func() {
		for {
			log.Info("info")
		}
	}()

	time.Sleep(time.Second)
	remoteCfg := localCfg.SetEnableRemote(true).
		SetAddr([]string{"127.0.0.1:31082"}).
		SetModule("example").
		SetMaxSize(100 << 20)
	log = New(remoteCfg, zap.DebugLevel)

	time.Sleep(time.Minute)
}

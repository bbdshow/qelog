package wrapzap

import (
	"testing"
	"time"

	"go.uber.org/zap"
)

const (
	_URL = "http://127.0.0.1:31081/v1/receive/packet"
)

func TestNewWrapZap(t *testing.T) {

	cfg := NewConfig("./log/logger.log", true, _URL, "test")
	cfg.WriteRemote.MaxPacket = 1024

	zapLog := NewWrapZap(cfg, zap.DebugLevel)

	zapLog.Debug("Debug", zap.String("k", "v"), zap.String("l", RandString(1024)))
	zapLog.Info("Info", zap.String("k", "v"), zap.String("k1", "v1"))
	zapLog.Warn("Warn", zap.String("k", "v"))
	zapLog.Error("Error", zap.String("k", "v"))
	zapLog.DPanic("DPanic", zap.String("k", "v"))

	time.Sleep(3 * time.Second)

	zapLog.Debug("Debug", zap.String("k", "v"))
	zapLog.Info("Info", zap.String("k", "v"))
	zapLog.Warn("Warn", zap.String("k", "v"))
	zapLog.Error("Error", zap.String("k", "v"))
	zapLog.DPanic("DPanic", zap.String("k", "v"))

	time.Sleep(5 * time.Second)

	zapLog.Error("Sync", zap.String("最后写入", "v"))

	zapLog.Sync()
}

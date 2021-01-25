package qezap

import (
	"testing"
	"time"

	"go.uber.org/zap"
)

func TestLogger_AppendWriter(t *testing.T) {
	cfg := NewConfig(nil, "")
	w := New(cfg, zap.DebugLevel)
	g := 5
	for g > 0 {
		g--
		go func() {
			for {
				w.Debug("debug", zap.Int64("nsec", time.Now().UnixNano()))
			}
		}()
	}
	time.Sleep(50 * time.Millisecond)

	w.Config().SetEnableRemote(true).
		SetModule("example").SetAddr([]string{"127.0.0.1:31082"})

	w.AppendWriter(NewWriteRemote(w.Config()))

	time.Sleep(5 * time.Second)
}

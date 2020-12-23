package main

import (
	"testing"

	"github.com/huzhongqing/qelog/qezap"

	"go.uber.org/zap"
)

func BenchmarkQzzap(b *testing.B) {
	addrs := []string{"127.0.0.1:31082"}
	cfg := qezap.NewConfig(addrs, "example")

	// 如果设置 false，可以 addrs = nil
	// cfg.SetEnableRemote(false)

	// 如果对默认配置不满足，可直接设置
	qzLog := qezap.New(cfg, zap.DebugLevel)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		qzLog.Info("benchmark", zap.Int("index", i))
	}
	qzLog.Sync()
	// BenchmarkQzzap-8   	   23943	     48908 ns/op
}

func BenchmarkQzzapHTTP(b *testing.B) {
	addrs := []string{"http://127.0.0.1:31081/v1/receiver/packet"}
	cfg := qezap.NewConfig(addrs, "example")
	cfg.SetHTTPTransport()
	// 如果设置 false，可以 addrs = nil
	// cfg.SetEnableRemote(false)

	// 如果对默认配置不满足，可直接设置
	qzLog := qezap.New(cfg, zap.DebugLevel)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		qzLog.Info("benchmark", zap.Int("index", i))
	}
	qzLog.Sync()
	// BenchmarkQzzapHTTP-8   	   19798	     57735 ns/op
}

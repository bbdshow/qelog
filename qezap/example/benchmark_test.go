package main

import (
	"testing"

	"github.com/huzhongqing/qelog/qezap"

	"go.uber.org/zap"
)

func BenchmarkQezap(b *testing.B) {
	addrs := []string{"127.0.0.1:31082"}
	cfg := qezap.NewConfig(addrs, "example")
	cfg.WriteRemote.MaxConcurrent = 10000
	// 如果设置 false，可以 addrs = nil
	//cfg.SetEnableRemote(false)
	// 如果对默认配置不满足，可直接设置
	qeLog := qezap.New(cfg, zap.DebugLevel)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		qeLog.Info("benchmark", zap.Int("index", i))
	}
	qeLog.Sync()
	// BenchmarkQezap-8   	   32242	     38522 ns/op
	// 关闭远程传输 性能相差5倍 ...
	// BenchmarkQezap-8   	  160515	      6305 ns/op
}

func BenchmarkQezapHTTP(b *testing.B) {
	addrs := []string{"http://127.0.0.1:31081/v1/receiver/packet"}
	cfg := qezap.NewConfig(addrs, "example")
	cfg.SetHTTPTransport()
	// 如果设置 false，可以 addrs = nil
	// cfg.SetEnableRemote(false)
	//cfg.WriteRemote.MaxConcurrent = 200
	//cfg.WriteRemote.MaxPacket = 500 * 1024
	// 如果对默认配置不满足，可直接设置
	qeLog := qezap.New(cfg, zap.DebugLevel)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		qeLog.Info("benchmark", zap.Int("index", i))
	}
	qeLog.Sync()
	// BenchmarkQzzapHTTP-8   	   19798	     57735 ns/op
}

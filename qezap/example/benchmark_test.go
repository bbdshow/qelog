package main

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/huzhongqing/qelog/qezap"

	"go.uber.org/zap"
)

func BenchmarkQezapRpc(b *testing.B) {
	addrs := []string{"127.0.0.1:31082"}
	cfg := qezap.NewConfig(addrs, "example")
	cfg.WriteRemote.MaxConcurrent = 50
	// 如果设置 false，可以 addrs = nil
	//cfg.SetEnableRemote(false)
	// 如果对默认配置不满足，可直接设置
	qeLog := qezap.New(cfg, zap.DebugLevel)
	time.Sleep(2 * time.Second)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		qeLog.Info(strconv.Itoa(i), zap.Int("index", i))
	}
	fmt.Println(b.N)
	qeLog.Sync()
	// BenchmarkQezapRpc-8   	  101323	     11634 ns/op
	// 关闭远程传输 性能相差接近2倍， 是因为 zap.core 要生成两份数据。
	// BenchmarkQezap-8   	  160515	      6305 ns/op
}

func BenchmarkQezapHTTP(b *testing.B) {
	addrs := []string{"http://127.0.0.1:31081/v1/receiver/packet"}
	cfg := qezap.NewConfig(addrs, "example")
	cfg.SetHTTPTransport()
	cfg.WriteRemote.MaxConcurrent = 50
	// 如果设置 false，可以 addrs = nil
	// cfg.SetEnableRemote(false)
	//cfg.WriteRemote.MaxConcurrent = 200
	//cfg.WriteRemote.MaxPacket = 500 * 1024
	// 如果对默认配置不满足，可直接设置
	qeLog := qezap.New(cfg, zap.DebugLevel)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		qeLog.Info(strconv.Itoa(i), zap.Int("index", i))
	}
	qeLog.Sync()
	// BenchmarkQezapHTTP-8   	   89986	     12213 ns/op
}

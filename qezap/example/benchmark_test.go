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
	//addrs := []string{"192.168.10.114:31082"}
	cfg := qezap.NewConfig(addrs, "benchmark")
	qeLog := qezap.New(cfg, zap.DebugLevel)
	time.Sleep(time.Second)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		qeLog.Info(strconv.Itoa(i), zap.Int("index", i))
	}
	fmt.Println(b.N)
	qeLog.Sync()
}

func BenchmarkQezapHTTP(b *testing.B) {
	//addrs := []string{"http://127.0.0.1:31081/v1/receiver/packet"}
	addrs := []string{"http://192.168.10.114:31081/v1/receiver/packet"}
	cfg := qezap.NewConfig(addrs, "benchmark")
	cfg.SetHTTPTransport()
	qeLog := qezap.New(cfg, zap.DebugLevel)
	time.Sleep(time.Second)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		qeLog.Info(strconv.Itoa(i), zap.Int("index", i))
	}
	fmt.Println(b.N)
	qeLog.Sync()
}

func BenchmarkTraceID(b *testing.B) {
	var tid qezap.TraceID
	for i := 0; i < b.N; i++ {
		tid.New()
	}
}

func BenchmarkQezapRpcWrite(b *testing.B) {
	addrs := []string{"127.0.0.1:31082"}
	cfg := qezap.NewConfig(addrs, "benchmark")

	writeRemote := qezap.NewWriteRemote(cfg)

	time.Sleep(time.Second)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c := fmt.Sprintf(`{"_level":"INFO","_time":1609944069261.2573,"_caller":"example/main.go:39","_func":"main.loopWriteLogging","_short":"%d","_traceid":"160994406926125730025244977","val":823537}`, i)
		writeRemote.Write([]byte(c))
	}
	fmt.Println(b.N)
	writeRemote.Sync()
}

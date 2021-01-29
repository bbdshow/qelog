package benchmark

import (
	crand "crypto/rand"
	"fmt"
	"io"
	"math/rand"
	"testing"
	"time"

	"github.com/huzhongqing/qelog/qezap"

	"go.uber.org/zap"
)

func BenchmarkQezapRpc(b *testing.B) {
	addrs := []string{"127.0.0.1:31082"}
	//addrs := []string{"192.168.10.114:31082"}
	//qezap.EnableRelease()
	cfg := qezap.NewConfig(addrs, "benchmark")
	qeLog := qezap.New(cfg, zap.DebugLevel)

	time.Sleep(time.Second)

	simpleWrite(qeLog, b)
	//write(qeLog, b)

	qeLog.Sync()
}

func BenchmarkQezapHTTP(b *testing.B) {
	addrs := []string{"http://127.0.0.1:31081/v1/receiver/packet"}
	//addrs := []string{"http://192.168.10.114:31081/v1/receiver/packet"}

	cfg := qezap.NewConfig(addrs, "benchmark")
	cfg.SetHTTPTransport()
	qeLog := qezap.New(cfg, zap.DebugLevel)

	time.Sleep(time.Second)

	write(qeLog, b)

	qeLog.Sync()
}

func simpleWrite(qeLog *qezap.Logger, b *testing.B) {
	countSize := 0
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := "index"
		countSize += 2*len("") + 128
		switch i % 4 {
		case 0:
			qeLog.Debug(key, zap.String("val", key))
		case 1:
			qeLog.Info(key, zap.String("val", key))
		case 2:
			qeLog.Warn(key, zap.String("val", key))
		case 3:
			qeLog.Error(key, zap.String("val", key))
		}

	}
	fmt.Println(b.N, "条数", "总数据大小约", countSize/1024, "KB")
}

func write(qeLog *qezap.Logger, b *testing.B) {
	keysMap := map[int]string{}

	for i := 0; i < 500; i++ {
		b := make([]byte, rand.Int31n(100))
		_, _ = io.ReadFull(crand.Reader, b[:])
		keysMap[i] = string(b)
	}
	countSize := 0
	b.ResetTimer()

	// 这样后台会创建不同的索引与值,模拟后端写入比较真实的情况
	for i := 0; i < b.N; i++ {
		key := keysMap[i%500]
		countSize += 2*len(key) + 128
		switch i % 4 {
		case 0:
			qeLog.Debug(key, zap.String("val", key))
		case 1:
			qeLog.Info(key, zap.String("val", key))
		case 2:
			qeLog.Warn(key, zap.String("val", key))
		case 3:
			qeLog.Error(key, zap.String("val", key))
		}

	}
	fmt.Println(b.N, "条数", "总数据大小约", countSize/1024, "KB")
}

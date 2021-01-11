package main

import (
	"context"
	"flag"
	"fmt"
	"math/rand"
	"runtime"
	"time"

	"github.com/huzhongqing/qelog/qezap"
	"go.uber.org/zap"
)

var (
	addr   string
	module string
)

func main() {
	flag.StringVar(&addr, "a", "127.0.0.1:31082", "remote grpc addr (default: 127.0.0.1:31082)")
	flag.StringVar(&module, "m", "benchmark", "module name,please register this module (default: benchmark)")
	flag.Parse()

	cfg := qezap.NewConfig([]string{addr}, module)
	cfg.SetFilename("./data/log/logger.log")
	cfg.SetEnableRemote(true)
	qelog := qezap.New(cfg, zap.DebugLevel)

	initMem := readMem()
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Minute)

	size := 4 * 1024
	tick := time.NewTicker(5 * time.Millisecond)
	exit := false

	sumSize := int64(size)
	count := 0
	for range tick.C {
		select {
		case <-ctx.Done():
			exit = true
			break
		default:
			sumSize += int64(size)
			count++
			if count%1000 == 0 {
				fmt.Println(memStatsString(readMem()))
			}
			buff := make([]byte, size)
			_, _ = rand.Read(buff)
			qelog.Debug("mem use", zap.ByteString("val", buff))
		}
		if exit {
			break
		}
	}
	fmt.Printf(`
init mem: %s,
over mem: %s,
sum size: %d,
sum count: %d,
`, memStatsString(initMem), memStatsString(readMem()), sumSize, count)

}

func readMem() *runtime.MemStats {
	stats := &runtime.MemStats{}
	runtime.ReadMemStats(stats)
	return stats
}

func memStatsString(v *runtime.MemStats) string {
	return fmt.Sprintf("Sys: %d, TotalAlloc: %d, Mallocs: %d, Frees:%d,HeapIdle: %d,HeapInuse: %d, HeapReleased: %d,GCCPUFraction:%f",
		v.Sys, v.TotalAlloc, v.Mallocs, v.Frees, v.HeapIdle, v.HeapInuse, v.HeapReleased, v.GCCPUFraction)
}

package main

import (
	"context"
	"flag"
	"fmt"
	"math/rand"
	"runtime"
	"time"

	"github.com/bbdshow/qelog/qezap"
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
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Minute)

	tick := time.NewTicker(10 * time.Millisecond)
	exit := false
	printNum := 0
	sumSize := int64(0)
	count := 0
	for range tick.C {
		select {
		case <-ctx.Done():
			exit = true
			break
		default:
			size := rand.Int63n(1024)
			sumSize += size
			count++
			if count%1000 == 0 {
				printNum++
				fmt.Println(printNum, memStatsString(readMem()))
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

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/bbdshow/qelog/qezap"
	"go.uber.org/zap"
)

var (
	addr     string
	module   string
	duration string
	size     int
)

// qezap client memory used test
func main() {
	flag.StringVar(&addr, "a", "127.0.0.1:31082", "remote grpc addr (default: 127.0.0.1:31082)")
	flag.StringVar(&module, "m", "benchmark", "module name,please register this module (default: benchmark)")
	flag.StringVar(&duration, "t", "30s", "benchmark time duration (default: 30s)")
	flag.IntVar(&size, "s", 256, "text size,single log total byte calc = (size+100) 100 extra char")
	flag.Parse()

	randStr := ""
	for i := 0; i < size; i++ {
		randStr += "a"
	}

	d, err := time.ParseDuration(duration)
	if err != nil {
		log.Fatal(err.Error())
	}

	lg := qezap.New(
		qezap.WithFilename("./log/memory_use.log"),
		qezap.WithAddrsAndModuleName([]string{addr}, module))

	initMem := readMem()
	ctx, cancel := context.WithTimeout(context.Background(), d)
	defer cancel()

	exit := false
	sumSize := int64(0)
	count := 0
	for {
		select {
		case <-ctx.Done():
			exit = true
			break
		default:
			sumSize += int64(len(randStr) + 100)
			count++
			if count%100000 == 0 {
				log.Println(memStatsString(readMem()))
			}
			switch count % 4 {
			case 0:
				lg.Debug("memory used debug", zap.String("txt", randStr))
			case 1:
				lg.Info("memory used info", zap.String("txt", randStr))
			case 2:
				lg.Warn("memory used warn", zap.String("txt", randStr))
			case 3:
				lg.Error("memory used error", zap.String("txt", randStr))
			}
		}
		if exit {
			break
		}
	}
	_ = lg.Close()

	fmt.Printf(`
Duration: %s,
Process init memory used: %s,
Process exec memory used: %s,
Total transport data size ≈: %d MB,
Total log count: %d,
Single log avg operation latency: %s
`, d.String(), memStatsString(initMem), memStatsString(readMem()), sumSize>>20, count, time.Duration(d.Nanoseconds()/int64(count)).String())
}

func readMem() *runtime.MemStats {
	stats := &runtime.MemStats{}
	runtime.ReadMemStats(stats)
	return stats
}

func memStatsString(v *runtime.MemStats) string {
	return fmt.Sprintf("Sys: %d MB, TotalAlloc: %d MB, Mallocs: %d, Frees:%d,HeapIdle: %d,HeapInuse: %d, HeapReleased: %d,GCCPUFraction:%f",
		v.Sys>>20, v.TotalAlloc>>20, v.Mallocs, v.Frees, v.HeapIdle, v.HeapInuse, v.HeapReleased, v.GCCPUFraction)
}

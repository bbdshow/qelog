package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"go.uber.org/zap"

	"github.com/bbdshow/qelog/qezap"
)

var (
	addr     string
	module   string
	duration string
	size     int
)

func main() {
	flag.StringVar(&addr, "a", "127.0.0.1:31082", "remote grpc addr (default: 127.0.0.1:31082)")
	flag.StringVar(&module, "mn", "benchmark", "module name,please register this module (default: benchmark)")
	flag.StringVar(&duration, "t", "3s", "benchmark time duration (default: 3s)")
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
		qezap.WithFilename("./log/benchmark.log"),
		qezap.WithAddrsAndModuleName([]string{addr}, module))

	count := 0
	debug := 0
	info := 0
	warn := 0
	errC := 0

	ctx, cancel := context.WithTimeout(context.Background(), d)
	defer cancel()

	exit := false
	for {
		select {
		case <-ctx.Done():
			exit = true
			break
		default:
			switch count % 4 {
			case 0:
				debug++
				lg.Debug("benchmark debug", zap.String("txt", randStr))
			case 1:
				info++
				lg.Info("benchmark info", zap.String("txt", randStr))
			case 2:
				warn++
				lg.Warn("benchmark warn", zap.String("txt", randStr))
			case 3:
				errC++
				lg.Error("benchmark error", zap.String("txt", randStr))
			}
			count++
		}
		if exit {
			break
		}
	}
	// keep sync over
	lg.Close()

	fmt.Printf(`
When the default IO capacity is exceeded, the backup file is written.
Benchmark duration: %s
Grpc addr: %s
Log count: %d
Level debug count: %d
Level info count: %d
Level warn count: %d
Level error count: %d
Single log size %d byte,avg operation latency: %s
Total data size ≈ : %d MB
`, d.String(), addr, count, debug, info, warn, errC, size+100,
		time.Duration(d.Nanoseconds()/int64(count)).String(), (count*(size+100))>>20)
}

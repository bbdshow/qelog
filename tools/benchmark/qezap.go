package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"go.uber.org/zap"

	"github.com/huzhongqing/qelog/qezap"
)

var (
	addr     string
	module   string
	duration string
)

func main() {
	flag.StringVar(&addr, "a", "127.0.0.1:31082", "remote grpc addr (default: 127.0.0.1:31082)")
	flag.StringVar(&module, "m", "benchmark", "module name,please register this module (default: benchmark)")
	flag.StringVar(&duration, "t", "1s", "benchmark time duration (default: 1s)")
	flag.Parse()

	d, err := time.ParseDuration(duration)
	if err != nil {
		log.Fatal(err.Error())
	}
	cfg := qezap.NewConfig([]string{addr}, module)
	cfg.SetFilename("./data/log/qelogger.log")
	qelog := qezap.New(cfg, zap.DebugLevel)

	count := 0
	debug := 0
	info := 0
	warn := 0
	errC := 0

	ctx, _ := context.WithTimeout(context.Background(), d)
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
				qelog.Debug("benchmark debug", zap.Int("count", count))
			case 1:
				info++
				qelog.Info("benchmark info", zap.Int("count", count))
			case 2:
				warn++
				qelog.Warn("benchmark warn", zap.Int("count", count))
			case 3:
				errC++
				qelog.Error("benchmark error", zap.Int("count", count))
			}
			count++
		}
		if exit {
			break
		}
	}
	qelog.Sync()
	fmt.Printf(`
When the default IO capacity is exceeded, the backup file is written.
benchmark duration: %s
grpc addr: %s
all count: %d
debug count: %d
info count: %d
warn count: %d
error count: %d
operation delay: %s
`, d.String(), addr, count, debug, info, warn, errC, time.Duration(d.Nanoseconds()/int64(count)).String())
}

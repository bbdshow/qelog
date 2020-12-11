package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/huzhongqing/qelog/cmd/receiver"

	"github.com/huzhongqing/qelog/config"
)

func main() {
	cfg := config.InitConfig("./config/config.toml")

	exec := "receiver"

	if len(os.Args) > 1 {
		exec = os.Args[1]
	}
	var finalRelease func() error

	switch exec {
	case "receiver":
		process, err := receiver.New(cfg)
		if err != nil {
			panic(err)
		}

		go func() {
			if err := process.Run(); err != nil {
				log.Fatalln(err)
			}
		}()
		finalRelease = process.Close
	default:
		panic("unknown exec")
	}

	signalAccept()

	if err := finalRelease(); err != nil {
		log.Fatalln("finalRelease", err)
	}
}

func signalAccept() {
	// 不同的信号量不同的处理方式
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		si := <-ch
		switch si {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			time.Sleep(time.Second)
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}

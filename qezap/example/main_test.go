package main

import (
	"strconv"
	"testing"
	"time"

	"github.com/huzhongqing/qelog/qezap"
	"go.uber.org/zap"
)

func TestQezapLoopWrite(t *testing.T) {
	// 持续写一段时间
	addrs := []string{"127.0.0.1:31082"}
	cfg := qezap.NewConfig(addrs, "example")
	cfg.SetMaxConcurrent(100)
	// 如果设置 false，可以 addrs = nil
	//cfg.SetEnableRemote(false)
	// 如果对默认配置不满足，可直接设置
	qeLog := qezap.New(cfg, zap.DebugLevel)
	time.Sleep(time.Second)
	go func() {
		for i := 0; i < 10000; i++ {
			qeLog.Info(strconv.Itoa(i), zap.Int("index", i))
		}
	}()
	go func() {
		for i := 0; i < 10000; i++ {
			qeLog.Warn(strconv.Itoa(i), zap.Int("index", i))
		}
	}()
	go func() {
		for i := 0; i < 10000; i++ {
			qeLog.Error(strconv.Itoa(i), zap.Int("index", i))
		}
	}()
	go func() {
		for i := 0; i < 10000; i++ {
			qeLog.Error(strconv.Itoa(i), zap.Int("index", i))
		}
	}()
	go func() {
		for i := 0; i < 10000; i++ {
			qeLog.Error(strconv.Itoa(i), zap.Int("index", i))
		}
	}()
	time.Sleep(10 * time.Second)
	qeLog.Sync()
}

func TestQezapWrite(t *testing.T) {
	// 持续写一段时间
	addrs := []string{"127.0.0.1:31082"}
	cfg := qezap.NewConfig(addrs, "example")

	// 如果设置 false，可以 addrs = nil
	//cfg.SetEnableRemote(false)
	// 如果对默认配置不满足，可直接设置
	qeLog := qezap.New(cfg, zap.DebugLevel)
	time.Sleep(2 * time.Second)

	for i := 0; i < 10; i++ {
		qeLog.Info(strconv.Itoa(i), zap.Int("index", i))
	}
	go func() {
		count := 0
		for range time.Tick(3 * time.Second) {
			count++
			qeLog.Info(strconv.Itoa(count), zap.Int("index", 1))
		}
	}()
	time.Sleep(5 * time.Minute)
	qeLog.Sync()
}

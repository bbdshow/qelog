package main

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"github.com/huzhongqing/qelog/qezap"

	"go.uber.org/zap"
)

func main() {
	//loopWriteLogging()
	writeLogging()
	//multiModuleLogging()
}

func loopWriteLogging() {
	addrs := []string{"127.0.0.1:31082"}
	cfg := qezap.NewConfig(addrs, "example")
	cfg.SetFilename("./data/log/example.log")
	qeLog := qezap.New(cfg, zap.DebugLevel)
	s := time.Now()
	count := 0
	for {
		time.Sleep(10 * time.Millisecond)
		count++
		if count > 10000000 {
			break
		}
		ctx := context.Background()
		ctx = qeLog.WithTraceID(ctx)
		val := rand.Int63n(10000000)
		shrot := strconv.Itoa(rand.Intn(100000))
		switch val % 4 {
		case 1:
			qeLog.Info(shrot, qeLog.TraceIDField(ctx), zap.Int64("val", val))
			qeLog.Warn(shrot, qeLog.TraceIDField(ctx), zap.Int64("val", val), qeLog.ConditionOne(shrot))
			qeLog.Error(shrot, qeLog.TraceIDField(ctx), zap.Int64("val", val), qeLog.ConditionOne(shrot), qeLog.ConditionTwo(shrot))
		case 2:
			qeLog.Warn(shrot, qeLog.TraceIDField(ctx), zap.Int64("val", val))
			qeLog.Error(shrot, qeLog.TraceIDField(ctx), zap.Int64("val", val), qeLog.ConditionOne(shrot), qeLog.ConditionTwo(shrot))
		case 3:
			qeLog.Error(shrot, qeLog.TraceIDField(ctx), zap.Int64("val", val))
		default:
			qeLog.Debug(shrot, qeLog.TraceIDField(ctx), zap.Int64("val", val))
		}
	}
	fmt.Println(time.Now().Sub(s))
	time.Sleep(30 * time.Minute)
}

func writeLogging() {
	//addrs := []string{"http://127.0.0.1:31081/v1/receiver/packet"}
	addrs := []string{"127.0.0.1:31082"}
	cfg := qezap.NewConfig(addrs, "example")
	//cfg.SetHTTPTransport()
	// 如果设置 false，可以 addrs = nil
	// cfg.SetEnableRemote(false)

	// 如果对默认配置不满足，可直接设置
	cfg.WriteRemote.MaxPacket = 512

	qeLog := qezap.New(cfg, zap.DebugLevel)

	qeLog.Debug("Debug", zap.String("k", "v"), zap.String("num", "1234567890"))
	qeLog.Info("Info", zap.String("k", "v"), zap.String("k1", "v1"))

	qeLog.Warn("Warn", zap.String("k", "v"),
		qeLog.ConditionOne("默认条件查询1"),
		qeLog.ConditionTwo("默认条件查询2, 当有条件1，在配合条件2，查询更快"),
		qeLog.ConditionThree("与2同理，我是条件3"))

	ctx := context.Background()
	ctx = qeLog.WithTraceID(ctx)
	qeLog.Info("teceid", qeLog.TraceIDField(ctx))

	qeLog.Error("Error", zap.String("k", "v"))
	qeLog.DPanic("DPanic", zap.String("k", "v"))

	// 在这之前，还未到默认发包时间，也不满足缓存容量，所以，这些信息是缓存在本地的。
	time.Sleep(2 * time.Second)
	//  满足默认发包时间了，所以日志已经发送走了。
	qeLog.Error("Alarm", zap.String("info", "测试一条报警信息"))
	qeLog.Error("Sync", zap.String("结束最后写入", "final"))
	// sync 执行后，缓存在本地的日志，将全部发送
	qeLog.Sync()
	time.Sleep(time.Minute)
	qeLog.Fatal("Fatal", zap.String("这个Fatal, 也是能写进去的哟", "Fatal"))
	fmt.Println("never print")
}

func multiModuleLogging() {
	addrs := []string{"127.0.0.1:31082"}

	exp := qezap.New(qezap.NewConfig(addrs, "example"), zap.DebugLevel)
	exp2 := qezap.New(qezap.NewConfig(addrs, "example2"), zap.DebugLevel)

	wg := sync.WaitGroup{}
	wg.Add(2)
	count := 10
	go func(c int) {
		for c > 0 {
			c--
			exp.Debug("example", exp.ConditionOne(strconv.Itoa(c)))
		}
		wg.Done()
	}(count)
	go func(c int) {
		for c > 0 {
			c--
			exp2.Debug("example2", exp.ConditionOne(strconv.Itoa(c)))
		}
		wg.Done()
	}(count)
	wg.Wait()

	exp.Sync()
	exp2.Sync()

	//time.Sleep(3 * time.Second)
}

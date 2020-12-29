package main

import (
	"fmt"
	"time"

	"github.com/huzhongqing/qelog/qezap"

	"go.uber.org/zap"
)

func main() {
	addrs := []string{"127.0.0.1:31082"}
	cfg := qezap.NewConfig(addrs, "example")

	// 如果设置 false，可以 addrs = nil
	// cfg.SetEnableRemote(false)

	// 如果对默认配置不满足，可直接设置
	cfg.WriteRemote.MaxPacket = 256

	qeLog := qezap.New(cfg, zap.DebugLevel)

	qeLog.Debug("Debug", zap.String("k", "v"), zap.String("num", "1234567890"))
	qeLog.Info("Info", zap.String("k", "v"), zap.String("k1", "v1"))

	qeLog.Warn("Warn", zap.String("k", "v"),
		qezap.NewCondition().StringFiled("默认条件查询1"),
		qezap.NewCondition().Two().StringFiled("默认条件查询2, 当有条件1，在配合条件2，查询更快"),
		qezap.NewCondition().Three().StringFiled("与2同理，我是条件3"))

	qeLog.Error("Error", zap.String("k", "v"))
	qeLog.DPanic("DPanic", zap.String("k", "v"))

	// 在这之前，还未到默认发包时间，也不满足缓存容量，所以，这些信息是缓存在本地的。
	time.Sleep(2 * time.Second)
	//  满足默认发包时间了，所以日志已经发送走了。
	qeLog.Error("Alarm", zap.String("info", "测试一条报警信息"))
	qeLog.Error("Sync", zap.String("结束最后写入", "final"))
	// sync 执行后，缓存在本地的日志，将全部发送
	qeLog.Sync()
	qeLog.Fatal("Fatal", zap.String("这个Fatal, 也是能写进去的哟", "Fatal"))
	fmt.Println("never print")
}

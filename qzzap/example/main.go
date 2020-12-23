package main

import (
	"fmt"
	"time"

	"github.com/huzhongqing/qzzap"

	"go.uber.org/zap"
)

func main() {
	addrs := []string{"127.0.0.1:31082"}
	cfg := qzzap.NewConfig(addrs, "example")

	// 如果设置 false，可以 addrs = nil
	// cfg.SetEnableRemote(false)

	// 如果对默认配置不满足，可直接设置
	cfg.WriteRemote.MaxPacket = 256

	qzLog := qzzap.New(cfg, zap.DebugLevel)

	qzLog.Debug("Debug", zap.String("k", "v"), zap.String("num", "1234567890"))
	qzLog.Info("Info", zap.String("k", "v"), zap.String("k1", "v1"))

	qzLog.Warn("Warn", zap.String("k", "v"),
		qzzap.NewCondition().StringFiled("默认条件查询1"),
		qzzap.NewCondition().Two().StringFiled("默认条件查询2, 当有条件1，在配合条件2，查询更快"),
		qzzap.NewCondition().Three().StringFiled("与2同理，我是条件3"))

	qzLog.Error("Error", zap.String("k", "v"))
	qzLog.DPanic("DPanic", zap.String("k", "v"))

	// 在这之前，还未到默认发包时间，也不满足缓存容量，所以，这些信息是缓存在本地的。
	time.Sleep(2 * time.Second)
	//  满足默认发包时间了，所以日志已经发送走了。

	qzLog.Error("Sync", zap.String("结束最后写入", "final"))
	// sync 执行后，缓存在本地的日志，将全部发送
	qzLog.Sync()
	qzLog.Fatal("Fatal", zap.String("这个Fatal, 也是能写进去的哟", "Fatal"))
	fmt.Println("never print")
}

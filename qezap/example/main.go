package main

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap/zapcore"

	"github.com/bbdshow/qelog/qezap"

	"go.uber.org/zap"
)

/*
// Qezap 支持远端的 Zap, 是Zap的超集
var Qezap *qezap.Logger

func init() {
	// 注册本地Log，不具有远程写入能力
	Qezap = qezap.New(qezap.NewConfig(nil, ""), zap.DebugLevel)
}

type Config struct {
	Addr     []string `defval:""`
	Module   string   `defval:""`
	Filename string   `defval:"./log/logger.log"`
	Level    int      `defval:"-1"` // -1=debug 0=info ...
}

func InitQezap(cfg *Config) {
	_ = Qezap.Close()
	Qezap = qezap.New(qezap.NewConfig(cfg.Addr, cfg.Module).SetFilename(cfg.Filename), zapcore.Level(cfg.Level))
}
*/
var qelog *qezap.Logger

func init() {
	// 一般使用上面的注释代码去init日志就行
	cfg := qezap.NewConfig([]string{"127.0.0.1:31082"}, "example")
	// 设置每一次发送远端的包大小
	cfg.SetMaxPacketSize(64 << 10)
	// 设置本地日志文件保存最大时间
	cfg.SetMaxAge(30 * 24 * time.Hour)
	// config 具体设置可查看响应的方法

	cfg.SetFilename("./log/qelogger.log")

	qelog = qezap.New(cfg, zap.DebugLevel, zap.AddStacktrace(zap.ErrorLevel))
	// 测试，等待后台建立好 gRPC 连接
	time.Sleep(time.Second)
}

func main() {

	// 普通用法
	qelog.Debug("Debug", zap.String("val", "i am string field"))

	// 动态修改日志等级
	qelog.SetEnabledLevel(zapcore.InfoLevel)

	qelog.Debug("Debug", zap.String("val", "should not be output"))

	// 携带条件查询, 条件必需前置设置，只能 1 或 1,2 不能 2,3 这样后台不会提供查询
	qelog.Error("condition example", qelog.ConditionOne("userid"), qelog.ConditionTwo("0001"), qelog.ConditionThree("phone"))

	// 携带 TraceID 打印到日志
	// 这是初始上下文
	ctx := context.Background()
	// 已经携带好 TraceID
	ctx = qelog.WithTraceID(ctx)
	// 会获取 ctx 的 TraceID
	qelog.Warn("have trace id field", zap.String("withCtx", "warn"), qelog.FieldTraceID(ctx))
	qelog.Error("have trace id field", zap.String("withCtx", "error"), qelog.FieldTraceID(ctx))

	// 还可以获取 ctx 里面的 TraceID 写入到 Response Header 等
	tid := qelog.TraceID(ctx)
	fmt.Println(tid.Hex())

	replaceZapLogger := qelog.Logger
	replaceZapLogger.Info("这种方式，可以替换之前项目用的 zap.Logger")

	w := qelog.NewWriter(zap.InfoLevel, "writer")
	w.Write([]byte("io.Writer impl"))

	qelog.DPanic("last message")
	if err := qelog.Sync(); err != nil {
		fmt.Println(err)
	}
	if err := qelog.Close(); err != nil {
		fmt.Println(err)
	}
}

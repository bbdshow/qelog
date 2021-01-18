package main

import (
	"context"
	"fmt"
	"io"
	"time"

	"go.uber.org/zap/zapcore"

	"github.com/huzhongqing/qelog/qezap"

	"go.uber.org/zap"
)

var qelog *qezap.Logger

func init() {
	cfg := qezap.NewConfig([]string{"127.0.0.1:31082"}, "example")
	// 设置每一次发送远端的包大小
	cfg.SetMaxPacketSize(64 << 10)
	// 设置本地日志文件保存最大时间
	cfg.SetMaxAge(30 * 24 * time.Hour)
	// config 具体设置可查看响应的方法

	qelog = qezap.New(cfg, zap.DebugLevel)
	// 测试，等待后台建立好 gRPC 连接
	time.Sleep(time.Second)
}

func main() {
	// 普通用法
	qelog.Debug("Debug", zap.String("val", "i am string field"))
	// {"_level":"DEBUG","_time":1610616377872.952,"_caller":"example/main.go:33","_func":"main.main","_short":"Debug","val":"i am string field"}

	// 动态修改日志等级
	qelog.SetEnabledLevel(zapcore.InfoLevel)

	qelog.Debug("Debug", zap.String("val", "should not be output"))

	// 携带条件查询, 条件必需前置设置，只能 1 或 1,2 不能 2,3 这样后台不会提供查询
	qelog.Error("condition example", qelog.ConditionOne("userid"), qelog.ConditionTwo("0001"), qelog.ConditionThree("phone"))
	// {"_level":"ERROR","_time":1610616377872.952,"_caller":"example/main.go:36","_func":"main.main","_short":"condition example","_condition1":"userid","_condition2":"0001","_condition3":"phone"}

	// 携带 TraceID 打印到日志
	// 这是初始上下文
	ctx := context.Background()
	// 已经携带好 TraceID
	ctx = qelog.WithTraceID(ctx)
	// 会获取 ctx 的 TraceID
	qelog.WarnWithCtx(ctx, "have trace id field", zap.String("withCtx", "warn"))
	qelog.ErrorWithCtx(ctx, "have trace id field", zap.String("withCtx", "error"))
	// {"_level":"WARN","_time":1610674696407.182,"_caller":"qezap/qezap.go:242","_func":"github.com/huzhongqing/qelog/qezap.(*Logger).encoderWithCtx","_short":"have trace id field","withCtx":"warn","_traceid":"165a441a480f69e89916bf32"}
	// {"_level":"ERROR","_time":1610674696407.182,"_caller":"qezap/qezap.go:244","_func":"github.com/huzhongqing/qelog/qezap.(*Logger).encoderWithCtx","_short":"have trace id field","withCtx":"error","_traceid":"165a441a480f69e89916bf32"}

	// 还可以获取 ctx 里面的 TraceID 写入到 Response Header 等
	tid := qelog.MustGetTraceID(ctx)
	fmt.Println(tid.Hex())

	// 用于替换需要 io.Writer 接口的其他组件
	// writer 复用，可以设置作为 io.Writer 输出的 prefix 和 level
	ginDefaultW := qelog.Clone()
	ginDefaultW.SetWriteLevel(zapcore.InfoLevel)
	ginDefaultW.SetWritePrefix("GinDefaultWriter")

	replaceGinLogger(ginDefaultW)
	// {"_level":"INFO","_time":1610616377873.9443,"_caller":"qezap/qezap.go:154","_func":"github.com/huzhongqing/qelog/qezap.(*Logger).Write","_short":"GinDefaultWriter","val":"gin out writer"}

	ginDefaultErrorW := qelog.Clone()
	ginDefaultErrorW.SetWriteLevel(zapcore.ErrorLevel)
	ginDefaultErrorW.SetWritePrefix("GinDefaultErrorWriter")

	replaceGinLogger(ginDefaultErrorW)

	// {"_level":"ERROR","_time":1610616377875.9512,"_caller":"qezap/qezap.go:154","_func":"github.com/huzhongqing/qelog/qezap.(*Logger).Write","_short":"GinDefaultErrorWriter","val":"gin out writer"}

	qelog.DPanic("last message")
	if err := qelog.Sync(); err != nil {
		fmt.Println(err)

	}
}

func replaceGinLogger(w io.Writer) {
	// 这里可以替换掉gin默认的输出文件
	// gin.DefaultWriter = w
	w.Write([]byte("gin out writer"))
}

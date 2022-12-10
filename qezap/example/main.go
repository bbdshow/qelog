package main

import (
	"context"
	"errors"
	"time"

	"github.com/bbdshow/qelog/qezap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	// 500MB rotate,keep 7days
	local := qezap.New(
		qezap.WithFilename("./log/local.log"),
		qezap.WithRotateMaxSizeAge(500<<20, 7*24*time.Hour))
	defer local.Close()

	local.Debug("Debug", zap.String("val", "only written local file"))

	// runtime change logger level
	local.SetEnabledLevel(zapcore.InfoLevel)
	local.Debug("Debug", zap.String("val", "this msg,not written to file"))

	// support remote storage
	multi := qezap.New(qezap.WithAddrsAndModuleName([]string{"127.0.0.1:31082"}, "demo"))
	defer multi.Close()

	multi.Info("local fs and remote storage will be written")

	// extension field
	// context bind traceID
	ctx := multi.WithTraceID(context.Background())
	// multi.FieldTraceID(ctx): get traceId from context, generate zap.Field written this log
	// admin manager can use traceId, find 'trace warn' and 'trace error' log
	multi.Warn("trace warn", zap.String("xx", "xx"), multi.FieldTraceID(ctx))
	multi.Error("trace error", zap.String("xx", "xx"), multi.FieldTraceID(ctx))

	// if we need manger filtering log, we provide multi condition field bind to log
	multi.Info("this msg as init filtering")
	multi.Info("this msg as init filtering", multi.ConditionOne("first condition"),
		multi.ConditionTwo("second condition"), multi.ConditionThree("third condition"))

	// if we need to io.Writer
	// gin.RecoveryWithWriter(ioWrite)
	ioWrite := multi.NewLevelWriter(zapcore.InfoLevel, "GIN logger output")
	ioWrite.Write([]byte("logger impl io.Writer, eg: gin.RecoveryWithWriter(ioWrite)"))

	// how to replace GO Logger
	// use zap.Sugar()
	slg := multi.Sugar()
	slg.Info("This info log")
	slg.Errorf("have err: %s", errors.New("mock error").Error())
}

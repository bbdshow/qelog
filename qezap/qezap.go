package qezap

import (
	"context"

	"go.uber.org/multierr"

	"github.com/huzhongqing/qelog/api/types"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var mode = Debug

const (
	Debug = iota
	Release
)

func EnableRelease() {
	mode = Release
}

type Logger struct {
	*zap.Logger
	cfg         *Config
	atomicLevel *zap.AtomicLevel

	WritePrefix string
	WriteLevel  zapcore.Level
}

func NewOneEncoderMultiWriterCore(enc zapcore.Encoder, level *zap.AtomicLevel, multiW []zapcore.WriteSyncer) *oneEncoderMultiWriter {
	return &oneEncoderMultiWriter{
		AtomicLevel: level,
		enc:         enc,
		multiW:      multiW,
	}
}

// 支持动态修改等级，一次编码，多处写入
type oneEncoderMultiWriter struct {
	*zap.AtomicLevel
	enc    zapcore.Encoder
	multiW []zapcore.WriteSyncer
}

func (mw *oneEncoderMultiWriter) With(fields []zap.Field) zapcore.Core {
	clone := mw.clone()
	for i := range fields {
		fields[i].AddTo(clone.enc)
	}
	return clone
}

func (mw *oneEncoderMultiWriter) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if mw.Enabled(ent.Level) {
		return ce.AddCore(ent, mw)
	}
	return ce
}

func (mw *oneEncoderMultiWriter) Write(ent zapcore.Entry, fields []zap.Field) error {
	buf, err := mw.enc.EncodeEntry(ent, fields)
	if err != nil {
		return err
	}

	for _, w := range mw.multiW {
		_, err = w.Write(buf.Bytes())
		if err != nil {
			err = multierr.Append(err, err)
			continue
		}
	}
	buf.Free()

	if ent.Level > zapcore.ErrorLevel {
		// Since we may be crashing the program, sync the output. Ignore Sync
		// errors, pending a clean solution to issue #370.
		mw.Sync()
	}
	return err
}

func (mw *oneEncoderMultiWriter) Sync() error {
	var err error
	for i := range mw.multiW {
		err = multierr.Append(err, mw.multiW[i].Sync())
	}
	return err
}

func (mw *oneEncoderMultiWriter) clone() *oneEncoderMultiWriter {
	return &oneEncoderMultiWriter{
		AtomicLevel: mw.AtomicLevel,
		enc:         mw.enc.Clone(),
		multiW:      mw.multiW,
	}
}

func New(cfg *Config, level zapcore.Level) *Logger {
	if err := cfg.Validate(); err != nil {
		panic(err)
	}
	atomicLevel := zap.NewAtomicLevelAt(level)

	var core zapcore.Core
	if mode == Release {
		// 一次编码 多次写入
		multiW := make([]zapcore.WriteSyncer, 0)
		multiW = append(multiW, NewWriteSync(cfg))

		if cfg.EnableRemote {
			multiW = append(multiW, NewWriteRemote(cfg))
		}
		core = NewOneEncoderMultiWriterCore(jsonEncoder(), &atomicLevel, multiW)
	} else {
		localW := NewWriteSync(cfg)
		localCore := zapcore.NewCore(consoleEncoder(), localW, &atomicLevel)
		cores := []zapcore.Core{localCore}
		if cfg.EnableRemote {
			cores = append(cores, zapcore.NewCore(jsonEncoder(), NewWriteRemote(cfg), &atomicLevel))
		}
		core = zapcore.NewTee(cores...)
	}

	return &Logger{cfg: cfg, atomicLevel: &atomicLevel,
		Logger: zap.New(core,
			zap.AddCaller(),
			zap.AddCallerSkip(2),
			zap.AddStacktrace(zap.DPanicLevel))}
}

// 可动态修改日志等级
func (log *Logger) SetEnabledLevel(lvl zapcore.Level) *Logger {
	log.atomicLevel.SetLevel(lvl)
	return log
}

// 暴露Write方法，用于替换使用  io.Writer 接口的地方
func (log *Logger) Write(b []byte) (n int, err error) {
	ec := log.Check(log.WriteLevel, log.WritePrefix)
	ec.Write(zap.String("val", string(b)))
	return len(b), nil
}

func (log *Logger) SetWriteLevel(lvl zapcore.Level) *Logger {
	log.WriteLevel = lvl
	return log
}

func (log *Logger) SetWritePrefix(s string) *Logger {
	log.WritePrefix = s
	return log
}

func (log *Logger) Clone() *Logger {
	return &Logger{
		Logger:      log.Logger,
		WritePrefix: "",
		WriteLevel:  0,
	}
}

func (log *Logger) ConditionOne(v string) zap.Field {
	return zap.String(types.EncoderConditionOneKey, v)
}

func (log *Logger) ConditionTwo(v string) zap.Field {
	return zap.String(types.EncoderConditionTwoKey, v)
}

func (log *Logger) ConditionThree(v string) zap.Field {
	return zap.String(types.EncoderConditionThreeKey, v)
}

func (log *Logger) WithTraceID(ctx context.Context) context.Context {
	return context.WithValue(ctx, types.EncoderTraceIDKey, types.NewTraceID())
}

func (log *Logger) Debug(msg string, fields ...zap.Field) {
	log.encoderWithCtx(zapcore.DebugLevel, nil, msg, fields...)
}

func (log *Logger) Info(msg string, fields ...zap.Field) {
	log.encoderWithCtx(zapcore.InfoLevel, nil, msg, fields...)
}

func (log *Logger) Warn(msg string, fields ...zap.Field) {
	log.encoderWithCtx(zapcore.WarnLevel, nil, msg, fields...)
}

func (log *Logger) Error(msg string, fields ...zap.Field) {
	log.encoderWithCtx(zapcore.ErrorLevel, nil, msg, fields...)
}

func (log *Logger) DPanic(msg string, fields ...zap.Field) {
	log.encoderWithCtx(zapcore.DPanicLevel, nil, msg, fields...)
}

func (log *Logger) Panic(msg string, fields ...zap.Field) {
	log.encoderWithCtx(zapcore.PanicLevel, nil, msg, fields...)
}

func (log *Logger) Fatal(msg string, fields ...zap.Field) {
	log.encoderWithCtx(zapcore.FatalLevel, nil, msg, fields...)
}

// 用于把上下文的一些信息打入日志
func (log *Logger) DebugWithCtx(ctx context.Context, msg string, fields ...zap.Field) {
	log.encoderWithCtx(zapcore.DebugLevel, ctx, msg, fields...)
}

func (log *Logger) InfoWithCtx(ctx context.Context, msg string, fields ...zap.Field) {
	log.encoderWithCtx(zapcore.InfoLevel, ctx, msg, fields...)
}

func (log *Logger) WarnWithCtx(ctx context.Context, msg string, fields ...zap.Field) {
	log.encoderWithCtx(zapcore.WarnLevel, ctx, msg, fields...)
}

func (log *Logger) ErrorWithCtx(ctx context.Context, msg string, fields ...zap.Field) {
	log.encoderWithCtx(zapcore.ErrorLevel, ctx, msg, fields...)
}

func (log *Logger) DPanicWithCtx(ctx context.Context, msg string, fields ...zap.Field) {
	log.encoderWithCtx(zapcore.DPanicLevel, ctx, msg, fields...)
}

func (log *Logger) PanicWithCtx(ctx context.Context, msg string, fields ...zap.Field) {
	log.encoderWithCtx(zapcore.PanicLevel, ctx, msg, fields...)
}

func (log *Logger) FatalWithCtx(ctx context.Context, msg string, fields ...zap.Field) {
	log.encoderWithCtx(zapcore.FatalLevel, ctx, msg, fields...)
}

func (log *Logger) MustGetTraceID(ctx context.Context) types.TraceID {
	val := ctx.Value(types.EncoderTraceIDKey)
	id, ok := val.(types.TraceID)
	if ok {
		return id
	}
	return types.NilTraceID
}

func (log *Logger) encoderWithCtx(level zapcore.Level, ctx context.Context, msg string, fields ...zap.Field) {
	if ctx != nil {
		id := log.MustGetTraceID(ctx)
		if !id.IsZero() {
			fields = append(fields, zap.String(types.EncoderTraceIDKey, id.Hex()))
		}
	}
	switch level {
	case zapcore.DebugLevel:
		log.Logger.Debug(msg, fields...)
	case zapcore.InfoLevel:
		log.Logger.Info(msg, fields...)
	case zapcore.WarnLevel:
		log.Logger.Warn(msg, fields...)
	case zapcore.ErrorLevel:
		log.Logger.Error(msg, fields...)
	case zapcore.DPanicLevel:
		log.Logger.DPanic(msg, fields...)
	case zapcore.PanicLevel:
		log.Logger.Panic(msg, fields...)
	case zapcore.FatalLevel:
		log.Logger.Fatal(msg, fields...)
	}
}

func (log *Logger) Config() *Config {
	return log.cfg
}

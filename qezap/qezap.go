package qezap

import (
	"context"

	"go.uber.org/multierr"

	"github.com/huzhongqing/qelog/api/types"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.Logger
	core        *oneEncoderMultiWriter
	WritePrefix string
	WriteLevel  zapcore.Level
}

func NewOneEncoderMultiWriterCore(enc zapcore.Encoder, level zap.AtomicLevel, multiW []zapcore.WriteSyncer) *oneEncoderMultiWriter {
	return &oneEncoderMultiWriter{
		AtomicLevel: level,
		enc:         enc,
		multiW:      multiW,
	}
}

// 支持动态修改等级，一次编码，多处写入
type oneEncoderMultiWriter struct {
	zap.AtomicLevel
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

func (mw *oneEncoderMultiWriter) SetEnabledLevel(lvl zapcore.Level) *oneEncoderMultiWriter {
	mw.AtomicLevel.SetLevel(lvl)
	return mw
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
	enc := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		MessageKey:       types.EncoderMessageKey,
		LevelKey:         types.EncoderLevelKey,
		TimeKey:          types.EncoderTimeKey,
		NameKey:          types.EncoderNameKey,
		CallerKey:        types.EncoderCallerKey,
		FunctionKey:      types.EncoderFunctionKey,
		StacktraceKey:    types.EncoderStacktraceKey,
		LineEnding:       "",
		EncodeLevel:      zapcore.CapitalLevelEncoder,
		EncodeTime:       zapcore.EpochMillisTimeEncoder,
		EncodeDuration:   zapcore.SecondsDurationEncoder,
		EncodeCaller:     zapcore.ShortCallerEncoder,
		ConsoleSeparator: zapcore.DefaultLineEnding,
	})

	multiW := make([]zapcore.WriteSyncer, 0)
	multiW = append(multiW, NewWriteSync(cfg))

	if cfg.EnableRemote {
		multiW = append(multiW, NewWriteRemote(cfg))
	}

	core := NewOneEncoderMultiWriterCore(enc, atomicLevel, multiW)

	return &Logger{core: core, Logger: zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.DPanicLevel))}
}

// 可动态修改日志等级
func (log *Logger) SetEnabledLevel(lvl zapcore.Level) *Logger {
	log.core.SetEnabledLevel(lvl)
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
		log.Debug(msg, fields...)
	case zapcore.InfoLevel:
		log.Info(msg, fields...)
	case zapcore.WarnLevel:
		log.Warn(msg, fields...)
	case zapcore.ErrorLevel:
		log.Error(msg, fields...)
	case zapcore.DPanicLevel:
		log.DPanic(msg, fields...)
	case zapcore.PanicLevel:
		log.Panic(msg, fields...)
	case zapcore.FatalLevel:
		log.Fatal(msg, fields...)
	}
}

package qezap

import (
	"context"
	"fmt"

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

	localW  *WriteSync
	remoteW *WriteRemote

	core            zapcore.Core
	ioWriterMessage string
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

func New(cfg *Config, level zapcore.Level, options ...zap.Option) *Logger {
	if err := cfg.Validate(); err != nil {
		panic(err)
	}
	atomicLevel := zap.NewAtomicLevelAt(level)

	log := &Logger{
		cfg:         cfg,
		atomicLevel: &atomicLevel,
		localW:      nil,
		remoteW:     nil,
	}

	if mode == Release {
		// 一次编码 多次写入
		multiW := make([]zapcore.WriteSyncer, 0)
		localW := NewWriteSync(cfg)
		log.localW = localW

		multiW = append(multiW, localW)

		if cfg.EnableRemote {
			remoteW := NewWriteRemote(cfg)
			log.remoteW = remoteW
			multiW = append(multiW, remoteW)
		}
		log.core = NewOneEncoderMultiWriterCore(jsonEncoder(), &atomicLevel, multiW)
	} else {
		localW := NewWriteSync(cfg)
		log.localW = localW
		localCore := zapcore.NewCore(consoleEncoder(), localW, &atomicLevel)
		cores := []zapcore.Core{localCore}
		if cfg.EnableRemote {
			remoteW := NewWriteRemote(cfg)
			log.remoteW = remoteW
			cores = append(cores, zapcore.NewCore(jsonEncoder(), remoteW, &atomicLevel))
		}
		log.core = zapcore.NewTee(cores...)
	}

	// 设置默认的 options, caller 设置最前面
	opts := make([]zap.Option, 0)
	opts = append(opts, zap.AddCaller())
	opts = append(opts, zap.AddStacktrace(zap.DPanicLevel))

	opts = append(opts, options...)

	log.Logger = zap.New(log.core, opts...)

	return log
}

type Writer struct {
	level   zapcore.Level
	message string
	log     *zap.Logger
}

func (w *Writer) Write(b []byte) (n int, err error) {
	if ce := w.log.Check(w.level, w.message); ce != nil {
		ce.Write(zap.String("content", string(b)))
	}
	return
}

func (log *Logger) NewWriter(level zapcore.Level, message string) *Writer {
	l := zap.New(log.core, zap.AddCaller(), zap.AddCallerSkip(1))
	return &Writer{
		level:   level,
		message: message,
		log:     l,
	}
}

// 可动态修改日志等级
func (log *Logger) SetEnabledLevel(lvl zapcore.Level) *Logger {
	log.atomicLevel.SetLevel(lvl)
	return log
}

func (log *Logger) ConditionOne(v interface{}) zap.Field {
	return ConditionOne(fmt.Sprintf("%v", v))
}

func (log *Logger) ConditionTwo(v interface{}) zap.Field {
	return ConditionTwo(fmt.Sprintf("%v", v))
}

func (log *Logger) ConditionThree(v interface{}) zap.Field {
	return ConditionThree(fmt.Sprintf("%v", v))
}

func (log *Logger) WithTraceID(ctx context.Context) context.Context {
	return WithTraceID(ctx)
}

func (log *Logger) FieldTraceID(ctx context.Context) zap.Field {
	return FieldTraceID(ctx)
}

func (log *Logger) TraceID(ctx context.Context) types.TraceID {
	return TraceID(ctx)
}

func (log *Logger) Config() *Config {
	return log.cfg
}

func (log *Logger) Close() error {
	var err error
	if log.localW != nil {
		err = multierr.Append(err, log.localW.Close())
	}
	if log.remoteW != nil {
		err = multierr.Append(err, log.remoteW.Close())
	}
	return err
}

func ConditionOne(v string) zap.Field {
	return zap.String(types.EncoderConditionOneKey, v)
}

func ConditionTwo(v string) zap.Field {
	return zap.String(types.EncoderConditionTwoKey, v)
}

func ConditionThree(v string) zap.Field {
	return zap.String(types.EncoderConditionThreeKey, v)
}

func WithTraceID(ctx context.Context) context.Context {
	return context.WithValue(ctx, types.EncoderTraceIDKey, types.NewTraceID())
}

func FieldTraceID(ctx context.Context) zap.Field {
	return zap.String(types.EncoderTraceIDKey, TraceID(ctx).Hex())
}

func TraceID(ctx context.Context) types.TraceID {
	val := ctx.Value(types.EncoderTraceIDKey)
	id, ok := val.(types.TraceID)
	if ok {
		return id
	}
	return types.NilTraceID
}

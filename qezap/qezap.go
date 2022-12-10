package qezap

import (
	"context"
	"errors"
	"fmt"

	"github.com/bbdshow/qelog/api/types"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.Logger
	opt *options
	// logger level
	lvl *zap.AtomicLevel
	// writer
	local  *WriteLocal
	remote *WriteRemote

	// zap core
	core zapcore.Core
}

// New wrap zap.Logger, impl multi writer
// local writer written to fs.
// remote writer, rely on GRPC or HTTP protocol written remote storage
// you can also only select local fs.
func New(opts ...Option) *Logger {
	lg := &Logger{opt: defaultOptions()}

	for _, opt := range opts {
		opt.apply(lg.opt)
	}

	initCoreWriter(lg)

	// setting zap options, caller first
	zapOpts := []zap.Option{zap.AddCaller(), zap.AddStacktrace(zap.DPanicLevel)}
	zapOpts = append(zapOpts, lg.opt.Zap...)

	lg.Logger = zap.New(lg.core, zapOpts...)

	return lg
}

func initCoreWriter(lg *Logger) {
	lvl := zap.NewAtomicLevelAt(lg.opt.Level)
	lg.lvl = &lvl

	multiWriter := make([]zapcore.WriteSyncer, 0)
	lg.local = NewWriteLocal(lg.opt.Local)
	multiWriter = append(multiWriter, lg.local)
	if lg.opt.EnableRemote {
		lg.remote = NewWriteRemote(lg.opt.Remote)
		multiWriter = append(multiWriter, lg.remote)
	}
	if lg.opt.Mode == ModeRelease {
		// one encoder
		lg.core = newOneEncoderMultiWriterCore(jsonEncoder(), lg.lvl, multiWriter)
		return
	}
	// two ways encoder
	localCore := zapcore.NewCore(consoleEncoder(), lg.local, lg.lvl)
	cores := []zapcore.Core{localCore}
	if lg.remote != nil {
		cores = append(cores, zapcore.NewCore(jsonEncoder(), lg.remote, lg.lvl))
	}
	lg.core = zapcore.NewTee(cores...)
}

// Close  will be sync logger data and close file handle.
func (lg *Logger) Close() error {
	var err error
	if lg.local != nil {
		err = multierr.Append(err, lg.local.Close())
	}
	if lg.remote != nil {
		err = multierr.Append(err, lg.remote.Close())
	}
	return err
}

// Writer Rely on zap.Logger to achieve IO write,use in io.Writer scenarios
// data will be written local | remote log file
type Writer struct {
	level zapcore.Level
	msg   string
	log   *zap.Logger
}

// Write rely on zap core write, actual data is written DATA field
func (w *Writer) Write(b []byte) (n int, err error) {
	str := string(b)
	if ce := w.log.Check(w.level, w.msg); ce != nil {
		ce.Write(zap.String("DATA", str))
		return len(b), nil
	}
	return 0, errors.New("not found Write, check level setting")
}

// NewLevelWriter  returns io.Writer impl
// fixed msg for remote storage index performance
func (lg *Logger) NewLevelWriter(level zapcore.Level, msg string) *Writer {
	l := zap.New(lg.core, zap.AddCaller(), zap.AddCallerSkip(1))
	return &Writer{
		level: level,
		msg:   msg,
		log:   l,
	}
}

// SetEnabledLevel runtime change logger level
func (lg *Logger) SetEnabledLevel(lvl zapcore.Level) *Logger {
	lg.lvl.SetLevel(lvl)
	return lg
}

// ConditionOne internal field extension, used for first condition filtering
func (lg *Logger) ConditionOne(v interface{}) zap.Field {
	return ConditionOne(fmt.Sprintf("%v", v))
}

// ConditionTwo internal field extension, used for second condition filtering
func (lg *Logger) ConditionTwo(v interface{}) zap.Field {
	return ConditionTwo(fmt.Sprintf("%v", v))
}

// ConditionThree internal field extension, used for third condition filtering
func (lg *Logger) ConditionThree(v interface{}) zap.Field {
	return ConditionThree(fmt.Sprintf("%v", v))
}

// FieldTraceID internal field extension, used for trace context filtering
func (lg *Logger) FieldTraceID(ctx context.Context) zap.Field {
	return FieldTraceID(ctx)
}

// WithTraceID generate traceID(bson._id), set value in context,relation context
func (lg *Logger) WithTraceID(ctx context.Context) context.Context {
	return WithTraceID(ctx)
}

// TraceID get traceId from context
func (lg *Logger) TraceID(ctx context.Context) types.TraceID {
	return TraceID(ctx)
}

// ConditionOne wrap internal field
func ConditionOne(v string) zap.Field {
	return zap.String(types.EncoderConditionOneKey, v)
}

// ConditionTwo wrap internal field
func ConditionTwo(v string) zap.Field {
	return zap.String(types.EncoderConditionTwoKey, v)
}

// ConditionThree wrap internal field
func ConditionThree(v string) zap.Field {
	return zap.String(types.EncoderConditionThreeKey, v)
}

// FieldTraceID wrap internal field
func FieldTraceID(ctx context.Context) zap.Field {
	return zap.String(types.EncoderTraceIDKey, TraceID(ctx).Hex())
}

// WithTraceID traceId setting in context
func WithTraceID(ctx context.Context) context.Context {
	return context.WithValue(ctx, types.EncoderTraceIDKey, types.NewTraceID())
}

// TraceID from context
func TraceID(ctx context.Context) types.TraceID {
	val := ctx.Value(types.EncoderTraceIDKey)
	id, ok := val.(types.TraceID)
	if ok {
		return id
	}
	return types.NilTraceID
}

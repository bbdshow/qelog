package qezap

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path"
	"strconv"
	"sync/atomic"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var _logFilename = "./log/logger.log"

type Config struct {
	WriteSync WriteSyncConfig

	EnableRemote bool
	WriteRemote  WriteRemoteConfig
}

func NewConfig(addrs []string, moduleName string) *Config {
	cfg := &Config{
		WriteSync:    NewWriteSyncConfig(_logFilename),
		EnableRemote: true,
		WriteRemote:  NewWriteRemoteConfig(addrs, moduleName),
	}
	return cfg
}

func (cfg *Config) SetFilename(filename string) *Config {
	dir := path.Dir(filename)
	cfg.WriteSync.Filename = filename
	cfg.WriteRemote.RemoteFailedBackup = path.Join(dir, "backup", "backup.log")
	return cfg
}

func (cfg *Config) SetEnableRemote(enable bool) *Config {
	cfg.EnableRemote = enable
	return cfg
}

func (cfg *Config) SetHTTPTransport() *Config {
	cfg.WriteRemote.Transport = "http"
	return cfg
}

type Logger struct {
	*zap.Logger
	WritePrefix string
	WriteLevel  zapcore.Level
}

func New(cfg *Config, level zapcore.Level) *Logger {
	if err := cfg.WriteSync.Validate(); err != nil {
		panic(err)
	}
	if cfg.EnableRemote {
		if err := cfg.WriteRemote.Validate(); err != nil {
			panic(err)
		}
	}

	prodEncCfg := zap.NewProductionEncoderConfig()
	prodEncCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	localEnc := zapcore.NewConsoleEncoder(prodEncCfg)
	localCore := zapcore.NewCore(localEnc, NewWriteSync(cfg.WriteSync), level)

	var core zapcore.Core

	if cfg.EnableRemote {
		remoteEnc := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
			MessageKey:       "_short",
			LevelKey:         "_level",
			TimeKey:          "_time",
			NameKey:          "_name",
			CallerKey:        "_caller",
			FunctionKey:      "_func",
			StacktraceKey:    "_stack",
			LineEnding:       "",
			EncodeLevel:      zapcore.CapitalLevelEncoder,
			EncodeTime:       zapcore.EpochMillisTimeEncoder,
			EncodeDuration:   zapcore.SecondsDurationEncoder,
			EncodeCaller:     zapcore.ShortCallerEncoder,
			ConsoleSeparator: zapcore.DefaultLineEnding,
		})

		remoteCore := zapcore.NewCore(remoteEnc, NewWriteRemote(cfg.WriteRemote), level)

		core = zapcore.NewTee(localCore, remoteCore)
	} else {
		core = localCore
	}

	return &Logger{Logger: zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.DPanicLevel))}
}

// 暴露Write方法，用于替换使用  io.Writer 接口的地方
func (log *Logger) Write(b []byte) (n int, err error) {
	ec := log.Check(log.WriteLevel, log.WritePrefix)
	ec.Write(zap.String("write_value", string(b)))
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
	return zap.String("_condition1", v)
}

func (log *Logger) ConditionTwo(v string) zap.Field {
	return zap.String("_condition2", v)
}

func (log *Logger) ConditionThree(v string) zap.Field {
	return zap.String("_condition3", v)
}

func (log *Logger) WithTraceID(ctx context.Context) context.Context {
	return context.WithValue(ctx, "_traceid", new(TraceID).New())
}

func (log *Logger) TraceIDField(ctx context.Context) zap.Field {
	id := ""
	val := ctx.Value("_traceid")
	tid, ok := val.(TraceID)
	if ok {
		id = tid.String()
	}
	return zap.String("_traceid", id)
}

var _pidString = func() string {
	pid := os.Getpid()
	return fmt.Sprintf("%05d", pid)
}()

var _incInt64 int64 = 0

type TraceID string

// [nsec:19]
func (tid TraceID) New() TraceID {
	buff := bytes.Buffer{}
	nsec := time.Now().UnixNano()
	nsecStr := strconv.FormatInt(nsec, 10)

	buff.WriteString(nsecStr)
	buff.WriteString(_pidString)
	buff.WriteString(strconv.FormatInt(atomic.AddInt64(&_incInt64, 1), 10))
	return TraceID(buff.String())
}

func (tid TraceID) Time() time.Time {
	if tid != "" {
		nsec, _ := strconv.ParseInt(string(tid[:19]), 10, 64)
		return time.Unix(0, nsec)
	}
	return time.Unix(0, 0)
}

func (tid TraceID) String() string {
	return string(tid)
}

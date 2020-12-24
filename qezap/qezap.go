package qezap

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	WriteSync WriteSyncConfig

	EnableRemote bool
	WriteRemote  WriteRemoteConfig
}

func NewConfig(addrs []string, moduleName string) *Config {
	defaultFilename := "./log/qelogger.log"
	cfg := &Config{
		WriteSync:    NewWriteSyncConfig(defaultFilename),
		EnableRemote: true,
		WriteRemote:  NewWriteRemoteConfig(addrs, moduleName),
	}
	return cfg
}

func (cfg *Config) SetFilename(filename string) *Config {
	cfg.WriteSync.Filename = filename
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

type Condition struct {
	key string
}

// 生成一个可以条件查询的字段
func NewCondition() *Condition {
	return &Condition{}
}

func (c *Condition) setKey(k string) *Condition {
	c.key = k
	return c
}
func (c *Condition) One() *Condition {
	return c.setKey("_condition1")
}
func (c *Condition) Two() *Condition {
	return c.setKey("_condition2")
}
func (c *Condition) Three() *Condition {
	return c.setKey("_condition3")
}
func (c *Condition) StringFiled(val string) zap.Field {
	if c.key == "" {
		c.One()
	}
	return zap.String(c.key, val)
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
	ec.Write(zap.Any("write_value", b))
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

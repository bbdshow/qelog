package qzzap

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
	defaultFilename := "./log/qzlogger.log"
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
	return c.setKey("_c1")
}
func (c *Condition) Two() *Condition {
	return c.setKey("_c2")
}
func (c *Condition) Three() *Condition {
	return c.setKey("_c3")
}
func (c *Condition) StringFiled(val string) zap.Field {
	if c.key == "" {
		c.One()
	}
	return zap.String(c.key, val)
}

func NewWrapZap(cfg *Config, level zapcore.Level) *zap.Logger {
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
			EncodeName:       zapcore.FullNameEncoder,
			ConsoleSeparator: zapcore.DefaultLineEnding,
		})

		remoteCore := zapcore.NewCore(remoteEnc, NewWriteRemote(cfg.WriteRemote), level)

		core = zapcore.NewTee(localCore, remoteCore)
	} else {
		core = localCore
	}

	return zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.DPanicLevel))
}

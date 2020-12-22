package wrapzap

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

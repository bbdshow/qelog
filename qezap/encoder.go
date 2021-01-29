package qezap

import (
	"github.com/huzhongqing/qelog/api/types"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func jsonEncoder() zapcore.Encoder {
	return zapcore.NewJSONEncoder(zapcore.EncoderConfig{
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
}

func consoleEncoder() zapcore.Encoder {
	return zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
}

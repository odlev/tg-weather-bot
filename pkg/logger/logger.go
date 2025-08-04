// Package logger is a nice package
package logger

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func SetupLogger() *zap.SugaredLogger {
	config := zap.Config{
		Level: zap.NewAtomicLevelAt(zapcore.DebugLevel),
		Encoding: "console",
		OutputPaths: []string{"stdout"},

		EncoderConfig: zapcore.EncoderConfig{
			TimeKey: "time",
			LevelKey: "level",
			MessageKey: "msg",
			EncodeLevel: zapcore.CapitalColorLevelEncoder,
			EncodeTime: customTimeEncoder,
			LineEnding: "\n",
			CallerKey: zapcore.OmitKey,
			FunctionKey: "function",
		},
	}

	log, err := config.Build()
	if err != nil {
		panic(err)
	}
	sugar := log.Sugar()
	return sugar
}

func customTimeEncoder(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
	encoder.AppendString(t.Format("2006/01/02 15:04:05"))
}

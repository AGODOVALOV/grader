package zaplogger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZLoggerConfig struct {
	Level    string
	Encoding string
	Name     string
}

func NewLogger(c ZLoggerConfig) (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()
	cfg.Level = zap.NewAtomicLevelAt(parseZapLevel(c.Level))
	cfg.DisableStacktrace = true
	cfg.Encoding = c.Encoding
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	l, err := cfg.Build(
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zap.ErrorLevel),
	)
	if err != nil {
		return nil, err
	}
	return l, nil
}

func parseZapLevel(l string) zapcore.Level {
	level, err := zapcore.ParseLevel(l)
	if err != nil {
		fmt.Println(err)
		return zapcore.InfoLevel
	}
	return level
}

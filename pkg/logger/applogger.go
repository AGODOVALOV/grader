package logger

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	"github.com/AGODOVALOV/grader/pkg/logger/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	zaplogger "github.com/AGODOVALOV/grader/pkg/logger/zap"
)

const logKey = "logger"

type AppLogger struct {
	Logger *zap.Logger
	level  string
}

var (
	GLogger          AppLogger
	onceLoggerCreate sync.Once
)

func NewAppLogger(cfg config.Config) (Logger, error) {
	var (
		err error
		l   AppLogger
	)
	onceLoggerCreate.Do(func() {
		l.Logger, err = zaplogger.NewLogger(zaplogger.ZLoggerConfig{
			Level:    cfg.Level,
			Encoding: cfg.Encoding,
			Name:     cfg.Name,
		})
		l.level = cfg.Level
	})
	GLogger = l
	return l, err
}

type Logger interface {
	Debug(ctx context.Context, action string, message string, extraFields ...map[string]string)

	Info(ctx context.Context, action string, message string, extraFields ...map[string]string)

	Warn(ctx context.Context, action string, message string, extraFields ...map[string]string)

	Error(ctx context.Context, action string, message string, extraFields ...map[string]string)

	Fatal(ctx context.Context, action string, message string, extraFields ...map[string]string)

	Panic(ctx context.Context, action string, message string, extraFields ...map[string]string)

	GetLogger() *zap.Logger
}

func (z AppLogger) Debug(ctx context.Context, action string, message string, extraFields ...map[string]string) {
	var ef map[string]string

	if len(extraFields) > 0 {
		ef = extraFields[0]
	}

	defer z.Logger.Sync()
	fields := z.getExtraFields(ctx, ef)

	z.Logger.Debug(getMessage(ctx, action, message), fields...)
}

func (z AppLogger) Info(ctx context.Context, action string, message string, extraFields ...map[string]string) {
	var ef map[string]string

	if len(extraFields) > 0 {
		ef = extraFields[0]
	}

	defer z.Logger.Sync()
	fields := z.getExtraFields(ctx, ef)

	z.Logger.Info(getMessage(ctx, action, message), fields...)
}

func (z AppLogger) Warn(ctx context.Context, action string, message string, extraFields ...map[string]string) {
	var ef map[string]string

	if len(extraFields) > 0 {
		ef = extraFields[0]
	}

	defer z.Logger.Sync()
	fields := z.getExtraFields(ctx, ef)

	z.Logger.Warn(getMessage(ctx, action, message), fields...)
}

func (z AppLogger) Error(ctx context.Context, action string, message string, extraFields ...map[string]string) {
	var ef map[string]string

	if len(extraFields) > 0 {
		ef = extraFields[0]
	}

	defer z.Logger.Sync()
	fields := z.getExtraFields(ctx, ef)

	z.Logger.Error(getMessage(ctx, action, message), fields...)
}

func (z AppLogger) Fatal(ctx context.Context, action string, message string, extraFields ...map[string]string) {
	var ef map[string]string

	if len(extraFields) > 0 {
		ef = extraFields[0]
	}

	defer z.Logger.Sync()
	fields := z.getExtraFields(ctx, ef)

	z.Logger.Fatal(getMessage(ctx, action, message), fields...)
}

func (z AppLogger) Panic(ctx context.Context, action string, message string, extraFields ...map[string]string) {
	var ef map[string]string

	if len(extraFields) > 0 {
		ef = extraFields[0]
	}

	defer z.Logger.Sync()
	fields := z.getExtraFields(ctx, ef)

	z.Logger.Panic(getMessage(ctx, action, message), fields...)
}

func getMessage(_ context.Context, action, message string) string {
	return fmt.Sprintf("%s - %s", action, message)
}

func (z AppLogger) getExtraFields(_ context.Context, extraFields map[string]string) []zapcore.Field {
	var zapFields []zapcore.Field

	for k, v := range extraFields {
		zapFields = append(
			zapFields,
			zap.String(k, v),
		)
	}

	return zapFields
}

func (z AppLogger) GetLogger() *zap.Logger {
	return z.Logger
}

func Z(ctx context.Context) Logger {
	if ctx == nil {
		return GLogger
	}

	ctxLogger, ok := ctx.Value(logKey).(Logger)
	if !ok || ctxLogger == nil {
		return GLogger
	}
	return ctxLogger
}

func CtxWithLogger(ctx context.Context, l Logger) context.Context {
	return context.WithValue(ctx, logKey, l)
}

func NewCtxLoggerWithExtraFields(ctx context.Context, extra map[string]string) context.Context {
	if len(extra) == 0 {
		return ctx
	}

	newAppLogger := AppLogger{}
	fields := make([]zap.Field, 0, len(extra))

	currentLogger := Z(ctx).GetLogger()

	fields = AppLogger.getExtraFields(newAppLogger, ctx, extra)

	newZapLogger := currentLogger.With(fields...)

	newCtx := CtxWithLogger(ctx, AppLogger{
		Logger: newZapLogger,
		level:  strconv.Itoa(int(currentLogger.Level())),
	})

	return newCtx
}

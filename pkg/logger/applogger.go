// Package logger provides logging functionality.
package logger

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/AGODOVALOV/grader/pkg/logger/config"
	zaplogger "github.com/AGODOVALOV/grader/pkg/logger/zap"
)

type contextKeyLogger string

const logKey contextKeyLogger = "logger"

var (
	// GLogger is a global logger.
	//nolint:gochecknoglobals // global variable is ok
	GLogger *AppLogger

	//nolint:gochecknoglobals // global variable is ok
	onceLoggerCreate sync.Once
)

// AppLogger is a struct that represents the logger.
type AppLogger struct {
	logger *zap.Logger
	level  string
}

// Logger defines a logging interface with support for various log levels and contextual data.
// Debug writes debug-level logs with additional contextual information.
// Info writes informational logs with a specified action and message.
// Warn records warning-level logs with optional extra fields for context.
// Error outputs error-level logs for significant issues in the application.
// Fatal logs a critical issue and terminates the application.
// Panic logs a critical issue and triggers a panic.
// GetLogger retrieves the underlying zap.Logger used by the Logger implementation.
type Logger interface {
	Debug(ctx context.Context, action, message string, extraFields ...map[string]string)

	Info(ctx context.Context, action, message string, extraFields ...map[string]string)

	Warn(ctx context.Context, action, message string, extraFields ...map[string]string)

	Error(ctx context.Context, action, message string, extraFields ...map[string]string)

	Fatal(ctx context.Context, action, message string, extraFields ...map[string]string)

	Panic(ctx context.Context, action, message string, extraFields ...map[string]string)

	GetLogger() *zap.Logger
}

// NewAppLogger creates a new logger.
func NewAppLogger(cfg config.Config) (*AppLogger, error) {
	var (
		err error
		l   *AppLogger
	)

	l = &AppLogger{}

	onceLoggerCreate.Do(func() {
		l.logger, err = zaplogger.NewLogger(zaplogger.ZLoggerConfig{
			Level:    cfg.Level,
			Encoding: cfg.Encoding,
			Name:     cfg.Name,
		})
		l.level = cfg.Level
	})

	if err != nil {
		return nil, err
	}

	GLogger = l
	return l, nil
}

// Debug logs a debug message.
func (z AppLogger) Debug(ctx context.Context, action, message string, extraFields ...map[string]string) {
	var ef map[string]string

	if len(extraFields) > 0 {
		ef = extraFields[0]
	}

	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			//fmt.Println(err)
		}
	}(z.logger)
	fields := z.getExtraFields(ctx, ef)

	z.logger.Debug(getMessage(ctx, action, message), fields...)
}

// Info logs an info message.
func (z AppLogger) Info(ctx context.Context, action, message string, extraFields ...map[string]string) {
	var ef map[string]string

	if len(extraFields) > 0 {
		ef = extraFields[0]
	}

	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			//fmt.Println(err)
		}
	}(z.logger)
	fields := z.getExtraFields(ctx, ef)

	z.logger.Info(getMessage(ctx, action, message), fields...)
}

// Warn logs a warning message.
func (z AppLogger) Warn(ctx context.Context, action, message string, extraFields ...map[string]string) {
	var ef map[string]string

	if len(extraFields) > 0 {
		ef = extraFields[0]
	}

	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			//fmt.Println(err)
		}
	}(z.logger)
	fields := z.getExtraFields(ctx, ef)

	z.logger.Warn(getMessage(ctx, action, message), fields...)
}

// Error logs an error message.
func (z AppLogger) Error(ctx context.Context, action, message string, extraFields ...map[string]string) {
	var ef map[string]string

	if len(extraFields) > 0 {
		ef = extraFields[0]
	}

	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			//fmt.Println(err)
		}
	}(z.logger)
	fields := z.getExtraFields(ctx, ef)

	z.logger.Error(getMessage(ctx, action, message), fields...)
}

// Fatal logs a fatal message.
func (z AppLogger) Fatal(ctx context.Context, action, message string, extraFields ...map[string]string) {
	var ef map[string]string

	if len(extraFields) > 0 {
		ef = extraFields[0]
	}

	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			//fmt.Println(err)
		}
	}(z.logger)
	fields := z.getExtraFields(ctx, ef)

	z.logger.Fatal(getMessage(ctx, action, message), fields...)
}

// Panic logs a panic message.
func (z AppLogger) Panic(ctx context.Context, action, message string, extraFields ...map[string]string) {
	var ef map[string]string

	if len(extraFields) > 0 {
		ef = extraFields[0]
	}

	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			fmt.Println(err)
		}
	}(z.logger)
	fields := z.getExtraFields(ctx, ef)

	z.logger.Panic(getMessage(ctx, action, message), fields...)
}

func getMessage(_ context.Context, action, message string) string {
	return fmt.Sprintf("%s - %s", action, message)
}

func (AppLogger) getExtraFields(_ context.Context, extraFields map[string]string) []zapcore.Field {
	zapFields := make([]zapcore.Field, 0, len(extraFields))

	for k, v := range extraFields {
		zapFields = append(
			zapFields,
			zap.String(k, v),
		)
	}

	return zapFields
}

// GetLogger returns the logger.
func (z AppLogger) GetLogger() *zap.Logger {
	return z.logger
}

// Z returns the logger from the context.
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

// CtxWithLogger returns a new context with the logger.
func CtxWithLogger(ctx context.Context, l Logger) context.Context {
	return context.WithValue(ctx, logKey, l)
}

// NewCtxLoggerWithExtraFields returns a new context with the logger and extra fields.
func NewCtxLoggerWithExtraFields(ctx context.Context, extra map[string]string) context.Context {
	if len(extra) == 0 {
		return ctx
	}

	newAppLogger := AppLogger{}

	currentLogger := Z(ctx).GetLogger()

	fields := newAppLogger.getExtraFields(ctx, extra)

	newZapLogger := currentLogger.With(fields...)

	newCtx := CtxWithLogger(ctx, AppLogger{
		logger: newZapLogger,
		level:  strconv.Itoa(int(currentLogger.Level())),
	})

	return newCtx
}

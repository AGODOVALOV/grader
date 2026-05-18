package logger_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/AGODOVALOV/grader/pkg/logger"
	"github.com/AGODOVALOV/grader/pkg/logger/config"
)

func TestRequestScopeLogger(_ *testing.T) {
	// create logger
	z, err := logger.NewAppLogger(config.Config{
		Level:    "debug",
		Encoding: "json",
		Name:     "test logger",
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	// ctx
	ctx, cancel := context.WithCancel(context.Background())
	// ctx add logger
	ctx = logger.CtxWithLogger(ctx, z)
	defer cancel()

	newCtx := logger.NewCtxLoggerWithExtraFields(ctx, map[string]string{"messageID": "msq1"})

	logger.Z(ctx).Info(ctx, "test without extra", "test without extra")

	logger.Z(newCtx).Info(ctx, "test with extra", "test with extra")

	logger.Z(newCtx).Info(ctx, "test with extra", "test with extra", map[string]string{"extraField": "ExtraValue"})
}

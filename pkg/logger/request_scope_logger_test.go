package logger

import (
	"context"
	"fmt"
	"testing"

	"github.com/AGODOVALOV/grader/pkg/logger/config"
)

func TestRequestScopeLogger(t *testing.T) {
	// create logger
	z, err := NewAppLogger(config.Config{
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
	ctx = CtxWithLogger(ctx, z)
	defer cancel()

	newCtx := NewCtxLoggerWithExtraFields(ctx, map[string]string{"messageID": "msq1"})

	Z(ctx).Info(ctx, "test without extra", "test without extra")

	Z(newCtx).Info(ctx, "test with extra", "test with extra")

	Z(newCtx).Info(ctx, "test with extra", "test with extra", map[string]string{"extraField": "ExtraValue"})
}

package rate_limiter

import (
	"context"

	"github.com/AGODOVALOV/grader/pkg/rate_limiter/config"
	"github.com/AGODOVALOV/grader/pkg/rate_limiter/fixed_window"
	"github.com/AGODOVALOV/grader/pkg/rate_limiter/token_bucket"
)

const (
	LimiterFixedWindow = "fixed_window"
	LimiterTokenBucket = "token_bucket"
)

type Limiter interface {
	Allow() bool
}

type RateLimiter struct {
	Limiter
	cfg *config.Config
}

func NewRateLimiter(ctx context.Context, cfg *config.Config) *RateLimiter {
	var l Limiter

	switch cfg.Type {
	case LimiterFixedWindow:
		l = fixed_window.NewFixedWindowLimiter(ctx, cfg)
	case LimiterTokenBucket:
		l = token_bucket.NewTokenBucketLimiter(ctx, cfg)
	}

	return &RateLimiter{
		Limiter: l,
		cfg:     cfg,
	}
}

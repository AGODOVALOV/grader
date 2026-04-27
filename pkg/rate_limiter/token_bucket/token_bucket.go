package token_bucket

import (
	"context"
	"time"

	"github.com/AGODOVALOV/grader/pkg/rate_limiter/config"
)

type TokenBucketLimiter struct {
	tokenCh chan struct{}
}

func NewTokenBucketLimiter(ctx context.Context, cfg *config.Config) *TokenBucketLimiter {
	limiter := &TokenBucketLimiter{
		tokenCh: make(chan struct{}, cfg.MaxRequests),
	}

	for i := int32(0); i < cfg.MaxRequests; i++ {
		limiter.tokenCh <- struct{}{}
	}

	refuelInterval := cfg.Interval / time.Duration(cfg.MaxRequests)
	if refuelInterval <= 0 {
		refuelInterval = time.Nanosecond
	}

	go limiter.refuelPeriodic(ctx, refuelInterval)

	return limiter
}

func (l *TokenBucketLimiter) refuelPeriodic(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			select {
			case l.tokenCh <- struct{}{}:
			default:
			}
		}
	}
}

func (l *TokenBucketLimiter) Allow() bool {
	select {
	case <-l.tokenCh:
		return true
	default:
		return false
	}
}

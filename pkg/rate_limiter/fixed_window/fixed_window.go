package fixed_window

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/AGODOVALOV/grader/pkg/rate_limiter/config"
)

type FixedWindowLimiter struct {
	count int32
	limit int32
}

func NewFixedWindowLimiter(ctx context.Context, cfg *config.Config) *FixedWindowLimiter {
	limiter := &FixedWindowLimiter{
		count: 0,
		limit: cfg.MaxRequests,
	}

	go limiter.startPeriodCountRefresh(ctx, cfg.Interval)

	return limiter
}

func (l *FixedWindowLimiter) startPeriodCountRefresh(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			atomic.StoreInt32(&l.count, 0)
		}
	}
}

func (l *FixedWindowLimiter) Allow() bool {
	count := atomic.LoadInt32(&l.count)

	if count > int32(l.limit) {
		return false
	}

	for !atomic.CompareAndSwapInt32(&l.count, count, count+1) {
		count = atomic.LoadInt32(&l.count)
	}

	return count < int32(l.limit)
}

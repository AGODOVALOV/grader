package rate_limiter_test

import (
	"context"
	"testing"
	"time"

	"github.com/AGODOVALOV/grader/pkg/rate_limiter"
	"github.com/stretchr/testify/require"

	"github.com/AGODOVALOV/grader/pkg/rate_limiter/config"
)

func TestFixedWindowLimiter_Allow(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var (
		limiter rate_limiter.Limiter
		err     error
	)

	cfg := &config.Config{
		MaxRequests: 5,
		Interval:    1 * time.Second,
		Type:        rate_limiter.LimiterFixedWindow,
	}

	t.Run("Allow Good", func(t *testing.T) {
		limiter = rate_limiter.NewRateLimiter(ctx, cfg)
		require.NoError(t, err)
		for range 5 {
			require.True(t, limiter.Allow())
		}
	})

	t.Run("Allow Bad", func(t *testing.T) {
		limiter = rate_limiter.NewRateLimiter(ctx, cfg)
		require.NoError(t, err)
		for range 5 {
			require.True(t, limiter.Allow())
		}
		require.False(t, limiter.Allow())
	})

	t.Run("Allow Bad Then True", func(t *testing.T) {
		limiter = rate_limiter.NewRateLimiter(ctx, cfg)
		require.NoError(t, err)
		for range 5 {
			require.True(t, limiter.Allow())
		}
		require.False(t, limiter.Allow())

		time.Sleep(1 * time.Second)
		require.True(t, limiter.Allow())
	})
}

func TestTokenBucketLimiter_Allow(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var (
		limiter rate_limiter.Limiter
		err     error
	)

	cfg := &config.Config{
		MaxRequests: 5,
		Interval:    1 * time.Second,
		Type:        rate_limiter.LimiterTokenBucket,
	}

	t.Run("Allow Good", func(t *testing.T) {
		limiter = rate_limiter.NewRateLimiter(ctx, cfg)
		require.NoError(t, err)
		for range 5 {
			require.True(t, limiter.Allow())
		}
	})

	t.Run("Allow Bad", func(t *testing.T) {
		limiter = rate_limiter.NewRateLimiter(ctx, cfg)
		require.NoError(t, err)
		for range 5 {
			require.True(t, limiter.Allow())
		}
		require.False(t, limiter.Allow())
	})

	t.Run("Allow Bad Then True", func(t *testing.T) {
		limiter = rate_limiter.NewRateLimiter(ctx, cfg)
		require.NoError(t, err)
		for range 5 {
			require.True(t, limiter.Allow())
		}
		require.False(t, limiter.Allow())

		time.Sleep(1 * time.Second)
		require.True(t, limiter.Allow())
	})
}

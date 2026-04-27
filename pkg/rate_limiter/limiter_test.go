package rate_limiter

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/AGODOVALOV/grader/pkg/rate_limiter/config"
)

func TestFixedWindowLimiter_Allow(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var (
		limiter Limiter
		err     error
	)

	cfg := &config.Config{
		MaxRequests: 5,
		Interval:    1 * time.Second,
		Type:        LimiterFixedWindow,
	}

	t.Run("Allow Good", func(t *testing.T) {
		limiter = NewRateLimiter(ctx, cfg)
		require.NoError(t, err)
		for range 5 {
			require.Equal(t, true, limiter.Allow())
		}
	})

	t.Run("Allow Bad", func(t *testing.T) {
		limiter = NewRateLimiter(ctx, cfg)
		require.NoError(t, err)
		for range 5 {
			require.Equal(t, true, limiter.Allow())
		}
		require.Equal(t, false, limiter.Allow())
	})

	t.Run("Allow Bad Then True", func(t *testing.T) {
		limiter = NewRateLimiter(ctx, cfg)
		require.NoError(t, err)
		for range 5 {
			require.Equal(t, true, limiter.Allow())
		}
		require.Equal(t, false, limiter.Allow())

		time.Sleep(1 * time.Second)
		require.Equal(t, true, limiter.Allow())
	})
}

func TestTokenBucketLimiter_Allow(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var (
		limiter Limiter
		err     error
	)

	cfg := &config.Config{
		MaxRequests: 5,
		Interval:    1 * time.Second,
		Type:        LimiterTokenBucket,
	}

	t.Run("Allow Good", func(t *testing.T) {
		limiter = NewRateLimiter(ctx, cfg)
		require.NoError(t, err)
		for range 5 {
			require.Equal(t, true, limiter.Allow())
		}
	})

	t.Run("Allow Bad", func(t *testing.T) {
		limiter = NewRateLimiter(ctx, cfg)
		require.NoError(t, err)
		for range 5 {
			require.Equal(t, true, limiter.Allow())
		}
		require.Equal(t, false, limiter.Allow())
	})

	t.Run("Allow Bad Then True", func(t *testing.T) {
		limiter = NewRateLimiter(ctx, cfg)
		require.NoError(t, err)
		for range 5 {
			require.Equal(t, true, limiter.Allow())
		}
		require.Equal(t, false, limiter.Allow())

		time.Sleep(1 * time.Second)
		require.Equal(t, true, limiter.Allow())
	})
}

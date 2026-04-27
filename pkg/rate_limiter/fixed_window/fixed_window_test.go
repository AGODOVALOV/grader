package fixed_window

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

	var limiter *FixedWindowLimiter

	cfg := &config.Config{
		MaxRequests: 5,
		Interval:    1 * time.Second,
	}

	t.Run("Allow Good", func(t *testing.T) {
		limiter = NewFixedWindowLimiter(ctx, cfg)
		for range 5 {
			require.Equal(t, true, limiter.Allow())
		}
	})

	t.Run("Allow Bad", func(t *testing.T) {
		limiter = NewFixedWindowLimiter(ctx, cfg)
		for range 5 {
			require.Equal(t, true, limiter.Allow())
		}
		require.Equal(t, false, limiter.Allow())
	})

	t.Run("Allow Bad Then True", func(t *testing.T) {
		limiter = NewFixedWindowLimiter(ctx, cfg)
		for range 5 {
			require.Equal(t, true, limiter.Allow())
		}
		require.Equal(t, false, limiter.Allow())

		time.Sleep(1 * time.Second)
		require.Equal(t, true, limiter.Allow())
	})
}

package rate_limiter

import (
	"context"
	"sync"
	"time"

	"github.com/AGODOVALOV/grader/pkg/logger"
	messagingConfig "github.com/AGODOVALOV/grader/pkg/queue/config"
	"github.com/AGODOVALOV/grader/pkg/rate_limiter/config"
)

type LimiterManager struct {
	limiters        map[string]Limiter
	messagingConfig *messagingConfig.MessagingConfig
	mu              *sync.RWMutex
}

func NewLimiterManager(cfg *messagingConfig.MessagingConfig) *LimiterManager {
	return &LimiterManager{
		limiters:        make(map[string]Limiter),
		mu:              &sync.RWMutex{},
		messagingConfig: cfg,
	}
}

func (lm *LimiterManager) Get(ctx context.Context, queueKey string) Limiter {
	lm.mu.RLock()
	limiter, ok := lm.limiters[queueKey]
	lm.mu.RUnlock()

	if ok {
		return limiter
	}

	lm.mu.Lock()
	defer lm.mu.Unlock()

	if limiter, ok = lm.limiters[queueKey]; ok {
		return limiter
	}

	limiterCfg := lm.getLimiterCfg(ctx, queueKey)
	if limiterCfg == nil {
		return nil
	}

	limiter = NewRateLimiter(ctx, limiterCfg)
	lm.limiters[queueKey] = limiter

	return limiter
}

func (lm *LimiterManager) getLimiterCfg(ctx context.Context, queueKey string) *config.Config {
	for _, v := range lm.messagingConfig.Channels {
		if v.Name == queueKey {
			return &v.RateLimiter
		}
	}

	// use default params
	logger.Z(ctx).Warn(
		ctx,
		"rate_limiter.getlimitercfg",
		"no limiter config found for queue - use default params",
		map[string]string{"queue": queueKey})

	return &config.Config{
		MaxRequests: 100,
		Interval:    1 * time.Second,
	}
}

// Package config contains client server config
package config

import (
	"time"

	"github.com/AGODOVALOV/grader/pkg/rate_limiter/config"
)

// Config represents the configuration for the HTTP server.
type Config struct {
	Host            string         `mapstructure:"host"             validate:"required,hostname|ip"`
	Port            int            `mapstructure:"port"             validate:"required,gte=1,lte=65535"`
	ReadTimeout     time.Duration  `mapstructure:"read-timeout"     validate:"required,gt=0"`
	WriteTimeout    time.Duration  `mapstructure:"write-timeout"    validate:"required,gt=0"`
	IdleTimeout     time.Duration  `mapstructure:"idle-timeout"     validate:"required,gt=0"`
	ShutdownTimeout time.Duration  `mapstructure:"shutdown-timeout" validate:"required,gt=0"`
	RateLimiter     *config.Config `mapstructure:"rate_limiter"     validate:"required"`
}

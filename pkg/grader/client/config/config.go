package config

import "time"

type Config struct {
	URL                string        `mapstructure:"callback-url"      validate:"required,url"`
	Timeout            time.Duration `mapstructure:"timeout"           validate:"required,gt=0"`
	MaxIdleConnections int           `mapstructure:"max-idle-conns"    validate:"required,min=1"`
	IdleConnTimeout    time.Duration `mapstructure:"idle-conn-timeout" validate:"required,gt=0"`
	TaskDelay          time.Duration `mapstructure:"task-delay"        validate:"required,gt=0"`
	Retry              RetryCallback `mapstructure:"retry"             validate:"required"`
	JWT                JWT           `mapstructure:"token"             validate:"required"`
}

type RetryCallback struct {
	MaxAttempts int           `mapstructure:"max-attempts" validate:"required,gte=1"`
	Backoff     time.Duration `mapstructure:"backoff"      validate:"required,gt=0"`
	BackoffMax  time.Duration `mapstructure:"backoff-max"  validate:"required,gt=0"`
}

type JWT struct {
	JWTSecret string `mapstructure:"secret"   validate:"required,min=32"`
	Duration  string `mapstructure:"duration" validate:"required"`
}

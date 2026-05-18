package config

import "time"

type Config struct {
	MaxRequests int32         `mapstructure:"max-requests" validate:"required,gte=1"`
	Interval    time.Duration `mapstructure:"interval"     validate:"required,gte=1s"`
	Type        string        `mapstructure:"type"         validate:"required,oneof=fixed_window token_bucket"`
}

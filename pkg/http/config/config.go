package config

import "time"

type Config struct {
	Host            string        `mapstructure:"host"             validate:"required,hostname|ip"`
	Port            int           `mapstructure:"port"             validate:"required,gte=1,lte=65535"`
	ReadTimeout     time.Duration `mapstructure:"read_timeout"     validate:"required,gt=0"`
	WriteTimeout    time.Duration `mapstructure:"write_timeout"    validate:"required,gt=0"`
	IdleTimeout     time.Duration `mapstructure:"idle_timeout"     validate:"required,gt=0"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout" validate:"required,gt=0"`
}

package config

import "time"

type Config struct {
	URL       string       `mapstructure:"url" validate:"required,url"`
	Container string       `mapstructure:"container" validate:"required"`
	Workers   int          `mapstructure:"workers" validate:"required,gte=1"`
	Server    ServerConfig `mapstructure:"server" validate:"required"`
}

type ServerConfig struct {
	Host            string        `mapstructure:"host"             validate:"required,hostname|ip"`
	Port            int           `mapstructure:"port"             validate:"required,gte=1,lte=65535"`
	ReadTimeout     time.Duration `mapstructure:"read-timeout"     validate:"required,gt=0"`
	WriteTimeout    time.Duration `mapstructure:"write-timeout"    validate:"required,gt=0"`
	IdleTimeout     time.Duration `mapstructure:"idle-timeout"     validate:"required,gt=0"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown-timeout" validate:"required,gt=0"`
}

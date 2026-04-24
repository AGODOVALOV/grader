// Package config contains the configuration for the database.
package config

import "time"

// Config represents the configuration for the database.
type Config struct {
	Host     string `mapstructure:"host"      validate:"required,hostname|ip"`
	Port     int    `mapstructure:"port"      validate:"required,gte=1,lte=65535"`
	User     string `mapstructure:"user"      validate:"required"`
	Password string `mapstructure:"password"  validate:"required"`
	DBName   string `mapstructure:"db-name"   validate:"required"`
	SSLMode  string `mapstructure:"ssl-mode"  validate:"required,oneof=disable require verify-ca verify-full"`
	TimeZone string `mapstructure:"time-zone" validate:"required,timezone"`
	Pool     struct {
		MaxOpenConns    int           `mapstructure:"max-open-conns"    validate:"required,gte=1"`
		MaxIdleConns    int           `mapstructure:"max-idle-conns"    validate:"required,gte=1"`
		ConnMaxLifetime time.Duration `mapstructure:"conn-max-lifetime" validate:"required,gte=1s"`
		ConnMaxIdleTime time.Duration `mapstructure:"conn-max-idletime" validate:"required,gte=1s"`
	}
}

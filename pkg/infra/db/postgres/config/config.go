package config

import "time"

type Config struct {
	Host     string `mapstructure:"host"      validate:"required,hostname|ip"`
	Port     int    `mapstructure:"port"      validate:"required,gte=1,lte=65535"`
	User     string `mapstructure:"user"      validate:"required"`
	Password string `mapstructure:"password"  validate:"required"`
	DBName   string `mapstructure:"db_name"   validate:"required"`
	SSLMode  string `mapstructure:"ssl_mode"  validate:"required,oneof=disable require verify-ca verify-full"`
	TimeZone string `mapstructure:"time_zone" validate:"required,timezone"`
	Pool     struct {
		MaxOpenConns    int           `mapstructure:"max_open_conns" validate:"required,gte=1"`
		MaxIdleConns    int           `mapstructure:"max_idle_conns" validate:"required,gte=1"`
		ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime" validate:"required,gte=1s"`
		ConnMaxIdleTime time.Duration `mapstructure:"conn_max_idle_time" validate:"required,gte=1s"`
	} `                                                                                         yaml:"pool"`
}

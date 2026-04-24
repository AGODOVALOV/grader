// Package config contains the configuration for the token.
package config

// Config represents the configuration for the token.
type Config struct {
	JWTSecret string `mapstructure:"secret"   validate:"required,min=32"`
	Duration  string `mapstructure:"duration" validate:"required"`
}

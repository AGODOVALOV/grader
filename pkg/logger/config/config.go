// Package config provides configuration for the logger.
package config

// Config represents the configuration for the logger.
type Config struct {
	Level    string `mapstructure:"level"    validate:"required,oneof=debug info warn error"`
	Encoding string `mapstructure:"encoding"`
	Name     string `mapstructure:"name"`
}

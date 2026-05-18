// Package config contains the configuration for the S3 storage.
package config

// Config represents the configuration for the S3 storage.
type Config struct {
	Endpoint  string `mapstructure:"endpoint"   validate:"required"`
	AccessKey string `mapstructure:"access-key" validate:"required"`
	SecretKey string `mapstructure:"secret-key" validate:"required"`
	Bucket    string `mapstructure:"bucket"     validate:"required"`
	UseSSL    bool   `mapstructure:"use-ssl"`
}

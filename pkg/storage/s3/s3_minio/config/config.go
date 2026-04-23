package config

type Config struct {
	Endpoint  string `mapstructure:"endpoint" validate:"required"`
	AccessKey string `mapstructure:"access_key" validate:"required"`
	SecretKey string `mapstructure:"secret_key" validate:"required"`
	Bucket    string `mapstructure:"bucket" validate:"required"`
	UseSSL    bool   `mapstructure:"use_ssl"`
}

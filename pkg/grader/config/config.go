package config

type Config struct {
	URL       string `mapstructure:"url" validate:"required,url"`
	Container string `mapstructure:"container" validate:"required"`
}

package config

type Config struct {
	JWTSecret string `mapstructure:"secret" validate:"required,min=32"`
	Duration  string `mapstructure:"duration" validate:"required"`
}

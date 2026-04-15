package config

type Config struct {
	Level    string `mapstructure:"level"    validate:"required,oneof=debug info warn error"`
	Encoding string `mapstructure:"encoding"`
	Name     string `mapstructure:"name"`
}

package loader

import (
	"fmt"

	"github.com/spf13/viper"

	"github.com/AGODOVALOV/grader/pkg/config/config"
)

type ConfigLoader struct {
	Viper *viper.Viper
	path  string
	name  string
}

func NewConfigLoader(p string, n string) *ConfigLoader {
	return &ConfigLoader{
		Viper: viper.New(),
		path:  p,
		name:  n,
	}
}

func (loader *ConfigLoader) Load() (*config.Config, error) {
	loader.Viper.AddConfigPath(loader.path)
	loader.Viper.SetConfigName(loader.name)
	loader.Viper.SetConfigType("yml")
	loader.Viper.SetDefault("log.level", "info")
	loader.Viper.SetDefault("log.encoding", "json")

	if err := loader.Viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config: %w", err)
	}

	var cfg config.Config

	if err := loader.Viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to decode into struct: %w", err)
	}

	return &cfg, nil
}

func LoadConfig(path string, name string) (*config.Config, error) {
	v := viper.New()
	v.AddConfigPath(path)
	v.SetConfigName(name)
	v.SetConfigType("yml")
	v.SetDefault("log.level", "info")
	v.SetDefault("log.encoding", "json")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config: %w", err)
	}

	var cfg config.Config

	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to decode into struct: %w", err)
	}

	return &cfg, nil
}

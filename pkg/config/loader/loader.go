// Package loader provides functionality for loading configuration files.
package loader

import (
	"fmt"

	"github.com/spf13/viper"

	"github.com/AGODOVALOV/grader/pkg/config/config"
)

// ConfigLoader is responsible for loading and unmarshalling configuration files into a structured configuration object.
type ConfigLoader struct {
	Viper *viper.Viper
	path  string
	name  string
}

// NewConfigLoader creates a new ConfigLoader instance.
func NewConfigLoader(p, n string) *ConfigLoader {
	return &ConfigLoader{
		Viper: viper.New(),
		path:  p,
		name:  n,
	}
}

// Load loads and unmarshals the configuration file into a structured configuration object.
func (loader *ConfigLoader) Load() (*config.Config, error) {
	loader.Viper.AddConfigPath(loader.path)
	loader.Viper.SetConfigName(loader.name)
	loader.Viper.SetConfigType("yml")
	loader.Viper.SetDefault("log.level", "info")
	loader.Viper.SetDefault("log.encoding", "json")

	err := loader.Viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("error reading config: %w", err)
	}

	var cfg config.Config

	err = loader.Viper.Unmarshal(&cfg)
	if err != nil {
		return nil, fmt.Errorf("unable to decode into struct: %w", err)
	}

	return &cfg, nil
}

// LoadConfig reads a configuration file from the given path and name, unmarshals it into a Config struct, and returns it.
func LoadConfig(path, name string) (*config.Config, error) {
	v := viper.New()
	v.AddConfigPath(path)
	v.SetConfigName(name)
	v.SetConfigType("yml")
	v.SetDefault("log.level", "info")
	v.SetDefault("log.encoding", "json")

	err := v.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("error reading config: %w", err)
	}

	var cfg config.Config

	err = v.Unmarshal(&cfg)
	if err != nil {
		return nil, fmt.Errorf("unable to decode into struct: %w", err)
	}

	return &cfg, nil
}

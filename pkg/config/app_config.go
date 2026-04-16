// Package config create application config
package config

import (
	"context"
	"os"

	"github.com/AGODOVALOV/grader/pkg/config/config"
	"github.com/AGODOVALOV/grader/pkg/config/loader"
	"github.com/AGODOVALOV/grader/pkg/config/service/validator"
	"github.com/AGODOVALOV/grader/pkg/config/service/watcher"
)

const (
	envConfigPath     = "GRADER_CONFIG_PATH"
	envConfigFileName = "GRADER_CONFIG_FILE_NAME"
	defaultPath       = "./"
	defaultFileName   = "config.yml"
)

// Config manages the application's configuration, including loading, validating, watching, and reloading configuration data.
type Config struct {
	appConfig *config.Config
	loader    *loader.ConfigLoader
	validator *validator.Validator
	watcher   *watcher.Watcher
}

// NewConfig creates a new Config instance.
func NewConfig(path, name string) *Config {
	return &Config{
		appConfig: &config.Config{},
		loader:    loader.NewConfigLoader(path, name),
		validator: validator.NewValidator(),
		watcher:   watcher.NewWatcher(),
	}
}

// Load loads the application configuration from the specified path and name.
func (cfg *Config) Load() error {
	appConfig, err := cfg.loader.Load()
	if err != nil {
		return err
	}

	err = cfg.validator.Validate(appConfig)
	if err != nil {
		return err
	}

	cfg.appConfig = appConfig
	return nil
}

// ReLoad reloads the application configuration from the specified path and name.
func (cfg *Config) ReLoad() (*config.Config, error) {
	newCfg := &config.Config{}
	err := cfg.loader.Viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	err = cfg.loader.Viper.Unmarshal(newCfg)
	if err != nil {
		return nil, err
	}

	err = cfg.validator.Validate(newCfg)
	if err != nil {
		return nil, err
	}

	return newCfg, nil
}

// GetConfig returns the application configuration.
func (cfg *Config) GetConfig() *config.Config {
	return cfg.appConfig
}

// Watch watches for changes in the application configuration and notifies the provided channel when a change is detected.
func (cfg *Config) Watch(ctx context.Context, newConfigSignal chan<- struct{}) {
	cfg.watcher.Watch(ctx, cfg.loader.Viper, newConfigSignal)
}

// GetApplicationConfig returns the application configuration.
func GetApplicationConfig() (*Config, error) {
	configPath := os.Getenv(envConfigPath)
	if configPath == "" {
		configPath = defaultPath
	}

	configFileName := os.Getenv(envConfigFileName)
	if configFileName == "" {
		configFileName = defaultFileName
	}

	cfg := NewConfig(configPath, configFileName)

	err := cfg.Load()
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

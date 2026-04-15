package config

import (
	"context"
	"os"

	"github.com/AGODOVALOV/grader/pkg/config/config"
	appCfg "github.com/AGODOVALOV/grader/pkg/config/config"
	"github.com/AGODOVALOV/grader/pkg/config/loader"
	"github.com/AGODOVALOV/grader/pkg/config/service/validator"
	"github.com/AGODOVALOV/grader/pkg/config/service/watcher"
)

const (
	ENVConfigPath     = "QPROC_CONFIG_PATH"
	ENVConfigFileName = "QPROC_CONFIG_FILE_NAME"
	defaultPath       = "./"
	defaultFileName   = "config.yml"
)

type Config struct {
	appConfig *config.Config
	loader    *loader.ConfigLoader
	validator *validator.Validator
	watcher   *watcher.Watcher
}

func NewConfig(path string, name string) *Config {
	return &Config{
		appConfig: &config.Config{},
		loader:    loader.NewConfigLoader(path, name),
		validator: validator.NewValidator(),
		watcher:   watcher.NewWatcher(),
	}
}

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

func (cfg *Config) ReLoad() (*appCfg.Config, error) {
	newCfg := &appCfg.Config{}
	if err := cfg.loader.Viper.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := cfg.loader.Viper.Unmarshal(newCfg); err != nil {
		return nil, err
	}

	err := cfg.validator.Validate(newCfg)
	if err != nil {
		return nil, err
	}

	return newCfg, nil
}

func (cfg *Config) GetConfig() *config.Config {
	return cfg.appConfig
}

func (cfg *Config) Watch(ctx context.Context, newConfigSignal chan<- struct{}) {
	cfg.watcher.Watch(ctx, cfg.loader.Viper, newConfigSignal)
}

func GetApplicationConfig() (*Config, error) {
	configPath := os.Getenv(ENVConfigPath)
	if configPath == "" {
		configPath = defaultPath
	}

	configFileName := os.Getenv(ENVConfigFileName)
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

// Package config provides configuration for the application.
package config

import (
	client "github.com/AGODOVALOV/grader/pkg/client/config"
	dbconfig "github.com/AGODOVALOV/grader/pkg/infra/db/postgres/config"
	loggerconfig "github.com/AGODOVALOV/grader/pkg/logger/config"
)

// Config represents the top-level configuration structure for the application.
type Config struct {
	App struct {
		Name string `mapstructure:"name" validate:"required"`
	} `mapstructure:"app" validate:"required"`

	Log loggerconfig.Config `mapstructure:"log" validate:"required"`

	WebServer client.Config `mapstructure:"web-server" validate:"required"`

	DB dbconfig.Config `mapstructure:"postgres" validate:"required"`
}

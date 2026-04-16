// Package config provides configuration for the application.
package config

import (
	loggerconfig "github.com/AGODOVALOV/grader/pkg/logger/config"
	webconfig "github.com/AGODOVALOV/grader/pkg/webserver/config"
)

// Config represents the top-level configuration structure for the application.
type Config struct {
	App struct {
		Name string `mapstructure:"name" validate:"required"`
	} `mapstructure:"app" validate:"required"`

	Log loggerconfig.Config `mapstructure:"log" validate:"required"`

	WebServer webconfig.Config `mapstructure:"web-server" validate:"required"`
}

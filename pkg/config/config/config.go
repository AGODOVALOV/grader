package config

import (
	webServerConfig "github.com/AGODOVALOV/grader/pkg/http/config"
	loggerConfig "github.com/AGODOVALOV/grader/pkg/logger/config"
)

type Config struct {
	App struct {
		Name string `mapstructure:"name" validate:"required"`
	} `mapstructure:"app" validate:"required"`

	Log loggerConfig.Config `mapstructure:"log" validate:"required"`

	WebServer webServerConfig.Config `mapstructure:"web_server" validate:"required"`

	//MsgQueue queue.Config `mapstructure:"message_queue" validate:"required"`
	//
	//WP workerpool.Config `mapstructure:"worker_pool" validate:"required"`
	//
	//Client client.Config `mapstructure:"http_client" validate:"required"`
	//
	//Server server.Config `mapstructure:"http_server" validate:"required"`
	//
	//Metrics metrics.Config `mapstructure:"metrics" validate:"required"`
}

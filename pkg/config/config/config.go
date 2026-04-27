// Package config provides configuration for the application.
package config

import (
	client "github.com/AGODOVALOV/grader/pkg/client/config"
	graderconfig "github.com/AGODOVALOV/grader/pkg/grader/config"
	loggerconfig "github.com/AGODOVALOV/grader/pkg/logger/config"
	messagequeue "github.com/AGODOVALOV/grader/pkg/queue/config"
	dbconfig "github.com/AGODOVALOV/grader/pkg/storage/db/postgres/config"
	filestorageconfig "github.com/AGODOVALOV/grader/pkg/storage/s3/s3minio/config"
	tokenconfig "github.com/AGODOVALOV/grader/pkg/token/config"
)

// Config represents the top-level configuration structure for the application.
type Config struct {
	App struct {
		Name string `mapstructure:"name" validate:"required"`
	} `mapstructure:"app" validate:"required"`

	Log loggerconfig.Config `mapstructure:"log" validate:"required"`

	WebServer client.Config `mapstructure:"web-server" validate:"required"`

	DB dbconfig.Config `mapstructure:"postgres" validate:"required"`

	FileStorage filestorageconfig.Config `mapstructure:"s3-minio" validate:"required"`

	Token tokenconfig.Config `mapstructure:"token" validate:"required"`

	Grader graderconfig.Config `mapstructure:"grader" validate:"required"`

	MsgQueue messagequeue.Config `mapstructure:"message_queue" validate:"required"`
}

// Package config contains the configuration for the queue.
package config

import ratelimiterconfig "github.com/AGODOVALOV/grader/pkg/rate_limiter/config"

// Config contains the configuration for the queue.
type Config struct {
	Broker    BrokerConfig    `mapstructure:"broker"    validate:"required"`
	Messaging MessagingConfig `mapstructure:"messaging" validate:"required"`
}

// BrokerConfig represents the configuration for the messaging broker and its associated parameters.
type BrokerConfig struct {
	Active string       `mapstructure:"active"   validate:"required,oneof=rabbit kafka postgres"`
	Rabbit RabbitConfig `mapstructure:"rabbit"   validate:"required_if=Active rabbit"`
}

// RabbitConfig represents the configuration for RabbitMQ.
type RabbitConfig struct {
	URL          string `mapstructure:"url"           validate:"required,amqpuri"`
	ExchangeType string `mapstructure:"exchange_type" validate:"required,oneof=direct topic"`
	Durable      bool   `mapstructure:"durable"       validate:"required"`
	Prefetch     int    `mapstructure:"prefetch"      validate:"required,gte=1"`
}

// MessagingConfig represents the configuration for messaging, containing one or more queue message channels.
type MessagingConfig struct {
	Channels []QueueMsgChannel `mapstructure:"channels" validate:"required,min=1,dive"`
}

// QueueMsgChannel represents a channel for sending messages to a queue.
type QueueMsgChannel struct {
	Name        string                   `mapstructure:"name"         validate:"required,alphanum"`
	Target      string                   `mapstructure:"target"       validate:"required,url"`
	APIPath     string                   `mapstructure:"api_path"     validate:"required,startswith=/"`
	RateLimiter ratelimiterconfig.Config `mapstructure:"rate_limiter"`
}

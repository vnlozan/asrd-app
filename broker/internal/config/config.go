package config

import (
	"broker/internal/utils"
)

type MailerConfig struct {
	ConnectionURL string
}

type LoggerConfig struct {
	ConnectionURL string
}

type AuthConfig struct {
	ConnectionURL string
}

type RabbitMQConfig struct {
	ConnectionURL string
}

type Config struct {
	MailerConfig   *MailerConfig
	LoggerConfig   *LoggerConfig
	AuthConfig     *AuthConfig
	RabbitMQConfig *RabbitMQConfig
	Port           string
}

func NewConfig() *Config {
	return &Config{
		Port:           utils.String("PORT", "8080"),
		RabbitMQConfig: &RabbitMQConfig{ConnectionURL: utils.String("RABBITMQ_URL", "amqp://guest:guest@rabbitmq")},
		MailerConfig:   &MailerConfig{ConnectionURL: utils.String("MAILER_URL", "http://mailer:8080")},
		LoggerConfig:   &LoggerConfig{ConnectionURL: utils.String("LOGGER_URL", "http://logger:8080")},
		AuthConfig:     &AuthConfig{ConnectionURL: utils.String("AUTH_URL", "http://auth:8080")},
	}
}

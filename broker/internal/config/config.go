package config

import "broker/internal/utils"

type Config struct {
	Port      string
	MailerURL string
	LoggerURL string
	AuthURL   string
}

func NewConfig() *Config {
	return &Config{
		Port:      utils.String("PORT", "8080"),
		MailerURL: utils.String("MAILER_URL", "http://mailer:8080"),
		LoggerURL: utils.String("LOGGER_URL", "http://logger:8080"),
		AuthURL:   utils.String("AUTH_URL", "http://auth:8080"),
	}
}

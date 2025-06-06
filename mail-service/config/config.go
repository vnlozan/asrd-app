package config

import "mailer-service/internal/utils"

type Config struct {
	Port string
}

func NewConfig() *Config {
	return &Config{
		Port: utils.String("PORT", "8080"),
	}
}

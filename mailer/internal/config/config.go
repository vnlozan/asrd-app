package config

import (
	"mailer/internal/utils"
)

type MailerClient struct {
	Domain      string
	Host        string
	Port        int
	Username    string
	Password    string
	Encryption  string
	FromName    string
	FromAddress string
}

type Config struct {
	Port         string
	MailerClient MailerClient
}

func NewConfig() *Config {
	return &Config{
		Port: utils.String("PORT", "8080"),
		MailerClient: MailerClient{
			Domain:      utils.String("MAIL_DOMAIN", "localhost"),
			Host:        utils.String("MAIL_HOST", "mailhog"),
			Port:        utils.Int("MAIL_PORT", 1025),
			Username:    utils.String("MAIL_USERNAME", "none"),
			Password:    utils.String("MAIL_PASSWORD", ""),
			Encryption:  utils.String("MAIL_ENCRYPTION", ""),
			FromName:    utils.String("FROM_NAME", "John Smith"),
			FromAddress: utils.String("FROM_ADDRESS", "john.smith@example.com"),
		},
	}
}

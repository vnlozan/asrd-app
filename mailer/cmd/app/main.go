package main

import (
	"mailer/internal/config"
	"mailer/internal/controller"
	"mailer/internal/repo/client"
	"mailer/internal/server"
	"mailer/internal/service"
)

func main() {
	config := config.NewConfig()
	mailerClient := client.NewMailerClient(
		config.MailerClient.Domain,
		config.MailerClient.Host,
		config.MailerClient.Port,
		config.MailerClient.Username,
		config.MailerClient.Password,
		config.MailerClient.Encryption,
		config.MailerClient.FromName,
		config.MailerClient.FromAddress,
	)
	mailerService := service.NewMailerService(mailerClient)
	mailController := controller.NewMailController(mailerService)

	server := server.NewServer(config, mailController)
	server.Start()
}

package main

import (
	"mailer/config"
	"mailer/internal/controller"
	"mailer/internal/repo/client"
	"mailer/internal/server"
)

func main() {
	mailer := client.NewMailerClient()
	mailController := controller.MailController{Mailer: &mailer}

	server := server.NewServer(
		*config.NewConfig(),
		mailController,
	)
	server.Start()
}

package main

import (
	"mailer-service/config"
	"mailer-service/internal/controller"
	"mailer-service/internal/repo/client"
	"mailer-service/internal/server"
)

const webPort = "8080"

func main() {
	mailer := client.NewMailerClient()
	mailController := controller.MailController{Mailer: &mailer}

	server := server.NewServer(
		config.Config{Port: webPort},
		mailController,
	)
	server.Start()
}

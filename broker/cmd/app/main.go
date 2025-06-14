package main

import (
	"broker/internal/config"
	"broker/internal/controller"
	"broker/internal/infra/rabbitmq"
	"broker/internal/server"
	"broker/internal/service"
	"log"
)

func main() {
	config := config.NewConfig()

	conn, err := rabbitmq.Dial(config.RabbitMQConfig.ConnectionURL)
	if err != nil {
		log.Fatalf("cannot start: %v", err)
	}
	defer conn.Close()

	emitter := rabbitmq.NewEmitter(conn)
	brokerService := service.NewBrokerService(config, conn, emitter)
	brokerController := controller.NewBrokerController(brokerService)

	server := server.NewServer(config, brokerController)
	server.Start()
}

package rabbitmq

import (
	"log"
	"math"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Dial tries to connect to RabbitMQ with exponential back-off
// and returns an established *amqp.Connection.
func Dial(url string) (*amqp.Connection, error) {
	var attempt int64
	for {
		conn, err := amqp.Dial(url)
		if err == nil {
			log.Println("connected to RabbitMQ")
			return conn, nil
		}

		if attempt > 5 {
			log.Printf("RabbitMQ not ready after %d attempts: %v", attempt, err)
			return nil, err
		}

		sleep := time.Duration(math.Pow(float64(attempt+1), 2)) * time.Second
		log.Printf("RabbitMQ not ready, backing off %s...", sleep)
		time.Sleep(sleep)
		attempt++
	}
}

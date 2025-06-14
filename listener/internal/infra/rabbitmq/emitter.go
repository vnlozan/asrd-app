package rabbitmq

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Emitter публикует события в RabbitMQ.
type Emitter struct {
	conn *amqp.Connection
}

// NewEmitter создаёт новый Emitter.
func NewEmitter(conn *amqp.Connection) *Emitter {
	return &Emitter{conn: conn}
}

// SetupChannel гарантирует, что exchange объявлен.
func (e *Emitter) SetupChannel() error {
	ch, err := e.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	return declareExchange(ch)
}

// Publish публикует текстовое сообщение event в routingKey (severity).
func (e *Emitter) Publish(event, routingKey string) error {
	ch, err := e.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	log.Printf("publish to exchange=%s rk=%s", exchangeName, routingKey)

	return ch.Publish(
		exchangeName,
		routingKey,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(event),
		},
	)
}

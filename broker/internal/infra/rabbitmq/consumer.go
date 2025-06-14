package rabbitmq

import (
	"context"
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Consumer читает сообщения из RabbitMQ и передаёт их обработчику.
type Consumer struct {
	conn *amqp.Connection
}

// NewConsumer создаёт новый Consumer.
func NewConsumer(conn *amqp.Connection) *Consumer {
	return &Consumer{conn: conn}
}

// SetupChannel гарантирует, что exchange объявлен.
func (c *Consumer) SetupChannel() error {
	ch, err := c.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	return declareExchange(ch)
}

// Listen подписывается на topics и вызывает handler для каждого сообщения.
// Метод блокирует текущую горутину, пока ctx не будет отменён.
func (c *Consumer) Listen(ctx context.Context, topics []string, handler func(Payload)) error {
	ch, err := c.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := declareRandomQueue(ch)
	if err != nil {
		return err
	}

	for _, rk := range topics {
		if err = ch.QueueBind(
			q.Name,
			rk,
			exchangeName,
			false,
			nil,
		); err != nil {
			return err
		}
	}

	msgs, err := ch.Consume(
		q.Name,
		"",    // consumer name
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,
	)
	if err != nil {
		return err
	}

	log.Printf("waiting for messages [exchange=%s queue=%s]", exchangeName, q.Name)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case d := <-msgs:
			var p Payload
			if err := json.Unmarshal(d.Body, &p); err == nil {
				go handler(p)
			}
		}
	}
}

package rabbitmq

import amqp "github.com/rabbitmq/amqp091-go"

const exchangeName = "logs_topic"
const exchangeKind = "topic"

func declareExchange(ch *amqp.Channel) error {
	return ch.ExchangeDeclare(
		exchangeName,
		exchangeKind,
		true,  // durable
		false, // auto-delete
		false, // internal
		false, // no-wait
		nil,   // args
	)
}

func declareRandomQueue(ch *amqp.Channel) (amqp.Queue, error) {
	return ch.QueueDeclare(
		"",    // let the server generate a name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // args
	)
}

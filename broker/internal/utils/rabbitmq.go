package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	amqp "github.com/rabbitmq/amqp091-go"
)

/* EVENT */
func declareExchange(ch *amqp.Channel) error {
	return ch.ExchangeDeclare(
		"logs_topic", // name
		"topic",      // type
		true,         // durable?
		false,        // auto-deleted?
		false,        // internal?
		false,        // no wait?
		nil,          // arguments?
	)
}

func declareRandomQueue(ch *amqp.Channel) (amqp.Queue, error) {
	return ch.QueueDeclare(
		"",    // name?
		false, // durable?
		false, // delete when unused?
		true,  // exclusive?
		false, // no wait?
		nil,   // arguments?
	)
}

/* CONSUMER */
type Payload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

type Consumer struct {
	rabbitMQConn *amqp.Connection
}

func NewConsumer(conn *amqp.Connection) *Consumer {
	return &Consumer{rabbitMQConn: conn}
}

func (c *Consumer) SetupChannel() error {
	channel, err := c.rabbitMQConn.Channel()
	if err != nil {
		return err
	}
	return declareExchange(channel)
}

func (c *Consumer) Listen(topics []string) error {
	ch, err := c.rabbitMQConn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := declareRandomQueue(ch)
	if err != nil {
		return err
	}

	for _, s := range topics {
		err = ch.QueueBind(
			q.Name,
			s,
			"logs_topic",
			false,
			nil,
		)

		if err != nil {
			return err
		}
	}

	messages, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	forever := make(chan bool)
	go func() {
		for d := range messages {
			var payload Payload
			_ = json.Unmarshal(d.Body, &payload)

			go handlePayload(payload)
		}
	}()

	log.Printf("Waiting for message [Exchange, Queue] [logs_topic, %s]\n", q.Name)

	<-forever

	return nil
}

func handlePayload(payload Payload) {
	switch payload.Name {
	case "log", "event":
		err := logEvent(payload)
		if err != nil {
			log.Println(err)
		}
	case "auth":
		//
	default:
		err := logEvent(payload)
		if err != nil {
			log.Println(err)
		}
	}
}

func logEvent(entry Payload) error {
	jsonData, _ := json.MarshalIndent(entry, "", "\t")

	logServiceURL := fmt.Sprintf("%s/log", "http://logger:8080") // TODO

	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		return err
	}

	return nil
}

/* EMITTER */
type Emitter struct {
	connection *amqp.Connection
}

func NewEventEmitter(conn *amqp.Connection) *Emitter {
	return &Emitter{connection: conn}
}

func (e *Emitter) SetupChannel() error {
	channel, err := e.connection.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	return declareExchange(channel)
}

func (e *Emitter) PublishToChannel(event string, severity string) error {
	channel, err := e.connection.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	log.Println("Pushing to channel")

	err = channel.Publish(
		"logs_topic",
		severity,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(event),
		},
	)
	if err != nil {
		return err
	}
	return nil
}

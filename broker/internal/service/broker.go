package service

import (
	"broker/internal/config"
	"broker/internal/dto"
	"broker/internal/infra/rabbitmq"
	"broker/internal/utils"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/rabbitmq/amqp091-go"
)

type BrokerService struct {
	config             *config.Config
	rabbitMQConnection *amqp091.Connection
	eventEmitter       *rabbitmq.Emitter
}

func NewBrokerService(config *config.Config, rabbitMQConnection *amqp091.Connection, eventEmitter *rabbitmq.Emitter) IBrokerService {
	return &BrokerService{config: config, rabbitMQConnection: rabbitMQConnection, eventEmitter: eventEmitter}
}

func (s *BrokerService) LogItem(entry dto.LogPayload) error {
	err := s.pushToQueue(entry.Name, entry.Data)
	if err != nil {
		return err
	}

	return nil
}

func (s *BrokerService) Authenticate(credentials dto.AuthPayload) (any, error) {
	jsonData, _ := json.MarshalIndent(credentials, "", "\t")

	request, err := http.NewRequest("POST", fmt.Sprintf("%s/authenticate", s.config.AuthConfig.ConnectionURL), bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusUnauthorized {
		return nil, errors.New("invalid credentials")
	} else if response.StatusCode != http.StatusAccepted {
		return nil, errors.New("error calling auth service")
	}

	var jsonFromService utils.JsonResponse
	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		return nil, err
	}

	if jsonFromService.Error {
		return nil, err
	}

	return jsonFromService.Data, nil
}

func (s *BrokerService) SendMail(msg dto.MailPayload) error {
	jsonData, _ := json.MarshalIndent(msg, "", "\t")

	mailServiceURL := fmt.Sprintf("%s/send", s.config.MailerConfig.ConnectionURL)

	request, err := http.NewRequest("POST", mailServiceURL, bytes.NewBuffer(jsonData))
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
		return errors.New("error calling mail service")
	}

	return nil
}

func (s *BrokerService) pushToQueue(name, msg string) error {
	err := s.eventEmitter.SetupChannel()
	if err != nil {
		return err
	}

	payload := dto.LogPayload{
		Name: name,
		Data: msg,
	}

	j, _ := json.MarshalIndent(&payload, "", "\t")

	return s.eventEmitter.Publish(string(j), "log.INFO")
}

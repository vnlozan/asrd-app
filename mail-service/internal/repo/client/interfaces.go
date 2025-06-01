package client

import (
	"mailer-service/internal/dto"
)

type IMailerClient interface {
	SendSMTPMessage(msg dto.Message) error
}

package client

import (
	"mailer/internal/dto"
)

type IMailerClient interface {
	SendSMTPMessage(msg dto.Message) error
}

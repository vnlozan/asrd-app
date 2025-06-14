package client

import (
	"mailer/internal/dto"
)

type IMailerClient interface {
	SendSMTPMessage(msg dto.MailMessage) error
}

package service

import (
	"mailer/internal/dto"
	"mailer/internal/repo/client"
)

type MailerService struct {
	MailerClient client.IMailerClient
}

func NewMailerService(mailerClient client.IMailerClient) IMailerService {
	return &MailerService{
		MailerClient: mailerClient,
	}
}

func (s *MailerService) SendMessage(msg dto.MailMessage) error {
	return s.MailerClient.SendSMTPMessage(msg)
}

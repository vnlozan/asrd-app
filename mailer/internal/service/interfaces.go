package service

import "mailer/internal/dto"

type IMailerService interface {
	SendMessage(msg dto.MailMessage) error
}

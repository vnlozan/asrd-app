package service

import "broker/internal/dto"

type IBrokerService interface {
	Authenticate(credentials dto.AuthPayload) (any, error)
	SendMail(msg dto.MailPayload) error
	LogItem(entry dto.LogPayload) error
}

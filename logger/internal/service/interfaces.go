package service

import "logger/internal/dto"

type ILoggerService interface {
	AddOneLog(logEntry dto.LogEntry) error
}

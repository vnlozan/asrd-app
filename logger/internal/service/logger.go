package service

import (
	"context"
	"logger/internal/dto"
	"logger/internal/repo"
)

type LoggerService struct {
	loggerStorage repo.ILoggerStorage
}

func NewLoggerService(loggerStorage repo.ILoggerStorage) ILoggerService {
	return &LoggerService{loggerStorage: loggerStorage}
}

func (s *LoggerService) AddOneLog(logEntry dto.LogEntry) error {
	return s.loggerStorage.InsertOne(context.Background(), logEntry)
}

package event

import (
	"context"
	"fmt"
	"listener/internal/infra/loggerclient"
	"listener/internal/infra/rabbitmq"
)

type LogHandler struct {
	logger *loggerclient.Client
}

func NewLogHandler(logger *loggerclient.Client) *LogHandler {
	return &LogHandler{logger: logger}
}

func (h *LogHandler) Handle(ctx context.Context, p rabbitmq.Payload) error {
	switch p.Name {
	case "log", "event":
		return h.logger.Log(ctx, loggerclient.Log{Name: p.Name, Data: p.Data})
	case "auth":
		return nil
	default:
		return fmt.Errorf("unknown payload type: %s", p.Name)
	}
}

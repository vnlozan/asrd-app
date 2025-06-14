package event

import (
	"context"
	"listener/internal/infra/rabbitmq"
)

type Handler interface {
	Handle(ctx context.Context, p rabbitmq.Payload) error
}

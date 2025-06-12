package repo

import (
	"context"
	"logger/internal/dto"
)

type ILoggerStorage interface {
	InsertOne(ctx context.Context, entry dto.LogEntry) error
	UpdateOne(ctx context.Context, entry dto.LogEntry) (MatchedCount int64, ModifiedCount int64, UpsertedCount int64, UpsertedID interface{}, err error)
	SelectAll(ctx context.Context) ([]*dto.LogEntry, error)
	SelectOne(ctx context.Context, id string) (*dto.LogEntry, error)
	DropCollection(ctx context.Context) error
}

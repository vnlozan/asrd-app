package repo

import (
	"auth/internal/dto"
	"context"
)

type IUserStorage interface {
	SelectAll(ctx context.Context) ([]*dto.User, error)
	SelectOneByEmail(ctx context.Context, email string) (*dto.User, error)
	SelectOne(ctx context.Context, id int) (*dto.User, error)
	UpdateOne(ctx context.Context, user *dto.User) error
	UpdateOnePassword(ctx context.Context, id int, encryptedPassword string) error
	DeleteOneByID(ctx context.Context, id int) error
	InsertOne(ctx context.Context, user dto.User) (int, error)
}

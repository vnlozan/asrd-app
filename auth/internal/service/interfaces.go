package service

import (
	"auth/internal/dto"
	"context"
)

type IAuthService interface {
	Authenticate(ctx context.Context, auth dto.AuthRequest) (*dto.User, error)
}

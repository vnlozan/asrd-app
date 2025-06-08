package service

import (
	"auth/internal/dto"
	"auth/internal/repo"
	"auth/internal/utils"
	"context"
	"errors"
)

type AuthService struct {
	userStorage repo.IUserStorage
}

func NewAuthService(userStorage repo.IUserStorage) IAuthService {
	return &AuthService{userStorage: userStorage}
}

func (a *AuthService) Authenticate(ctx context.Context, auth dto.AuthRequest) (*dto.User, error) {
	user, err := a.userStorage.SelectOneByEmail(ctx, auth.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	valid, err := utils.PasswordMatches(user.Password, auth.Password)
	if err != nil || !valid {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

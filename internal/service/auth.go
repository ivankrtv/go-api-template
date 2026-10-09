package service

import (
	"go-api-template/internal/domain/token"
	"go-api-template/internal/domain/user"
	"time"
)

type AuthService struct {
}

func NewAuthService() *AuthService {
	return &AuthService{}
}

func (a *AuthService) CreateTokenPair(user user.User) (token.AuthTokenPair, error) {
	// TODO: implement logic
	return token.AuthTokenPair{
		AuthToken:             "token",
		AuthTokenExpiredAt:    time.Now().Add(15 * time.Minute),
		RefreshToken:          "refresh_token",
		RefreshTokenExpiredAt: time.Now().AddDate(0, 1, 0),
	}, nil
}

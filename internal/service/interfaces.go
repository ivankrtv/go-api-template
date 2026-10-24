package service

import (
	"context"
	"go-api-template/internal/domain/token"
	"go-api-template/internal/domain/user"
)

type UserRepositoryInterface interface {
	SaveUser(ctx context.Context, user user.User) error
}

type AuthRepositoryInterface interface {
	SaveTokenPair(ctx context.Context, tokenPair token.AuthTokenPair) error
}

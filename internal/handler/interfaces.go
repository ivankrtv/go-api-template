package handler

import (
	"context"
	"go-api-template/internal/domain/token"
	"go-api-template/internal/domain/user"
)

type UserService interface {
	Registration(ctx context.Context, emailAsString string, password string, username string) (user.User, error)
}

type AuthService interface {
	CreateTokenPair(user user.User) (token.AuthTokenPair, error)
}

package service

import (
	"context"
	"go-api-template/internal/domain/user"
)

type UserRepositoryInterface interface {
	SaveUser(ctx context.Context, user user.User) error
}

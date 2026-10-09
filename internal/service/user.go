package service

import (
	"context"
	"errors"
	"fmt"
	"go-api-template/internal/domain/user"
)

type UserService struct {
	repository UserRepositoryInterface
}

func NewUserService(repo UserRepositoryInterface) *UserService {
	return &UserService{repository: repo}
}

func (s *UserService) Registration(
	ctx context.Context,
	emailAsString string,
	password string,
	username string,
) (user.User, error) {
	email := user.Email(emailAsString)
	isEmailValid := email.Validate()

	if !isEmailValid {
		return user.User{}, user.ErrInvalidEmail
	}

	usr := user.New(email, password, username)

	err := s.repository.SaveUser(ctx, usr)
	if err != nil {
		if errors.Is(err, user.ErrEmailAlreadyTaken) {
			return user.User{}, err
		}
		return user.User{}, fmt.Errorf("fail to save user: %w", err)
	}

	return usr, nil
}

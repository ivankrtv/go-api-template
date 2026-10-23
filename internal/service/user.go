package service

import (
	"context"
	"errors"
	"fmt"
	"go-api-template/internal/domain/user"

	"github.com/rs/zerolog/log"
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
	l := log.With().Any("request_id", ctx.Value("X-Request-ID")).Logger()

	l.Info().Msg("test")

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
		l.Error().Err(err).Msg("fail to save user")
		return user.User{}, fmt.Errorf("fail to save user: %w", err)
	}

	return usr, nil
}

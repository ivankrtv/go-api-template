package service

import (
	"context"
	"errors"
	"fmt"
	"go-api-template/internal/domain/user"
	"go-api-template/internal/utils"

	"github.com/rs/zerolog"
)

type UserService struct {
	logger     *zerolog.Logger
	repository UserRepositoryInterface
}

func NewUserService(logger *zerolog.Logger, repo UserRepositoryInterface) *UserService {
	return &UserService{
		logger:     logger,
		repository: repo,
	}
}

func (s *UserService) Registration(
	ctx context.Context,
	emailAsString string,
	password string,
	username string,
) (user.User, error) {
	l := s.logger.With().Str(utils.RequestID, utils.GetRequestIDFromCtx(ctx)).Logger()

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

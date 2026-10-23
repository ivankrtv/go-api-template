package service

import (
	"context"
	mocks "go-api-template/gen/mocks"
	"go-api-template/internal/domain/user"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func newUserService(t *testing.T) *UserService {
	nopLogger := zerolog.Nop()
	return NewUserService(&nopLogger, mocks.NewMockUserRepositoryInterface(t))
}

func TestUserService_Registration(t *testing.T) {
	ctx := context.TODO()

	type args struct {
		email    string
		password string
		username string
	}

	type want struct {
		err  string
		user user.User
	}

	tests := []struct {
		name       string
		args       args
		setupMocks func(s *UserService)
		want       want
	}{
		{
			name: "Success",
			args: args{
				email:    "test@test.ru",
				password: "Qwerty123",
				username: "test",
			},
			setupMocks: func(s *UserService) {
				db := s.repository.(*mocks.MockUserRepositoryInterface)

				db.EXPECT().SaveUser(ctx, mock.Anything).Return(nil)
			},
			want: want{
				user: user.User{
					Email:    user.Email("test@test.ru"),
					UserName: "test",
				},
			},
		},
		{
			name: "fail invalid email",
			args: args{
				email:    "test_invalid_email",
				password: "Qwerty123",
				username: "test",
			},
			setupMocks: func(s *UserService) {
			},
			want: want{
				err: "invalid_email",
			},
		},
		{
			name: "fail email already taken",
			args: args{
				email:    "test@test.ru",
				password: "Qwerty123",
				username: "test",
			},
			setupMocks: func(s *UserService) {
				db := s.repository.(*mocks.MockUserRepositoryInterface)

				db.EXPECT().SaveUser(ctx, mock.Anything).Return(user.ErrEmailAlreadyTaken)
			},
			want: want{
				err: "email_already_taken",
			},
		},
	}

	is := require.New(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			serv := newUserService(t)
			tt.setupMocks(serv)

			got, err := serv.Registration(ctx, tt.args.email, tt.args.password, tt.args.username)

			if tt.want.err != "" {
				is.Error(err)
				is.ErrorContains(err, tt.want.err)
				return
			}

			is.NotEmpty(got.ID)
			is.NotEmpty(got.CreatedAt)
			is.Equal(tt.want.user.Email, got.Email)
			is.Equal(tt.want.user.UserName, got.UserName)
		})
	}
}

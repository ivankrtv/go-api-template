package repository

import (
	"context"
	"go-api-template/internal/domain/user"
	"go-api-template/internal/testutils"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestUserRepository_SaveUser(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		user    user.User
		wantErr string
	}{
		{
			name: "Success",
			user: user.User{
				ID:        uuid.New(),
				Email:     user.Email("test@test.ru"),
				Password:  "Qwerty123",
				UserName:  "username",
				CreatedAt: time.Now(),
			},
		},
		{
			name: "FailEmailAlreadyTaken",
			user: user.User{
				ID:        uuid.New(),
				Email:     user.Email("taken_email@test.ru"),
				Password:  "Qwerty123",
				UserName:  "username",
				CreatedAt: time.Now(),
			},
			wantErr: "email_already_taken",
		},
	}

	is := require.New(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := testutils.NewTestDB(t)
			testutils.LoadFixtures(t, db, "testdata/user.sql")
			repo := NewUserRepository(db)

			err := repo.SaveUser(ctx, tt.user)
			if tt.wantErr != "" {
				is.Error(err)
				is.ErrorContains(err, tt.wantErr)
				return
			}

			is.NoError(err)
		})
	}
}

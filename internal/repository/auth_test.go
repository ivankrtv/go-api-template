package repository

import (
	"context"
	"go-api-template/internal/domain/token"
	"go-api-template/internal/testutils"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestAuthRepository_SaveTokenPair(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name  string
		token token.AuthTokenPair
	}{
		{
			name: "Success",
			token: token.AuthTokenPair{
				ID:                    uuid.New(),
				AccessToken:           "test_access_token",
				AccessTokenExpiredAt:  time.Now().Add(15 * time.Minute),
				RefreshToken:          "test_refresh_token",
				RefreshTokenExpiredAt: time.Now().Add(43200 * time.Minute),
			},
		},
	}

	is := require.New(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := testutils.NewTestDB(t)
			repo := NewAuthRepository(db)

			err := repo.SaveTokenPair(ctx, tt.token)
			is.NoError(err)
		})
	}
}

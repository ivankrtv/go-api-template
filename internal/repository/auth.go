package repository

import (
	"context"
	"fmt"
	"go-api-template/internal/database"
	"go-api-template/internal/domain/token"
)

type AuthRepository struct {
	db *database.DB
}

func NewAuthRepository(db *database.DB) *AuthRepository {
	return &AuthRepository{
		db: db,
	}
}

func (r *AuthRepository) SaveTokenPair(ctx context.Context, tokenPair token.AuthTokenPair) error {
	query := `
		INSERT INTO token (id, user_id, access_token, access_token_expires_at, refresh_token, refresh_token_expires_at) 
			VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.Primary.Exec(ctx, query,
		tokenPair.ID,
		tokenPair.User.ID,
		tokenPair.AccessToken,
		tokenPair.AccessTokenExpiredAt,
		tokenPair.RefreshToken,
		tokenPair.RefreshTokenExpiredAt,
	)

	if err != nil {
		return fmt.Errorf("fail save token pair into db: %w", err)
	}

	return nil
}

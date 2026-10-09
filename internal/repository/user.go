package repository

import (
	"context"
	"errors"
	"fmt"
	"go-api-template/internal/database"
	"go-api-template/internal/domain/user"
	"go-api-template/internal/utils/pg"

	"github.com/jackc/pgx/v5/pgconn"
)

type UserRepository struct {
	db *database.DB
}

func NewUserRepository(db *database.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (ur *UserRepository) SaveUser(ctx context.Context, usr user.User) error {
	query := `
		INSERT INTO "user" (id, email, password, username, created_at) VALUES ($1, $2, $3, $4, $5)
	`

	_, err := ur.db.Primary.Exec(ctx, query,
		usr.ID,
		usr.Email,
		usr.Password,
		usr.UserName,
		usr.CreatedAt,
	)

	if err != nil {
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) {
			if pgError.Code == pg.ErrUniqueConstrainViolation {
				return user.ErrEmailAlreadyTaken
			}
		}

		return fmt.Errorf("fail save user into db: %w", err)
	}

	return nil
}

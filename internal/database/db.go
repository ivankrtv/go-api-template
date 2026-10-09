package database

import (
	"context"
	"go-api-template/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	Primary *pgxpool.Pool
	Replica *pgxpool.Pool
}

func NewDB(ctx context.Context, cfg config.Database) (*DB, error) {
	primary, err := pgxpool.New(ctx, cfg.PrimaryDSN)
	if err != nil {
		return nil, err
	}

	replica, err := pgxpool.New(ctx, cfg.ReplicaDSN)
	if err != nil {
		primary.Close()
		return nil, err
	}

	return &DB{
		Primary: primary,
		Replica: replica,
	}, nil
}

func (db *DB) Close() {
	db.Primary.Close()
	db.Replica.Close()
}

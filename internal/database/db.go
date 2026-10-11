package database

import (
	"context"
	"fmt"
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

	err = primary.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("fail connect to primary db: %w", err)
	}

	replica, err := pgxpool.New(ctx, cfg.ReplicaDSN)
	if err != nil {
		primary.Close()
		return nil, err
	}

	err = replica.Ping(ctx)
	if err != nil {
		primary.Close()
		return nil, fmt.Errorf("fail connect to primary db: %w", err)
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

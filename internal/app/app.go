package app

import (
	"context"
	"fmt"
	"go-api-template/internal/config"
	"go-api-template/internal/database"
	"go-api-template/internal/handler"
	"go-api-template/internal/repository"
	"go-api-template/internal/service"
	"go-api-template/migrations"

	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

type App struct {
	db *database.DB

	userService *service.UserService
	authService *service.AuthService

	userRepository *repository.UserRepository
}

func New(ctx context.Context, cfg config.Config) (*App, error) {
	db, err := database.NewDB(ctx, cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("fail create database: %w", err)
	}

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)

	authService := service.NewAuthService()

	return &App{
		db:             db,
		userService:    userService,
		userRepository: userRepo,
		authService:    authService,
	}, nil
}

func (a *App) NewHandler() *handler.Handler {
	return handler.New(a.userService, a.authService)
}

func (a *App) UpMigrations() error {
	goose.SetBaseFS(migrations.FS)

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	db := stdlib.OpenDBFromPool(a.db.Primary)
	defer db.Close()

	if err := goose.Up(db, "."); err != nil {
		return err
	}

	return nil
}

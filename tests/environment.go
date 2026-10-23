package tests

import (
	transport "go-api-template/gen/oapi"
	"go-api-template/internal/handler"
	"go-api-template/internal/repository"
	"go-api-template/internal/service"
	"go-api-template/internal/testutils"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

func initHandler(t *testing.T) *handler.Handler {
	nopLogger := zerolog.Nop()
	db := testutils.NewTestDB(t)
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(&nopLogger, userRepo)

	authService := service.NewAuthService()

	return handler.New(userService, authService)
}

func setupApp(t *testing.T) *fiber.App {
	handler := initHandler(t)
	app := fiber.New()
	transport.RegisterHandlers(app, handler)
	return app
}

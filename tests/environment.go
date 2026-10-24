package tests

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	transport "go-api-template/gen/oapi"
	"go-api-template/internal/config"
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

	authRepo := repository.NewAuthRepository(db)
	authService := service.NewAuthService(initAuthConfig(), authRepo)

	return handler.New(userService, authService)
}

func setupApp(t *testing.T) *fiber.App {
	handler := initHandler(t)
	app := fiber.New()
	transport.RegisterHandlers(app, handler)
	return app
}

func initAuthConfig() *config.Auth {
	privateKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}

	privBytes := x509.MarshalPKCS1PrivateKey(privateKey)

	pubBytes := x509.MarshalPKCS1PublicKey(&privateKey.PublicKey)

	return &config.Auth{
		Issuer:                        "testapp",
		TokenPrivateKeyPEM:            string(privBytes),
		TokenPublicKeyPEM:             string(pubBytes),
		AccessTokenLifeTimeInMinutes:  15,
		RefreshTokenLifeTimeInMinutes: 43200,
	}
}

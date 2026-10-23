package main

import (
	"context"
	"fmt"
	transport "go-api-template/gen/oapi"
	"go-api-template/internal/app"
	"go-api-template/internal/config"
	"go-api-template/internal/metrics"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

func main() {
	ctx := context.Background()

	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatalf("fail init config: %s", err.Error())
	}

	logLevel, err := zerolog.ParseLevel(cfg.LogLevel)
	if err != nil {
		log.Fatalf("invalid log level: %s", err.Error())
	}
	logger := zerolog.New(os.Stderr).Level(logLevel).With().Timestamp().Logger()

	application, err := app.New(ctx, cfg, &logger)
	if err != nil {
		log.Fatalf("fail init app: %s", err.Error())
	}

	err = application.UpMigrations()
	if err != nil {
		log.Fatalf("fail up migrations: %s", err.Error())
	}

	h := application.NewHandler()

	server := fiber.New(fiber.Config{
		AppName: cfg.ServiceName,
	})

	server.Use(fiberzerolog.New(fiberzerolog.Config{
		Logger: &logger,
	}))

	server.Use("/metrics", metrics.NewMetricsMiddleware(cfg.ServiceName))

	transport.RegisterHandlers(server, h)

	go func() {
		host := fmt.Sprintf(":%d", cfg.HTTPPort)
		if err := server.Listen(host); err != nil {
			log.Fatalf("fail start server: %s", err.Error())
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	err = server.Shutdown()
	if err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	application.Close()
}

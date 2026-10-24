package main

import (
	"context"
	transport "go-api-template/gen/oapi"
	"go-api-template/internal/app"
	"go-api-template/internal/config"
	"go-api-template/internal/utils/logger"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
)

func main() {
	ctx := context.Background()

	if os.Getenv("ENVIRONMENT_NAME") == "" {
		err := godotenv.Load()
		if err != nil {
			log.Fatalf("fail parse environment vars from .env")
		}
	}

	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatalf("fail init config: %s", err.Error())
	}

	l, err := logger.InitLogger(cfg.LogLevel)
	if err != nil {
		log.Fatalf("fail init logger: %s", err.Error())
	}

	application, err := app.New(ctx, cfg, l)
	if err != nil {
		log.Fatalf("fail init app: %s", err.Error())
	}

	err = application.UpMigrations()
	if err != nil {
		log.Fatalf("fail up migrations: %s", err.Error())
	}

	h := application.NewHandler()

	server := app.NewServer(&cfg, l)

	transport.RegisterHandlers(server.GetRouter(), h)

	go func() {
		if err := server.Run(); err != nil {
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

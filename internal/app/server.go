package app

import (
	"context"
	"fmt"
	"go-api-template/internal/config"
	"go-api-template/internal/metrics"
	"go-api-template/internal/utils"

	fiberzerolog "github.com/gofiber/contrib/v3/zerolog"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type Server struct {
	cfg    *config.Config
	logger *zerolog.Logger
	srv    *fiber.App
}

func NewServer(cfg *config.Config, logger *zerolog.Logger) *Server {
	srv := fiber.New(fiber.Config{
		AppName: cfg.ServiceName,
	})

	srv.Use(func(c fiber.Ctx) error {
		requestID := c.Get(utils.HTTPRequestID)
		if requestID == "" {
			requestID = uuid.New().String()
		}

		c.SetContext(context.WithValue(c.Context(), utils.RequestIDCtxKey, requestID))

		return c.Next()
	})

	srv.Use(fiberzerolog.New(fiberzerolog.Config{
		GetLogger: func(c fiber.Ctx) zerolog.Logger {
			requestID := c.Context().Value(utils.RequestIDCtxKey)
			return logger.With().Any(utils.RequestID, requestID).Logger()
		},
		Levels: []zerolog.Level{zerolog.InfoLevel, zerolog.ErrorLevel},
	}))

	srv.Use("/metrics", metrics.NewMetricsMiddleware(cfg.ServiceName))

	return &Server{
		cfg:    cfg,
		logger: logger,
		srv:    srv,
	}
}

func (s *Server) GetRouter() *fiber.App {
	return s.srv
}

func (s *Server) Run() error {
	host := fmt.Sprintf(":%d", s.cfg.HTTPPort)
	return s.srv.Listen(host)
}

func (s *Server) Shutdown() error {
	return s.srv.Shutdown()
}

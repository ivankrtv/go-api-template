package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	UserService UserService
	AuthService AuthService
}

func New(
	userService UserService,
	authService AuthService,
) *Handler {
	return &Handler{
		UserService: userService,
		AuthService: authService,
	}
}

func (h *Handler) Health(c *fiber.Ctx) error {
	if _, err := c.Status(http.StatusOK).Write([]byte("ok")); err != nil {
		return err
	}
	return nil
}

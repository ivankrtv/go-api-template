package handler

import (
	"errors"
	transport "go-api-template/gen/oapi"
	httpErrors "go-api-template/internal/errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) Registration(c *fiber.Ctx) error {
	var request transport.RegistrationRequest
	err := c.BodyParser(&request)
	if err != nil {
		if errors.Is(err, fiber.ErrUnprocessableEntity) {
			httpErrors.BadRequestError(c, errors.New("empty_request"))
			return nil
		}
		httpErrors.InternalServerError(c, err)
		return nil
	}

	usr, err := h.UserService.Registration(c.Context(), request.Email, request.Password, request.Username)
	if err != nil {
		httpErrors.BadRequestError(c, err)
		return nil
	}

	tokenPair, err := h.AuthService.CreateTokenPair(usr)
	if err != nil {
		httpErrors.InternalServerError(c, err)
		return nil
	}

	resp := transport.RegistrationResponse{
		AuthToken: &tokenPair.AuthToken,
	}
	if err = c.Status(http.StatusOK).JSON(resp); err != nil {
		return err
	}

	return nil
}

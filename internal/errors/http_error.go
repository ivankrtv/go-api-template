package errors

import (
	"go-api-template/internal/utils"
	"net/http"

	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog/log"
)

type HttpError struct {
	Error string `json:"error"`
	Type  string `json:"type"`
}

func InternalServerError(c fiber.Ctx, err error) HttpError {
	requestID := utils.GetRequestIDFromCtx(c.Context())
	l := log.With().Timestamp().Str(utils.RequestID, requestID).Logger()
	l.Error().Err(err).Msg("Internal server error")

	httpErr := HttpError{
		Type:  "internal_server_error",
		Error: "Internal server error",
	}

	c.Status(http.StatusInternalServerError).JSON(httpErr) //nolint:errcheck

	return httpErr
}

func BadRequestError(c fiber.Ctx, err error) HttpError {
	httpErr := HttpError{
		Type:  "Bad request",
		Error: err.Error(),
	}

	c.Status(http.StatusBadRequest).JSON(httpErr) //nolint:errcheck

	return httpErr
}

func UnprocessableEntity(c fiber.Ctx, err error) HttpError {
	httpErr := HttpError{
		Type:  "Unprocessable entity",
		Error: err.Error(),
	}

	c.Status(http.StatusUnprocessableEntity).JSON(httpErr) //nolint:errcheck

	return httpErr
}

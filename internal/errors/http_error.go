package errors

import (
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type HttpError struct {
	Error string `json:"error"`
	Type  string `json:"type"`
}

func InternalServerError(c *fiber.Ctx, err error) HttpError {
	log.Printf("Internal error: %s", err.Error())

	httpErr := HttpError{
		Type:  "internal_server_error",
		Error: "Internal server error",
	}

	c.Status(http.StatusInternalServerError).JSON(httpErr) //nolint:errcheck

	return httpErr
}

func BadRequestError(c *fiber.Ctx, err error) HttpError {
	httpErr := HttpError{
		Type:  "Bad request",
		Error: err.Error(),
	}

	c.Status(http.StatusBadRequest).JSON(httpErr) //nolint:errcheck

	return httpErr
}

func UnprocessableEntity(c *fiber.Ctx, err error) HttpError {
	httpErr := HttpError{
		Type:  "Unprocessable entity",
		Error: err.Error(),
	}

	c.Status(http.StatusUnprocessableEntity).JSON(httpErr) //nolint:errcheck

	return httpErr
}

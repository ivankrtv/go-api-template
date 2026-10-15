package validator

import (
	transport "go-api-template/gen/oapi"
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

func ValidateRegistration(req transport.RegistrationRequest) error {
	validation.ErrRequired = validation.ErrRequired.SetMessage("the value is required")

	return validation.ValidateStruct(&req,
		validation.Field(&req.Username, validation.Required, validation.Length(3, 50)),
		validation.Field(&req.Email, validation.Required, is.Email),
		validation.Field(&req.Password, validation.Required, validation.Match(regexp.MustCompile("^[a-zA-Z0-9]"))),
	)
}

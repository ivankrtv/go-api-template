package user

import "errors"

var (
	ErrInvalidEmail      = errors.New("invalid_email")
	ErrEmailAlreadyTaken = errors.New("email_already_taken")
)

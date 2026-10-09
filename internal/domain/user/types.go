package user

import (
	"net/mail"
)

type Email string

func (e Email) Validate() bool {
	_, err := mail.ParseAddress(string(e))

	return err == nil
}

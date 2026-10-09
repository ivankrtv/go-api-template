package user

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID
	Email     Email
	Password  string
	UserName  string
	CreatedAt time.Time
}

func New(email Email, password string, userName string) User {
	id, _ := uuid.NewV7()
	return User{
		ID:        id,
		Email:     email,
		Password:  password,
		UserName:  userName,
		CreatedAt: time.Now(),
	}
}

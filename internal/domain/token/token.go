package token

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID uuid.UUID
}

type AuthTokenPair struct {
	ID                    uuid.UUID
	AccessToken           string
	AccessTokenExpiredAt  time.Time
	RefreshToken          string
	RefreshTokenExpiredAt time.Time
	User                  User
}

package token

import "time"

type AuthTokenPair struct {
	AuthToken             string
	AuthTokenExpiredAt    time.Time
	RefreshToken          string
	RefreshTokenExpiredAt time.Time
}

package service

import (
	"context"
	"crypto/rsa"
	"fmt"
	"go-api-template/internal/config"
	"go-api-template/internal/domain/token"
	"go-api-template/internal/domain/user"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AccessClaims struct {
	UserID uuid.UUID `json:"user_id"`
	jwt.RegisteredClaims
}

type RefreshClaims struct {
	jwt.RegisteredClaims
}

type AuthService struct {
	cfg  *config.Auth
	repo AuthRepositoryInterface
}

func NewAuthService(cfg *config.Auth, repo AuthRepositoryInterface) *AuthService {
	return &AuthService{
		cfg:  cfg,
		repo: repo,
	}
}

func (a *AuthService) CreateTokenPair(ctx context.Context, user user.User) (token.AuthTokenPair, error) {
	accessTokenExpiresAt := time.Now().Add(time.Duration(a.cfg.AccessTokenLifeTimeInMinutes) * time.Minute)
	refreshTokenExpiresAt := time.Now().Add(time.Duration(a.cfg.RefreshTokenLifeTimeInMinutes) * time.Minute)

	accessClaims := AccessClaims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessTokenExpiresAt),
			Issuer:    a.cfg.Issuer,
		},
	}

	refreshClaims := RefreshClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshTokenExpiresAt),
			Issuer:    a.cfg.Issuer,
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodRS256, accessClaims)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodRS256, refreshClaims)

	signKey, err := a.getSignKey()
	if err != nil {
		return token.AuthTokenPair{}, fmt.Errorf("fail get sign key: %w", err)
	}

	signedAccessToken, err := accessToken.SignedString(signKey)
	if err != nil {
		return token.AuthTokenPair{}, fmt.Errorf("fail sign access token: %w", err)
	}

	signedRefreshToken, err := refreshToken.SignedString(signKey)
	if err != nil {
		return token.AuthTokenPair{}, fmt.Errorf("fail sign access token: %w", err)
	}

	tokenPair := token.AuthTokenPair{
		ID:                    uuid.New(),
		AccessToken:           signedAccessToken,
		AccessTokenExpiredAt:  accessTokenExpiresAt,
		RefreshToken:          signedRefreshToken,
		RefreshTokenExpiredAt: refreshTokenExpiresAt,
		User: token.User{
			ID: user.ID,
		},
	}

	err = a.repo.SaveTokenPair(ctx, tokenPair)
	if err != nil {
		return token.AuthTokenPair{}, fmt.Errorf("fail to save token pair in db: %w", err)
	}

	return tokenPair, nil
}

func (a *AuthService) getSignKey() (*rsa.PrivateKey, error) {
	signKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(a.cfg.TokenPrivateKeyPEM))
	if err != nil {
		return nil, fmt.Errorf("fail parse auth private key pem: %w", err)
	}

	return signKey, nil
}

//nolint:unused // tmp
func (a *AuthService) getVerifyKey() (*rsa.PublicKey, error) {
	verifyKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(a.cfg.TokenPublicKeyPEM))
	if err != nil {
		return nil, fmt.Errorf("fail parse auth public key pem: %w", err)
	}

	return verifyKey, nil
}

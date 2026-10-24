package service

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"go-api-template/internal/config"
	"go-api-template/internal/domain/user"
	"testing"

	mocks "go-api-template/gen/mocks"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const (
	testIssuer            = "test-issuer"
	testAccessTTLMinutes  = 15
	testRefreshTTLMinutes = 1440
)

func generateKeyPairPEM(t *testing.T) (privateKeyPEM, publicKeyPEM string) {
	t.Helper()

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	privateKeyPEMBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	publicKeyDER, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	require.NoError(t, err)

	publicKeyPEMBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyDER,
	})

	return string(privateKeyPEMBytes), string(publicKeyPEMBytes)
}

func TestAuthService_CreateTokenPair(t *testing.T) {
	t.Parallel()

	privateKeyPEM, publicKeyPEM := generateKeyPairPEM(t)

	validCfg := &config.Auth{
		Issuer:                        testIssuer,
		TokenPrivateKeyPEM:            privateKeyPEM,
		TokenPublicKeyPEM:             publicKeyPEM,
		AccessTokenLifeTimeInMinutes:  testAccessTTLMinutes,
		RefreshTokenLifeTimeInMinutes: testRefreshTTLMinutes,
	}

	invalidKeyCfg := *validCfg
	invalidKeyCfg.TokenPrivateKeyPEM = "not-a-pem"

	tests := []struct {
		name      string
		config    *config.Auth
		setupMock func(m *mocks.MockAuthRepositoryInterface)
		wantErr   string
	}{
		{
			name:   "Success",
			config: validCfg,
			setupMock: func(m *mocks.MockAuthRepositoryInterface) {
				m.EXPECT().SaveTokenPair(mock.Anything, mock.Anything).Return(nil)
			},
		},
		{
			name:      "InvalidSignKeyPEM_Fail",
			config:    &invalidKeyCfg,
			setupMock: func(m *mocks.MockAuthRepositoryInterface) {},
			wantErr:   "fail get sign key",
		},
		{
			name:   "SaveTokenPairError_Fail",
			config: validCfg,
			setupMock: func(m *mocks.MockAuthRepositoryInterface) {
				m.EXPECT().SaveTokenPair(mock.Anything, mock.Anything).Return(errors.New("db down"))
			},
			wantErr: "fail to save token pair in db",
		},
	}

	is := require.New(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockRepo := mocks.NewMockAuthRepositoryInterface(t)
			tt.setupMock(mockRepo)

			svc := NewAuthService(tt.config, mockRepo)

			userID := uuid.New()
			got, err := svc.CreateTokenPair(context.TODO(), user.User{ID: userID})

			if tt.wantErr != "" {
				is.ErrorContains(err, tt.wantErr)
				return
			}

			is.NoError(err)

			publicKey, parseErr := jwt.ParseRSAPublicKeyFromPEM([]byte(tt.config.TokenPublicKeyPEM))
			is.NoError(parseErr)

			accessParsed, parseErr := jwt.ParseWithClaims(got.AccessToken, &AccessClaims{}, func(_ *jwt.Token) (any, error) {
				return publicKey, nil
			}, jwt.WithValidMethods([]string{jwt.SigningMethodRS256.Alg()}))
			is.NoError(parseErr)

			accessClaims, ok := accessParsed.Claims.(*AccessClaims)
			is.True(ok)
			is.Equal(userID, accessClaims.UserID)
		})
	}
}

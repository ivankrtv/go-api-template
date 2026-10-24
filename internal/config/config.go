package config

import "github.com/caarlos0/env/v11"

type Database struct {
	PrimaryDSN string `env:"DATABASE_PRIMARY_DSN,required,notEmpty"`
	ReplicaDSN string `env:"DATABASE_REPLICA_DSN,required,notEmpty"`
}

type Auth struct {
	Issuer                        string `env:"AUTH_TOKEN_ISSUER,required,notEmpty"`
	TokenPublicKeyPEM             string `env:"AUTH_TOKEN_PUBLIC_KEY_PEM,required,notEmpty"`
	TokenPrivateKeyPEM            string `env:"AUTH_TOKEN_PRIVATE_KEY_PEM,required,notEmpty"`
	AccessTokenLifeTimeInMinutes  int    `env:"AUTH_ACCESS_TOKEN_LIFE_TIME_IN_MINUTES,required,notEmpty"`
	RefreshTokenLifeTimeInMinutes int    `env:"AUTH_REFRESH_TOKEN_LIFE_TIME_IN_MINUTES,required,notEmpty"`
}

type Config struct {
	ServiceName string `env:"SERVICE_NAME" envDefault:"service"`
	HTTPPort    int    `env:"HTTP_PORT,required"`
	LogLevel    string `env:"LOG_LEVEL" envDefault:"info"`
	Auth        Auth
	Database    Database
}

func InitConfig() (Config, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

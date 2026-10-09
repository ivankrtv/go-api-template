package config

import "github.com/caarlos0/env/v11"

type Database struct {
	PrimaryDSN string `env:"DATABASE_PRIMARY_DSN,required,notEmpty"`
	ReplicaDSN string `env:"DATABASE_REPLICA_DSN,required,notEmpty"`
}
type Config struct {
	Database Database
}

func InitConfig() (Config, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

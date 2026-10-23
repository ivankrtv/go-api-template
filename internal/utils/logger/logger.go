package logger

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
)

func InitLogger(logLevel string) (*zerolog.Logger, error) {
	lvl, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		return nil, fmt.Errorf("fail parse log level: %w", err)
	}
	l := zerolog.New(os.Stderr).Level(lvl).With().Timestamp().Logger()

	return &l, nil
}

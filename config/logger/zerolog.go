package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

func NewLogger() zerolog.Logger {
	zerolog.TimeFieldFormat = time.RFC3339
	appLogger := zerolog.New(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	}).With().Timestamp().Logger()

	return appLogger
}

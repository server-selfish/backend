package logger

import (
	"io"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/server-selfish/backend/config/otel"
	"github.com/spf13/viper"
)

func NewLogger(logs *otel.Logs) zerolog.Logger {
	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.stdout.enabled", true)
	viper.SetDefault("log.stdout.pretty", true)

	zerolog.TimeFieldFormat = time.RFC3339

	writers := make([]io.Writer, 0, 2)
	if viper.GetBool("log.stdout.enabled") {
		if viper.GetBool("log.stdout.pretty") {
			// ConsoleWriter parses the JSON event and pretty-prints it,
			// so stdout keeps its current human-readable format.
			writers = append(writers, zerolog.ConsoleWriter{
				Out:        os.Stdout,
				TimeFormat: time.RFC3339,
			})
		} else {
			writers = append(writers, os.Stdout)
		}
	}
	if logs.Enabled() {
		writers = append(writers, newOTELWriter(logs))
	}
	if len(writers) == 0 {
		writers = append(writers, io.Discard)
	}

	appLogger := zerolog.New(zerolog.MultiLevelWriter(writers...)).
		With().Timestamp().Logger().
		Level(parseLevel(viper.GetString("log.level")))

	if warning := logs.Warning(); warning != "" {
		appLogger.Warn().Msg(warning)
	}

	return appLogger
}

func parseLevel(level string) zerolog.Level {
	if strings.TrimSpace(level) == "" {
		return zerolog.InfoLevel
	}
	if l, err := zerolog.ParseLevel(level); err == nil {
		return l
	}
	return zerolog.InfoLevel
}

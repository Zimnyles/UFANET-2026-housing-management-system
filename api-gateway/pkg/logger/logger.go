package logger

import (
	"os"
	"strings"

	"github.com/rs/zerolog"
)

func NewLogger(serviceName, logLevel string) *zerolog.Logger {
	zerolog.SetGlobalLevel(parseLogLevel(logLevel))

	l := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).
		With().
		Timestamp().
		Str("service", serviceName).
		Logger()

	return &l
}

func parseLogLevel(level string) zerolog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return zerolog.DebugLevel
	case "warn":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	default:
		return zerolog.InfoLevel
	}
}

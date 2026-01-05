package config

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// SetupLogger configures the global logger based on log level
func SetupLogger(level string) {
	// Parse log level
	var zeroLogLevel zerolog.Level
	switch level {
	case "debug":
		zeroLogLevel = zerolog.DebugLevel
	case "info":
		zeroLogLevel = zerolog.InfoLevel
	case "warn":
		zeroLogLevel = zerolog.WarnLevel
	case "error":
		zeroLogLevel = zerolog.ErrorLevel
	default:
		zeroLogLevel = zerolog.InfoLevel
	}

	// Set global log level
	zerolog.SetGlobalLevel(zeroLogLevel)

	// Use console writer for human-readable output
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
}

package logger

import (
	"log/slog"
	"os"
	"strings"
)

// New initializes and returns a new slog.Logger.
// The log level is determined by the `LOG_LEVEL` environment variable.
func New(logLevelStr string) *slog.Logger {
	var level slog.Level

	switch strings.ToUpper(logLevelStr) {
	case "DEBUG":
		level = slog.LevelDebug
	case "INFO":
		level = slog.LevelInfo
	case "WARN":
		level = slog.LevelWarn
	case "ERROR":
		level = slog.LevelError
	default:
		level = slog.LevelInfo // Default level
	}

	opts := &slog.HandlerOptions{
		Level: level,
	}

	handler := slog.NewJSONHandler(os.Stdout, opts)
	return slog.New(handler)
}

package logger

import (
	"log/slog"
	"os"

	"github.com/M-kos/hitalent_test/internal/config"
)

type Logger struct {
	slog.Logger
}

func New(config *config.Config) *Logger {
	level := slog.LevelDebug

	if config != nil && config.LogLevel != "" {
		switch config.LogLevel {
		case "debug":
			level = slog.LevelDebug
		case "info":
			level = slog.LevelInfo
		case "warn":
			level = slog.LevelWarn
		case "error":
			level = slog.LevelError
		default:
			level = slog.LevelInfo
		}
	}

	opts := &slog.HandlerOptions{
		Level: level,
	}

	handler := slog.NewJSONHandler(os.Stdout, opts)

	return slog.New(handler)
}

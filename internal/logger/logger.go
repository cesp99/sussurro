package logger

import (
	"log/slog"
	"os"
)

func Init(level string) *slog.Logger {
	var lvl slog.Level
	switch level {
	case "debug":
		lvl = slog.LevelDebug
	case "info":
		lvl = slog.LevelInfo
	case "warn":
		lvl = slog.LevelWarn
	case "error":
		lvl = slog.LevelError
	default:
		lvl = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: lvl,
	}

	// Use JSON handler for structured logging, good for production/parsing
	// For local development, TextHandler might be easier to read, but JSON is safer generally.
	// Given this is a desktop app, maybe TextHandler is better for stdout debugging?
	// I'll stick to JSON for now as it's "infrastructure".
	handler := slog.NewJSONHandler(os.Stdout, opts)
	logger := slog.New(handler)
	slog.SetDefault(logger)

	return logger
}

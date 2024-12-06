package logger

import (
	"log/slog"
	"os"
	"strings"
)

func SetLogger(logLevel string) {
	logLevel = strings.ToUpper(logLevel)
	level := slog.LevelDebug

	switch logLevel {
	case "DEBUG":
		level = slog.LevelDebug
	case "INFO":
		level = slog.LevelInfo
	case "WARN":
		level = slog.LevelWarn
	case "ERROR":
		level = slog.LevelError
	}

	opts := &slog.HandlerOptions{Level: level}
	handler := slog.NewJSONHandler(os.Stdout, opts)
	logger := slog.New(handler)
	slog.SetDefault(logger)
}

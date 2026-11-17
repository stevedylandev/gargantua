package main

import (
	"log/slog"
	"os"
	"strings"
)

var logger = newLogger()

func newLogger() *slog.Logger {
	opts := &slog.HandlerOptions{Level: levelFromEnv()}

	var handler slog.Handler
	if strings.EqualFold(envString("GARGANTUA_LOG_FORMAT", "text"), "json") {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	return slog.New(handler)
}

func levelFromEnv() slog.Level {
	switch strings.ToUpper(envString("GARGANTUA_LOG_LEVEL", "info")) {
	case "DEBUG":
		return slog.LevelDebug
	case "WARN":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func fatal(msg string, args ...any) {
	logger.Error(msg, args...)
	os.Exit(1)
}

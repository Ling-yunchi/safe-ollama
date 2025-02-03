package utils

import (
	"log/slog"
	"os"
	"safe-ollama/config"
)

func InitLogger() *slog.Logger {
	textOpts := &slog.HandlerOptions{
		Level: config.LoggingLevel,
	}
	var handler slog.Handler
	switch config.LoggingOutput {
	case "stdout":
		handler = slog.NewTextHandler(os.Stdout, textOpts)
	case "stderr":
		handler = slog.NewTextHandler(os.Stderr, textOpts)
	default:
		logFile, err := os.OpenFile(config.LoggingOutput, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			slog.Warn("Failed to open log file, log to stdout", "error", err)
			handler = slog.NewTextHandler(os.Stdout, textOpts)
		}
		handler = slog.NewTextHandler(logFile, textOpts)
	}
	logger := slog.New(handler)
	slog.SetDefault(logger)
	return logger
}

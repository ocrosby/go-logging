package main

import (
	"log/slog"
	"os"

	"github.com/ocrosby/go-logging/pkg/logging"
)

func main() {
	logger := logging.NewSlogTextLogger(logging.InfoLevel)

	logger.Info("Application started with slog backend")
	logger.Debug("This won't appear because level is INFO")
	logger.Warn("This is a warning")
	logger.Error("This is an error")

	logger = logger.WithField("service", "slog-example")
	logger.Info("Logger with static field")

	logger = logger.WithFields(map[string]interface{}{
		"version": "2.0.0",
		"env":     "production",
	})
	logger.Info("Logger with multiple fields")

	ctx := logging.NewContextWithTrace()
	logger.InfoContext(ctx, "Message with trace ID from context")

	jsonLogger := logging.NewSlogJSONLogger(logging.DebugLevel)
	jsonLogger.Debug("JSON formatted log with slog")
	jsonLogger.Info("Another JSON log")

	customHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	})
	customLogger := logging.NewWithHandler(customHandler)
	customLogger.Info("Using custom slog handler")
	customLogger.Debug("Debug message with custom handler")

	customLogger.Fluent().Info().
		Str("user", "john_doe").
		Int("attempts", 3).
		Msg("Login attempt")
}

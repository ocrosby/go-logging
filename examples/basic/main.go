package main

import (
	"github.com/ocrosby/go-logging/pkg/logging"
)

func main() {
	logger := logging.NewWithLevel(logging.InfoLevel)

	logger.Info("Application started")
	logger.Debug("This won't appear because level is INFO")
	logger.Warn("This is a warning")
	logger.Error("This is an error")

	logger = logger.WithField("service", "example-app")
	logger.Info("Logger with static field")

	logger = logger.WithFields(map[string]interface{}{
		"version": "1.0.0",
		"env":     "development",
	})
	logger.Info("Logger with multiple fields")

	ctx := logging.NewContextWithTrace()
	logger.InfoContext(ctx, "Message with trace ID from context")
}

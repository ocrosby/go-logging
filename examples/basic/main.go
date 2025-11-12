package main

import (
	"github.com/ocrosby/go-logging/pkg/logging"
)

func main() {
	// Super simple - just works!
	logger := logging.NewSimple()

	logger.Info("Application started")
	logger.Debug("This won't appear because level is INFO")
	logger.Warn("This is a warning")
	logger.Error("This is an error")

	// Want JSON? Easy!
	jsonLogger := logging.NewEasyJSON()
	jsonLogger.Info("This will be formatted as JSON")

	// Need a different level? Simple!
	debugLogger := logging.NewEasyJSONWithLevel(logging.DebugLevel)
	debugLogger.Debug("This debug message will now appear")

	// Progressive configuration with fluent builder
	complexLogger := logging.NewEasyBuilder().
		Level(logging.InfoLevel).
		JSON().
		WithFile().
		Field("service", "basic-example").
		Fields(map[string]interface{}{
			"version": "1.0.0",
			"env":     "development",
		}).
		Build()

	complexLogger.Info("Logger with multiple configured options")

	// Environment-based configuration
	envLogger := logging.NewFromEnvSimple()
	envLogger.Info("This logger is configured from environment variables")

	// Context still works the same
	ctx := logging.NewContextWithTrace()
	logger.InfoContext(ctx, "Message with trace ID from context")
}

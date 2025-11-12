// Configuration Examples
// This example demonstrates various ways to configure the logging library

package main

import (
	"github.com/ocrosby/go-logging/pkg/logging"
)

func main() {
	// Basic configuration with builder pattern
	logger := logging.NewEasyBuilder().
		Level(logging.DebugLevel).
		JSON().
		Build()

	logger.Debug("Debug message with JSON formatting")

	// Configuration with static fields
	appLogger := logging.NewEasyBuilder().
		Level(logging.InfoLevel).
		Field("service", "my-app").
		Field("version", "1.0.0").
		Fields(map[string]any{
			"environment": "production",
			"region":      "us-east-1",
		}).
		Build()

	appLogger.Info("Application started with context")

	// File output configuration (includes file info in logs)
	fileLogger := logging.NewEasyBuilder().
		Level(logging.InfoLevel).
		JSON().
		WithFile().
		Field("component", "file-writer").
		Build()

	fileLogger.Info("This logs with file information")

	// Environment-based configuration
	// Set LOG_LEVEL=debug, LOG_FORMAT=json in your environment
	envLogger := logging.NewFromEnvSimple()
	envLogger.Info("Configuration from environment variables")

	// Different log levels using shortcuts
	criticalLogger := logging.NewEasyBuilder().
		Critical().
		JSON().
		Field("target", "alerts").
		Build()

	criticalLogger.Error("Only critical and error messages will show")
	criticalLogger.Info("This won't show (below critical level)")
}

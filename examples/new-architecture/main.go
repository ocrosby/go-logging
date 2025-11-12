package main

import (
	"context"
	"os"

	"github.com/ocrosby/go-logging/pkg/logging"
)

func main() {
	// Example 1: Using the new unified logger with old-style factory methods
	logger := logging.NewWithLevel(logging.InfoLevel)

	// Core Logger interface - works with any logger
	logger.Log(logging.InfoLevel, "This uses the core Logger interface")
	logger.LogContext(context.Background(), logging.InfoLevel, "This includes context")

	// Use level-specific methods directly (now part of Logger interface)
	logger.Info("This uses the convenience Info method")
	logger.Debug("This debug message won't show (level is Info)")

	// Type assertion to access configurable methods
	if configLogger, ok := logger.(logging.ConfigurableLogger); ok {
		configLogger.SetLevel(logging.DebugLevel)
	}

	// Use fluent interface directly (now part of Logger interface)
	logger.Fluent().Info().
		Str("service", "example").
		Int("version", 1).
		Msg("This uses the fluent interface")

	// Example 2: Using new configuration structure
	config := logging.NewLoggerConfig().
		WithLevel(logging.InfoLevel).
		WithJSONFormat().
		Build()

	jsonLogger := logging.NewWithLoggerConfig(config)

	jsonLogger.Info("This will be formatted as JSON")

	// WithField returns a Logger with level methods included
	loggerWithField := jsonLogger.WithField("user_id", 123)
	loggerWithField.Info("With structured fields")

	// Example 3: Using formatters and outputs directly
	formatterConfig := logging.NewFormatterConfig().WithJSONFormat().Build()
	formatter := logging.NewJSONFormatter(formatterConfig)

	entry := logging.LogEntry{
		Level:   logging.InfoLevel,
		Message: "Direct formatter usage",
		Fields: map[string]interface{}{
			"component": "formatter-demo",
		},
	}

	data, err := formatter.Format(entry)
	if err == nil {
		output := logging.NewWriterOutput(os.Stdout)
		_ = output.Write(data)
	}

	// Example 4: Using handler registry
	handlers := logging.ListHandlers()
	for _, name := range handlers {
		jsonLogger.Info("Available handler: %s", name)
	}
}

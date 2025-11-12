package main

import (
	"github.com/ocrosby/go-logging/pkg/logging"
)

func main() {
	// Quickest way to get started - one line!
	logger := logging.NewSimple()

	// Log at different levels
	logger.Debug("Debug message - won't show (level is INFO by default)")
	logger.Info("Application starting...")
	logger.Warn("This is a warning")
	logger.Error("This is an error")

	// Want JSON output instead?
	jsonLogger := logging.NewEasyJSON()
	jsonLogger.Info("This message appears as JSON")

	// Need debug level?
	debugLogger := logging.NewEasyJSONWithLevel(logging.DebugLevel)
	debugLogger.Debug("Now debug messages appear!")
}

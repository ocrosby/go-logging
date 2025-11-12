// YAML Configuration Example
// This example demonstrates how to use YAML configuration files
// to control logging setup declaratively.

package main

import (
	"fmt"
	"log"

	"github.com/ocrosby/go-logging/pkg/logging"
)

func main() {
	fmt.Println("=== YAML Configuration Example ===")

	// Example 1: Load from YAML file
	fmt.Println("1. Loading from development config...")
	devLogger, err := logging.NewFromYAMLFile("../../configs/development.yaml")
	if err != nil {
		log.Printf("Error loading development config: %v", err)
	} else {
		devLogger.Debug("Development logger loaded successfully")
		devLogger.Info("This is an info message in development mode")
	}

	fmt.Println()

	// Example 2: Load from production config
	fmt.Println("2. Loading from production config...")
	prodLogger, err := logging.NewFromYAMLFile("../../configs/production.yaml")
	if err != nil {
		log.Printf("Error loading production config: %v", err)
	} else {
		prodLogger.Info("Production logger loaded successfully")
		prodLogger.Warn("This is a warning in production mode")
	}

	fmt.Println()

	// Example 3: Embedded YAML configuration
	fmt.Println("3. Using embedded YAML configuration...")
	embeddedConfig := `
level: info
format: json
include_file: true
include_time: true
static_fields:
  application: yaml-example
  environment: embedded
output:
  type: stdout
redact_patterns:
  - "password=[^\\s]*"
`

	embeddedLogger, err := logging.LoadFromYAMLString(embeddedConfig)
	if err != nil {
		log.Printf("Error loading embedded config: %v", err)
	} else {
		embeddedLogger.Info("Embedded logger created successfully")
		embeddedLogger.Error("This error log includes static fields")
	}

	fmt.Println()

	// Example 4: Environment-based configuration
	fmt.Println("4. Using environment-based configuration...")
	// This falls back to simple logger if LOG_CONFIG_FILE not set
	envLogger := logging.NewFromYAMLEnv("LOG_CONFIG_FILE")
	envLogger.Info("Environment logger (falls back to simple if no env var set)")

	fmt.Println()

	// Example 5: Structured logging with YAML config
	fmt.Println("5. Structured logging with YAML configuration...")
	structuredLogger, err := logging.NewFromYAMLFile("../../configs/microservice.yaml")
	if err != nil {
		log.Printf("Error loading microservice config: %v", err)
	} else {
		// The YAML config includes static fields, so they appear automatically
		structuredLogger.Info("User action performed",
			"user_id", 12345,
			"action", "login",
			"ip_address", "192.168.1.100",
			"success", true,
		)

		structuredLogger.Error("Database connection failed",
			"database", "postgres",
			"host", "db.example.com",
			"error", "connection timeout",
			"retry_count", 3,
		)
	}

	fmt.Println()

	// Example 6: Demonstrate preset usage
	fmt.Println("6. Using preset configurations...")

	presets := []string{"minimal", "debug"}
	for _, preset := range presets {
		presetConfig := fmt.Sprintf(`preset: %s
static_fields:
  preset_demo: %s
`, preset, preset)

		presetLogger, err := logging.LoadFromYAMLString(presetConfig)
		if err != nil {
			log.Printf("Error with preset %s: %v", preset, err)
			continue
		}

		fmt.Printf("   Using %s preset:\n", preset)
		presetLogger.Info("Message with preset configuration")
		if preset == "debug" {
			presetLogger.Debug("Debug message (only shows with debug preset)")
		}
	}

	fmt.Println("\n=== Example completed ===")
}

// Structured Logging Examples
// This example demonstrates how to add structured data to log messages

package main

import (
	"github.com/ocrosby/go-logging/pkg/logging"
)

func main() {
	// Create a JSON logger for structured logging
	logger := logging.NewEasyJSON()

	// Basic structured logging with key-value pairs
	logger.Info("User login",
		"user_id", 12345,
		"email", "user@example.com",
		"ip", "192.168.1.100",
		"success", true,
	)

	// Logger with static fields - great for service context
	serviceLogger := logging.NewEasyBuilder().
		JSON().
		Field("service", "user-service").
		Field("version", "1.2.0").
		Field("environment", "production").
		Build()

	serviceLogger.Info("Processing request",
		"request_id", "req-abc123",
		"method", "POST",
		"path", "/api/users",
		"duration_ms", 42,
	)

	// More complex structured data
	serviceLogger.Warn("Rate limit exceeded",
		"user_id", 67890,
		"current_requests", 1000,
		"limit", 500,
		"window", "1h",
		"action", "blocked",
		"metadata", map[string]any{
			"source":    "rate_limiter",
			"algorithm": "sliding_window",
		},
	)

	// Error logging with structured context
	serviceLogger.Error("Database connection failed",
		"database", "postgres",
		"host", "db.example.com",
		"port", 5432,
		"error", "connection timeout",
		"retry_count", 3,
		"last_success", "2023-01-15T10:30:00Z",
	)

	// Using advanced configuration for payment logging
	paymentLogger := logging.NewEasyBuilder().
		Level(logging.InfoLevel).
		JSON().
		Field("component", "payment").
		Build()

	// Payment processing with structured data
	paymentLogger.Info("Payment processed",
		"transaction_id", "txn-xyz789",
		"amount", 99.99,
		"currency", "USD",
		"payment_method", "credit_card",
		"merchant_id", "merchant_123",
		"customer_id", "customer_456",
		"status", "completed",
	)
}

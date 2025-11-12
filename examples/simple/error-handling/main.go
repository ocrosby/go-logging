// Error Handling Examples
// This example demonstrates best practices for logging errors and failures

package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ocrosby/go-logging/pkg/logging"
)

func main() {
	logger := logging.NewEasyJSON()

	// Basic error logging
	err := processPayment("invalid-card")
	if err != nil {
		logger.Error("Payment processing failed",
			"error", err.Error(),
			"user_id", "user_123",
			"amount", 99.99,
		)
	}

	// Error with structured context
	ctx := logging.NewContextWithTrace()
	ctx = logging.WithRequestID(ctx, "req_payment_001")

	err = validateUser(ctx, "invalid_user")
	if err != nil {
		logger.ErrorContext(ctx, "User validation failed",
			"error", err.Error(),
			"attempted_user", "invalid_user",
			"validation_step", "existence_check",
		)
	}

	// Different error scenarios
	demonstrateErrorScenarios(logger)

	// Recovery from panics
	demonstratePanicRecovery(logger)

	// Warning vs Error distinction
	demonstrateWarningVsError(logger)
}

func processPayment(cardNumber string) error {
	if cardNumber == "invalid-card" {
		return errors.New("invalid credit card number")
	}
	return nil
}

func validateUser(ctx context.Context, userID string) error {
	if userID == "invalid_user" {
		return fmt.Errorf("user %s does not exist", userID)
	}
	return nil
}

func demonstrateErrorScenarios(logger logging.Logger) {
	ctx := logging.NewContextWithTrace()

	// Network error
	logger.ErrorContext(ctx, "Database connection failed",
		"error", "connection timeout after 30s",
		"database", "postgresql",
		"host", "db.example.com",
		"port", 5432,
		"retry_count", 3,
		"severity", "high",
	)

	// Validation error
	logger.WarnContext(ctx, "Invalid input received",
		"field", "email",
		"value", "not-an-email",
		"error", "invalid email format",
		"user_id", "user_456",
		"action_taken", "rejected_request",
	)

	// Business logic error
	logger.ErrorContext(ctx, "Insufficient funds",
		"user_id", "user_789",
		"requested_amount", 1000.00,
		"available_balance", 50.00,
		"error", "transaction denied",
		"account_type", "checking",
	)

	// Third-party service error
	logger.ErrorContext(ctx, "External API call failed",
		"service", "payment_processor",
		"endpoint", "https://api.payments.com/charge",
		"http_status", 503,
		"error", "service temporarily unavailable",
		"retry_after", "60s",
	)
}

func demonstratePanicRecovery(logger logging.Logger) {
	defer func() {
		if r := recover(); r != nil {
			logger.Error("Panic recovered",
				"panic", fmt.Sprintf("%v", r),
				"function", "demonstratePanicRecovery",
				"timestamp", time.Now().Unix(),
			)
		}
	}()

	// Simulate a panic
	var nilPointer *string
	logger.Info("About to cause a panic...")
	_ = *nilPointer // This will panic
}

func demonstrateWarningVsError(logger logging.Logger) {
	ctx := logging.NewContextWithTrace()

	// Warning: Something unusual but handled
	logger.WarnContext(ctx, "Slow query detected",
		"query", "SELECT * FROM users WHERE age > 25",
		"duration_ms", 5000,
		"threshold_ms", 1000,
		"table", "users",
		"rows_affected", 10000,
		"action_taken", "query_completed",
	)

	// Error: Something that prevents normal operation
	logger.ErrorContext(ctx, "Cache miss caused fallback to database",
		"cache_key", "user_profile_123",
		"fallback_duration_ms", 250,
		"cache_status", "unavailable",
		"error", "redis connection lost",
		"impact", "increased_latency",
	)

	// Critical: System-level issue
	logger.Error("Disk space critically low",
		"available_space_gb", 0.5,
		"threshold_gb", 1.0,
		"partition", "/var/log",
		"alert_level", "critical",
		"action_required", "immediate_cleanup",
	)
}

// Context-Aware Logging Examples
// This example demonstrates how to use context for request tracing and user information

package main

import (
	"context"

	"github.com/ocrosby/go-logging/pkg/logging"
)

func main() {
	logger := logging.NewEasyJSON()

	// Basic context logging - automatically generates trace IDs
	ctx := logging.NewContextWithTrace()
	logger.InfoContext(ctx, "Request started")
	logger.InfoContext(ctx, "Processing user data")
	logger.InfoContext(ctx, "Request completed")

	// Custom context with user and request information
	userCtx := logging.WithRequestID(ctx, "req_abc123")
	userCtx = logging.WithCorrelationID(userCtx, "user_12345")

	logger.InfoContext(userCtx, "User action logged",
		"action", "profile_update",
		"fields_changed", []string{"email", "name"},
	)

	// Different operations with same context maintain trace
	processOrder(userCtx, logger)
	sendNotification(userCtx, logger)

	// Nested operations preserve context
	handleUserRequest(logger)
}

func processOrder(ctx context.Context, logger logging.Logger) {
	logger.InfoContext(ctx, "Processing order",
		"order_id", "order_789",
		"amount", 49.99,
		"status", "processing",
	)

	// Simulate some processing steps
	logger.InfoContext(ctx, "Validating payment method")
	logger.InfoContext(ctx, "Checking inventory")
	logger.InfoContext(ctx, "Order confirmed",
		"estimated_delivery", "2-3 days",
	)
}

func sendNotification(ctx context.Context, logger logging.Logger) {
	logger.InfoContext(ctx, "Sending notification",
		"type", "email",
		"template", "order_confirmation",
		"recipient", "user@example.com",
	)
}

func handleUserRequest(logger logging.Logger) {
	// Create request context
	ctx := logging.NewContextWithTrace()
	ctx = logging.WithRequestID(ctx, "req_def456")

	logger.InfoContext(ctx, "Handling user request")

	// Pass context through the call chain
	authenticateUser(ctx, logger)
	authorizeUser(ctx, logger)
	processUserData(ctx, logger)

	logger.InfoContext(ctx, "Request handling complete")
}

func authenticateUser(ctx context.Context, logger logging.Logger) {
	logger.InfoContext(ctx, "Authenticating user",
		"method", "jwt",
		"token_valid", true,
	)
}

func authorizeUser(ctx context.Context, logger logging.Logger) {
	userCtx := logging.WithCorrelationID(ctx, "user_67890")
	logger.InfoContext(userCtx, "Checking user permissions",
		"resource", "/api/users/profile",
		"method", "PUT",
		"authorized", true,
	)
}

func processUserData(ctx context.Context, logger logging.Logger) {
	logger.InfoContext(ctx, "Processing user data",
		"operation", "profile_update",
		"fields", []string{"name", "email", "preferences"},
		"duration_ms", 45,
	)
}

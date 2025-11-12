package main

import (
	"errors"

	"github.com/ocrosby/go-logging/pkg/logging"
)

func main() {
	// Much simpler setup - just one line!
	logger := logging.NewEasyBuilder().
		Debug(). // Set to debug level
		JSON().  // Use JSON format
		Build()

	// Everything else stays the same
	logger.Fluent().Info().
		Str("service", "fluent-example").
		Str("version", "1.0.0").
		Msg("Application started with fluent interface")

	logger.Fluent().Debug().
		Int("user_id", 12345).
		Str("username", "john_doe").
		Bool("active", true).
		Msg("User details")

	err := errors.New("connection timeout")
	logger.Fluent().Error().
		Err(err).
		Str("host", "db.example.com").
		Int("port", 5432).
		Msg("Database connection failed")

	ctx := logging.NewContextWithTrace()
	ctx = logging.WithRequestID(ctx, "req-456")

	logger.Fluent().Info().
		Ctx(ctx).
		Str("operation", "fetch_user").
		Msgf("Processing request for user %s", "john_doe")
}

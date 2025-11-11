package main

import "context"

func main() {
	logger := InitializeLogger()

	logger.Info("Application started using Wire DI")

	logger.Fluent().Info().
		Str("service", "di-example").
		Str("version", "1.0.0").
		Msg("Demonstrating dependency injection")

	ctx := context.Background()
	logger.InfoContext(ctx, "Logging with context")

	userLogger := logger.WithField("user_id", 12345)
	userLogger.Info("User-specific logging")
}

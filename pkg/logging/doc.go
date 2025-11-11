// Package logging provides a configurable, production-ready logging library
// for Go with support for structured logging, request tracing, and fluent interfaces.
//
// # Features
//
//   - Multiple log levels: TRACE, DEBUG, INFO, WARN, ERROR, CRITICAL
//   - Multiple output formats: Text and JSON
//   - Fluent interface for expressive logging
//   - Request tracing with trace IDs, request IDs, and correlation IDs
//   - Structured logging with contextual fields
//   - HTTP middleware for automatic request tracing
//   - Sensitive data redaction
//   - Thread-safe concurrent logging
//   - Environment variable configuration
//
// # Quick Start
//
// Basic logging:
//
//	logger := logging.NewWithLevel(logging.InfoLevel)
//	logger.Info("Application started")
//	logger.Error("An error occurred: %v", err)
//
// Structured logging:
//
//	logger := logger.WithFields(map[string]interface{}{
//		"service": "api",
//		"version": "1.0.0",
//	})
//	logger.Info("Server starting")
//
// Fluent interface:
//
//	logger.Fluent().Info().
//		Str("service", "api").
//		Int("port", 8080).
//		Msg("Server started")
//
// Request tracing:
//
//	ctx := logging.NewContextWithTrace()
//	logger.InfoContext(ctx, "Processing request")
//
// # Configuration
//
// Using the builder pattern:
//
//	config := logging.NewConfig().
//		WithLevel(logging.DebugLevel).
//		WithJSONFormat().
//		WithStaticField("service", "my-app").
//		Build()
//	logger := logging.NewStandardLogger(config)
//
// From environment variables:
//
//	logger := logging.NewFromEnvironment()
//
// Supported environment variables:
//   - LOG_LEVEL: TRACE, DEBUG, INFO, WARN, ERROR, CRITICAL
//   - LOG_FORMAT: json (text is default)
//
// # HTTP Middleware
//
// Automatic request tracing:
//
//	mux := http.NewServeMux()
//	mux.HandleFunc("/api", handler)
//	handler := logging.TracingMiddleware(logger)(mux)
//	http.ListenAndServe(":8080", handler)
//
// The middleware automatically:
//   - Generates trace IDs for requests
//   - Extracts X-Trace-ID, X-Request-ID, X-Correlation-ID headers
//   - Logs request start and completion
//   - Adds trace ID to response headers
//
// # Design Principles
//
// This library follows SOLID principles:
//   - Single Responsibility: Each component has a focused purpose
//   - Open/Closed: Extensible through interfaces
//   - Liskov Substitution: Logger interface is easily mockable
//   - Interface Segregation: Clean, minimal interfaces
//   - Dependency Injection: Configuration via builder pattern
package logging

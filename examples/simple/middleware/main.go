// HTTP Middleware Logging Examples
// This example demonstrates how to integrate logging with HTTP middleware

package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ocrosby/go-logging/pkg/logging"
)

func main() {
	// Create a logger for the HTTP server
	logger := logging.NewEasyBuilder().
		Level(logging.InfoLevel).
		JSON().
		Field("service", "web-server").
		Field("version", "1.0.0").
		Build()

	// Create HTTP handlers
	mux := http.NewServeMux()

	// Simple endpoints
	mux.HandleFunc("/health", healthHandler(logger))
	mux.HandleFunc("/api/users", usersHandler(logger))
	mux.HandleFunc("/api/orders", ordersHandler(logger))

	// Chain middleware: logging -> auth -> CORS -> handler
	handler := loggingMiddleware(logger)(
		authMiddleware(logger)(
			corsMiddleware(logger)(mux),
		),
	)

	logger.Info("Starting HTTP server",
		"port", 8080,
		"endpoints", []string{"/health", "/api/users", "/api/orders"},
	)

	if err := http.ListenAndServe(":8080", handler); err != nil {
		logger.Error("Server failed to start",
			"error", err.Error(),
			"port", 8080,
		)
	}
}

// HTTP request logging middleware
func loggingMiddleware(logger logging.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Create context with trace information
			ctx := logging.NewContextWithTrace()
			ctx = logging.WithRequestID(ctx, fmt.Sprintf("req_%d", time.Now().UnixNano()))

			// Add context to request
			r = r.WithContext(ctx)

			// Log incoming request
			logger.InfoContext(ctx, "HTTP request started",
				"method", r.Method,
				"path", r.URL.Path,
				"user_agent", r.UserAgent(),
				"remote_addr", r.RemoteAddr,
				"content_length", r.ContentLength,
			)

			// Wrap response writer to capture status code
			wrapped := &responseWriter{ResponseWriter: w, statusCode: 200}

			// Process request
			next.ServeHTTP(wrapped, r)

			// Log request completion
			duration := time.Since(start)
			logger.InfoContext(ctx, "HTTP request completed",
				"method", r.Method,
				"path", r.URL.Path,
				"status_code", wrapped.statusCode,
				"duration_ms", duration.Milliseconds(),
				"response_size", wrapped.bytesWritten,
			)
		})
	}
}

// Authentication middleware
func authMiddleware(logger logging.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			// Simple auth check (in real apps, validate JWT, etc.)
			authHeader := r.Header.Get("Authorization")

			if r.URL.Path != "/health" && authHeader == "" {
				logger.WarnContext(ctx, "Unauthorized request",
					"path", r.URL.Path,
					"method", r.Method,
					"reason", "missing_auth_header",
				)

				w.WriteHeader(http.StatusUnauthorized)
				fmt.Fprint(w, `{"error": "unauthorized"}`)
				return
			}

			if authHeader != "" {
				// Extract user info (simplified)
				userID := "user_123" // In reality, decode from token
				ctx = logging.WithCorrelationID(ctx, userID)
				r = r.WithContext(ctx)

				logger.InfoContext(ctx, "Request authenticated",
					"user_id", userID,
					"path", r.URL.Path,
				)
			}

			next.ServeHTTP(w, r)
		})
	}
}

// CORS middleware
func corsMiddleware(logger logging.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			origin := r.Header.Get("Origin")
			if origin != "" {
				logger.InfoContext(ctx, "CORS request detected",
					"origin", origin,
					"method", r.Method,
					"path", r.URL.Path,
				)

				w.Header().Set("Access-Control-Allow-Origin", "*")
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			}

			if r.Method == "OPTIONS" {
				logger.InfoContext(ctx, "CORS preflight request handled")
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// Health check handler
func healthHandler(logger logging.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		logger.InfoContext(ctx, "Health check requested")

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"status": "healthy", "timestamp": "`+time.Now().Format(time.RFC3339)+`"}`)
	}
}

// Users API handler
func usersHandler(logger logging.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		switch r.Method {
		case "GET":
			logger.InfoContext(ctx, "Fetching users list",
				"endpoint", "/api/users",
			)
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, `{"users": [{"id": 1, "name": "John"}, {"id": 2, "name": "Jane"}]}`)

		case "POST":
			logger.InfoContext(ctx, "Creating new user",
				"endpoint", "/api/users",
				"content_type", r.Header.Get("Content-Type"),
			)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			fmt.Fprint(w, `{"id": 3, "name": "New User", "created": true}`)

		default:
			logger.WarnContext(ctx, "Method not allowed",
				"method", r.Method,
				"endpoint", "/api/users",
			)
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprint(w, `{"error": "method not allowed"}`)
		}
	}
}

// Orders API handler
func ordersHandler(logger logging.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Simulate some processing time
		time.Sleep(50 * time.Millisecond)

		logger.InfoContext(ctx, "Processing order request",
			"endpoint", "/api/orders",
			"method", r.Method,
			"processing_time_ms", 50,
		)

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"orders": [{"id": "ord_123", "status": "completed"}]}`)
	}
}

// Response writer wrapper to capture status code and bytes written
type responseWriter struct {
	http.ResponseWriter
	statusCode   int
	bytesWritten int
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func (rw *responseWriter) Write(data []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(data)
	rw.bytesWritten += n
	return n, err
}

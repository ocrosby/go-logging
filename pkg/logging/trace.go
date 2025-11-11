// Package logging provides structured logging with request tracing support.
package logging

import (
	"context"
	"crypto/rand"
	"fmt"
)

type contextKey string

const (
	// TraceIDKey is the context key for trace identifiers that follow requests through the system.
	TraceIDKey contextKey = "trace_id"
	// RequestIDKey is the context key for unique request identifiers.
	RequestIDKey contextKey = "request_id"
	// CorrelationKey is the context key for correlation identifiers linking related requests.
	CorrelationKey contextKey = "correlation_id"
)

// NewTraceID generates a new unique trace identifier using UUID v4.
// Use this to create trace IDs for tracking requests through your system.
//
// Example:
//
//	traceID := logging.NewTraceID()
//	ctx := logging.WithTraceID(context.Background(), traceID)
func NewTraceID() string {
	var uuid [16]byte
	_, err := rand.Read(uuid[:])
	if err != nil {
		panic(err)
	}
	uuid[6] = (uuid[6] & 0x0f) | 0x40
	uuid[8] = (uuid[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:])
}

// WithTraceID returns a new context with the trace ID attached.
// The trace ID can be retrieved later with GetTraceID.
//
// Example:
//
//	ctx := logging.WithTraceID(r.Context(), "trace-123")
//	logger.InfoContext(ctx, "Processing request")
func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, TraceIDKey, traceID)
}

// GetTraceID retrieves the trace ID from the context.
// Returns the trace ID and true if present, empty string and false otherwise.
func GetTraceID(ctx context.Context) (string, bool) {
	traceID, ok := ctx.Value(TraceIDKey).(string)
	return traceID, ok
}

// WithRequestID returns a new context with the request ID attached.
// The request ID can be retrieved later with GetRequestID.
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, RequestIDKey, requestID)
}

// GetRequestID retrieves the request ID from the context.
// Returns the request ID and true if present, empty string and false otherwise.
func GetRequestID(ctx context.Context) (string, bool) {
	requestID, ok := ctx.Value(RequestIDKey).(string)
	return requestID, ok
}

// WithCorrelationID returns a new context with the correlation ID attached.
// The correlation ID can be retrieved later with GetCorrelationID.
// Use correlation IDs to link related requests across different services.
func WithCorrelationID(ctx context.Context, correlationID string) context.Context {
	return context.WithValue(ctx, CorrelationKey, correlationID)
}

// GetCorrelationID retrieves the correlation ID from the context.
// Returns the correlation ID and true if present, empty string and false otherwise.
func GetCorrelationID(ctx context.Context) (string, bool) {
	correlationID, ok := ctx.Value(CorrelationKey).(string)
	return correlationID, ok
}

// NewContextWithTrace creates a new context with an automatically generated trace ID.
// This is a convenience function equivalent to WithTraceID(context.Background(), NewTraceID()).
//
// Example:
//
//	ctx := logging.NewContextWithTrace()
//	logger.InfoContext(ctx, "Starting new operation")
func NewContextWithTrace() context.Context {
	return WithTraceID(context.Background(), NewTraceID())
}

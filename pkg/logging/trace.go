package logging

import (
	"context"
	"github.com/google/uuid"
)

type contextKey string

const (
	TraceIDKey     contextKey = "trace_id"
	RequestIDKey   contextKey = "request_id"
	CorrelationKey contextKey = "correlation_id"
)

func NewTraceID() string {
	return uuid.New().String()
}

func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, TraceIDKey, traceID)
}

func GetTraceID(ctx context.Context) (string, bool) {
	traceID, ok := ctx.Value(TraceIDKey).(string)
	return traceID, ok
}

func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, RequestIDKey, requestID)
}

func GetRequestID(ctx context.Context) (string, bool) {
	requestID, ok := ctx.Value(RequestIDKey).(string)
	return requestID, ok
}

func WithCorrelationID(ctx context.Context, correlationID string) context.Context {
	return context.WithValue(ctx, CorrelationKey, correlationID)
}

func GetCorrelationID(ctx context.Context) (string, bool) {
	correlationID, ok := ctx.Value(CorrelationKey).(string)
	return correlationID, ok
}

func NewContextWithTrace() context.Context {
	return WithTraceID(context.Background(), NewTraceID())
}

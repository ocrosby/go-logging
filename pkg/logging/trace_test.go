package logging

import (
	"context"
	"testing"
)

func TestNewTraceID(t *testing.T) {
	id1 := NewTraceID()
	id2 := NewTraceID()

	if id1 == "" {
		t.Error("TraceID should not be empty")
	}

	if id1 == id2 {
		t.Error("TraceIDs should be unique")
	}
}

func TestWithTraceID(t *testing.T) {
	ctx := context.Background()
	traceID := "test-trace-id"

	ctx = WithTraceID(ctx, traceID)

	retrievedID, ok := GetTraceID(ctx)
	if !ok {
		t.Error("TraceID should be present in context")
	}

	if retrievedID != traceID {
		t.Errorf("Expected traceID %s, got %s", traceID, retrievedID)
	}
}

func TestWithRequestID(t *testing.T) {
	ctx := context.Background()
	requestID := "test-request-id"

	ctx = WithRequestID(ctx, requestID)

	retrievedID, ok := GetRequestID(ctx)
	if !ok {
		t.Error("RequestID should be present in context")
	}

	if retrievedID != requestID {
		t.Errorf("Expected requestID %s, got %s", requestID, retrievedID)
	}
}

func TestWithCorrelationID(t *testing.T) {
	ctx := context.Background()
	correlationID := "test-correlation-id"

	ctx = WithCorrelationID(ctx, correlationID)

	retrievedID, ok := GetCorrelationID(ctx)
	if !ok {
		t.Error("CorrelationID should be present in context")
	}

	if retrievedID != correlationID {
		t.Errorf("Expected correlationID %s, got %s", correlationID, retrievedID)
	}
}

func TestNewContextWithTrace(t *testing.T) {
	ctx := NewContextWithTrace()

	traceID, ok := GetTraceID(ctx)
	if !ok {
		t.Error("TraceID should be present in context")
	}

	if traceID == "" {
		t.Error("TraceID should not be empty")
	}
}

func TestGetTraceID_NotPresent(t *testing.T) {
	ctx := context.Background()

	_, ok := GetTraceID(ctx)
	if ok {
		t.Error("TraceID should not be present in empty context")
	}
}

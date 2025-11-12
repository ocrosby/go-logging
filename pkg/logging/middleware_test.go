package logging

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestTracingMiddleware(t *testing.T) {
	buf := &bytes.Buffer{}
	config := NewLoggerConfig().
		WithLevel(InfoLevel).
		WithWriter(buf).
		WithTextFormat().
		Build()
	logger := NewWithLoggerConfig(config)

	// Create a test handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check that context has trace ID
		traceID, ok := GetTraceID(r.Context())
		if !ok || traceID == "" {
			t.Error("expected trace ID in context")
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("test response"))
	})

	// Wrap with tracing middleware
	middleware := TracingMiddleware(logger)
	wrappedHandler := middleware(handler)

	// Create test request
	req := httptest.NewRequest("GET", "/test", nil)
	recorder := httptest.NewRecorder()

	// Execute request
	wrappedHandler.ServeHTTP(recorder, req)

	// Check response
	if recorder.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", recorder.Code)
	}

	// Check that logging occurred
	output := buf.String()
	if !strings.Contains(output, "Request started") && !strings.Contains(output, "Request completed") {
		t.Errorf("expected request logging in output, got: %s", output)
	}
}

func TestTracingMiddleware_WithExistingTraceID(t *testing.T) {
	buf := &bytes.Buffer{}
	config := NewLoggerConfig().
		WithLevel(InfoLevel).
		WithWriter(buf).
		WithTextFormat().
		Build()
	logger := NewWithLoggerConfig(config)

	existingTraceID := "existing-trace-123"

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceID, ok := GetTraceID(r.Context())
		if !ok || traceID != existingTraceID {
			t.Errorf("expected trace ID %s, got %s", existingTraceID, traceID)
		}
		w.WriteHeader(http.StatusOK)
	})

	middleware := TracingMiddleware(logger)
	wrappedHandler := middleware(handler)

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set(HeaderTraceID, existingTraceID)
	recorder := httptest.NewRecorder()

	wrappedHandler.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", recorder.Code)
	}
}

func TestTracingMiddleware_WithRequestAndCorrelationIDs(t *testing.T) {
	buf := &bytes.Buffer{}
	config := NewLoggerConfig().
		WithLevel(InfoLevel).
		WithWriter(buf).
		WithTextFormat().
		Build()
	logger := NewWithLoggerConfig(config)

	requestID := "req-123"
	correlationID := "corr-456"

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request ID
		gotRequestID, ok := GetRequestID(r.Context())
		if !ok || gotRequestID != requestID {
			t.Errorf("expected request ID %s, got %s", requestID, gotRequestID)
		}

		// Check correlation ID
		gotCorrelationID, ok := GetCorrelationID(r.Context())
		if !ok || gotCorrelationID != correlationID {
			t.Errorf("expected correlation ID %s, got %s", correlationID, gotCorrelationID)
		}

		w.WriteHeader(http.StatusOK)
	})

	middleware := TracingMiddleware(logger)
	wrappedHandler := middleware(handler)

	req := httptest.NewRequest("POST", "/api/test", nil)
	req.Header.Set(HeaderRequestID, requestID)
	req.Header.Set(HeaderCorrelationID, correlationID)
	recorder := httptest.NewRecorder()

	wrappedHandler.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", recorder.Code)
	}
}

func TestRequestLogger(t *testing.T) {
	buf := &bytes.Buffer{}
	config := NewLoggerConfig().
		WithLevel(InfoLevel).
		WithWriter(buf).
		WithTextFormat().
		Build()
	logger := NewWithLoggerConfig(config)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte("created"))
	})

	middleware := RequestLogger(logger, "User-Agent", "Content-Type")
	wrappedHandler := middleware(handler)

	req := httptest.NewRequest("POST", "/api/users", strings.NewReader(`{"name":"test"}`))
	req.Header.Set("User-Agent", "test-client/1.0")
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	wrappedHandler.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", recorder.Code)
	}

	output := buf.String()

	// Should log request
	if !strings.Contains(output, "HTTP Request") {
		t.Errorf("expected HTTP Request log message, got: %s", output)
	}
}

func TestRequestLogger_NilLogger(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// This will panic with nil logger - expected behavior
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic with nil logger")
		}
	}()

	middleware := RequestLogger(nil)
	wrappedHandler := middleware(handler)

	req := httptest.NewRequest("GET", "/test", nil)
	recorder := httptest.NewRecorder()

	wrappedHandler.ServeHTTP(recorder, req)
}

func TestResponseWriter_WriteHeader(t *testing.T) {
	recorder := httptest.NewRecorder()
	rw := &responseWriter{
		ResponseWriter: recorder,
		statusCode:     0,
	}

	rw.WriteHeader(http.StatusNotFound)

	if rw.statusCode != http.StatusNotFound {
		t.Errorf("expected status code %d, got %d", http.StatusNotFound, rw.statusCode)
	}

	if recorder.Code != http.StatusNotFound {
		t.Errorf("expected recorder status %d, got %d", http.StatusNotFound, recorder.Code)
	}
}

func TestResponseWriter_Write(t *testing.T) {
	recorder := httptest.NewRecorder()
	rw := &responseWriter{
		ResponseWriter: recorder,
	}

	data := []byte("test response data")
	n, err := rw.Write(data)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if n != len(data) {
		t.Errorf("expected %d bytes written, got %d", len(data), n)
	}

	if rw.written != int64(len(data)) {
		t.Errorf("expected written count %d, got %d", len(data), rw.written)
	}

	if recorder.Body.String() != string(data) {
		t.Errorf("expected body %s, got %s", string(data), recorder.Body.String())
	}
}

func TestMiddlewareConstants(t *testing.T) {
	// Test that middleware header constants are defined
	if HeaderTraceID == "" {
		t.Error("HeaderTraceID should not be empty")
	}

	if HeaderRequestID == "" {
		t.Error("HeaderRequestID should not be empty")
	}

	if HeaderCorrelationID == "" {
		t.Error("HeaderCorrelationID should not be empty")
	}

	// Verify expected values
	if HeaderTraceID != "X-Trace-ID" {
		t.Errorf("expected HeaderTraceID to be 'X-Trace-ID', got %s", HeaderTraceID)
	}

	if HeaderRequestID != "X-Request-ID" {
		t.Errorf("expected HeaderRequestID to be 'X-Request-ID', got %s", HeaderRequestID)
	}

	if HeaderCorrelationID != "X-Correlation-ID" {
		t.Errorf("expected HeaderCorrelationID to be 'X-Correlation-ID', got %s", HeaderCorrelationID)
	}
}

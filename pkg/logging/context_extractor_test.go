package logging

import (
	"context"
	"log/slog"
	"testing"
)

func TestTraceContextExtractor(t *testing.T) {
	extractor := TraceContextExtractor()

	ctx := context.Background()
	ctx = WithTraceID(ctx, "trace-123")
	ctx = WithRequestID(ctx, "req-456")
	ctx = WithCorrelationID(ctx, "corr-789")

	attrs := extractor.Extract(ctx)

	if len(attrs) != 3 {
		t.Errorf("Expected 3 attributes, got %d", len(attrs))
	}

	foundTrace := false
	foundRequest := false
	foundCorrelation := false

	for _, attr := range attrs {
		switch attr.Key {
		case "trace_id":
			if attr.Value.String() == "trace-123" {
				foundTrace = true
			}
		case "request_id":
			if attr.Value.String() == "req-456" {
				foundRequest = true
			}
		case "correlation_id":
			if attr.Value.String() == "corr-789" {
				foundCorrelation = true
			}
		}
	}

	if !foundTrace {
		t.Error("Expected to find trace_id attribute")
	}
	if !foundRequest {
		t.Error("Expected to find request_id attribute")
	}
	if !foundCorrelation {
		t.Error("Expected to find correlation_id attribute")
	}
}

func TestTraceContextExtractor_EmptyContext(t *testing.T) {
	extractor := TraceContextExtractor()
	ctx := context.Background()

	attrs := extractor.Extract(ctx)

	if len(attrs) != 0 {
		t.Errorf("Expected 0 attributes, got %d", len(attrs))
	}
}

func TestStringContextExtractor(t *testing.T) {
	const userKey ContextKey = "user"
	extractor := StringContextExtractor(userKey, "username")

	ctx := context.WithValue(context.Background(), userKey, "john_doe")
	attrs := extractor.Extract(ctx)

	if len(attrs) != 1 {
		t.Fatalf("Expected 1 attribute, got %d", len(attrs))
	}

	if attrs[0].Key != "username" {
		t.Errorf("Expected key 'username', got '%s'", attrs[0].Key)
	}

	if attrs[0].Value.String() != "john_doe" {
		t.Errorf("Expected value 'john_doe', got '%s'", attrs[0].Value.String())
	}
}

func TestStringContextExtractor_MissingValue(t *testing.T) {
	const userKey ContextKey = "user"
	extractor := StringContextExtractor(userKey, "username")

	ctx := context.Background()
	attrs := extractor.Extract(ctx)

	if len(attrs) != 0 {
		t.Errorf("Expected 0 attributes, got %d", len(attrs))
	}
}

func TestIntContextExtractor(t *testing.T) {
	const ageKey ContextKey = "age"
	extractor := IntContextExtractor(ageKey, "user_age")

	ctx := context.WithValue(context.Background(), ageKey, 30)
	attrs := extractor.Extract(ctx)

	if len(attrs) != 1 {
		t.Fatalf("Expected 1 attribute, got %d", len(attrs))
	}

	if attrs[0].Key != "user_age" {
		t.Errorf("Expected key 'user_age', got '%s'", attrs[0].Key)
	}

	if attrs[0].Value.Int64() != 30 {
		t.Errorf("Expected value 30, got %d", attrs[0].Value.Int64())
	}
}

func TestInt64ContextExtractor(t *testing.T) {
	const idKey ContextKey = "user_id"
	extractor := Int64ContextExtractor(idKey, "id")

	ctx := context.WithValue(context.Background(), idKey, int64(123456789))
	attrs := extractor.Extract(ctx)

	if len(attrs) != 1 {
		t.Fatalf("Expected 1 attribute, got %d", len(attrs))
	}

	if attrs[0].Key != "id" {
		t.Errorf("Expected key 'id', got '%s'", attrs[0].Key)
	}

	if attrs[0].Value.Int64() != 123456789 {
		t.Errorf("Expected value 123456789, got %d", attrs[0].Value.Int64())
	}
}

func TestBoolContextExtractor(t *testing.T) {
	const flagKey ContextKey = "is_admin"
	extractor := BoolContextExtractor(flagKey, "admin")

	ctx := context.WithValue(context.Background(), flagKey, true)
	attrs := extractor.Extract(ctx)

	if len(attrs) != 1 {
		t.Fatalf("Expected 1 attribute, got %d", len(attrs))
	}

	if attrs[0].Key != "admin" {
		t.Errorf("Expected key 'admin', got '%s'", attrs[0].Key)
	}

	if attrs[0].Value.Bool() != true {
		t.Errorf("Expected value true, got %v", attrs[0].Value.Bool())
	}
}

func TestCustomContextExtractor(t *testing.T) {
	const dataKey ContextKey = "custom_data"
	extractor := CustomContextExtractor(dataKey, "data")

	type CustomData struct {
		Name  string
		Value int
	}

	data := CustomData{Name: "test", Value: 42}
	ctx := context.WithValue(context.Background(), dataKey, data)
	attrs := extractor.Extract(ctx)

	if len(attrs) != 1 {
		t.Fatalf("Expected 1 attribute, got %d", len(attrs))
	}

	if attrs[0].Key != "data" {
		t.Errorf("Expected key 'data', got '%s'", attrs[0].Key)
	}
}

func TestCompositeContextExtractor(t *testing.T) {
	const (
		userKey ContextKey = "user"
		ageKey  ContextKey = "age"
	)

	composite := NewCompositeContextExtractor(
		TraceContextExtractor(),
		StringContextExtractor(userKey, "username"),
		IntContextExtractor(ageKey, "user_age"),
	)

	ctx := context.Background()
	ctx = WithTraceID(ctx, "trace-123")
	ctx = context.WithValue(ctx, userKey, "john_doe")
	ctx = context.WithValue(ctx, ageKey, 30)

	attrs := composite.Extract(ctx)

	if len(attrs) != 3 {
		t.Errorf("Expected 3 attributes, got %d", len(attrs))
	}

	foundTrace := false
	foundUser := false
	foundAge := false

	for _, attr := range attrs {
		switch attr.Key {
		case "trace_id":
			foundTrace = true
		case "username":
			foundUser = true
		case "user_age":
			foundAge = true
		}
	}

	if !foundTrace {
		t.Error("Expected to find trace_id attribute")
	}
	if !foundUser {
		t.Error("Expected to find username attribute")
	}
	if !foundAge {
		t.Error("Expected to find user_age attribute")
	}
}

func TestCompositeContextExtractor_Add(t *testing.T) {
	const userKey ContextKey = "user"

	composite := NewCompositeContextExtractor()
	composite.Add(TraceContextExtractor())
	composite.Add(StringContextExtractor(userKey, "username"))

	ctx := context.Background()
	ctx = WithTraceID(ctx, "trace-123")
	ctx = context.WithValue(ctx, userKey, "john_doe")

	attrs := composite.Extract(ctx)

	if len(attrs) != 2 {
		t.Errorf("Expected 2 attributes, got %d", len(attrs))
	}
}

func BenchmarkTraceContextExtractor(b *testing.B) {
	extractor := TraceContextExtractor()
	ctx := context.Background()
	ctx = WithTraceID(ctx, "trace-123")
	ctx = WithRequestID(ctx, "req-456")
	ctx = WithCorrelationID(ctx, "corr-789")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = extractor.Extract(ctx)
	}
}

func BenchmarkCompositeContextExtractor(b *testing.B) {
	const (
		userKey ContextKey = "user"
		ageKey  ContextKey = "age"
	)

	composite := NewCompositeContextExtractor(
		TraceContextExtractor(),
		StringContextExtractor(userKey, "username"),
		IntContextExtractor(ageKey, "user_age"),
	)

	ctx := context.Background()
	ctx = WithTraceID(ctx, "trace-123")
	ctx = context.WithValue(ctx, userKey, "john_doe")
	ctx = context.WithValue(ctx, ageKey, 30)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = composite.Extract(ctx)
	}
}

func ExampleTraceContextExtractor() {
	extractor := TraceContextExtractor()

	ctx := context.Background()
	ctx = WithTraceID(ctx, "trace-123")
	ctx = WithRequestID(ctx, "req-456")

	attrs := extractor.Extract(ctx)
	for _, attr := range attrs {
		_ = slog.String(attr.Key, attr.Value.String())
	}
}

func ExampleCompositeContextExtractor() {
	const userKey ContextKey = "user"

	composite := NewCompositeContextExtractor(
		TraceContextExtractor(),
		StringContextExtractor(userKey, "username"),
	)

	ctx := context.Background()
	ctx = WithTraceID(ctx, "trace-123")
	ctx = context.WithValue(ctx, userKey, "john_doe")

	attrs := composite.Extract(ctx)
	for _, attr := range attrs {
		_ = slog.Any(attr.Key, attr.Value)
	}
}

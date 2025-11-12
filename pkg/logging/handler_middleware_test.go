package logging

import (
	"bytes"
	"context"
	"log/slog"
	"regexp"
	"strings"
	"testing"
	"time"
)

func TestMiddlewareHandler_ChainExecution(t *testing.T) {
	var buf bytes.Buffer
	handler := slog.NewTextHandler(&buf, nil)

	var calls []string

	middleware1 := handlerMiddlewareFunc(func(ctx context.Context, record slog.Record, next HandlerFunc) error {
		calls = append(calls, "middleware1-before")
		err := next(ctx, record)
		calls = append(calls, "middleware1-after")
		return err
	})

	middleware2 := handlerMiddlewareFunc(func(ctx context.Context, record slog.Record, next HandlerFunc) error {
		calls = append(calls, "middleware2-before")
		err := next(ctx, record)
		calls = append(calls, "middleware2-after")
		return err
	})

	mh := NewMiddlewareHandler(handler, middleware1, middleware2)

	record := slog.NewRecord(time.Now(), slog.LevelInfo, "test message", 0)
	err := mh.Handle(context.Background(), record)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expectedOrder := []string{
		"middleware1-before",
		"middleware2-before",
		"middleware2-after",
		"middleware1-after",
	}

	if len(calls) != len(expectedOrder) {
		t.Fatalf("Expected %d calls, got %d", len(expectedOrder), len(calls))
	}

	for i, expected := range expectedOrder {
		if calls[i] != expected {
			t.Errorf("Call %d: expected '%s', got '%s'", i, expected, calls[i])
		}
	}
}

func TestTimestampMiddleware(t *testing.T) {
	var buf bytes.Buffer
	handler := slog.NewJSONHandler(&buf, nil)

	mh := NewMiddlewareHandler(handler, TimestampMiddleware())

	record := slog.NewRecord(time.Now(), slog.LevelInfo, "test message", 0)
	err := mh.Handle(context.Background(), record)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, `"time"`) {
		t.Error("Expected timestamp in output")
	}
}

func TestContextExtractorMiddleware(t *testing.T) {
	var buf bytes.Buffer
	handler := slog.NewJSONHandler(&buf, nil)

	extractor := TraceContextExtractor()
	mh := NewMiddlewareHandler(handler, ContextExtractorMiddleware(extractor))

	ctx := WithTraceID(context.Background(), "trace-123")
	record := slog.NewRecord(time.Now(), slog.LevelInfo, "test message", 0)

	err := mh.Handle(ctx, record)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "trace-123") {
		t.Errorf("Expected trace ID in output, got: %s", output)
	}
}

func TestLevelFilterMiddleware(t *testing.T) {
	var buf bytes.Buffer
	handler := slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})

	mh := NewMiddlewareHandler(handler, LevelFilterMiddleware(slog.LevelWarn))

	debugRecord := slog.NewRecord(time.Now(), slog.LevelDebug, "debug message", 0)
	err := mh.Handle(context.Background(), debugRecord)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	warnRecord := slog.NewRecord(time.Now(), slog.LevelWarn, "warn message", 0)
	err = mh.Handle(context.Background(), warnRecord)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	output := buf.String()
	if strings.Contains(output, "debug message") {
		t.Error("Debug message should be filtered out")
	}
	if !strings.Contains(output, "warn message") {
		t.Error("Warn message should be present")
	}
}

func TestSamplingMiddleware(t *testing.T) {
	var buf bytes.Buffer
	handler := slog.NewTextHandler(&buf, nil)

	mh := NewMiddlewareHandler(handler, SamplingMiddleware(2))

	for i := 0; i < 10; i++ {
		record := slog.NewRecord(time.Now(), slog.LevelInfo, "test message", 0)
		_ = mh.Handle(context.Background(), record)
	}

	output := buf.String()
	count := strings.Count(output, "test message")

	if count != 5 {
		t.Errorf("Expected 5 messages (sampled every 2), got %d", count)
	}
}

func TestStaticFieldsMiddleware(t *testing.T) {
	var buf bytes.Buffer
	handler := slog.NewJSONHandler(&buf, nil)

	fields := map[string]interface{}{
		"service": "test-service",
		"version": "1.0.0",
	}

	mh := NewMiddlewareHandler(handler, StaticFieldsMiddleware(fields))

	record := slog.NewRecord(time.Now(), slog.LevelInfo, "test message", 0)
	err := mh.Handle(context.Background(), record)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, `"service":"test-service"`) {
		t.Error("Expected service field in output")
	}
	if !strings.Contains(output, `"version":"1.0.0"`) {
		t.Error("Expected version field in output")
	}
}

func TestRedactionMiddleware(t *testing.T) {
	var buf bytes.Buffer
	handler := slog.NewTextHandler(&buf, nil)

	redactor := NewRegexRedactor(
		regexp.MustCompile(`password=\w+`),
		"password=***REDACTED***",
	)

	mh := NewMiddlewareHandler(handler, RedactionMiddleware(redactor))

	record := slog.NewRecord(time.Now(), slog.LevelInfo, "user logged in with password=secret123", 0)
	err := mh.Handle(context.Background(), record)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	output := buf.String()
	if strings.Contains(output, "secret123") {
		t.Error("Password should be redacted")
	}
	if !strings.Contains(output, "***REDACTED***") {
		t.Error("Expected redacted placeholder")
	}
}

func TestLoggingMiddleware(t *testing.T) {
	var buf bytes.Buffer
	handler := slog.NewTextHandler(&buf, nil)

	var beforeCalled, afterCalled bool

	onBefore := func(ctx context.Context, record slog.Record) {
		beforeCalled = true
	}

	onAfter := func(ctx context.Context, record slog.Record, err error) {
		afterCalled = true
	}

	mh := NewMiddlewareHandler(handler, NewLoggingMiddleware(onBefore, onAfter))

	record := slog.NewRecord(time.Now(), slog.LevelInfo, "test message", 0)
	err := mh.Handle(context.Background(), record)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if !beforeCalled {
		t.Error("Expected onBefore to be called")
	}
	if !afterCalled {
		t.Error("Expected onAfter to be called")
	}
}

func TestMetricsMiddleware(t *testing.T) {
	var buf bytes.Buffer
	handler := slog.NewTextHandler(&buf, nil)

	var recordedLevel slog.Level
	recordMetric := func(level slog.Level) {
		recordedLevel = level
	}

	mh := NewMiddlewareHandler(handler, MetricsMiddleware(recordMetric))

	record := slog.NewRecord(time.Now(), slog.LevelWarn, "test message", 0)
	err := mh.Handle(context.Background(), record)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if recordedLevel != slog.LevelWarn {
		t.Errorf("Expected level WARN, got %v", recordedLevel)
	}
}

func TestMiddlewareHandler_WithAttrs(t *testing.T) {
	var buf bytes.Buffer
	handler := slog.NewJSONHandler(&buf, nil)

	mh := NewMiddlewareHandler(handler, TimestampMiddleware())
	mhWithAttrs := mh.WithAttrs([]slog.Attr{
		slog.String("service", "test"),
	})

	record := slog.NewRecord(time.Now(), slog.LevelInfo, "test message", 0)
	err := mhWithAttrs.Handle(context.Background(), record)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, `"service":"test"`) {
		t.Error("Expected service attribute in output")
	}
}

func TestMiddlewareHandler_WithGroup(t *testing.T) {
	var buf bytes.Buffer
	handler := slog.NewJSONHandler(&buf, nil)

	mh := NewMiddlewareHandler(handler, TimestampMiddleware())
	mhWithGroup := mh.WithGroup("app")

	record := slog.NewRecord(time.Now(), slog.LevelInfo, "test message", 0)
	record.AddAttrs(slog.String("version", "1.0"))
	err := mhWithGroup.Handle(context.Background(), record)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, `"app":{"version":"1.0"}`) {
		t.Errorf("Expected grouped attribute in output, got: %s", output)
	}
}

func TestMultipleMiddlewares(t *testing.T) {
	var buf bytes.Buffer
	handler := slog.NewJSONHandler(&buf, nil)

	fields := map[string]interface{}{
		"service": "test",
	}

	mh := NewMiddlewareHandler(
		handler,
		TimestampMiddleware(),
		StaticFieldsMiddleware(fields),
		ContextExtractorMiddleware(TraceContextExtractor()),
	)

	ctx := WithTraceID(context.Background(), "trace-123")
	record := slog.NewRecord(time.Now(), slog.LevelInfo, "test message", 0)

	err := mh.Handle(ctx, record)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, `"service":"test"`) {
		t.Error("Expected service field")
	}
	if !strings.Contains(output, "trace-123") {
		t.Error("Expected trace ID")
	}
	if !strings.Contains(output, `"time"`) {
		t.Error("Expected timestamp")
	}
}

func BenchmarkMiddlewareHandler_NoMiddleware(b *testing.B) {
	var buf bytes.Buffer
	handler := slog.NewTextHandler(&buf, nil)
	mh := NewMiddlewareHandler(handler)

	ctx := context.Background()
	record := slog.NewRecord(time.Now(), slog.LevelInfo, "benchmark message", 0)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = mh.Handle(ctx, record)
	}
}

func BenchmarkMiddlewareHandler_SingleMiddleware(b *testing.B) {
	var buf bytes.Buffer
	handler := slog.NewTextHandler(&buf, nil)
	mh := NewMiddlewareHandler(handler, TimestampMiddleware())

	ctx := context.Background()
	record := slog.NewRecord(time.Now(), slog.LevelInfo, "benchmark message", 0)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = mh.Handle(ctx, record)
	}
}

func BenchmarkMiddlewareHandler_MultipleMiddlewares(b *testing.B) {
	var buf bytes.Buffer
	handler := slog.NewTextHandler(&buf, nil)

	mh := NewMiddlewareHandler(
		handler,
		TimestampMiddleware(),
		StaticFieldsMiddleware(map[string]interface{}{"service": "test"}),
		ContextExtractorMiddleware(TraceContextExtractor()),
	)

	ctx := WithTraceID(context.Background(), "trace-123")
	record := slog.NewRecord(time.Now(), slog.LevelInfo, "benchmark message", 0)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = mh.Handle(ctx, record)
	}
}

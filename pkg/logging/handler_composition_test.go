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

func TestNewConditionalHandler(t *testing.T) {
	buf := &bytes.Buffer{}
	baseHandler := slog.NewTextHandler(buf, nil)

	// Create conditional handler that only logs errors
	condition := func(ctx context.Context, record slog.Record) bool {
		return record.Level >= slog.LevelError
	}

	handler := NewConditionalHandler(baseHandler, condition)

	// Test that condition works
	ctx := context.Background()

	// This should be logged (error level)
	errorRecord := slog.NewRecord(time.Now(), slog.LevelError, "error message", 0)
	_ = handler.Handle(ctx, errorRecord)

	// This should not be logged (info level)
	infoRecord := slog.NewRecord(time.Now(), slog.LevelInfo, "info message", 0)
	_ = handler.Handle(ctx, infoRecord)

	output := buf.String()
	if !strings.Contains(output, "error message") {
		t.Error("expected error message to be logged")
	}
	if strings.Contains(output, "info message") {
		t.Error("expected info message to be filtered out")
	}
}

func TestConditionalHandler_Enabled(t *testing.T) {
	buf := &bytes.Buffer{}
	baseHandler := slog.NewTextHandler(buf, nil)

	condition := func(ctx context.Context, record slog.Record) bool {
		return record.Level >= slog.LevelError
	}

	handler := NewConditionalHandler(baseHandler, condition)
	ctx := context.Background()

	// Test enabled for different levels
	if !handler.Enabled(ctx, slog.LevelError) {
		t.Error("expected handler to be enabled for error level")
	}

	if !handler.Enabled(ctx, slog.LevelInfo) {
		t.Error("expected handler to report enabled (actual filtering happens in Handle)")
	}
}

func TestConditionalHandler_WithAttrs(t *testing.T) {
	buf := &bytes.Buffer{}
	baseHandler := slog.NewTextHandler(buf, nil)

	condition := func(ctx context.Context, record slog.Record) bool {
		return true
	}

	handler := NewConditionalHandler(baseHandler, condition)

	// Test WithAttrs
	attrs := []slog.Attr{slog.String("key", "value")}
	newHandler := handler.WithAttrs(attrs)

	if newHandler == nil {
		t.Error("expected WithAttrs to return a handler")
	}
}

func TestConditionalHandler_WithGroup(t *testing.T) {
	buf := &bytes.Buffer{}
	baseHandler := slog.NewTextHandler(buf, nil)

	condition := func(ctx context.Context, record slog.Record) bool {
		return true
	}

	handler := NewConditionalHandler(baseHandler, condition)

	// Test WithGroup
	newHandler := handler.WithGroup("testgroup")

	if newHandler == nil {
		t.Error("expected WithGroup to return a handler")
	}
}

func TestNewBufferedHandler(t *testing.T) {
	buf := &bytes.Buffer{}
	baseHandler := slog.NewTextHandler(buf, nil)

	handler := NewBufferedHandler(baseHandler, 2)
	ctx := context.Background()

	// Add one record - should not be flushed yet
	record1 := slog.NewRecord(time.Now(), slog.LevelInfo, "message 1", 0)
	_ = handler.Handle(ctx, record1)

	if buf.Len() > 0 {
		t.Error("expected no output before buffer is full")
	}

	// Add second record - should trigger flush
	record2 := slog.NewRecord(time.Now(), slog.LevelInfo, "message 2", 0)
	_ = handler.Handle(ctx, record2)

	output := buf.String()
	if !strings.Contains(output, "message 1") || !strings.Contains(output, "message 2") {
		t.Error("expected both messages after buffer flush")
	}
}

func TestBufferedHandler_Flush(t *testing.T) {
	buf := &bytes.Buffer{}
	baseHandler := slog.NewTextHandler(buf, nil)

	handler := NewBufferedHandler(baseHandler, 10)
	ctx := context.Background()

	// Add record but don't fill buffer
	record := slog.NewRecord(time.Now(), slog.LevelInfo, "test message", 0)
	_ = handler.Handle(ctx, record)

	if buf.Len() > 0 {
		t.Error("expected no output before manual flush")
	}

	// Manual flush
	_ = handler.Flush(ctx)

	output := buf.String()
	if !strings.Contains(output, "test message") {
		t.Error("expected message after manual flush")
	}
}

func TestNewAsyncHandler(t *testing.T) {
	buf := &bytes.Buffer{}
	baseHandler := slog.NewTextHandler(buf, nil)

	handler := NewAsyncHandler(baseHandler, 10)
	ctx := context.Background()

	// Add a record
	record := slog.NewRecord(time.Now(), slog.LevelInfo, "async message", 0)
	_ = handler.Handle(ctx, record)

	// Give async processing time
	time.Sleep(10 * time.Millisecond)

	// Close to ensure all messages are processed
	handler.Close()

	output := buf.String()
	if !strings.Contains(output, "async message") {
		t.Error("expected async message to be processed")
	}
}

func TestNewHandlerBuilder(t *testing.T) {
	buf := &bytes.Buffer{}
	baseHandler := slog.NewTextHandler(buf, nil)

	builder := NewHandlerBuilder(baseHandler)

	if builder == nil {
		t.Error("expected builder to be created")
	}
}

func TestHandlerBuilder_WithMiddleware(t *testing.T) {
	buf := &bytes.Buffer{}
	baseHandler := slog.NewTextHandler(buf, nil)

	builder := NewHandlerBuilder(baseHandler)

	// Create a simple middleware that adds a field
	middleware := handlerMiddlewareFunc(func(ctx context.Context, record slog.Record, next HandlerFunc) error {
		record.AddAttrs(slog.String("middleware", "applied"))
		return next(ctx, record)
	})

	newBuilder := builder.WithMiddleware(middleware)

	if newBuilder == nil {
		t.Error("expected WithMiddleware to return builder")
	}
}

func TestHandlerBuilder_WithTimestamp(t *testing.T) {
	buf := &bytes.Buffer{}
	baseHandler := slog.NewTextHandler(buf, nil)

	builder := NewHandlerBuilder(baseHandler)
	newBuilder := builder.WithTimestamp()

	if newBuilder == nil {
		t.Error("expected WithTimestamp to return builder")
	}
}

func TestHandlerBuilder_WithStaticFields(t *testing.T) {
	buf := &bytes.Buffer{}
	baseHandler := slog.NewTextHandler(buf, nil)

	builder := NewHandlerBuilder(baseHandler)
	fields := map[string]interface{}{
		"service": "test-service",
		"version": "1.0.0",
	}

	newBuilder := builder.WithStaticFields(fields)

	if newBuilder == nil {
		t.Error("expected WithStaticFields to return builder")
	}
}

func TestHandlerBuilder_WithContextExtractor(t *testing.T) {
	buf := &bytes.Buffer{}
	baseHandler := slog.NewTextHandler(buf, nil)

	builder := NewHandlerBuilder(baseHandler)
	extractor := TraceContextExtractor()

	newBuilder := builder.WithContextExtractor(extractor)

	if newBuilder == nil {
		t.Error("expected WithContextExtractor to return builder")
	}
}

func TestHandlerBuilder_WithTraceContext(t *testing.T) {
	buf := &bytes.Buffer{}
	baseHandler := slog.NewTextHandler(buf, nil)

	builder := NewHandlerBuilder(baseHandler)
	newBuilder := builder.WithTraceContext()

	if newBuilder == nil {
		t.Error("expected WithTraceContext to return builder")
	}
}

func TestHandlerBuilder_WithLevelFilter(t *testing.T) {
	buf := &bytes.Buffer{}
	baseHandler := slog.NewTextHandler(buf, nil)

	builder := NewHandlerBuilder(baseHandler)
	newBuilder := builder.WithLevelFilter(slog.LevelWarn)

	if newBuilder == nil {
		t.Error("expected WithLevelFilter to return builder")
	}
}

func TestHandlerBuilder_WithRedaction(t *testing.T) {
	buf := &bytes.Buffer{}
	baseHandler := slog.NewTextHandler(buf, nil)

	builder := NewHandlerBuilder(baseHandler)
	pattern := regexp.MustCompile(`password=\w+`)
	redactor := NewRegexRedactor(pattern, "password=***")

	newBuilder := builder.WithRedaction(redactor)

	if newBuilder == nil {
		t.Error("expected WithRedaction to return builder")
	}
}

func TestHandlerBuilder_WithSampling(t *testing.T) {
	buf := &bytes.Buffer{}
	baseHandler := slog.NewTextHandler(buf, nil)

	builder := NewHandlerBuilder(baseHandler)
	newBuilder := builder.WithSampling(2) // Every 2nd message

	if newBuilder == nil {
		t.Error("expected WithSampling to return builder")
	}
}

func TestHandlerBuilder_Build(t *testing.T) {
	buf := &bytes.Buffer{}
	baseHandler := slog.NewTextHandler(buf, nil)

	builder := NewHandlerBuilder(baseHandler)
	handler := builder.Build()

	if handler == nil {
		t.Error("expected Build to return a handler")
	}

	// Test that built handler works
	ctx := context.Background()
	record := slog.NewRecord(time.Now(), slog.LevelInfo, "test message", 0)
	err := handler.Handle(ctx, record)

	if err != nil {
		t.Errorf("expected no error from handler, got: %v", err)
	}
}

func TestMultiHandlerBuilder(t *testing.T) {
	buf1 := &bytes.Buffer{}
	buf2 := &bytes.Buffer{}
	handler1 := slog.NewTextHandler(buf1, nil)
	handler2 := slog.NewTextHandler(buf2, nil)

	builder := MultiHandlerBuilder(handler1, handler2)

	if builder == nil {
		t.Error("expected MultiHandlerBuilder to return builder")
	}
}

func TestConditionalHandlerBuilder(t *testing.T) {
	buf := &bytes.Buffer{}
	baseHandler := slog.NewTextHandler(buf, nil)
	condition := func(ctx context.Context, record slog.Record) bool {
		return record.Level >= slog.LevelWarn
	}

	builder := ConditionalHandlerBuilder(baseHandler, condition)

	if builder == nil {
		t.Error("expected ConditionalHandlerBuilder to return builder")
	}
}

func TestBufferedHandlerBuilder(t *testing.T) {
	buf := &bytes.Buffer{}
	baseHandler := slog.NewTextHandler(buf, nil)

	builder := BufferedHandlerBuilder(baseHandler, 5)

	if builder == nil {
		t.Error("expected BufferedHandlerBuilder to return builder")
	}
}

func TestAsyncHandlerBuilder(t *testing.T) {
	buf := &bytes.Buffer{}
	baseHandler := slog.NewTextHandler(buf, nil)

	builder := AsyncHandlerBuilder(baseHandler, 10)

	if builder == nil {
		t.Error("expected AsyncHandlerBuilder to return builder")
	}
}

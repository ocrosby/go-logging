package logging

import (
	"bytes"
	"context"
	"log/slog"
	"testing"
	"time"
)

func TestNewJSONLogger_Factory(t *testing.T) {
	logger := NewJSONLogger(InfoLevel)
	if logger == nil {
		t.Fatal("expected JSON logger to be created")
	}

	// Test that it can log
	logger.Info("test JSON message")
}

func TestFluentMissingMethods(t *testing.T) {
	buf := &bytes.Buffer{}
	config := NewLoggerConfig().
		WithLevel(InfoLevel).
		WithWriter(buf).
		WithJSONFormat().
		Build()
	logger := NewWithLoggerConfig(config)

	// Test Field method
	logger.Fluent().Info().
		Field("custom", map[string]interface{}{"nested": "value"}).
		Msg("field test")

	// Test Bool method
	logger.Fluent().Info().
		Bool("flag", true).
		Msg("bool test")

	output := buf.String()
	if !contains(output, "field test") || !contains(output, "bool test") {
		t.Error("expected fluent methods to work")
	}
}

func TestMultiHandler_Basic(t *testing.T) {
	buf1 := &bytes.Buffer{}
	buf2 := &bytes.Buffer{}

	handler1 := slog.NewJSONHandler(buf1, nil)
	handler2 := slog.NewTextHandler(buf2, nil)

	multiHandler := NewMultiHandler(handler1, handler2)
	if multiHandler == nil {
		t.Fatal("expected multi handler to be created")
	}

	// Test Enabled
	if !multiHandler.Enabled(context.Background(), slog.LevelInfo) {
		t.Error("expected handler to be enabled")
	}

	// Test Handle
	record := slog.NewRecord(time.Now(), slog.LevelInfo, "test message", 0)
	err := multiHandler.Handle(context.Background(), record)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Both buffers should have output
	if buf1.Len() == 0 || buf2.Len() == 0 {
		t.Error("expected output in both handlers")
	}
}

func TestMultiHandler_WithAttrsAndGroup(t *testing.T) {
	buf1 := &bytes.Buffer{}
	handler1 := slog.NewJSONHandler(buf1, nil)

	multiHandler := NewMultiHandler(handler1)

	// Test WithAttrs
	attrs := []slog.Attr{slog.String("key", "value")}
	newHandler := multiHandler.WithAttrs(attrs)
	if newHandler == nil {
		t.Error("expected WithAttrs to return handler")
	}

	// Test WithGroup
	groupHandler := multiHandler.WithGroup("testgroup")
	if groupHandler == nil {
		t.Error("expected WithGroup to return handler")
	}
}

func TestMultiHandler_AddRemove(t *testing.T) {
	buf1 := &bytes.Buffer{}
	buf2 := &bytes.Buffer{}

	handler1 := slog.NewJSONHandler(buf1, nil)
	handler2 := slog.NewTextHandler(buf2, nil)

	multiHandler := NewMultiHandler()

	// Add handlers
	multiHandler.AddHandler(handler1)
	multiHandler.AddHandler(handler2)

	// Remove a handler
	multiHandler.RemoveHandler(handler1)

	// Test that only handler2 receives logs
	record := slog.NewRecord(time.Now(), slog.LevelInfo, "test after remove", 0)
	err := multiHandler.Handle(context.Background(), record)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Only buf2 should have output
	if buf1.Len() > 0 {
		t.Error("expected buf1 to be empty after removal")
	}

	if buf2.Len() == 0 {
		t.Error("expected buf2 to have output")
	}
}

func TestNewTextLogger_Factory(t *testing.T) {
	logger := NewTextLogger(WarnLevel)
	if logger == nil {
		t.Fatal("expected text logger to be created")
	}

	// Test that it can log
	logger.Warn("test text message")
}

func TestFromEnvironment_Config(t *testing.T) {
	// This tests the environment configuration
	config := NewConfig().FromEnvironment().Build()
	if config == nil {
		t.Fatal("expected config to be created from environment")
	}
}

func TestCallerMiddleware_Basic(t *testing.T) {
	buf := &bytes.Buffer{}
	baseHandler := slog.NewJSONHandler(buf, nil)

	middleware := CallerMiddleware(2) // Skip 2 frames
	handler := NewMiddlewareHandler(baseHandler, middleware)

	record := slog.NewRecord(time.Now(), slog.LevelInfo, "caller test", 0)
	err := handler.Handle(context.Background(), record)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Should not crash
	if buf.Len() == 0 {
		t.Error("expected some output")
	}
}

package logging

import (
	"bytes"
	"context"
	"log/slog"
	"testing"
)

func TestNewUnifiedLogger_WithNilConfig(t *testing.T) {
	logger := NewUnifiedLogger(nil, nil)
	if logger == nil {
		t.Fatal("expected logger to be created with nil config")
	}

	// Should use defaults and be able to log
	logger.Info("test message")
}

func TestUnifiedLogger_WithField_Chaining(t *testing.T) {
	buf := &bytes.Buffer{}
	config := NewLoggerConfig().
		WithLevel(InfoLevel).
		WithWriter(buf).
		WithJSONFormat().
		Build()
	logger := NewUnifiedLogger(config, nil)

	// Test chaining multiple WithField calls
	enhancedLogger := logger.WithField("key1", "value1").WithField("key2", "value2")
	enhancedLogger.Info("test message")

	output := buf.String()
	if !contains(output, "key1") || !contains(output, "key2") {
		t.Errorf("expected fields in output, got: %s", output)
	}
}

func TestUnifiedLogger_WithFields_MergeExisting(t *testing.T) {
	buf := &bytes.Buffer{}
	config := NewLoggerConfig().
		WithLevel(InfoLevel).
		WithWriter(buf).
		WithJSONFormat().
		Build()
	logger := NewUnifiedLogger(config, nil)

	// Add initial field
	logger1 := logger.WithField("existing", "field")

	// Add multiple fields - should merge
	fields := map[string]interface{}{
		"new1": "value1",
		"new2": 42,
	}
	logger2 := logger1.WithFields(fields)

	logger2.Info("test message")

	output := buf.String()
	if !contains(output, "existing") || !contains(output, "new1") || !contains(output, "new2") {
		t.Errorf("expected all fields in output, got: %s", output)
	}
}

func TestUnifiedLogger_SetLevel_GetLevel(t *testing.T) {
	logger := NewUnifiedLogger(nil, nil).(*unifiedLogger)

	// Test default level
	if logger.GetLevel() != InfoLevel {
		t.Errorf("expected default level %v, got %v", InfoLevel, logger.GetLevel())
	}

	// Test setting level
	logger.SetLevel(DebugLevel)
	if logger.GetLevel() != DebugLevel {
		t.Errorf("expected level %v after set, got %v", DebugLevel, logger.GetLevel())
	}
}

func TestUnifiedLogger_IsLevelEnabled(t *testing.T) {
	logger := NewUnifiedLogger(nil, nil)

	// Assuming default level is Info
	if !logger.IsLevelEnabled(InfoLevel) {
		t.Error("expected info level to be enabled")
	}

	if !logger.IsLevelEnabled(WarnLevel) {
		t.Error("expected warn level to be enabled")
	}

	if logger.IsLevelEnabled(DebugLevel) {
		t.Error("expected debug level to be disabled with info level logger")
	}
}

func TestUnifiedLogger_Fluent(t *testing.T) {
	buf := &bytes.Buffer{}
	config := NewLoggerConfig().
		WithLevel(InfoLevel).
		WithWriter(buf).
		WithJSONFormat().
		Build()
	logger := NewUnifiedLogger(config, nil)

	// Test fluent interface
	logger.Fluent().Info().
		Str("key1", "value1").
		Int("key2", 42).
		Msg("fluent test message")

	output := buf.String()
	if !contains(output, "fluent test message") {
		t.Error("expected fluent message in output")
	}

	if !contains(output, "key1") || !contains(output, "key2") {
		t.Error("expected fluent fields in output")
	}
}

func TestUnifiedLogger_LogAndLogContext(t *testing.T) {
	buf := &bytes.Buffer{}
	config := NewLoggerConfig().
		WithLevel(TraceLevel).
		WithWriter(buf).
		WithTextFormat().
		Build()
	logger := NewUnifiedLogger(config, nil)

	// Test Log method
	logger.Log(WarnLevel, "warn message %s", "arg")

	// Test LogContext method
	ctx := WithTraceID(context.Background(), "trace-123")
	logger.LogContext(ctx, ErrorLevel, "error message %s", "arg")

	output := buf.String()
	if !contains(output, "warn message arg") {
		t.Error("expected warn message in output")
	}

	if !contains(output, "error message arg") {
		t.Error("expected error message in output")
	}
}

func TestUnifiedLogger_LevelMethods(t *testing.T) {
	buf := &bytes.Buffer{}
	config := NewLoggerConfig().
		WithLevel(TraceLevel).
		WithWriter(buf).
		WithTextFormat().
		Build()
	logger := NewUnifiedLogger(config, nil)

	// Test all level methods
	logger.Trace("trace %s", "message")
	logger.Debug("debug %s", "message")
	logger.Info("info %s", "message")
	logger.Warn("warn %s", "message")
	logger.Error("error %s", "message")
	logger.Critical("critical %s", "message")

	output := buf.String()

	levels := []string{"trace", "debug", "info", "warn", "error", "critical"}
	for _, level := range levels {
		if !contains(output, level+" message") {
			t.Errorf("expected %s message in output", level)
		}
	}
}

func TestUnifiedLogger_ContextMethods(t *testing.T) {
	buf := &bytes.Buffer{}
	config := NewLoggerConfig().
		WithLevel(TraceLevel).
		WithWriter(buf).
		WithTextFormat().
		Build()
	logger := NewUnifiedLogger(config, nil)

	ctx := WithTraceID(context.Background(), "ctx-123")

	// Test all context methods
	logger.TraceContext(ctx, "trace context %s", "message")
	logger.DebugContext(ctx, "debug context %s", "message")
	logger.InfoContext(ctx, "info context %s", "message")
	logger.WarnContext(ctx, "warn context %s", "message")
	logger.ErrorContext(ctx, "error context %s", "message")
	logger.CriticalContext(ctx, "critical context %s", "message")

	output := buf.String()

	contextLevels := []string{"trace", "debug", "info", "warn", "error", "critical"}
	for _, level := range contextLevels {
		if !contains(output, level+" context message") {
			t.Errorf("expected %s context message in output, got: %s", level, output)
		}
	}
}

func TestUnifiedLogger_WithSlog(t *testing.T) {
	buf := &bytes.Buffer{}
	handler := slog.NewJSONHandler(buf, nil)

	config := NewLoggerConfig().
		WithLevel(InfoLevel).
		WithHandler(handler).
		UseSlog(true).
		Build()

	logger := NewUnifiedLogger(config, nil)
	logger.Info("slog test message")

	// Just verify it doesn't crash - slog logging is harder to test output
}

func TestUnifiedLogger_ThreadSafety(t *testing.T) {
	buf := &bytes.Buffer{}
	config := NewLoggerConfig().
		WithLevel(InfoLevel).
		WithWriter(buf).
		WithTextFormat().
		Build()
	logger := NewUnifiedLogger(config, nil)

	// Test concurrent access
	done := make(chan bool)

	// Concurrent field setting
	go func() {
		for i := 0; i < 10; i++ {
			logger.WithField("field1", i).Info("message 1")
		}
		done <- true
	}()

	go func() {
		for i := 0; i < 10; i++ {
			logger.WithField("field2", i).Info("message 2")
		}
		done <- true
	}()

	// Wait for both goroutines
	<-done
	<-done

	// Should not crash and should have output
	if buf.Len() == 0 {
		t.Error("expected some output from concurrent logging")
	}
}

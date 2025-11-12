package logging

import (
	"bytes"
	"context"
	"log/slog"
	"strings"
	"testing"
)

func TestSlogLogger_Levels(t *testing.T) {
	var buf bytes.Buffer
	handler := slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	logger := NewWithHandler(handler)

	logger.Debug("debug message")
	logger.Info("info message")
	logger.Warn("warn message")
	logger.Error("error message")

	output := buf.String()
	if !strings.Contains(output, "debug message") {
		t.Errorf("Expected debug message in output, got: %s", output)
	}
	if !strings.Contains(output, "info message") {
		t.Errorf("Expected info message in output, got: %s", output)
	}
	if !strings.Contains(strings.ToUpper(output), "WARN") {
		t.Errorf("Expected WARN level in output, got: %s", output)
	}
	if !strings.Contains(strings.ToUpper(output), "ERROR") {
		t.Errorf("Expected ERROR level in output, got: %s", output)
	}
}

func TestSlogLogger_WithFields(t *testing.T) {
	var buf bytes.Buffer
	handler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	logger := NewWithHandler(handler)

	logger = logger.WithField("service", "test")
	logger.Info("message with field")

	output := buf.String()
	if !strings.Contains(output, `"service":"test"`) {
		t.Errorf("Expected service field in output, got: %s", output)
	}
	if !strings.Contains(output, "message with field") {
		t.Errorf("Expected message in output, got: %s", output)
	}
}

func TestSlogLogger_WithMultipleFields(t *testing.T) {
	var buf bytes.Buffer
	handler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	logger := NewWithHandler(handler)

	fields := map[string]interface{}{
		"service": "test",
		"version": "1.0.0",
		"env":     "production",
	}
	logger = logger.WithFields(fields)
	logger.Info("message with multiple fields")

	output := buf.String()
	if !strings.Contains(output, `"service":"test"`) {
		t.Errorf("Expected service field in output")
	}
	if !strings.Contains(output, `"version":"1.0.0"`) {
		t.Errorf("Expected version field in output")
	}
	if !strings.Contains(output, `"env":"production"`) {
		t.Errorf("Expected env field in output")
	}
}

func TestSlogLogger_Context(t *testing.T) {
	var buf bytes.Buffer
	handler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	logger := NewWithHandler(handler)

	ctx := WithRequestID(context.Background(), "req-123")
	logger.InfoContext(ctx, "message with context")

	output := buf.String()
	if !strings.Contains(output, "req-123") {
		t.Errorf("Expected request ID in output, got: %s", output)
	}
}

func TestSlogLogger_LevelFiltering(t *testing.T) {
	var buf bytes.Buffer
	handler := slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelWarn,
	})
	logger := NewWithHandler(handler)

	logger.Debug("debug message")
	logger.Info("info message")
	logger.Warn("warn message")
	logger.Error("error message")

	output := buf.String()
	if strings.Contains(output, "debug message") {
		t.Errorf("Debug message should be filtered out")
	}
	if strings.Contains(output, "info message") {
		t.Errorf("Info message should be filtered out")
	}
	if !strings.Contains(output, "warn message") {
		t.Errorf("Expected warn message in output")
	}
	if !strings.Contains(output, "error message") {
		t.Errorf("Expected error message in output")
	}
}

func TestSlogLogger_Formatting(t *testing.T) {
	var buf bytes.Buffer
	handler := slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	logger := NewWithHandler(handler)

	logger.Info("user %s logged in", "john")
	logger.Error("failed with error: %v", "connection timeout")

	output := buf.String()
	if !strings.Contains(output, "user john logged in") {
		t.Errorf("Expected formatted message in output, got: %s", output)
	}
	if !strings.Contains(output, "failed with error: connection timeout") {
		t.Errorf("Expected formatted error message in output, got: %s", output)
	}
}

func TestSlogLogger_Fluent(t *testing.T) {
	var buf bytes.Buffer
	handler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	logger := NewWithHandler(handler)

	logger.Fluent().Info().
		Str("user", "john").
		Int("age", 30).
		Msg("User information")

	output := buf.String()
	if !strings.Contains(output, `"user":"john"`) {
		t.Errorf("Expected user field in output, got: %s", output)
	}
	if !strings.Contains(output, `"age":30`) {
		t.Errorf("Expected age field in output, got: %s", output)
	}
	if !strings.Contains(output, "User information") {
		t.Errorf("Expected message in output, got: %s", output)
	}
}

func TestSlogLogger_CustomLevels(t *testing.T) {
	var buf bytes.Buffer
	handler := slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.Level(-8), // Allow trace level (our custom level)
	})
	logger := NewWithHandler(handler)

	logger.Trace("trace message")
	logger.Critical("critical message")

	output := buf.String()
	if !strings.Contains(output, "trace message") {
		t.Errorf("Expected trace message in output, got: %s", output)
	}
	if !strings.Contains(output, "critical message") {
		t.Errorf("Expected critical message in output, got: %s", output)
	}
}

func TestSlogLogger_IsLevelEnabled(t *testing.T) {
	handler := slog.NewTextHandler(&bytes.Buffer{}, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	logger := NewWithHandler(handler)

	if logger.IsLevelEnabled(TraceLevel) {
		t.Error("TraceLevel should not be enabled")
	}
	if logger.IsLevelEnabled(DebugLevel) {
		t.Error("DebugLevel should not be enabled")
	}
	if !logger.IsLevelEnabled(InfoLevel) {
		t.Error("InfoLevel should be enabled")
	}
	if !logger.IsLevelEnabled(WarnLevel) {
		t.Error("WarnLevel should be enabled")
	}
	if !logger.IsLevelEnabled(ErrorLevel) {
		t.Error("ErrorLevel should be enabled")
	}
	if !logger.IsLevelEnabled(CriticalLevel) {
		t.Error("CriticalLevel should be enabled")
	}
}

func TestSlogLogger_GetSetLevel(t *testing.T) {
	handler := slog.NewTextHandler(&bytes.Buffer{}, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	logger := NewWithHandler(handler)

	// Test setting to DebugLevel
	logger.SetLevel(DebugLevel)
	if logger.GetLevel() != DebugLevel {
		t.Errorf("Expected level to be DebugLevel after SetLevel, got %v", logger.GetLevel())
	}

	// Test setting to ErrorLevel
	logger.SetLevel(ErrorLevel)
	if logger.GetLevel() != ErrorLevel {
		t.Errorf("Expected level to be ErrorLevel after SetLevel, got %v", logger.GetLevel())
	}
}

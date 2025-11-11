package logging

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"
)

func TestStandardLogger_Levels(t *testing.T) {
	buf := &bytes.Buffer{}
	config := NewConfig().
		WithLevel(InfoLevel).
		WithOutput(buf).
		Build()

	logger := NewStandardLogger(config)

	logger.Trace("trace message")
	logger.Debug("debug message")
	logger.Info("info message")
	logger.Warn("warn message")
	logger.Error("error message")

	output := buf.String()

	if strings.Contains(output, "trace message") {
		t.Error("Trace message should not appear when level is Info")
	}
	if strings.Contains(output, "debug message") {
		t.Error("Debug message should not appear when level is Info")
	}
	if !strings.Contains(output, "info message") {
		t.Error("Info message should appear when level is Info")
	}
	if !strings.Contains(output, "warn message") {
		t.Error("Warn message should appear when level is Info")
	}
	if !strings.Contains(output, "error message") {
		t.Error("Error message should appear when level is Info")
	}
}

func TestStandardLogger_JSONFormat(t *testing.T) {
	buf := &bytes.Buffer{}
	config := NewConfig().
		WithLevel(InfoLevel).
		WithOutput(buf).
		WithJSONFormat().
		Build()

	logger := NewStandardLogger(config)
	logger.Info("test message")

	output := buf.String()

	var entry map[string]interface{}
	if err := json.Unmarshal([]byte(output), &entry); err != nil {
		t.Fatalf("Failed to parse JSON output: %v", err)
	}

	if entry["level"] != "INFO" {
		t.Errorf("Expected level INFO, got %v", entry["level"])
	}
	if entry["message"] != "test message" {
		t.Errorf("Expected message 'test message', got %v", entry["message"])
	}
	if _, ok := entry["timestamp"]; !ok {
		t.Error("Expected timestamp in JSON output")
	}
}

func TestStandardLogger_WithFields(t *testing.T) {
	buf := &bytes.Buffer{}
	config := NewConfig().
		WithLevel(InfoLevel).
		WithOutput(buf).
		WithJSONFormat().
		Build()

	logger := NewStandardLogger(config)
	logger = logger.WithField("user_id", "12345")
	logger.Info("test message")

	output := buf.String()

	var entry map[string]interface{}
	if err := json.Unmarshal([]byte(output), &entry); err != nil {
		t.Fatalf("Failed to parse JSON output: %v", err)
	}

	if entry["user_id"] != "12345" {
		t.Errorf("Expected user_id 12345, got %v", entry["user_id"])
	}
}

func TestStandardLogger_WithMultipleFields(t *testing.T) {
	buf := &bytes.Buffer{}
	config := NewConfig().
		WithLevel(InfoLevel).
		WithOutput(buf).
		WithJSONFormat().
		Build()

	logger := NewStandardLogger(config)
	logger = logger.WithFields(map[string]interface{}{
		"request_id": "req-123",
		"user_id":    "user-456",
	})
	logger.Info("test message")

	output := buf.String()

	var entry map[string]interface{}
	if err := json.Unmarshal([]byte(output), &entry); err != nil {
		t.Fatalf("Failed to parse JSON output: %v", err)
	}

	if entry["request_id"] != "req-123" {
		t.Errorf("Expected request_id req-123, got %v", entry["request_id"])
	}
	if entry["user_id"] != "user-456" {
		t.Errorf("Expected user_id user-456, got %v", entry["user_id"])
	}
}

func TestStandardLogger_Context(t *testing.T) {
	buf := &bytes.Buffer{}
	config := NewConfig().
		WithLevel(InfoLevel).
		WithOutput(buf).
		WithJSONFormat().
		Build()

	logger := NewStandardLogger(config)
	ctx := WithRequestID(context.Background(), "req-789")
	logger.InfoContext(ctx, "test message")

	output := buf.String()

	var entry map[string]interface{}
	if err := json.Unmarshal([]byte(output), &entry); err != nil {
		t.Fatalf("Failed to parse JSON output: %v", err)
	}

	if entry["request_id"] != "req-789" {
		t.Errorf("Expected request_id req-789, got %v", entry["request_id"])
	}
}

func TestStandardLogger_SetLevel(t *testing.T) {
	buf := &bytes.Buffer{}
	config := NewConfig().
		WithLevel(InfoLevel).
		WithOutput(buf).
		Build()

	logger := NewStandardLogger(config)

	logger.Debug("should not appear")
	if buf.Len() > 0 {
		t.Error("Debug message should not appear when level is Info")
	}

	logger.SetLevel(DebugLevel)
	logger.Debug("should appear")

	if !strings.Contains(buf.String(), "should appear") {
		t.Error("Debug message should appear after setting level to Debug")
	}
}

func TestStandardLogger_IsLevelEnabled(t *testing.T) {
	config := NewConfig().
		WithLevel(InfoLevel).
		Build()

	logger := NewStandardLogger(config)

	if logger.IsLevelEnabled(DebugLevel) {
		t.Error("Debug should not be enabled when level is Info")
	}
	if !logger.IsLevelEnabled(InfoLevel) {
		t.Error("Info should be enabled when level is Info")
	}
	if !logger.IsLevelEnabled(ErrorLevel) {
		t.Error("Error should be enabled when level is Info")
	}
}

func TestStandardLogger_Formatting(t *testing.T) {
	buf := &bytes.Buffer{}
	config := NewConfig().
		WithLevel(InfoLevel).
		WithOutput(buf).
		Build()

	logger := NewStandardLogger(config)
	logger.Info("test %s %d", "message", 42)

	output := buf.String()
	if !strings.Contains(output, "test message 42") {
		t.Errorf("Expected formatted message, got %s", output)
	}
}

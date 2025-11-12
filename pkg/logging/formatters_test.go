package logging

import (
	"context"
	"encoding/json"
	"regexp"
	"strings"
	"testing"
	"time"
)

const testMessage = "test message"

func TestNewJSONFormatter(t *testing.T) {
	// Test with nil config
	formatter := NewJSONFormatter(nil)
	if formatter == nil {
		t.Fatal("expected formatter to be created with nil config")
	}

	// Test with custom config
	config := NewFormatterConfig().WithJSONFormat().Build()
	formatter2 := NewJSONFormatter(config)
	if formatter2 == nil {
		t.Fatal("expected formatter to be created with custom config")
		return
	}

	if formatter2.config != config {
		t.Error("expected formatter to use provided config")
	}
}

func TestJSONFormatter_Format(t *testing.T) {
	config := NewFormatterConfig().
		WithJSONFormat().
		IncludeTime(true).
		IncludeFile(false).
		Build()

	formatter := NewJSONFormatter(config)

	entry := LogEntry{
		Level:     InfoLevel,
		Message:   testMessage,
		Timestamp: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
		Fields: map[string]interface{}{
			"user_id": 123,
			"action":  "login",
		},
		Context: context.Background(),
	}

	data, err := formatter.Format(entry)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}

	// Check required fields
	if result["level"] != "INFO" {
		t.Errorf("expected level 'INFO', got %v", result["level"])
	}

	if result["message"] != testMessage {
		t.Errorf("expected message 'test message', got %v", result["message"])
	}

	if result["timestamp"] != "2023-01-01T12:00:00Z" {
		t.Errorf("expected timestamp '2023-01-01T12:00:00Z', got %v", result["timestamp"])
	}

	if result["user_id"] != float64(123) { // JSON numbers are float64
		t.Errorf("expected user_id 123, got %v", result["user_id"])
	}

	if result["action"] != "login" {
		t.Errorf("expected action 'login', got %v", result["action"])
	}
}

func TestJSONFormatter_Format_WithFile(t *testing.T) {
	config := NewFormatterConfig().
		WithJSONFormat().
		IncludeTime(false).
		IncludeFile(true).
		UseShortFile(true).
		Build()

	formatter := NewJSONFormatter(config)

	entry := LogEntry{
		Level:   InfoLevel,
		Message: testMessage,
		File:    "/path/to/test.go",
		Line:    42,
		Context: context.Background(),
	}

	data, err := formatter.Format(entry)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}

	if !strings.Contains(result["file"].(string), "test.go:42") {
		t.Errorf("expected file info to contain test.go:42, got %v", result["file"])
	}
}

func TestJSONFormatter_Format_WithRedaction(t *testing.T) {
	re := regexp.MustCompile(`password=\w+`)
	config := NewFormatterConfig().
		WithJSONFormat().
		AddRedactRegex(re).
		Build()

	formatter := NewJSONFormatter(config)

	entry := LogEntry{
		Level:   InfoLevel,
		Message: "user login with password=secret123",
		Context: context.Background(),
	}

	data, err := formatter.Format(entry)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}

	message := result["message"].(string)
	if strings.Contains(message, "secret123") {
		t.Error("expected password to be redacted")
	}

	if !strings.Contains(message, "[REDACTED]") {
		t.Error("expected [REDACTED] placeholder")
	}
}

func TestNewTextFormatter(t *testing.T) {
	// Test with nil config
	formatter := NewTextFormatter(nil)
	if formatter == nil {
		t.Fatal("expected formatter to be created with nil config")
	}

	// Test with custom config
	config := NewFormatterConfig().WithTextFormat().Build()
	formatter2 := NewTextFormatter(config)
	if formatter2 == nil {
		t.Fatal("expected formatter to be created with custom config")
	}
}

func TestTextFormatter_Format(t *testing.T) {
	config := NewFormatterConfig().
		WithTextFormat().
		IncludeTime(true).
		IncludeFile(false).
		Build()

	formatter := NewTextFormatter(config)

	entry := LogEntry{
		Level:     InfoLevel,
		Message:   testMessage,
		Timestamp: time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
		Fields: map[string]interface{}{
			"user_id": 123,
			"action":  "login",
		},
		Context: context.Background(),
	}

	data, err := formatter.Format(entry)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := string(data)

	// Check that components are present
	if !strings.Contains(output, "2023/01/01 12:00:00") {
		t.Error("expected timestamp in output")
	}

	if !strings.Contains(output, "[INFO]") {
		t.Error("expected level in output")
	}

	if !strings.Contains(output, testMessage) {
		t.Error("expected message in output")
	}

	if !strings.Contains(output, "user_id=123") {
		t.Error("expected user_id field in output")
	}

	if !strings.Contains(output, "action=login") {
		t.Error("expected action field in output")
	}
}

func TestTextFormatter_Format_WithFile(t *testing.T) {
	config := NewFormatterConfig().
		WithTextFormat().
		IncludeTime(false).
		IncludeFile(true).
		UseShortFile(true).
		Build()

	formatter := NewTextFormatter(config)

	entry := LogEntry{
		Level:   WarnLevel,
		Message: "warning message",
		File:    "/long/path/to/test.go",
		Line:    100,
		Context: context.Background(),
	}

	data, err := formatter.Format(entry)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := string(data)

	if !strings.Contains(output, "[WARN]") {
		t.Error("expected warn level in output")
	}

	if !strings.Contains(output, "test.go:100") {
		t.Error("expected short filename in output")
	}
}

func TestTextFormatter_Format_WithRedaction(t *testing.T) {
	re := regexp.MustCompile(`api_key=\w+`)
	config := NewFormatterConfig().
		WithTextFormat().
		AddRedactRegex(re).
		Build()

	formatter := NewTextFormatter(config)

	entry := LogEntry{
		Level:   ErrorLevel,
		Message: "API call failed with api_key=secret123",
		Context: context.Background(),
	}

	data, err := formatter.Format(entry)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := string(data)

	if strings.Contains(output, "secret123") {
		t.Error("expected API key to be redacted")
	}

	if !strings.Contains(output, "[REDACTED]") {
		t.Error("expected [REDACTED] placeholder")
	}
}

func TestNewConsoleFormatter(t *testing.T) {
	// Test with nil config
	formatter := NewConsoleFormatter(nil, true)
	if formatter == nil {
		t.Fatal("expected formatter to be created with nil config")
		return
	}

	if !formatter.useColors {
		t.Error("expected colors to be enabled")
	}

	// Test with custom config and no colors
	config := NewFormatterConfig().WithTextFormat().Build()
	formatter2 := NewConsoleFormatter(config, false)
	if formatter2 == nil {
		t.Fatal("expected formatter to be created")
		return
	}

	if formatter2.useColors {
		t.Error("expected colors to be disabled")
	}
}

func TestConsoleFormatter_Format(t *testing.T) {
	config := NewFormatterConfig().
		WithTextFormat().
		IncludeTime(true).
		IncludeFile(false).
		Build()

	formatter := NewConsoleFormatter(config, false) // No colors for easier testing

	entry := LogEntry{
		Level:     WarnLevel,
		Message:   "warning message",
		Timestamp: time.Date(2023, 1, 1, 15, 30, 45, 0, time.UTC),
		Fields: map[string]interface{}{
			"component": "auth",
		},
		Context: context.Background(),
	}

	data, err := formatter.Format(entry)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := string(data)

	// Check components
	if !strings.Contains(output, "15:30:45") {
		t.Error("expected time format HH:MM:SS in output")
	}

	if !strings.Contains(output, "[WARN]") {
		t.Error("expected warn level in output")
	}

	if !strings.Contains(output, "warning message") {
		t.Error("expected message in output")
	}

	if !strings.Contains(output, "component=auth") {
		t.Error("expected component field in output")
	}
}

func TestConsoleFormatter_Format_WithColors(t *testing.T) {
	config := NewFormatterConfig().
		WithTextFormat().
		IncludeTime(true).
		Build()

	formatter := NewConsoleFormatter(config, true) // Enable colors

	entry := LogEntry{
		Level:     ErrorLevel,
		Message:   "error message",
		Timestamp: time.Now(),
		Context:   context.Background(),
	}

	data, err := formatter.Format(entry)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := string(data)

	// Check that color codes are present
	if !strings.Contains(output, "\033[") {
		t.Error("expected ANSI color codes in output")
	}

	if !strings.Contains(output, "\033[31m") { // Red for error
		t.Error("expected red color code for error level")
	}

	if !strings.Contains(output, "\033[0m") { // Reset
		t.Error("expected color reset code")
	}
}

func TestConsoleFormatter_LevelColors(t *testing.T) {
	formatter := NewConsoleFormatter(nil, true)

	expectedColors := map[Level]string{
		TraceLevel:    "\033[36m", // Cyan
		DebugLevel:    "\033[37m", // White
		InfoLevel:     "\033[32m", // Green
		WarnLevel:     "\033[33m", // Yellow
		ErrorLevel:    "\033[31m", // Red
		CriticalLevel: "\033[35m", // Magenta
	}

	for level, expectedColor := range expectedColors {
		if color, exists := formatter.levelColors[level]; !exists || color != expectedColor {
			t.Errorf("expected color %s for level %v, got %s", expectedColor, level, color)
		}
	}
}

func TestFormattersWithMinimalFields(t *testing.T) {
	entry := LogEntry{
		Level:   InfoLevel,
		Message: "minimal message",
		Context: context.Background(),
	}

	// Test JSON formatter
	jsonFormatter := NewJSONFormatter(nil)
	jsonData, err := jsonFormatter.Format(entry)
	if err != nil {
		t.Errorf("JSON formatter error: %v", err)
	}

	if !strings.Contains(string(jsonData), "minimal message") {
		t.Error("expected message in JSON output")
	}

	// Test Text formatter
	textFormatter := NewTextFormatter(nil)
	textData, err := textFormatter.Format(entry)
	if err != nil {
		t.Errorf("Text formatter error: %v", err)
	}

	if !strings.Contains(string(textData), "minimal message") {
		t.Error("expected message in text output")
	}

	// Test Console formatter
	consoleFormatter := NewConsoleFormatter(nil, false)
	consoleData, err := consoleFormatter.Format(entry)
	if err != nil {
		t.Errorf("Console formatter error: %v", err)
	}

	if !strings.Contains(string(consoleData), "minimal message") {
		t.Error("expected message in console output")
	}
}

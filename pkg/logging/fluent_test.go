package logging

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"
)

func TestFluentInterface_BasicUsage(t *testing.T) {
	buf := &bytes.Buffer{}
	config := NewConfig().
		WithLevel(InfoLevel).
		WithOutput(buf).
		WithJSONFormat().
		Build()

	redactorChain := ProvideRedactorChain(config)
	logger := NewStandardLogger(config, redactorChain)

	logger.Fluent().Info().
		Str("user", "john").
		Int("age", 30).
		Msg("User logged in")

	output := buf.String()

	var entry map[string]interface{}
	if err := json.Unmarshal([]byte(output), &entry); err != nil {
		t.Fatalf("Failed to parse JSON output: %v", err)
	}

	if entry["user"] != "john" {
		t.Errorf("Expected user john, got %v", entry["user"])
	}

	if age, ok := entry["age"].(float64); !ok || int(age) != 30 {
		t.Errorf("Expected age 30, got %v", entry["age"])
	}
}

func TestFluentInterface_WithTraceID(t *testing.T) {
	buf := &bytes.Buffer{}
	config := NewConfig().
		WithLevel(InfoLevel).
		WithOutput(buf).
		WithJSONFormat().
		Build()

	redactorChain := ProvideRedactorChain(config)
	logger := NewStandardLogger(config, redactorChain)

	logger.Fluent().Info().
		TraceID("trace-123").
		Msg("Test message")

	output := buf.String()

	var entry map[string]interface{}
	if err := json.Unmarshal([]byte(output), &entry); err != nil {
		t.Fatalf("Failed to parse JSON output: %v", err)
	}

	if entry["trace_id"] != "trace-123" {
		t.Errorf("Expected trace_id trace-123, got %v", entry["trace_id"])
	}
}

func TestFluentInterface_WithContext(t *testing.T) {
	buf := &bytes.Buffer{}
	config := NewConfig().
		WithLevel(InfoLevel).
		WithOutput(buf).
		WithJSONFormat().
		Build()

	redactorChain := ProvideRedactorChain(config)
	logger := NewStandardLogger(config, redactorChain)

	ctx := WithTraceID(context.Background(), "ctx-trace-456")
	ctx = WithRequestID(ctx, "req-789")

	logger.Fluent().Info().
		Ctx(ctx).
		Msg("Test message")

	output := buf.String()

	var entry map[string]interface{}
	if err := json.Unmarshal([]byte(output), &entry); err != nil {
		t.Fatalf("Failed to parse JSON output: %v", err)
	}

	if entry["trace_id"] != "ctx-trace-456" {
		t.Errorf("Expected trace_id ctx-trace-456, got %v", entry["trace_id"])
	}

	if entry["request_id"] != "req-789" {
		t.Errorf("Expected request_id req-789, got %v", entry["request_id"])
	}
}

func TestFluentInterface_WithError(t *testing.T) {
	buf := &bytes.Buffer{}
	config := NewConfig().
		WithLevel(InfoLevel).
		WithOutput(buf).
		WithJSONFormat().
		Build()

	redactorChain := ProvideRedactorChain(config)
	logger := NewStandardLogger(config, redactorChain)

	testErr := errors.New("test error")

	logger.Fluent().Error().
		Err(testErr).
		Msg("An error occurred")

	output := buf.String()

	var entry map[string]interface{}
	if err := json.Unmarshal([]byte(output), &entry); err != nil {
		t.Fatalf("Failed to parse JSON output: %v", err)
	}

	if entry["error"] != "test error" {
		t.Errorf("Expected error 'test error', got %v", entry["error"])
	}
}

func TestFluentInterface_Msgf(t *testing.T) {
	buf := &bytes.Buffer{}
	config := NewConfig().
		WithLevel(InfoLevel).
		WithOutput(buf).
		Build()

	redactorChain := ProvideRedactorChain(config)
	logger := NewStandardLogger(config, redactorChain)

	logger.Fluent().Info().
		Str("name", "test").
		Msgf("User %s logged in", "john")

	output := buf.String()

	if !strings.Contains(output, "User john logged in") {
		t.Errorf("Expected formatted message, got %s", output)
	}
}

func TestFluentInterface_MultipleFields(t *testing.T) {
	buf := &bytes.Buffer{}
	config := NewConfig().
		WithLevel(InfoLevel).
		WithOutput(buf).
		WithJSONFormat().
		Build()

	redactorChain := ProvideRedactorChain(config)
	logger := NewStandardLogger(config, redactorChain)

	fields := map[string]interface{}{
		"field1": "value1",
		"field2": 123,
		"field3": true,
	}

	logger.Fluent().Info().
		Fields(fields).
		Msg("Test message")

	output := buf.String()

	var entry map[string]interface{}
	if err := json.Unmarshal([]byte(output), &entry); err != nil {
		t.Fatalf("Failed to parse JSON output: %v", err)
	}

	if entry["field1"] != "value1" {
		t.Errorf("Expected field1 value1, got %v", entry["field1"])
	}

	if field2, ok := entry["field2"].(float64); !ok || int(field2) != 123 {
		t.Errorf("Expected field2 123, got %v", entry["field2"])
	}

	if entry["field3"] != true {
		t.Errorf("Expected field3 true, got %v", entry["field3"])
	}
}

func TestFluentInterface_AllLevels(t *testing.T) {
	buf := &bytes.Buffer{}
	config := NewConfig().
		WithLevel(TraceLevel).
		WithOutput(buf).
		WithJSONFormat().
		Build()
	redactorChain := ProvideRedactorChain(config)
	logger := NewStandardLogger(config, redactorChain)

	levels := []struct {
		name  string
		entry *FluentEntry
	}{
		{"TRACE", logger.Fluent().Trace()},
		{"DEBUG", logger.Fluent().Debug()},
		{"INFO", logger.Fluent().Info()},
		{"WARN", logger.Fluent().Warn()},
		{"ERROR", logger.Fluent().Error()},
		{"CRITICAL", logger.Fluent().Critical()},
	}

	for _, level := range levels {
		buf.Reset()
		level.entry.Msg("test message")

		output := buf.String()
		var entry map[string]interface{}
		if err := json.Unmarshal([]byte(output), &entry); err != nil {
			t.Fatalf("Failed to parse JSON for %s: %v", level.name, err)
		}

		if entry["level"] != level.name {
			t.Errorf("Expected level %s, got %v", level.name, entry["level"])
		}
	}
}

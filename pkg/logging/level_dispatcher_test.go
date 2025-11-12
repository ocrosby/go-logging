package logging

import (
	"bytes"
	"context"
	"testing"
)

func TestNewLevelDispatcher(t *testing.T) {
	logger := NewTextLogger(InfoLevel)

	dispatcher := NewLevelDispatcher(logger)
	if dispatcher == nil {
		t.Fatal("expected dispatcher to be created")
	}
}

func TestLevelDispatcher_DispatchMethods(t *testing.T) {
	buf := &bytes.Buffer{}
	config := NewLoggerConfig().
		WithLevel(TraceLevel).
		WithWriter(buf).
		WithTextFormat().
		Build()
	logger := NewWithLoggerConfig(config)

	dispatcher := NewLevelDispatcher(logger)

	// Test all dispatch methods
	dispatcher.DispatchTrace("trace message")
	dispatcher.DispatchDebug("debug message")
	dispatcher.DispatchInfo("info message")
	dispatcher.DispatchWarn("warn message")
	dispatcher.DispatchError("error message")
	dispatcher.DispatchCritical("critical message")

	output := buf.String()

	messages := []string{"trace message", "debug message", "info message", "warn message", "error message", "critical message"}
	for _, msg := range messages {
		if !contains(output, msg) {
			t.Errorf("expected message '%s' in output", msg)
		}
	}
}

func TestLevelDispatcher_DispatchContextMethods(t *testing.T) {
	buf := &bytes.Buffer{}
	config := NewLoggerConfig().
		WithLevel(TraceLevel).
		WithWriter(buf).
		WithTextFormat().
		Build()
	logger := NewWithLoggerConfig(config)

	dispatcher := NewLevelDispatcher(logger)
	ctx := context.Background()

	// Test all context dispatch methods
	dispatcher.DispatchTraceContext(ctx, "trace context message")
	dispatcher.DispatchDebugContext(ctx, "debug context message")
	dispatcher.DispatchInfoContext(ctx, "info context message")
	dispatcher.DispatchWarnContext(ctx, "warn context message")
	dispatcher.DispatchErrorContext(ctx, "error context message")
	dispatcher.DispatchCriticalContext(ctx, "critical context message")

	output := buf.String()

	messages := []string{"trace context", "debug context", "info context", "warn context", "error context", "critical context"}
	for _, msg := range messages {
		if !contains(output, msg) {
			t.Errorf("expected message '%s' in output", msg)
		}
	}
}

func TestLoggerLevelMethods(t *testing.T) {
	buf := &bytes.Buffer{}
	config := NewLoggerConfig().
		WithLevel(TraceLevel).
		WithWriter(buf).
		WithTextFormat().
		Build()
	logger := NewWithLoggerConfig(config)

	levelMethods := &LoggerLevelMethods{}
	levelMethods.InitLevelMethods(logger)

	// Test initialized methods
	levelMethods.Trace("init trace message")
	levelMethods.Debug("init debug message")
	levelMethods.Info("init info message")
	levelMethods.Warn("init warn message")
	levelMethods.Error("init error message")
	levelMethods.Critical("init critical message")

	output := buf.String()

	messages := []string{"init trace", "init debug", "init info", "init warn", "init error", "init critical"}
	for _, msg := range messages {
		if !contains(output, msg) {
			t.Errorf("expected message '%s' in output", msg)
		}
	}
}

func TestLoggerLevelMethods_ContextMethods(t *testing.T) {
	buf := &bytes.Buffer{}
	config := NewLoggerConfig().
		WithLevel(TraceLevel).
		WithWriter(buf).
		WithTextFormat().
		Build()
	logger := NewWithLoggerConfig(config)

	levelMethods := &LoggerLevelMethods{}
	levelMethods.InitLevelMethods(logger)

	ctx := context.Background()

	// Test context methods
	levelMethods.TraceContext(ctx, "trace ctx message")
	levelMethods.DebugContext(ctx, "debug ctx message")
	levelMethods.InfoContext(ctx, "info ctx message")
	levelMethods.WarnContext(ctx, "warn ctx message")
	levelMethods.ErrorContext(ctx, "error ctx message")
	levelMethods.CriticalContext(ctx, "critical ctx message")

	output := buf.String()

	messages := []string{"trace ctx", "debug ctx", "info ctx", "warn ctx", "error ctx", "critical ctx"}
	for _, msg := range messages {
		if !contains(output, msg) {
			t.Errorf("expected message '%s' in output", msg)
		}
	}
}

func TestLevelDispatcher_LevelFiltering(t *testing.T) {
	buf := &bytes.Buffer{}
	config := NewLoggerConfig().
		WithLevel(WarnLevel).
		WithWriter(buf).
		WithTextFormat().
		Build()
	logger := NewWithLoggerConfig(config)

	dispatcher := NewLevelDispatcher(logger)

	// These should not be logged (below warn level)
	dispatcher.DispatchTrace("trace should not appear")
	dispatcher.DispatchDebug("debug should not appear")
	dispatcher.DispatchInfo("info should not appear")

	// These should be logged (at or above warn level)
	dispatcher.DispatchWarn("warn should appear")
	dispatcher.DispatchError("error should appear")
	dispatcher.DispatchCritical("critical should appear")

	output := buf.String()

	// Should not contain low-level messages
	if contains(output, "should not appear") {
		t.Error("expected low-level messages to be filtered out")
	}

	// Should contain high-level messages
	if !contains(output, "warn should appear") ||
		!contains(output, "error should appear") ||
		!contains(output, "critical should appear") {
		t.Error("expected high-level messages to be logged")
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || someContains(s, substr)))
}

func someContains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

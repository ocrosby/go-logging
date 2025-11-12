package logging

import (
	"bytes"
	"log/slog"
	"os"
	"strings"
	"testing"
)

func TestGetDefaultLogger(t *testing.T) {
	logger := GetDefaultLogger()
	if logger == nil {
		t.Fatal("expected default logger to be created")
	}

	// Should return same instance on subsequent calls
	logger2 := GetDefaultLogger()
	if logger != logger2 {
		t.Error("expected same logger instance")
	}
}

func TestSetDefaultLogger(t *testing.T) {
	originalLogger := GetDefaultLogger()
	defer SetDefaultLogger(originalLogger) // Restore

	customLogger := NewTextLogger(InfoLevel)
	SetDefaultLogger(customLogger)

	retrieved := GetDefaultLogger()
	if retrieved != customLogger {
		t.Error("expected custom logger to be set as default")
	}
}

func TestNew(t *testing.T) {
	buf := &bytes.Buffer{}

	logger := New(func(builder *ConfigBuilder) {
		builder.WithLevel(DebugLevel).WithOutput(buf)
	})

	if logger == nil {
		t.Fatal("expected logger to be created")
	}

	// Test that the logger works
	logger.Info("test message")
	output := buf.String()
	if !strings.Contains(output, "test message") {
		t.Error("expected log message to appear in output")
	}
}

func TestNewWithLoggerConfig(t *testing.T) {
	buf := &bytes.Buffer{}
	config := NewLoggerConfig().
		WithLevel(WarnLevel).
		WithWriter(buf).
		Build()

	logger := NewWithLoggerConfig(config)
	if logger == nil {
		t.Fatal("expected logger to be created")
	}

	// Test with nil config (should use defaults)
	logger2 := NewWithLoggerConfig(nil)
	if logger2 == nil {
		t.Fatal("expected logger with default config to be created")
	}
}

func TestNewFromEnvironment(t *testing.T) {
	// Set environment variable
	os.Setenv("LOG_LEVEL", "debug")
	defer os.Unsetenv("LOG_LEVEL")

	logger := NewFromEnvironment()
	if logger == nil {
		t.Fatal("expected logger to be created from environment")
	}

	// Check that it respects the environment level
	if !logger.IsLevelEnabled(DebugLevel) {
		t.Error("expected debug level to be enabled from environment")
	}
}

func TestNewWithLevel(t *testing.T) {
	logger := NewWithLevel(ErrorLevel)
	if logger == nil {
		t.Fatal("expected logger to be created")
	}

	if logger.IsLevelEnabled(InfoLevel) {
		t.Error("expected info level to be disabled with error level logger")
	}

	if !logger.IsLevelEnabled(ErrorLevel) {
		t.Error("expected error level to be enabled")
	}
}

func TestNewWithLevelString(t *testing.T) {
	logger := NewWithLevelString("warn")
	if logger == nil {
		t.Fatal("expected logger to be created")
	}

	if logger.IsLevelEnabled(InfoLevel) {
		t.Error("expected info level to be disabled with warn level")
	}

	if !logger.IsLevelEnabled(WarnLevel) {
		t.Error("expected warn level to be enabled")
	}
}

func TestNewJSONLogger(t *testing.T) {
	buf := &bytes.Buffer{}

	// Redirect output to buffer for testing
	config := NewConfig().WithLevel(InfoLevel).WithOutput(buf).WithJSONFormat().Build()
	redactor := ProvideRedactorChain(config)
	logger := ProvideLogger(config, redactor)

	logger.Info("test message")
	output := buf.String()

	// JSON output should contain key-value pairs
	if !strings.Contains(output, "test message") {
		t.Error("expected log message in output")
	}
}

func TestNewTextLogger(t *testing.T) {
	buf := &bytes.Buffer{}

	// Create text logger with buffer output
	config := NewConfig().WithLevel(InfoLevel).WithOutput(buf).WithTextFormat().Build()
	redactor := ProvideRedactorChain(config)
	logger := ProvideLogger(config, redactor)

	logger.Info("test message")
	output := buf.String()

	if !strings.Contains(output, "test message") {
		t.Error("expected log message in output")
	}
}

func TestNewWithHandler(t *testing.T) {
	buf := &bytes.Buffer{}
	handler := slog.NewTextHandler(buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	logger := NewWithHandler(handler)
	if logger == nil {
		t.Fatal("expected logger to be created with handler")
	}
}

func TestNewSlogJSONLogger(t *testing.T) {
	logger := NewSlogJSONLogger(InfoLevel)
	if logger == nil {
		t.Fatal("expected slog JSON logger to be created")
	}
}

func TestNewSlogTextLogger(t *testing.T) {
	logger := NewSlogTextLogger(InfoLevel)
	if logger == nil {
		t.Fatal("expected slog text logger to be created")
	}
}

func TestGlobalLoggingFunctions(t *testing.T) {
	// Save original logger
	originalLogger := GetDefaultLogger()
	defer SetDefaultLogger(originalLogger)

	buf := &bytes.Buffer{}

	// Set custom logger to capture output
	config := NewLoggerConfig().
		WithLevel(TraceLevel).
		WithWriter(buf).
		WithTextFormat().
		Build()
	customLogger := NewWithLoggerConfig(config)
	SetDefaultLogger(customLogger)

	// Test global functions
	Trace("trace message")
	Debug("debug message")
	Info("info message")
	Warn("warn message")
	Error("error message")
	Critical("critical message")

	output := buf.String()

	messages := []string{"trace message", "debug message", "info message", "warn message", "error message", "critical message"}
	for _, msg := range messages {
		if !strings.Contains(output, msg) {
			t.Errorf("expected message '%s' in output, got: %s", msg, output)
		}
	}
}

func TestShorthandLoggerFunctions(t *testing.T) {
	// Test shorthand logger getter functions
	loggers := []Logger{T(), D(), I(), E()}

	for i, logger := range loggers {
		if logger == nil {
			t.Errorf("expected shorthand logger %d to be non-nil", i)
		}

		if logger != GetDefaultLogger() {
			t.Errorf("expected shorthand logger %d to return default logger", i)
		}
	}
}

func TestIsLevelEnabled(t *testing.T) {
	// Set a trace level logger (lowest level, enables all)
	SetDefaultLogger(NewWithLevel(TraceLevel))

	if !IsDebugEnabled() {
		t.Error("expected debug to be enabled")
	}

	if !IsTraceEnabled() {
		t.Error("expected trace to be enabled with trace level")
	}

	// Set higher level logger
	SetDefaultLogger(NewWithLevel(ErrorLevel))

	if IsDebugEnabled() {
		t.Error("expected debug to be disabled with error level")
	}

	if IsTraceEnabled() {
		t.Error("expected trace to be disabled with error level")
	}
}

func TestMustGetEnv(t *testing.T) {
	// Test with existing environment variable
	os.Setenv("TEST_VAR", "test_value")
	defer os.Unsetenv("TEST_VAR")

	value := MustGetEnv("TEST_VAR")
	if value != "test_value" {
		t.Errorf("expected 'test_value', got '%s'", value)
	}
}

func TestMustGetEnv_Missing(t *testing.T) {
	// This is tricky to test since it calls os.Exit(1)
	// We'll test the setup but not the actual exit
	originalLogger := GetDefaultLogger()
	defer SetDefaultLogger(originalLogger)

	SetDefaultLogger(NewTextLogger(DebugLevel))

	// We can't easily test os.Exit, but we know it would be called
	// with missing environment variables
}

func TestNewStandardLogger(t *testing.T) {
	config := NewConfig().WithLevel(InfoLevel).Build()
	redactor := ProvideRedactorChain(config)

	logger := NewStandardLogger(config, redactor)
	if logger == nil {
		t.Fatal("expected standard logger to be created")
	}
}

func TestNewSlogLogger(t *testing.T) {
	config := NewConfig().WithLevel(InfoLevel).UseSlog(true).Build()
	redactor := ProvideRedactorChain(config)

	logger := NewSlogLogger(config, redactor)
	if logger == nil {
		t.Fatal("expected slog logger to be created")
	}
}

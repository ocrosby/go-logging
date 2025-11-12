package logging

import (
	"os"
	"testing"
)

func TestNewSimple(t *testing.T) {
	logger := NewSimple()
	if logger == nil {
		t.Fatal("NewSimple() returned nil")
	}
}

func TestNewEasyJSON(t *testing.T) {
	logger := NewEasyJSON()
	if logger == nil {
		t.Fatal("NewEasyJSON() returned nil")
	}
}

func TestNewEasyJSONWithLevel(t *testing.T) {
	logger := NewEasyJSONWithLevel(DebugLevel)
	if logger == nil {
		t.Fatal("NewEasyJSONWithLevel() returned nil")
	}
}

func TestNewFromEnvSimple(t *testing.T) {
	// Test with default environment (no env vars set)
	logger := NewFromEnvSimple()
	if logger == nil {
		t.Fatal("NewFromEnvSimple() returned nil")
	}

	// Test with environment variables set
	os.Setenv("LOG_LEVEL", "debug")
	os.Setenv("LOG_FORMAT", "json")
	os.Setenv("LOG_INCLUDE_FILE", "true")
	os.Setenv("LOG_INCLUDE_TIME", "false")

	logger2 := NewFromEnvSimple()
	if logger2 == nil {
		t.Fatal("NewFromEnvSimple() with env vars returned nil")
	}

	// Clean up
	os.Unsetenv("LOG_LEVEL")
	os.Unsetenv("LOG_FORMAT")
	os.Unsetenv("LOG_INCLUDE_FILE")
	os.Unsetenv("LOG_INCLUDE_TIME")
}

func TestNewEasyBuilder(t *testing.T) {
	logger := NewEasyBuilder().
		Level(DebugLevel).
		JSON().
		WithFile().
		Field("service", "test").
		Build()

	if logger == nil {
		t.Fatal("NewEasyBuilder().Build() returned nil")
	}
}

func TestEasyLoggerBuilderMethods(t *testing.T) {
	builder := NewEasyBuilder()

	// Test level methods
	builder = builder.Trace().Debug().Info().Warn().Error().Critical()
	if builder == nil {
		t.Fatal("Level methods returned nil")
	}

	// Test format methods
	builder = builder.JSON().Text()
	if builder == nil {
		t.Fatal("Format methods returned nil")
	}

	// Test option methods
	builder = builder.WithFile().WithoutTime()
	if builder == nil {
		t.Fatal("Option methods returned nil")
	}

	// Test field methods
	builder = builder.Field("key", "value").Fields(map[string]interface{}{
		"service": "test",
		"version": "1.0",
	})
	if builder == nil {
		t.Fatal("Field methods returned nil")
	}

	logger := builder.Build()
	if logger == nil {
		t.Fatal("Build() returned nil")
	}
}

func TestGetEnvLevel(t *testing.T) {
	// Test default
	level := getEnvLevel()
	if level != InfoLevel {
		t.Errorf("Expected InfoLevel for empty LOG_LEVEL, got %v", level)
	}

	// Test valid levels
	testCases := []struct {
		envValue string
		expected Level
	}{
		{"trace", TraceLevel},
		{"debug", DebugLevel},
		{"info", InfoLevel},
		{"warn", WarnLevel},
		{"error", ErrorLevel},
		{"critical", CriticalLevel},
		{"TRACE", TraceLevel},
		{"DEBUG", DebugLevel},
		{"invalid", InfoLevel}, // fallback to default
	}

	for _, tc := range testCases {
		os.Setenv("LOG_LEVEL", tc.envValue)
		level := getEnvLevel()
		if level != tc.expected {
			t.Errorf("Expected %v for LOG_LEVEL=%s, got %v", tc.expected, tc.envValue, level)
		}
		os.Unsetenv("LOG_LEVEL")
	}
}

func TestGetEnvFormat(t *testing.T) {
	// Test default
	format := getEnvFormat()
	if format != TextFormat {
		t.Errorf("Expected TextFormat for empty LOG_FORMAT, got %v", format)
	}

	// Test valid formats
	testCases := []struct {
		envValue string
		expected OutputFormat
	}{
		{"text", TextFormat},
		{"json", JSONFormat},
		{"TEXT", TextFormat},
		{"JSON", JSONFormat},
		{"invalid", TextFormat}, // fallback to default
	}

	for _, tc := range testCases {
		os.Setenv("LOG_FORMAT", tc.envValue)
		format := getEnvFormat()
		if format != tc.expected {
			t.Errorf("Expected %v for LOG_FORMAT=%s, got %v", tc.expected, tc.envValue, format)
		}
		os.Unsetenv("LOG_FORMAT")
	}
}

func TestGetEnvBool(t *testing.T) {
	testCases := []struct {
		envValue string
		defValue bool
		expected bool
	}{
		{"true", false, true},
		{"1", false, true},
		{"yes", false, true},
		{"TRUE", false, true},
		{"false", true, false},
		{"0", true, false},
		{"no", true, false},
		{"FALSE", true, false},
		{"", true, true},          // default when empty
		{"", false, false},        // default when empty
		{"invalid", true, true},   // default when invalid
		{"invalid", false, false}, // default when invalid
	}

	for _, tc := range testCases {
		os.Setenv("TEST_BOOL", tc.envValue)
		result := getEnvBool("TEST_BOOL", tc.defValue)
		if result != tc.expected {
			t.Errorf("Expected %v for env=%s, default=%v, got %v",
				tc.expected, tc.envValue, tc.defValue, result)
		}
		os.Unsetenv("TEST_BOOL")
	}
}

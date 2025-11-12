package logging

import (
	"testing"
)

func TestUntestedFactoryMethods(t *testing.T) {
	// Test NewWithLevelString
	logger := NewWithLevelString("debug")
	if logger == nil {
		t.Fatal("expected logger from level string")
	}

	// Test NewFromEnvironment
	logger2 := NewFromEnvironment()
	if logger2 == nil {
		t.Fatal("expected logger from environment")
	}

	// Test NewSlogJSONLogger
	logger3 := NewSlogJSONLogger(InfoLevel)
	if logger3 == nil {
		t.Fatal("expected slog JSON logger")
	}

	// Test NewSlogTextLogger
	logger4 := NewSlogTextLogger(InfoLevel)
	if logger4 == nil {
		t.Fatal("expected slog text logger")
	}
}

func TestMoreMissingMethods(t *testing.T) {
	// Test global shorthand functions
	T().Info("T() test")
	D().Debug("D() test")
	I().Info("I() test")
	E().Error("E() test")

	// Test level enabled functions
	_ = IsDebugEnabled()
	_ = IsTraceEnabled()
}

func TestFormatterEdgeCases(t *testing.T) {

	// Test formatters with edge cases
	config := NewFormatterConfig().
		WithJSONFormat().
		IncludeTime(false).
		IncludeFile(true).
		Build()

	formatter := NewJSONFormatter(config)

	// Test with empty fields
	entry := LogEntry{
		Level:   InfoLevel,
		Message: "test",
		Fields:  map[string]interface{}{},
	}

	_, err := formatter.Format(entry)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Test text formatter too
	textFormatter := NewTextFormatter(config)
	_, err = textFormatter.Format(entry)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestConfigurationEdgeCases(t *testing.T) {
	// Test invalid level strings
	_ = NewLoggerConfig().WithLevelString("invalid")
	// Should not crash and should use default level

	// Test environment variable reading
	config2 := NewLoggerConfig().FromEnvironment().Build()
	if config2 == nil {
		t.Error("expected config from environment")
	}
}

func TestProviderEdgeCases(t *testing.T) {
	// Test provider functions that might not be covered
	output := ProvideOutput()
	if output == nil {
		t.Error("expected output provider to return writer")
	}

	config := ProvideConfig()
	if config == nil {
		t.Error("expected config provider to return config")
	}

	loggerConfig := ProvideLoggerConfig()
	if loggerConfig == nil {
		t.Error("expected logger config provider to return config")
	}
}

func TestLevelStringConversions(t *testing.T) {
	levels := []Level{TraceLevel, DebugLevel, InfoLevel, WarnLevel, ErrorLevel, CriticalLevel}

	for _, level := range levels {
		str := level.String()
		if str == "" {
			t.Errorf("expected non-empty string for level %v", level)
		}

		// Test parsing back
		parsed, ok := ParseLevel(str)
		if !ok {
			t.Errorf("failed to parse level string %s", str)
		}

		if parsed != level {
			t.Errorf("expected %v after parsing %s, got %v", level, str, parsed)
		}
	}
}

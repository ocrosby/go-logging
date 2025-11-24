package logging

import (
	"bytes"
	"log/slog"
	"os"
	"regexp"
	"testing"
)

func TestNewCoreConfig(t *testing.T) {
	builder := NewCoreConfig()
	config := builder.Build()

	if config.Level != InfoLevel {
		t.Errorf("expected default level %v, got %v", InfoLevel, config.Level)
	}

	if config.StaticFields == nil {
		t.Error("expected static fields map to be initialized")
	}
}

func TestCoreConfigBuilder_WithLevel(t *testing.T) {
	builder := NewCoreConfig().WithLevel(DebugLevel)
	config := builder.Build()

	if config.Level != DebugLevel {
		t.Errorf("expected level %v, got %v", DebugLevel, config.Level)
	}
}

func TestCoreConfigBuilder_WithLevelString(t *testing.T) {
	tests := []struct {
		levelStr string
		expected Level
	}{
		{"debug", DebugLevel},
		{"info", InfoLevel},
		{"warn", WarnLevel},
		{"error", ErrorLevel},
		{"invalid", InfoLevel}, // Should keep default
	}

	for _, tt := range tests {
		builder := NewCoreConfig().WithLevelString(tt.levelStr)
		config := builder.Build()

		if config.Level != tt.expected {
			t.Errorf("for level string %s, expected %v, got %v", tt.levelStr, tt.expected, config.Level)
		}
	}
}

func TestCoreConfigBuilder_WithStaticField(t *testing.T) {
	builder := NewCoreConfig().WithStaticField("service", "test-service")
	config := builder.Build()

	if config.StaticFields["service"] != "test-service" {
		t.Error("expected static field to be set")
	}
}

func TestCoreConfigBuilder_WithStaticFields(t *testing.T) {
	fields := map[string]interface{}{
		"service": "test-service",
		"version": "1.0.0",
	}

	builder := NewCoreConfig().WithStaticFields(fields)
	config := builder.Build()

	for k, v := range fields {
		if config.StaticFields[k] != v {
			t.Errorf("expected field %s to be %v, got %v", k, v, config.StaticFields[k])
		}
	}
}

func TestNewFormatterConfig(t *testing.T) {
	builder := NewFormatterConfig()
	config := builder.Build()

	if config.Format != CommonLogFormat {
		t.Errorf("expected default format %v, got %v", CommonLogFormat, config.Format)
	}

	if !config.IncludeFile {
		t.Error("expected IncludeFile to be true by default")
	}

	if !config.IncludeTime {
		t.Error("expected IncludeTime to be true by default")
	}

	if !config.UseShortFile {
		t.Error("expected UseShortFile to be true by default")
	}

	if config.RedactPatterns == nil {
		t.Error("expected RedactPatterns to be initialized")
	}
}

func TestFormatterConfigBuilder_WithFormat(t *testing.T) {
	builder := NewFormatterConfig().WithFormat(JSONFormat)
	config := builder.Build()

	if config.Format != JSONFormat {
		t.Errorf("expected format %v, got %v", JSONFormat, config.Format)
	}
}

func TestFormatterConfigBuilder_WithJSONFormat(t *testing.T) {
	builder := NewFormatterConfig().WithJSONFormat()
	config := builder.Build()

	if config.Format != JSONFormat {
		t.Error("expected JSON format")
	}
}

func TestFormatterConfigBuilder_WithTextFormat(t *testing.T) {
	builder := NewFormatterConfig().WithTextFormat()
	config := builder.Build()

	if config.Format != TextFormat {
		t.Error("expected text format")
	}
}

func TestFormatterConfigBuilder_IncludeFile(t *testing.T) {
	builder := NewFormatterConfig().IncludeFile(false)
	config := builder.Build()

	if config.IncludeFile {
		t.Error("expected IncludeFile to be false")
	}
}

func TestFormatterConfigBuilder_IncludeTime(t *testing.T) {
	builder := NewFormatterConfig().IncludeTime(false)
	config := builder.Build()

	if config.IncludeTime {
		t.Error("expected IncludeTime to be false")
	}
}

func TestFormatterConfigBuilder_UseShortFile(t *testing.T) {
	builder := NewFormatterConfig().UseShortFile(false)
	config := builder.Build()

	if config.UseShortFile {
		t.Error("expected UseShortFile to be false")
	}
}

func TestFormatterConfigBuilder_AddRedactPattern(t *testing.T) {
	builder := NewFormatterConfig().AddRedactPattern(`password`)
	config := builder.Build()

	if len(config.RedactPatterns) != 1 {
		t.Errorf("expected 1 redact pattern, got %d", len(config.RedactPatterns))
	}

	if !config.RedactPatterns[0].MatchString("password123") {
		t.Error("expected pattern to match 'password123'")
	}
}

func TestFormatterConfigBuilder_AddRedactPattern_Invalid(t *testing.T) {
	builder := NewFormatterConfig().AddRedactPattern(`[invalid`)
	config := builder.Build()

	if len(config.RedactPatterns) != 0 {
		t.Error("expected invalid pattern to be ignored")
	}
}

func TestFormatterConfigBuilder_AddRedactRegex(t *testing.T) {
	re := regexp.MustCompile(`secret`)
	builder := NewFormatterConfig().AddRedactRegex(re)
	config := builder.Build()

	if len(config.RedactPatterns) != 1 {
		t.Error("expected 1 redact pattern")
	}

	if config.RedactPatterns[0] != re {
		t.Error("expected pattern to match provided regex")
	}
}

func TestNewOutputConfig(t *testing.T) {
	builder := NewOutputConfig()
	config := builder.Build()

	if config.Writer != os.Stdout {
		t.Error("expected default writer to be os.Stdout")
	}
}

func TestOutputConfigBuilder_WithWriter(t *testing.T) {
	buf := &bytes.Buffer{}
	builder := NewOutputConfig().WithWriter(buf)
	config := builder.Build()

	if config.Writer != buf {
		t.Error("expected writer to match provided buffer")
	}
}

func TestNewLoggerConfig(t *testing.T) {
	builder := NewLoggerConfig()

	if builder.config == nil {
		t.Fatal("expected config to be initialized")
	}

	if builder.config.Core == nil {
		t.Error("expected core config to be initialized")
	}

	if builder.config.Formatter == nil {
		t.Error("expected formatter config to be initialized")
	}

	if builder.config.Output == nil {
		t.Error("expected output config to be initialized")
	}
}

func TestLoggerConfigBuilder_WithCore(t *testing.T) {
	core := NewCoreConfig().WithLevel(ErrorLevel).Build()
	builder := NewLoggerConfig().WithCore(core)
	config := builder.Build()

	if config.Core.Level != ErrorLevel {
		t.Error("expected core config to be set")
	}
}

func TestLoggerConfigBuilder_WithFormatter(t *testing.T) {
	formatter := NewFormatterConfig().WithJSONFormat().Build()
	builder := NewLoggerConfig().WithFormatter(formatter)
	config := builder.Build()

	if config.Formatter.Format != JSONFormat {
		t.Error("expected formatter config to be set")
	}
}

func TestLoggerConfigBuilder_WithOutput(t *testing.T) {
	buf := &bytes.Buffer{}
	output := NewOutputConfig().WithWriter(buf).Build()
	builder := NewLoggerConfig().WithOutput(output)
	config := builder.Build()

	if config.Output.Writer != buf {
		t.Error("expected output config to be set")
	}
}

func TestLoggerConfigBuilder_WithHandler(t *testing.T) {
	handler := slog.NewTextHandler(os.Stdout, nil)
	builder := NewLoggerConfig().WithHandler(handler)
	config := builder.Build()

	if config.Handler != handler {
		t.Error("expected handler to be set")
	}
}

func TestLoggerConfigBuilder_UseSlog(t *testing.T) {
	builder := NewLoggerConfig().UseSlog(true)
	config := builder.Build()

	if !config.UseSlog {
		t.Error("expected UseSlog to be true")
	}
}

func TestLoggerConfigBuilder_FromEnvironment(t *testing.T) {
	// Set environment variable
	os.Setenv("LOG_LEVEL", "warn")
	defer os.Unsetenv("LOG_LEVEL")

	builder := NewLoggerConfig().FromEnvironment()
	config := builder.Build()

	if config.Core.Level != WarnLevel {
		t.Errorf("expected level from environment to be warn, got %v", config.Core.Level)
	}
}

func TestLoggerConfigBuilder_ChainedMethods(t *testing.T) {
	buf := &bytes.Buffer{}

	config := NewLoggerConfig().
		WithLevel(WarnLevel).
		WithWriter(buf).
		WithJSONFormat().
		Build()

	if config.Core.Level != WarnLevel {
		t.Error("expected level to be warn")
	}

	if config.Output.Writer != buf {
		t.Error("expected writer to be set")
	}

	if config.Formatter.Format != JSONFormat {
		t.Error("expected JSON format")
	}
}

func TestLoggerConfigBuilder_WithLevelString(t *testing.T) {
	builder := NewLoggerConfig().WithLevelString("debug")
	config := builder.Build()

	if config.Core.Level != DebugLevel {
		t.Error("expected debug level")
	}
}

func TestLoggerConfigBuilder_WithFormat(t *testing.T) {
	builder := NewLoggerConfig().WithFormat(JSONFormat)
	config := builder.Build()

	if config.Formatter.Format != JSONFormat {
		t.Error("expected JSON format")
	}
}

func TestLoggerConfigBuilder_WithTextFormat(t *testing.T) {
	builder := NewLoggerConfig().WithTextFormat()
	config := builder.Build()

	if config.Formatter.Format != TextFormat {
		t.Error("expected text format")
	}
}

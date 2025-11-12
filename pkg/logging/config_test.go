package logging

import (
	"bytes"
	"log/slog"
	"os"
	"regexp"
	"testing"
)

func TestConfig_ToLoggerConfig(t *testing.T) {
	buf := &bytes.Buffer{}
	re := regexp.MustCompile(`test`)
	staticFields := map[string]interface{}{
		"service": "test",
		"version": "1.0",
	}

	config := &Config{
		Level:          InfoLevel,
		Output:         buf,
		Format:         JSONFormat,
		IncludeFile:    true,
		IncludeTime:    true,
		UseShortFile:   true,
		RedactPatterns: []*regexp.Regexp{re},
		StaticFields:   staticFields,
		UseSlog:        true,
	}

	loggerConfig := config.ToLoggerConfig()

	if loggerConfig.Core.Level != InfoLevel {
		t.Errorf("expected level %v, got %v", InfoLevel, loggerConfig.Core.Level)
	}

	if loggerConfig.Output.Writer != buf {
		t.Error("expected output writer to match")
	}

	if loggerConfig.Formatter.Format != JSONFormat {
		t.Errorf("expected format %v, got %v", JSONFormat, loggerConfig.Formatter.Format)
	}

	if !loggerConfig.Formatter.IncludeFile {
		t.Error("expected IncludeFile to be true")
	}

	if !loggerConfig.Formatter.IncludeTime {
		t.Error("expected IncludeTime to be true")
	}

	if !loggerConfig.Formatter.UseShortFile {
		t.Error("expected UseShortFile to be true")
	}

	if len(loggerConfig.Formatter.RedactPatterns) != 1 {
		t.Errorf("expected 1 redact pattern, got %d", len(loggerConfig.Formatter.RedactPatterns))
	}

	if loggerConfig.Core.StaticFields["service"] != "test" {
		t.Error("expected static field service to be 'test'")
	}

	if !loggerConfig.UseSlog {
		t.Error("expected UseSlog to be true")
	}
}

func TestNewConfig(t *testing.T) {
	builder := NewConfig()
	if builder == nil {
		t.Fatal("expected builder to be created")
	}

	if builder.builder == nil {
		t.Fatal("expected internal builder to be created")
	}
}

func TestConfigBuilder_WithLevel(t *testing.T) {
	builder := NewConfig().WithLevel(WarnLevel)
	config := builder.Build()

	if config.Level != WarnLevel {
		t.Errorf("expected level %v, got %v", WarnLevel, config.Level)
	}
}

func TestConfigBuilder_WithLevelString(t *testing.T) {
	builder := NewConfig().WithLevelString("error")
	config := builder.Build()

	if config.Level != ErrorLevel {
		t.Errorf("expected level %v, got %v", ErrorLevel, config.Level)
	}
}

func TestConfigBuilder_WithOutput(t *testing.T) {
	buf := &bytes.Buffer{}
	builder := NewConfig().WithOutput(buf)
	config := builder.Build()

	if config.Output != buf {
		t.Error("expected output to match provided buffer")
	}
}

func TestConfigBuilder_WithFormat(t *testing.T) {
	builder := NewConfig().WithFormat(JSONFormat)
	config := builder.Build()

	if config.Format != JSONFormat {
		t.Errorf("expected format %v, got %v", JSONFormat, config.Format)
	}
}

func TestConfigBuilder_WithJSONFormat(t *testing.T) {
	builder := NewConfig().WithJSONFormat()
	config := builder.Build()

	if config.Format != JSONFormat {
		t.Errorf("expected JSON format, got %v", config.Format)
	}
}

func TestConfigBuilder_WithTextFormat(t *testing.T) {
	builder := NewConfig().WithTextFormat()
	config := builder.Build()

	if config.Format != TextFormat {
		t.Errorf("expected text format, got %v", config.Format)
	}
}

func TestConfigBuilder_IncludeFile(t *testing.T) {
	builder := NewConfig().IncludeFile(true)
	config := builder.Build()

	if !config.IncludeFile {
		t.Error("expected IncludeFile to be true")
	}
}

func TestConfigBuilder_IncludeTime(t *testing.T) {
	builder := NewConfig().IncludeTime(true)
	config := builder.Build()

	if !config.IncludeTime {
		t.Error("expected IncludeTime to be true")
	}
}

func TestConfigBuilder_UseShortFile(t *testing.T) {
	builder := NewConfig().UseShortFile(true)
	config := builder.Build()

	if !config.UseShortFile {
		t.Error("expected UseShortFile to be true")
	}
}

func TestConfigBuilder_AddRedactPattern(t *testing.T) {
	builder := NewConfig().AddRedactPattern(`api[_-]?key`)
	config := builder.Build()

	if len(config.RedactPatterns) != 1 {
		t.Errorf("expected 1 redact pattern, got %d", len(config.RedactPatterns))
	}

	if !config.RedactPatterns[0].MatchString("api_key") {
		t.Error("expected pattern to match api_key")
	}
}

func TestConfigBuilder_AddRedactPattern_InvalidRegex(t *testing.T) {
	builder := NewConfig().AddRedactPattern(`[invalid`)
	config := builder.Build()

	// Invalid regex should be ignored
	if len(config.RedactPatterns) != 0 {
		t.Errorf("expected 0 redact patterns for invalid regex, got %d", len(config.RedactPatterns))
	}
}

func TestConfigBuilder_AddRedactRegex(t *testing.T) {
	re := regexp.MustCompile(`password`)
	builder := NewConfig().AddRedactRegex(re)
	config := builder.Build()

	if len(config.RedactPatterns) != 1 {
		t.Errorf("expected 1 redact pattern, got %d", len(config.RedactPatterns))
	}

	if config.RedactPatterns[0] != re {
		t.Error("expected regex to match provided regex")
	}
}

func TestConfigBuilder_WithStaticField(t *testing.T) {
	builder := NewConfig().WithStaticField("service", "test-service")
	config := builder.Build()

	if config.StaticFields["service"] != "test-service" {
		t.Error("expected static field to be set")
	}
}

func TestConfigBuilder_WithStaticFields(t *testing.T) {
	fields := map[string]interface{}{
		"service": "test-service",
		"version": "1.0.0",
		"env":     "test",
	}

	builder := NewConfig().WithStaticFields(fields)
	config := builder.Build()

	for k, v := range fields {
		if config.StaticFields[k] != v {
			t.Errorf("expected static field %s to be %v, got %v", k, v, config.StaticFields[k])
		}
	}
}

func TestConfigBuilder_WithHandler(t *testing.T) {
	handler := slog.NewTextHandler(os.Stdout, nil)
	builder := NewConfig().WithHandler(handler)
	config := builder.Build()

	if config.Handler != handler {
		t.Error("expected handler to match provided handler")
	}
}

func TestConfigBuilder_UseSlog(t *testing.T) {
	builder := NewConfig().UseSlog(true)
	config := builder.Build()

	if !config.UseSlog {
		t.Error("expected UseSlog to be true")
	}
}

func TestConfigBuilder_FromEnvironment(t *testing.T) {
	// Set environment variable
	os.Setenv("LOG_LEVEL", "debug")
	defer os.Unsetenv("LOG_LEVEL")

	builder := NewConfig().FromEnvironment()
	config := builder.Build()

	if config.Level != DebugLevel {
		t.Errorf("expected level from environment to be debug, got %v", config.Level)
	}
}

func TestConfigBuilder_ChainedCalls(t *testing.T) {
	buf := &bytes.Buffer{}
	fields := map[string]interface{}{
		"service": "test",
	}

	config := NewConfig().
		WithLevel(WarnLevel).
		WithOutput(buf).
		WithJSONFormat().
		IncludeFile(true).
		IncludeTime(true).
		UseShortFile(true).
		AddRedactPattern(`secret`).
		WithStaticFields(fields).
		UseSlog(true).
		Build()

	if config.Level != WarnLevel {
		t.Errorf("expected level %v, got %v", WarnLevel, config.Level)
	}

	if config.Output != buf {
		t.Error("expected output to match")
	}

	if config.Format != JSONFormat {
		t.Error("expected JSON format")
	}

	if !config.IncludeFile || !config.IncludeTime || !config.UseShortFile {
		t.Error("expected file options to be true")
	}

	if len(config.RedactPatterns) != 1 {
		t.Error("expected 1 redact pattern")
	}

	if config.StaticFields["service"] != "test" {
		t.Error("expected static field to be set")
	}

	if !config.UseSlog {
		t.Error("expected UseSlog to be true")
	}
}

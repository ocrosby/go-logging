package logging

import (
	"bytes"
	"log/slog"
	"regexp"
	"testing"
)

func TestProvideConfig(t *testing.T) {
	config := ProvideConfig()
	if config == nil {
		t.Fatal("expected config to be provided")
		return
	}

	if config.Level != InfoLevel {
		t.Errorf("expected default level %v, got %v", InfoLevel, config.Level)
	}
}

func TestProvideConfigWithLevel(t *testing.T) {
	config := ProvideConfigWithLevel(DebugLevel)
	if config == nil {
		t.Fatal("expected config to be provided")
	}

	if config != nil && config.Level != DebugLevel {
		t.Errorf("expected level %v, got %v", DebugLevel, config.Level)
	}
}

func TestProvideLoggerConfig(t *testing.T) {
	config := ProvideLoggerConfig()
	if config == nil {
		t.Fatal("expected logger config to be provided")
	}

	if config != nil && config.Core.Level != InfoLevel {
		t.Errorf("expected level %v, got %v", InfoLevel, config.Core.Level)
	}
}

func TestProvideLoggerConfigWithLevel(t *testing.T) {
	config := ProvideLoggerConfigWithLevel(WarnLevel)
	if config == nil {
		t.Fatal("expected logger config to be provided")
	}

	if config != nil && config.Core.Level != WarnLevel {
		t.Errorf("expected level %v, got %v", WarnLevel, config.Core.Level)
	}
}

func TestProvideOutput(t *testing.T) {
	output := ProvideOutput()
	if output == nil {
		t.Fatal("expected output to be provided")
	}

	// Should be able to write to it
	_, err := output.Write([]byte("test"))
	if err != nil {
		t.Errorf("unexpected error writing to output: %v", err)
	}
}

func TestProvideRedactorChain(t *testing.T) {
	config := &Config{
		RedactPatterns: []*regexp.Regexp{
			regexp.MustCompile(`password=\w+`),
		},
	}

	chain := ProvideRedactorChain(config)
	if chain == nil {
		t.Fatal("expected redactor chain to be provided")
	}

	// Test redaction
	input := "login with password=secret123"
	output := chain.Redact(input)

	if output == input {
		t.Error("expected input to be redacted")
	}
}

func TestProvideRedactorChainWithPatterns(t *testing.T) {
	// Test with pattern
	pattern := regexp.MustCompile(`api_key=\w+`)

	chain := ProvideRedactorChainWithPatterns(pattern)
	if chain == nil {
		t.Fatal("expected redactor chain to be provided")
	}

	// Test that redaction works
	input := "request with api_key=secret123"
	output := chain.Redact(input)

	if output == input {
		t.Error("expected input to be redacted")
	}
}

func TestProvideRedactorChainFromLoggerConfig(t *testing.T) {
	config := &LoggerConfig{
		Formatter: &FormatterConfig{
			RedactPatterns: []*regexp.Regexp{
				regexp.MustCompile(`secret=\w+`),
			},
		},
	}

	chain := ProvideRedactorChainFromLoggerConfig(config)
	if chain == nil {
		t.Fatal("expected redactor chain to be provided")
	}

	// Test redaction
	input := "data with secret=hidden123"
	output := chain.Redact(input)

	if output == input {
		t.Error("expected input to be redacted")
	}
}

func TestProvideLogger(t *testing.T) {
	buf := &bytes.Buffer{}
	config := &Config{
		Level:  InfoLevel,
		Output: buf,
		Format: TextFormat,
	}
	chain := ProvideRedactorChain(config)

	logger := ProvideLogger(config, chain)
	if logger == nil {
		t.Fatal("expected logger to be provided")
	}

	// Test that logger works
	logger.Info("test message")
	if buf.Len() == 0 {
		t.Error("expected output to be written")
	}
}

func TestProvideLogger_WithSlog(t *testing.T) {
	buf := &bytes.Buffer{}
	handler := slog.NewTextHandler(buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	config := &Config{
		Level:   InfoLevel,
		Output:  buf,
		Format:  TextFormat,
		Handler: handler,
		UseSlog: true,
	}
	chain := ProvideRedactorChain(config)

	logger := ProvideLogger(config, chain)
	if logger == nil {
		t.Fatal("expected slog logger to be provided")
	}
}

func TestProvideLoggerFromConfig(t *testing.T) {
	buf := &bytes.Buffer{}
	config := &LoggerConfig{
		Core: &CoreConfig{
			Level: InfoLevel,
		},
		Formatter: &FormatterConfig{
			Format: TextFormat,
		},
		Output: &OutputConfig{
			Writer: buf,
		},
	}
	chain := ProvideRedactorChainFromLoggerConfig(config)

	logger := ProvideLoggerFromConfig(config, chain)
	if logger == nil {
		t.Fatal("expected logger to be provided")
	}

	// Test that logger works
	logger.Info("test message from new config")
	if buf.Len() == 0 {
		t.Error("expected output to be written")
	}
}

func TestProvideLoggerFromConfig_WithSlog(t *testing.T) {
	buf := &bytes.Buffer{}
	handler := slog.NewTextHandler(buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	config := &LoggerConfig{
		Core: &CoreConfig{
			Level: InfoLevel,
		},
		Formatter: &FormatterConfig{
			Format: TextFormat,
		},
		Output: &OutputConfig{
			Writer: buf,
		},
		Handler: handler,
		UseSlog: true,
	}
	chain := ProvideRedactorChainFromLoggerConfig(config)

	logger := ProvideLoggerFromConfig(config, chain)
	if logger == nil {
		t.Fatal("expected slog logger to be provided")
	}
}

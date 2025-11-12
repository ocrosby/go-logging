package logging

import (
	"os"
	"strings"
)

// NewSimple creates a logger with sensible defaults:
// - INFO level
// - Text format
// - Outputs to stdout
// This is the simplest way to get started with logging.
func NewSimple() Logger {
	return NewWithLevel(InfoLevel)
}

// NewEasyJSON creates a logger with JSON format and INFO level.
// Perfect for production environments or structured logging needs.
func NewEasyJSON() Logger {
	return NewEasyJSONWithLevel(InfoLevel)
}

// NewEasyJSONWithLevel creates a JSON logger with a specific level.
func NewEasyJSONWithLevel(level Level) Logger {
	config := NewLoggerConfig().
		WithLevel(level).
		WithJSONFormat().
		Build()
	return NewWithLoggerConfig(config)
}

// NewFromEnvSimple creates a logger configured from environment variables:
// - LOG_LEVEL: trace, debug, info, warn, error, critical (default: info)
// - LOG_FORMAT: text, json (default: text)
// - LOG_INCLUDE_FILE: true, false (default: false)
// - LOG_INCLUDE_TIME: true, false (default: true)
func NewFromEnvSimple() Logger {
	level := getEnvLevel()
	format := getEnvFormat()
	includeFile := getEnvBool("LOG_INCLUDE_FILE", false)
	includeTime := getEnvBool("LOG_INCLUDE_TIME", true)

	config := NewLoggerConfig().
		WithLevel(level).
		WithFormatter(NewFormatterConfig().
			WithFormat(format).
			IncludeFile(includeFile).
			IncludeTime(includeTime).
			Build()).
		Build()

	return NewWithLoggerConfig(config)
}

// NewEasyBuilder returns a fluent builder for progressive configuration.
// Use this when you need more control than the simple factory functions.
//
// Example:
//
//	logger := logging.NewEasyBuilder().
//	    Level(logging.DebugLevel).
//	    JSON().
//	    Field("service", "my-app").
//	    Build()
func NewEasyBuilder() *EasyLoggerBuilder {
	return &EasyLoggerBuilder{
		level:       InfoLevel,
		format:      TextFormat,
		includeFile: false,
		includeTime: true,
		fields:      make(map[string]interface{}),
	}
}

// EasyLoggerBuilder provides a fluent API for configuring loggers.
type EasyLoggerBuilder struct {
	level       Level
	format      OutputFormat
	includeFile bool
	includeTime bool
	fields      map[string]interface{}
}

// Level sets the minimum logging level.
func (b *EasyLoggerBuilder) Level(level Level) *EasyLoggerBuilder {
	b.level = level
	return b
}

// Trace sets the level to TRACE (most verbose).
func (b *EasyLoggerBuilder) Trace() *EasyLoggerBuilder {
	return b.Level(TraceLevel)
}

// Debug sets the level to DEBUG.
func (b *EasyLoggerBuilder) Debug() *EasyLoggerBuilder {
	return b.Level(DebugLevel)
}

// Info sets the level to INFO (default).
func (b *EasyLoggerBuilder) Info() *EasyLoggerBuilder {
	return b.Level(InfoLevel)
}

// Warn sets the level to WARN.
func (b *EasyLoggerBuilder) Warn() *EasyLoggerBuilder {
	return b.Level(WarnLevel)
}

// Error sets the level to ERROR.
func (b *EasyLoggerBuilder) Error() *EasyLoggerBuilder {
	return b.Level(ErrorLevel)
}

// Critical sets the level to CRITICAL (least verbose).
func (b *EasyLoggerBuilder) Critical() *EasyLoggerBuilder {
	return b.Level(CriticalLevel)
}

// JSON enables JSON formatting.
func (b *EasyLoggerBuilder) JSON() *EasyLoggerBuilder {
	b.format = JSONFormat
	return b
}

// Text enables text formatting (default).
func (b *EasyLoggerBuilder) Text() *EasyLoggerBuilder {
	b.format = TextFormat
	return b
}

// WithFile includes file and line information in log entries.
func (b *EasyLoggerBuilder) WithFile() *EasyLoggerBuilder {
	b.includeFile = true
	return b
}

// WithoutTime excludes timestamp from log entries.
func (b *EasyLoggerBuilder) WithoutTime() *EasyLoggerBuilder {
	b.includeTime = false
	return b
}

// Field adds a static field that will be included in all log entries.
func (b *EasyLoggerBuilder) Field(key string, value interface{}) *EasyLoggerBuilder {
	b.fields[key] = value
	return b
}

// Fields adds multiple static fields.
func (b *EasyLoggerBuilder) Fields(fields map[string]interface{}) *EasyLoggerBuilder {
	for k, v := range fields {
		b.fields[k] = v
	}
	return b
}

// Build creates the logger with the configured options.
func (b *EasyLoggerBuilder) Build() Logger {
	config := NewLoggerConfig().
		WithLevel(b.level).
		WithFormatter(NewFormatterConfig().
			WithFormat(b.format).
			IncludeFile(b.includeFile).
			IncludeTime(b.includeTime).
			Build()).
		WithCore(NewCoreConfig().
			WithLevel(b.level).
			WithStaticFields(b.fields).
			Build()).
		Build()

	return NewWithLoggerConfig(config)
}

// Helper functions for environment variable parsing

func getEnvLevel() Level {
	levelStr := strings.ToLower(os.Getenv("LOG_LEVEL"))
	if levelStr == "" {
		return InfoLevel
	}

	level, ok := ParseLevel(levelStr)
	if !ok {
		return InfoLevel
	}
	return level
}

func getEnvFormat() OutputFormat {
	formatStr := strings.ToLower(os.Getenv("LOG_FORMAT"))
	switch formatStr {
	case "json":
		return JSONFormat
	case "text", "":
		return TextFormat
	default:
		return TextFormat
	}
}

func getEnvBool(key string, defaultValue bool) bool {
	value := strings.ToLower(os.Getenv(key))
	switch value {
	case "true", "1", "yes":
		return true
	case "false", "0", "no":
		return false
	default:
		return defaultValue
	}
}

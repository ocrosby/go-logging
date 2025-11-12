package logging

import (
	"context"
	"time"
)

// Logger defines the complete logging interface with structured logging support.
// All implementations must be thread-safe and support concurrent usage.
//
// Example usage:
//
//	logger := logging.NewWithLevel(logging.InfoLevel)
//	logger.Info("Application started")
//	logger.WithField("user_id", 123).Info("User logged in")
type Logger interface {
	// Core logging methods
	Log(level Level, msg string, args ...interface{})
	LogContext(ctx context.Context, level Level, msg string, args ...interface{})

	// Field attachment methods (immutable)
	WithField(key string, value interface{}) Logger
	WithFields(fields map[string]interface{}) Logger

	// Level checking
	IsLevelEnabled(level Level) bool

	// Level-specific logging methods
	Trace(msg string, args ...interface{})
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Critical(msg string, args ...interface{})

	// Context-aware level-specific logging methods
	TraceContext(ctx context.Context, msg string, args ...interface{})
	DebugContext(ctx context.Context, msg string, args ...interface{})
	InfoContext(ctx context.Context, msg string, args ...interface{})
	WarnContext(ctx context.Context, msg string, args ...interface{})
	ErrorContext(ctx context.Context, msg string, args ...interface{})
	CriticalContext(ctx context.Context, msg string, args ...interface{})

	// Fluent interface support
	Fluent() FluentLogger

	// Configuration methods
	SetLevel(level Level)
	GetLevel() Level
}

// ConfigurableLogger allows runtime configuration changes.
type ConfigurableLogger interface {
	Logger

	// SetLevel dynamically changes the minimum log level.
	SetLevel(level Level)

	// GetLevel returns the current minimum log level.
	GetLevel() Level
}

// LogEntry represents a structured log entry with all its metadata.
type LogEntry struct {
	Timestamp time.Time
	Level     Level
	Message   string
	Fields    map[string]interface{}
	Context   context.Context
	File      string
	Line      int
}

// Formatter defines how log entries are converted to output format.
type Formatter interface {
	// Format converts a LogEntry to bytes for output.
	Format(entry LogEntry) ([]byte, error)
}

// Output defines where formatted log entries are written.
type Output interface {
	// Write outputs the formatted log data.
	Write(data []byte) error

	// Close cleanly shuts down the output (optional).
	Close() error
}

// BufferedOutputInterface extends Output with buffering capabilities.
type BufferedOutputInterface interface {
	Output

	// Flush ensures all buffered data is written.
	Flush() error
}

// AsyncOutputInterface extends Output with asynchronous writing capabilities.
type AsyncOutputInterface interface {
	Output

	// Stop gracefully shuts down async processing.
	Stop() error
}

// FluentLogger provides a fluent interface for building log entries
// with chained method calls. Each method returns the FluentEntry for
// further chaining until Msg() or Msgf() is called.
//
// Example:
//
//	logger.Fluent().Error().
//		Err(err).
//		Str("operation", "database_query").
//		Int("retry_count", 3).
//		Msg("Query failed after retries")
type FluentLogger interface {
	// Trace creates a new fluent entry at TRACE level.
	Trace() *FluentEntry

	// Debug creates a new fluent entry at DEBUG level.
	Debug() *FluentEntry

	// Info creates a new fluent entry at INFO level.
	Info() *FluentEntry

	// Warn creates a new fluent entry at WARN level.
	Warn() *FluentEntry

	// Error creates a new fluent entry at ERROR level.
	Error() *FluentEntry

	// Critical creates a new fluent entry at CRITICAL level.
	Critical() *FluentEntry
}

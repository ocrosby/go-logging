package logging

import "context"

// Logger defines the core logging interface with multiple log levels,
// structured logging support, and context awareness. All implementations
// must be thread-safe and support concurrent usage.
//
// Example usage:
//
//	logger := logging.NewWithLevel(logging.InfoLevel)
//	logger.Info("Application started")
//	logger.WithField("user_id", 123).Info("User logged in")
type Logger interface {
	// Trace logs a message at TRACE level with optional formatting arguments.
	// This is the most verbose level and should be used for detailed debugging.
	Trace(msg string, args ...interface{})

	// Debug logs a message at DEBUG level with optional formatting arguments.
	// Use this for diagnostic information useful during development.
	Debug(msg string, args ...interface{})

	// Info logs a message at INFO level with optional formatting arguments.
	// This is the default level for general informational messages.
	Info(msg string, args ...interface{})

	// Warn logs a message at WARN level with optional formatting arguments.
	// Use this for warning conditions that don't prevent operation.
	Warn(msg string, args ...interface{})

	// Error logs a message at ERROR level with optional formatting arguments.
	// Use this for error conditions that may affect functionality.
	Error(msg string, args ...interface{})

	// Critical logs a message at CRITICAL level with optional formatting arguments.
	// This is the least verbose level for critical conditions requiring immediate attention.
	Critical(msg string, args ...interface{})

	// TraceContext logs a message at TRACE level with context for trace/request ID propagation.
	TraceContext(ctx context.Context, msg string, args ...interface{})

	// DebugContext logs a message at DEBUG level with context for trace/request ID propagation.
	DebugContext(ctx context.Context, msg string, args ...interface{})

	// InfoContext logs a message at INFO level with context for trace/request ID propagation.
	InfoContext(ctx context.Context, msg string, args ...interface{})

	// WarnContext logs a message at WARN level with context for trace/request ID propagation.
	WarnContext(ctx context.Context, msg string, args ...interface{})

	// ErrorContext logs a message at ERROR level with context for trace/request ID propagation.
	ErrorContext(ctx context.Context, msg string, args ...interface{})

	// CriticalContext logs a message at CRITICAL level with context for trace/request ID propagation.
	CriticalContext(ctx context.Context, msg string, args ...interface{})

	// WithField returns a new Logger instance with an additional field attached.
	// The original logger is not modified (immutable pattern).
	//
	// Example:
	//	userLogger := logger.WithField("user_id", 123)
	//	userLogger.Info("User action")
	WithField(key string, value interface{}) Logger

	// WithFields returns a new Logger instance with multiple fields attached.
	// The original logger is not modified (immutable pattern).
	//
	// Example:
	//	contextLogger := logger.WithFields(map[string]interface{}{
	//		"service": "api",
	//		"version": "1.0.0",
	//	})
	WithFields(fields map[string]interface{}) Logger

	// IsLevelEnabled returns true if the given level will produce output.
	// Use this to avoid expensive operations when the level is disabled.
	//
	// Example:
	//	if logger.IsLevelEnabled(logging.DebugLevel) {
	//		logger.Debug("Expensive debug info: %v", computeExpensiveData())
	//	}
	IsLevelEnabled(level Level) bool

	// SetLevel dynamically changes the minimum log level.
	// Only messages at or above this level will be output.
	SetLevel(level Level)

	// GetLevel returns the current minimum log level.
	GetLevel() Level

	// Fluent returns a fluent interface for expressive chained logging.
	//
	// Example:
	//	logger.Fluent().Info().
	//		Str("service", "api").
	//		Int("user_id", 123).
	//		Msg("User logged in")
	Fluent() FluentLogger
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

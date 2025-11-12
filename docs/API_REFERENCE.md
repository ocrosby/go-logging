# API Reference

Complete API reference for the go-logging library with unified architecture.

## Table of Contents

- [Core Interfaces](#core-interfaces)
- [Logger Interface](#logger-interface)
- [Factory Functions](#factory-functions)
- [Configuration](#configuration)
- [Level System](#level-system)
- [Context Support](#context-support)
- [Fluent Interface](#fluent-interface)
- [Handler System](#handler-system)
- [Middleware](#middleware)
- [Async Processing](#async-processing)
- [Utilities](#utilities)

## Core Interfaces

### Logger Interface

The main Logger interface provides all logging capabilities in a unified interface:

```go
type Logger interface {
    // Core logging methods
    Log(level Level, msg string, args ...interface{})
    LogContext(ctx context.Context, level Level, msg string, args ...interface{})

    // Field attachment (immutable pattern)
    WithField(key string, value interface{}) Logger
    WithFields(fields map[string]interface{}) Logger

    // Level checking
    IsLevelEnabled(level Level) bool

    // Level-specific methods
    Trace(msg string, args ...interface{})
    Debug(msg string, args ...interface{})
    Info(msg string, args ...interface{})
    Warn(msg string, args ...interface{})
    Error(msg string, args ...interface{})
    Critical(msg string, args ...interface{})

    // Context-aware variants
    TraceContext(ctx context.Context, msg string, args ...interface{})
    DebugContext(ctx context.Context, msg string, args ...interface{})
    InfoContext(ctx context.Context, msg string, args ...interface{})
    WarnContext(ctx context.Context, msg string, args ...interface{})
    ErrorContext(ctx context.Context, msg string, args ...interface{})
    CriticalContext(ctx context.Context, msg string, args ...interface{})

    // Fluent interface
    Fluent() FluentLogger

    // Configuration
    SetLevel(level Level)
    GetLevel() Level
}
```

#### Method Details

**Core Methods**
- `Log(level, msg, args...)` - Log a message at specified level with optional formatting
- `LogContext(ctx, level, msg, args...)` - Log with context for tracing support

**Field Methods**
- `WithField(key, value) Logger` - Return new logger with additional field (immutable)
- `WithFields(fields) Logger` - Return new logger with multiple fields (immutable)

**Level Methods**
- `Trace/Debug/Info/Warn/Error/Critical(msg, args...)` - Level-specific logging
- `TraceContext/DebugContext/...Context(ctx, msg, args...)` - Context-aware variants

**Utility Methods**
- `IsLevelEnabled(level) bool` - Check if level would produce output
- `Fluent() FluentLogger` - Access fluent interface
- `SetLevel(level)/GetLevel() Level` - Runtime level management

### FluentLogger Interface

Provides expressive method chaining for log construction:

```go
type FluentLogger interface {
    // Level selection
    Trace() *FluentEntry
    Debug() *FluentEntry
    Info() *FluentEntry
    Warn() *FluentEntry
    Error() *FluentEntry
    Critical() *FluentEntry
}
```

### FluentEntry Methods

```go
type FluentEntry struct {
    // Field methods (chainable)
    Str(key, value string) *FluentEntry
    Int(key string, value int) *FluentEntry
    Int64(key string, value int64) *FluentEntry
    Float64(key string, value float64) *FluentEntry
    Bool(key string, value bool) *FluentEntry
    Err(err error) *FluentEntry
    Field(key string, value interface{}) *FluentEntry
    Fields(fields map[string]interface{}) *FluentEntry
    
    // Context methods
    Ctx(ctx context.Context) *FluentEntry
    TraceID(id string) *FluentEntry
    
    // Output methods (terminal)
    Msg(msg string)
    Msgf(format string, args ...interface{})
}
```

## Factory Functions

### Simple Creation

```go
// Create logger with specific level
func NewWithLevel(level Level) Logger

// Create from environment variables (LOG_LEVEL, LOG_FORMAT)
func NewFromEnvironment() Logger

// Create with slog backend
func NewSlogTextLogger(level Level) Logger
func NewSlogJSONLogger(level Level) Logger

// Create with custom slog handler
func NewWithHandler(handler slog.Handler) Logger

// Create JSON/text loggers
func NewJSONLogger(level Level) Logger
func NewTextLogger(level Level) Logger
```

### Configuration-Based Creation

```go
// Using new structured configuration
func NewWithLoggerConfig(config *LoggerConfig) Logger

// Using legacy configuration (backward compatible)
func NewStandardLogger(config *Config, redactorChain RedactorChainInterface) Logger
```

### Global Functions

```go
// Package-level convenience functions
func Trace(msg string, args ...interface{})
func Debug(msg string, args ...interface{})
func Info(msg string, args ...interface{})
func Warn(msg string, args ...interface{})
func Error(msg string, args ...interface{})
func Critical(msg string, args ...interface{})

// Global logger management
func GetDefaultLogger() Logger
func SetDefaultLogger(logger Logger)

// Level checking
func IsDebugEnabled() bool
func IsTraceEnabled() bool
```

## Configuration

### New Structured Configuration

```go
type LoggerConfig struct {
    Core      *CoreConfig
    Formatter *FormatterConfig
    Output    *OutputConfig
    Handler   slog.Handler
    UseSlog   bool
}

type CoreConfig struct {
    Level        Level
    StaticFields map[string]interface{}
}

type FormatterConfig struct {
    Format         OutputFormat
    IncludeFile    bool
    IncludeTime    bool
    UseShortFile   bool
    RedactPatterns []*regexp.Regexp
}

type OutputConfig struct {
    Writer io.Writer
}
```

### Configuration Builders

```go
// Core configuration
func NewCoreConfig() *CoreConfigBuilder
func (b *CoreConfigBuilder) WithLevel(level Level) *CoreConfigBuilder
func (b *CoreConfigBuilder) WithStaticField(key string, value interface{}) *CoreConfigBuilder
func (b *CoreConfigBuilder) WithStaticFields(fields map[string]interface{}) *CoreConfigBuilder
func (b *CoreConfigBuilder) Build() *CoreConfig

// Formatter configuration
func NewFormatterConfig() *FormatterConfigBuilder
func (b *FormatterConfigBuilder) WithFormat(format OutputFormat) *FormatterConfigBuilder
func (b *FormatterConfigBuilder) WithJSONFormat() *FormatterConfigBuilder
func (b *FormatterConfigBuilder) WithTextFormat() *FormatterConfigBuilder
func (b *FormatterConfigBuilder) IncludeFile(include bool) *FormatterConfigBuilder
func (b *FormatterConfigBuilder) IncludeTime(include bool) *FormatterConfigBuilder
func (b *FormatterConfigBuilder) UseShortFile(useShort bool) *FormatterConfigBuilder
func (b *FormatterConfigBuilder) AddRedactPattern(pattern string) *FormatterConfigBuilder
func (b *FormatterConfigBuilder) Build() *FormatterConfig

// Output configuration
func NewOutputConfig() *OutputConfigBuilder
func (b *OutputConfigBuilder) WithWriter(w io.Writer) *OutputConfigBuilder
func (b *OutputConfigBuilder) Build() *OutputConfig

// Complete logger configuration
func NewLoggerConfig() *LoggerConfigBuilder
func (b *LoggerConfigBuilder) WithCore(core *CoreConfig) *LoggerConfigBuilder
func (b *LoggerConfigBuilder) WithFormatter(formatter *FormatterConfig) *LoggerConfigBuilder
func (b *LoggerConfigBuilder) WithOutput(output *OutputConfig) *LoggerConfigBuilder
func (b *LoggerConfigBuilder) WithHandler(handler slog.Handler) *LoggerConfigBuilder
func (b *LoggerConfigBuilder) UseSlog(use bool) *LoggerConfigBuilder
func (b *LoggerConfigBuilder) FromEnvironment() *LoggerConfigBuilder
func (b *LoggerConfigBuilder) Build() *LoggerConfig
```

### Legacy Configuration (Backward Compatible)

```go
type Config struct {
    Level          Level
    Output         io.Writer
    Format         OutputFormat
    IncludeFile    bool
    IncludeTime    bool
    UseShortFile   bool
    RedactPatterns []*regexp.Regexp
    StaticFields   map[string]interface{}
    Handler        slog.Handler
    UseSlog        bool
}

func NewConfig() *ConfigBuilder
func (b *ConfigBuilder) WithLevel(level Level) *ConfigBuilder
func (b *ConfigBuilder) WithOutput(w io.Writer) *ConfigBuilder
func (b *ConfigBuilder) WithFormat(format OutputFormat) *ConfigBuilder
func (b *ConfigBuilder) WithJSONFormat() *ConfigBuilder
func (b *ConfigBuilder) WithTextFormat() *ConfigBuilder
func (b *ConfigBuilder) IncludeFile(include bool) *ConfigBuilder
func (b *ConfigBuilder) IncludeTime(include bool) *ConfigBuilder
func (b *ConfigBuilder) UseShortFile(useShort bool) *ConfigBuilder
func (b *ConfigBuilder) AddRedactPattern(pattern string) *ConfigBuilder
func (b *ConfigBuilder) WithStaticField(key string, value interface{}) *ConfigBuilder
func (b *ConfigBuilder) WithStaticFields(fields map[string]interface{}) *ConfigBuilder
func (b *ConfigBuilder) WithHandler(handler slog.Handler) *ConfigBuilder
func (b *ConfigBuilder) FromEnvironment() *ConfigBuilder
func (b *ConfigBuilder) Build() *Config
```

## Level System

### Level Constants

```go
const (
    TraceLevel    Level = iota // Most verbose
    DebugLevel
    InfoLevel                  // Default level
    WarnLevel
    ErrorLevel
    CriticalLevel              // Least verbose
)
```

### Level Operations

```go
// Convert string to level
func ParseLevel(level string) (Level, bool)

// Convert level to string
func (l Level) String() string

// Level comparison (levels are ordered)
if level >= WarnLevel {
    // Handle important messages
}
```

## Context Support

### Context Key Management

```go
type contextKey string

const (
    TraceIDKey       contextKey = "trace_id"
    RequestIDKey     contextKey = "request_id"
    CorrelationKey   contextKey = "correlation_id"
)
```

### Context Functions

```go
// Add values to context
func WithTraceID(ctx context.Context, traceID string) context.Context
func WithRequestID(ctx context.Context, requestID string) context.Context
func WithCorrelationID(ctx context.Context, correlationID string) context.Context

// Retrieve values from context
func GetTraceID(ctx context.Context) (string, bool)
func GetRequestID(ctx context.Context) (string, bool)
func GetCorrelationID(ctx context.Context) (string, bool)

// Utilities
func NewTraceID() string
func NewContextWithTrace() context.Context
```

## Handler System

### Handler Interfaces

```go
// Unified handler interface combining factory and middleware capabilities
type UnifiedHandlerInterface interface {
    slog.Handler
    
    // Factory methods
    Create(config interface{}) (slog.Handler, error)
    Name() string
    ConfigType() interface{}
    
    // Middleware support
    WithMiddleware(middleware ...HandlerMiddleware) slog.Handler
    
    // Lifecycle management
    Close() error
}

// Middleware interface
type HandlerMiddleware interface {
    Handle(ctx context.Context, record slog.Record, next HandlerFunc) error
}

type HandlerFunc func(context.Context, slog.Record) error
```

### Handler Composition

```go
// Multi-output handler
func NewMultiHandler(handlers ...slog.Handler) *MultiHandler

// Conditional handler
func NewConditionalHandler(handler slog.Handler, condition func(context.Context, slog.Record) bool) *ConditionalHandler

// Buffered handler
func NewBufferedHandler(handler slog.Handler, maxSize int) *BufferedHandler

// Async handler
func NewAsyncHandler(handler slog.Handler, queueSize int) *AsyncHandler

// Rotating handler
func NewRotatingHandler(handlers ...slog.Handler) *RotatingHandler

// Handler builder
func NewHandlerBuilder(handler slog.Handler) *HandlerBuilder
```

### Handler Builder Methods

```go
func (b *HandlerBuilder) WithTimestamp() *HandlerBuilder
func (b *HandlerBuilder) WithTraceContext() *HandlerBuilder
func (b *HandlerBuilder) WithStaticFields(fields map[string]interface{}) *HandlerBuilder
func (b *HandlerBuilder) WithLevelFilter(minLevel slog.Level) *HandlerBuilder
func (b *HandlerBuilder) WithSampling(rate int) *HandlerBuilder
func (b *HandlerBuilder) WithMiddleware(middleware ...HandlerMiddleware) *HandlerBuilder
func (b *HandlerBuilder) Build() slog.Handler
```

## Middleware

### Built-in Middleware

```go
// Timestamp middleware
func TimestampMiddleware() HandlerMiddleware

// Context extraction middleware
func ContextExtractorMiddleware(extractor ContextExtractor) HandlerMiddleware

// Level filtering middleware
func LevelFilterMiddleware(minLevel slog.Level) HandlerMiddleware

// Sampling middleware
func SamplingMiddleware(sampleRate int) HandlerMiddleware

// Static fields middleware
func StaticFieldsMiddleware(fields map[string]interface{}) HandlerMiddleware

// Redaction middleware
func RedactionMiddleware(redactor RedactorInterface) HandlerMiddleware

// Metrics middleware
func MetricsMiddleware(callback func(slog.Level)) HandlerMiddleware
```

### Context Extractors

```go
// Built-in extractors
func TraceContextExtractor() ContextExtractor
func StringContextExtractor(contextKey string, fieldName string) ContextExtractor
func IntContextExtractor(contextKey string, fieldName string) ContextExtractor
func Int64ContextExtractor(contextKey string, fieldName string) ContextExtractor
func BoolContextExtractor(contextKey string, fieldName string) ContextExtractor
func CustomContextExtractor(contextKey string, fieldName string) ContextExtractor

// Composite extractor
func NewCompositeContextExtractor(extractors ...ContextExtractor) *CompositeContextExtractor
```

## Async Processing

### AsyncWorker (Generic)

```go
type AsyncWorker[T any] struct {
    // Internal fields
}

type AsyncWorkerConfig[T any] struct {
    QueueSize  int
    Processor  func(T) error
    OnShutdown func() error
}

func NewAsyncWorker[T any](config AsyncWorkerConfig[T]) *AsyncWorker[T]

func (w *AsyncWorker[T]) Submit(item T) bool
func (w *AsyncWorker[T]) SubmitBlocking(item T) bool
func (w *AsyncWorker[T]) Stop() error
func (w *AsyncWorker[T]) IsClosed() bool
func (w *AsyncWorker[T]) QueueSize() int
func (w *AsyncWorker[T]) QueueCapacity() int
```

### AsyncOutput

```go
func NewAsyncOutput(output Output, queueSize int) *AsyncOutput

func (ao *AsyncOutput) Write(data []byte) error
func (ao *AsyncOutput) Stop() error
func (ao *AsyncOutput) Close() error
```

### AsyncHandler

```go
func NewAsyncHandler(handler slog.Handler, queueSize int) *AsyncHandler

func (h *AsyncHandler) Handle(ctx context.Context, record slog.Record) error
func (h *AsyncHandler) Close()
func (h *AsyncHandler) WithAttrs(attrs []slog.Attr) slog.Handler
func (h *AsyncHandler) WithGroup(name string) slog.Handler
```

## Utilities

### HTTP Utilities

```go
// Middleware functions
func TracingMiddleware(logger Logger) func(http.Handler) http.Handler
func RequestLogger(logger Logger, headers ...string) func(http.Handler) http.Handler

// HTTP logging helpers
func LogHTTPRequest(logger Logger, r *http.Request, headers []string)
func LogHTTPResponse(logger Logger, statusCode int, url string)

// URL utilities
func RedactedURL(url string) string
func RequestHeaders(r *http.Request, headersToPrint []string) string
```

### Redaction

```go
// Redactor interfaces
type RedactorInterface interface {
    Redact(input string) string
}

type RedactorChainInterface interface {
    RedactorInterface
    AddRedactor(redactor RedactorInterface) RedactorChainInterface
}

// Built-in redactors
func NewRegexRedactor(pattern *regexp.Regexp, replacement string) *RegexRedactor
func NewRedactorChain(patterns ...*regexp.Regexp) *RedactorChain

// Predefined redactors
func RedactAPIKeys(input string) string
```

### Output Types

```go
// Output interface
type Output interface {
    Write(data []byte) error
    Close() error
}

// Extended interfaces
type BufferedOutputInterface interface {
    Output
    Flush() error
}

type AsyncOutputInterface interface {
    Output
    Stop() error
}

// Built-in outputs
func NewFileOutput(filename string) *FileOutput
func NewRotatingFileOutput(pattern string, maxSize int64, maxAge time.Duration) *RotatingFileOutput
func NewConsoleOutput() *ConsoleOutput
func NewAsyncOutput(output Output, queueSize int) *AsyncOutput
```

### Environment Support

```go
// Environment variable names
const (
    EnvLogLevel  = "LOG_LEVEL"   // TRACE, DEBUG, INFO, WARN, ERROR, CRITICAL
    EnvLogFormat = "LOG_FORMAT"  // text, json
)

// Utility functions
func MustGetEnv(key string) string // Panics if not found
```

### Registry System

```go
// Handler registry for dynamic handler creation
func RegisterHandler(name string, factory NamedHandlerFactory)
func GetHandler(name string) (NamedHandlerFactory, bool)
func ListHandlers() []string
func CreateHandler(name string, config interface{}) (slog.Handler, error)
```

## Usage Examples

### Basic Usage

```go
// Simple logger creation
logger := logging.NewWithLevel(logging.InfoLevel)
logger.Info("Application started")

// With fields
serviceLogger := logger.WithFields(map[string]interface{}{
    "service": "api-gateway",
    "version": "2.1.0",
})
serviceLogger.Info("Service initialized")
```

### Fluent Interface

```go
logger.Fluent().Error().
    Str("operation", "payment").
    Int("user_id", 12345).
    Err(err).
    Msg("Payment processing failed")
```

### Context Usage

```go
ctx := logging.WithRequestID(context.Background(), "req-123")
logger.InfoContext(ctx, "Processing request")
```

### Configuration

```go
config := logging.NewLoggerConfig().
    WithCore(
        logging.NewCoreConfig().
            WithLevel(logging.DebugLevel).
            WithStaticField("service", "payment-api").
            Build(),
    ).
    WithFormatter(
        logging.NewFormatterConfig().
            WithJSONFormat().
            IncludeFile(true).
            Build(),
    ).
    Build()

logger := logging.NewWithLoggerConfig(config)
```

### Custom Handler

```go
customHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
    Level: slog.LevelDebug,
    AddSource: true,
})

logger := logging.NewWithHandler(customHandler)
logger.Debug("Debug message with source location")
```

This API reference covers all public interfaces and functions available in the go-logging library. The unified architecture ensures that all functionality is accessible through clean, consistent interfaces while maintaining backward compatibility.
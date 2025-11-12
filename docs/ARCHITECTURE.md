# Architecture Guide

This document provides a detailed overview of the go-logging library's unified architecture, design decisions, and implementation patterns.

## Table of Contents

- [Overview](#overview)
- [Core Components](#core-components)
- [Unified Logger System](#unified-logger-system)
- [Configuration Architecture](#configuration-architecture)
- [Async Processing](#async-processing)
- [Handler System](#handler-system)
- [Level Dispatch](#level-dispatch)
- [Context Management](#context-management)
- [Design Patterns](#design-patterns)
- [Performance Considerations](#performance-considerations)

## Overview

The go-logging library has been architected as a **unified logging system** that seamlessly bridges traditional Go logging with modern structured logging via slog. The architecture prioritizes:

- **Unified Interface**: Single Logger interface with all capabilities built-in
- **Backend Flexibility**: Transparent switching between standard and slog backends
- **Performance**: Optimized for high-throughput scenarios with async processing
- **Extensibility**: Clean interfaces for custom handlers and middleware
- **Backward Compatibility**: Legacy APIs continue to work unchanged

## Core Components

### 1. Unified Logger (`unified_logger.go`)

The heart of the system is the `unifiedLogger` struct, which implements the complete Logger interface:

```go
type unifiedLogger struct {
    mu            sync.RWMutex
    config        *LoggerConfig
    fields        map[string]interface{}
    textLoggers   map[Level]*log.Logger
    slogLogger    *slog.Logger
    discard       *log.Logger
    redactorChain RedactorChainInterface
}
```

**Key Features:**
- **Dual Backend Support**: Switches between standard Go logging and slog based on configuration
- **Thread-Safe**: All operations are protected with RWMutex for concurrent access
- **Field Management**: Immutable field attachment with efficient copying
- **Level-Aware Routing**: Automatically routes calls to appropriate backend

### 2. Logger Interface (`logger.go`)

The unified Logger interface consolidates all logging capabilities:

```go
type Logger interface {
    // Core methods
    Log(level Level, msg string, args ...interface{})
    LogContext(ctx context.Context, level Level, msg string, args ...interface{})
    
    // Field attachment (immutable)
    WithField(key string, value interface{}) Logger
    WithFields(fields map[string]interface{}) Logger
    
    // Level management
    IsLevelEnabled(level Level) bool
    SetLevel(level Level)
    GetLevel() Level
    
    // Level-specific methods (all included)
    Trace(msg string, args ...interface{})
    Debug(msg string, args ...interface{})
    Info(msg string, args ...interface{})
    Warn(msg string, args ...interface{})
    Error(msg string, args ...interface{})
    Critical(msg string, args ...interface{})
    
    // Context variants
    TraceContext(ctx context.Context, msg string, args ...interface{})
    DebugContext(ctx context.Context, msg string, args ...interface{})
    InfoContext(ctx context.Context, msg string, args ...interface{})
    WarnContext(ctx context.Context, msg string, args ...interface{})
    ErrorContext(ctx context.Context, msg string, args ...interface{})
    CriticalContext(ctx context.Context, msg string, args ...interface{})
    
    // Fluent interface
    Fluent() FluentLogger
}
```

**Design Benefits:**
- **No Type Assertions**: All methods available on every logger instance
- **Consistent API**: Same interface regardless of backend
- **Context First-Class**: Context support built into every level
- **Fluent Always Available**: No need to check for FluentCapable interface

## Unified Logger System

### Backend Selection Logic

The unified logger automatically selects the appropriate backend:

```go
func (ul *unifiedLogger) Log(ctx context.Context, level Level, msg string, args ...interface{}) {
    if !ul.isLevelEnabledInternal(level) {
        return
    }

    message := fmt.Sprintf(msg, args...)
    message = ul.redactorChain.Redact(message)

    if ul.config.UseSlog {
        ul.logSlog(ctx, level, message)
    } else if ul.config.Formatter.Format == JSONFormat {
        ul.logJSON(level, message, ctx)
    } else {
        ul.logText(level, message)
    }
}
```

### Level Checking with Backend Delegation

Level checking respects the backend's capabilities:

```go
func (ul *unifiedLogger) isLevelEnabledInternal(level Level) bool {
    // When using slog, delegate to the slog handler for level checking
    if ul.config.UseSlog && ul.slogLogger != nil {
        return ul.slogLogger.Enabled(context.Background(), ul.levelToSlog(level))
    }
    // For standard logging, use config level
    return level >= ul.config.Core.Level
}
```

This ensures that when using slog with custom handlers, the handler's level configuration is respected.

## Configuration Architecture

### Structured Configuration System

The new configuration system separates concerns into logical components:

```go
type LoggerConfig struct {
    Core      *CoreConfig      // Level, static fields
    Formatter *FormatterConfig // Format, redaction, file inclusion
    Output    *OutputConfig    // Writer configuration
    Handler   slog.Handler     // Optional slog handler
    UseSlog   bool            // Backend selection
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

### Builder Pattern with Composition

The configuration builders support method chaining:

```go
config := logging.NewLoggerConfig().
    WithCore(
        logging.NewCoreConfig().
            WithLevel(logging.DebugLevel).
            WithStaticField("service", "api").
            Build(),
    ).
    WithFormatter(
        logging.NewFormatterConfig().
            WithJSONFormat().
            IncludeFile(true).
            Build(),
    ).
    WithOutput(
        logging.NewOutputConfig().
            WithWriter(os.Stdout).
            Build(),
    ).
    Build()
```

### Backward Compatibility

The legacy `Config` type is preserved and delegates to the new system:

```go
// Old API still works
config := logging.NewConfig().
    WithLevel(logging.InfoLevel).
    WithJSONFormat().
    Build()

logger := logging.NewStandardLogger(config, redactorChain)
```

## Async Processing

### Generic Async Worker Pattern

The `AsyncWorker[T]` provides a reusable async processing pattern:

```go
type AsyncWorker[T any] struct {
    queue      chan T
    done       chan struct{}
    wg         sync.WaitGroup
    closed     bool
    mu         sync.Mutex
    processor  func(T) error
    onShutdown func() error
}
```

### Usage in Components

Both `AsyncOutput` and `AsyncHandler` use the generic worker:

```go
// AsyncOutput for high-throughput file writing
ao.worker = NewAsyncWorker(AsyncWorkerConfig[[]byte]{
    QueueSize: queueSize,
    Processor: func(data []byte) error {
        return ao.output.Write(data)
    },
})

// AsyncHandler for non-blocking slog handling
ah.worker = NewAsyncWorker(AsyncWorkerConfig[slog.Record]{
    QueueSize: queueSize,
    Processor: func(record slog.Record) error {
        return ah.handler.Handle(context.Background(), record)
    },
})
```

## Handler System

### Unified Handler Interface

The `UnifiedHandlerInterface` consolidates handler operations:

```go
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
```

### Handler Composition

The `HandlerCompositor` enables sophisticated handler combinations:

```go
compositor := logging.NewHandlerCompositor()

// Multi-output handler
multiHandler := compositor.
    Add(textHandler).
    Add(jsonHandler).
    Add(remoteHandler).
    Multi()

// Middleware chain
chainedHandler := compositor.
    Add(baseHandler).
    Chain(
        logging.TimestampMiddleware(),
        logging.ContextExtractorMiddleware(),
        logging.LevelFilterMiddleware(slog.LevelWarn),
    )
```

## Level Dispatch

### Unified Dispatch System

The `LevelDispatcher` centralizes level method routing:

```go
type LevelDispatcher struct {
    logger Logger
}

func (d *LevelDispatcher) DispatchInfo(msg string, args ...interface{}) {
    d.logger.Log(InfoLevel, msg, args...)
}
```

### Embedded Level Methods

The `LoggerLevelMethods` struct provides default implementations:

```go
type LoggerLevelMethods struct {
    dispatcher *LevelDispatcher
}

func (l *LoggerLevelMethods) InitLevelMethods(coreLogger Logger) {
    l.dispatcher = NewLevelDispatcher(coreLogger)
}

func (l *LoggerLevelMethods) Info(msg string, args ...interface{}) {
    l.dispatcher.DispatchInfo(msg, args...)
}
```

This pattern eliminates code duplication across logger implementations.

## Context Management

### Context Field Extraction

Context fields are extracted using the actual context key functions:

```go
func (ul *unifiedLogger) addContextFields(entry map[string]interface{}, ctx context.Context) {
    if requestID, ok := GetRequestID(ctx); ok && requestID != "" {
        entry["request_id"] = requestID
    }
    if traceID, ok := GetTraceID(ctx); ok && traceID != "" {
        entry["trace_id"] = traceID
    }
    if correlationID, ok := GetCorrelationID(ctx); ok && correlationID != "" {
        entry["correlation_id"] = correlationID
    }
}
```

### Context Key Management

Context keys are defined with proper typing:

```go
type contextKey string

const (
    TraceIDKey       contextKey = "trace_id"
    RequestIDKey     contextKey = "request_id" 
    CorrelationKey   contextKey = "correlation_id"
)

func WithRequestID(ctx context.Context, requestID string) context.Context {
    return context.WithValue(ctx, RequestIDKey, requestID)
}

func GetRequestID(ctx context.Context) (string, bool) {
    if val := ctx.Value(RequestIDKey); val != nil {
        if requestID, ok := val.(string); ok {
            return requestID, true
        }
    }
    return "", false
}
```

## Design Patterns

### 1. Unified Interface Pattern

Instead of multiple interfaces (`LevelLogger`, `ContextLogger`, `FluentCapable`), we provide a single comprehensive interface. This eliminates type assertions and provides a consistent API.

### 2. Backend Abstraction Pattern

The unified logger abstracts backend differences while preserving each backend's strengths:
- Standard logging for simple text output
- Slog for structured logging with custom handlers

### 3. Immutable Field Pattern

Field attachment creates new logger instances rather than modifying existing ones:

```go
func (ul *unifiedLogger) WithField(key string, value interface{}) Logger {
    ul.mu.RLock()
    newFields := make(map[string]interface{}, len(ul.fields)+1)
    for k, v := range ul.fields {
        newFields[k] = v
    }
    ul.mu.RUnlock()
    
    newFields[key] = value
    
    return &unifiedLogger{
        config:        ul.config,
        fields:        newFields,
        textLoggers:   ul.textLoggers,
        slogLogger:    ul.slogLogger,
        discard:       ul.discard,
        redactorChain: ul.redactorChain,
    }
}
```

### 4. Generic Worker Pattern

The `AsyncWorker[T]` uses Go generics to provide type-safe async processing for different data types while sharing the same underlying implementation.

### 5. Builder Pattern with Validation

Configuration builders validate inputs and provide sensible defaults:

```go
func NewCoreConfig() *CoreConfigBuilder {
    return &CoreConfigBuilder{
        config: &CoreConfig{
            Level:        InfoLevel,  // Sensible default
            StaticFields: make(map[string]interface{}),
        },
    }
}
```

## Performance Considerations

### Memory Efficiency

- **Field Copying**: Only copies fields when creating new loggers, not on every log call
- **Level Checking**: Fast path for disabled levels avoids formatting and processing
- **Async Processing**: Non-blocking writes for high-throughput scenarios

### Concurrency

- **RWMutex Usage**: Allows multiple concurrent reads while protecting writes
- **Lock Granularity**: Minimal lock scope to reduce contention
- **Async Workers**: Independent goroutines for background processing

### Benchmarking

The library includes comprehensive benchmarks:

```bash
go test ./pkg/logging -bench=. -benchmem
```

Key benchmark results show:
- Minimal overhead for level checking
- Efficient field attachment
- Comparable performance to other popular logging libraries

## Extensibility Points

### Custom Handlers

Implement the `slog.Handler` interface:

```go
type MyCustomHandler struct {
    // Implementation
}

func (h *MyCustomHandler) Enabled(ctx context.Context, level slog.Level) bool { ... }
func (h *MyCustomHandler) Handle(ctx context.Context, record slog.Record) error { ... }
func (h *MyCustomHandler) WithAttrs(attrs []slog.Attr) slog.Handler { ... }
func (h *MyCustomHandler) WithGroup(name string) slog.Handler { ... }
```

### Custom Middleware

Implement the `HandlerMiddleware` interface:

```go
type MyMiddleware struct{}

func (m *MyMiddleware) Handle(ctx context.Context, record slog.Record, next HandlerFunc) error {
    // Pre-processing
    err := next(ctx, record)
    // Post-processing
    return err
}
```

### Custom Outputs

Implement the `Output` interface:

```go
type MyOutput struct{}

func (o *MyOutput) Write(data []byte) error { ... }
func (o *MyOutput) Close() error { ... }
```

---

This architecture enables a powerful, flexible, and performant logging system that grows with your application's needs while maintaining simplicity for basic use cases.
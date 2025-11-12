# Advanced Features Guide

This document covers the advanced features of the go-logging library, including async processing, handler composition, custom middleware, and performance optimization techniques with the unified architecture.

## Table of Contents

1. [Unified Architecture Features](#unified-architecture-features)
2. [Async Processing](#async-processing)
3. [Context Value Extraction](#context-value-extraction)
4. [Handler Middleware](#handler-middleware)
5. [Handler Composition](#handler-composition)
6. [Performance Optimization](#performance-optimization)
7. [Custom Handlers](#custom-handlers)
8. [Performance Benchmarks](#performance-benchmarks)

## Unified Architecture Features

### Consolidated Logger Interface

The unified architecture eliminates the need for type assertions by providing all logging capabilities in a single interface:

```go
// Before: Required type assertions
var logger logging.Logger = getLogger()
if ll, ok := logger.(logging.LevelLogger); ok {
    ll.Info("Using level logger")
}

// After: Direct method calls
var logger logging.Logger = getLogger()
logger.Info("Direct access to all methods")
logger.InfoContext(ctx, "Context support built-in")
logger.Fluent().Info().Msg("Fluent interface always available")
logger.SetLevel(logging.DebugLevel) // Configuration methods included
```

### Automatic Backend Selection

The unified logger transparently switches between standard and slog backends:

```go
// Creates logger that automatically uses slog when beneficial
logger := logging.NewWithLevel(logging.InfoLevel)

// Or explicitly use slog backend
slogLogger := logging.NewSlogJSONLogger(logging.DebugLevel)

// Or use custom slog handler
customHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
    Level: slog.LevelDebug,
})
handlerLogger := logging.NewWithHandler(customHandler)

// All have the same unified interface
logger.Info("Standard backend")
slogLogger.Info("Slog backend") 
handlerLogger.Info("Custom handler backend")
```

### Level Dispatch System

Unified level method dispatch eliminates code duplication:

```go
// All level methods automatically delegate to core Log/LogContext
logger.Trace("Trace message")       // -> logger.Log(TraceLevel, msg)
logger.Debug("Debug message")       // -> logger.Log(DebugLevel, msg)
logger.InfoContext(ctx, "Info")     // -> logger.LogContext(ctx, InfoLevel, msg)
logger.ErrorContext(ctx, "Error")   // -> logger.LogContext(ctx, ErrorLevel, msg)
```

## Async Processing

### Generic AsyncWorker Pattern

The new `AsyncWorker[T]` provides type-safe async processing for any data type:

```go
// Create async worker for processing log data
worker := logging.NewAsyncWorker(logging.AsyncWorkerConfig[[]byte]{
    QueueSize: 1000,
    Processor: func(data []byte) error {
        return writeToRemoteEndpoint(data)
    },
    OnShutdown: func() error {
        return flushRemoteConnection()
    },
})

// Submit data for async processing
if !worker.Submit(logData) {
    // Queue full, handle fallback
    return writeDirectly(logData)
}

// Graceful shutdown
defer worker.Stop()
```

### AsyncOutput for High-Throughput Logging

```go
// Create async file output with buffering
fileOutput := &logging.FileOutput{
    Filename: "high-volume.log",
}

asyncOutput := logging.NewAsyncOutput(fileOutput, 5000) // 5000 item buffer
defer asyncOutput.Close()

// Use with any logger type
config := logging.NewLoggerConfig().
    WithOutput(
        logging.NewOutputConfig().
            WithWriter(asyncOutput).
            Build(),
    ).
    Build()

logger := logging.NewWithLoggerConfig(config)

// High-throughput logging - non-blocking writes
for i := 0; i < 100000; i++ {
    logger.Info("High volume log entry %d", i)
}

// Ensure all logs are written before exit
asyncOutput.Stop()
```

### AsyncHandler Integration

```go
// Create async slog handler wrapper using generic worker
baseHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
    Level: slog.LevelInfo,
})

asyncHandler := logging.NewAsyncHandler(baseHandler, 2000) // 2000 record buffer
defer asyncHandler.Close()

logger := logging.NewWithHandler(asyncHandler)

// Non-blocking slog operations with full unified interface
logger.Info("This won't block the caller")
logger.Fluent().Error().Str("service", "payment").Msg("Error occurred")
```

## Context Value Extraction

Extract values from context and add them as structured log fields.

### Built-in Extractors

#### TraceContextExtractor
Extracts trace, request, and correlation IDs:
```go
extractor := logging.TraceContextExtractor()
ctx := logging.WithTraceID(context.Background(), "trace-123")
attrs := extractor.Extract(ctx)
```

#### Type-Specific Extractors
```go
// String values
userExtractor := logging.StringContextExtractor("user_key", "username")

// Integer values
ageExtractor := logging.IntContextExtractor("age_key", "age")

// Int64 values
idExtractor := logging.Int64ContextExtractor("id_key", "user_id")

// Boolean values
flagExtractor := logging.BoolContextExtractor("flag_key", "is_admin")

// Any type
customExtractor := logging.CustomContextExtractor("data_key", "custom_data")
```

### Composite Extractor
Combine multiple extractors:
```go
composite := logging.NewCompositeContextExtractor(
    logging.TraceContextExtractor(),
    logging.StringContextExtractor("user", "username"),
    logging.IntContextExtractor("age", "user_age"),
)

attrs := composite.Extract(ctx)
```

## Handler Middleware

Chain middleware to modify log records before they're written.

### Available Middlewares

#### TimestampMiddleware
Add/update timestamps:
```go
middleware := logging.TimestampMiddleware()
```

#### ContextExtractorMiddleware
Extract context values:
```go
extractor := logging.TraceContextExtractor()
middleware := logging.ContextExtractorMiddleware(extractor)
```

#### LevelFilterMiddleware
Filter by minimum level:
```go
middleware := logging.LevelFilterMiddleware(slog.LevelWarn)
```

#### SamplingMiddleware
Sample logs at specified rate:
```go
middleware := logging.SamplingMiddleware(10) // Log every 10th message
```

#### StaticFieldsMiddleware
Add static fields to all logs:
```go
fields := map[string]interface{}{
    "service": "my-service",
    "version": "1.0.0",
}
middleware := logging.StaticFieldsMiddleware(fields)
```

#### RedactionMiddleware
Redact sensitive data:
```go
redactor := logging.NewRegexRedactor(
    regexp.MustCompile(`password=\w+`),
    "password=***",
)
middleware := logging.RedactionMiddleware(redactor)
```

#### MetricsMiddleware
Record logging metrics:
```go
middleware := logging.MetricsMiddleware(func(level slog.Level) {
    metrics.RecordLog(level.String())
})
```

### Using Middleware

```go
handler := slog.NewJSONHandler(os.Stdout, nil)
middlewareHandler := logging.NewMiddlewareHandler(
    handler,
    logging.TimestampMiddleware(),
    logging.StaticFieldsMiddleware(fields),
    logging.ContextExtractorMiddleware(extractor),
)

logger := logging.NewWithHandler(middlewareHandler)
```

### Custom Middleware

Implement the `HandlerMiddleware` interface:
```go
type CustomMiddleware struct{}

func (m *CustomMiddleware) Handle(ctx context.Context, record slog.Record, next logging.HandlerFunc) error {
    // Modify record before passing to next middleware
    return next(ctx, record)
}
```

## Handler Composition

Compose multiple handlers for advanced logging patterns.

### MultiHandler
Log to multiple destinations:
```go
stdoutHandler := slog.NewTextHandler(os.Stdout, nil)
fileHandler := slog.NewJSONHandler(file, nil)

multiHandler := logging.NewMultiHandler(stdoutHandler, fileHandler)
logger := logging.NewWithHandler(multiHandler)
```

### ConditionalHandler
Conditional logging based on record properties:
```go
handler := logging.NewConditionalHandler(baseHandler, func(ctx context.Context, record slog.Record) bool {
    return record.Level >= slog.LevelWarn
})
```

### BufferedHandler
Buffer logs for batch processing:
```go
bufferedHandler := logging.NewBufferedHandler(baseHandler, 100)
defer bufferedHandler.Flush(ctx)
```

### AsyncHandler
Non-blocking asynchronous logging:
```go
asyncHandler := logging.NewAsyncHandler(baseHandler, 1000)
defer asyncHandler.Close()
```

### RotatingHandler
Round-robin across multiple handlers:
```go
rotatingHandler := logging.NewRotatingHandler(
    handler1,
    handler2,
    handler3,
)
```

### HandlerBuilder
Fluent interface for building complex handlers:
```go
handler := logging.NewHandlerBuilder(baseHandler).
    WithTimestamp().
    WithTraceContext().
    WithStaticFields(map[string]interface{}{
        "service": "my-service",
    }).
    WithLevelFilter(slog.LevelInfo).
    WithSampling(10).
    Build()
```

## Performance Optimization

### Memory-Efficient Field Usage

The unified architecture optimizes field management through immutable patterns:

```go
type OptimizedService struct {
    baseLogger   logging.Logger
    serviceInfo  map[string]interface{}
}

func NewOptimizedService(logger logging.Logger) *OptimizedService {
    // Pre-create static fields once - works with any logger type
    serviceInfo := map[string]interface{}{
        "service":     "payment-processor",
        "version":     "2.1.0",
        "datacenter":  "us-east-1",
        "environment": "production",
    }
    
    return &OptimizedService{
        baseLogger:  logger.WithFields(serviceInfo),
        serviceInfo: serviceInfo,
    }
}

func (s *OptimizedService) ProcessPayment(userID int, amount float64) {
    // Efficient: Reuse base logger with static fields
    userLogger := s.baseLogger.WithField("user_id", userID)
    
    // All methods available without type assertions
    userLogger.Info("Processing payment of $%.2f", amount)
    userLogger.Debug("Validating payment details")
    userLogger.InfoContext(ctx, "Payment completed successfully")
}
```

### Level Checking Optimization

The unified logger provides optimized level checking that delegates to backend handlers when appropriate:

```go
func expensiveLogging(logger logging.Logger, data *ComplexData) {
    // Optimized level checking works with any backend
    if logger.IsLevelEnabled(logging.DebugLevel) {
        expensiveMetrics := computeExpensiveMetrics(data) // Expensive operation
        
        // Direct access to fluent interface
        logger.Fluent().Debug().
            Field("metrics", expensiveMetrics).
            Int("data_size", len(data.Items)).
            Msg("Detailed debug information")
    }
    
    // Always efficient - level checking is built-in
    logger.Info("Operation completed for %d items", len(data.Items))
}
```

### Async Pattern Performance

```go
// High-performance async logging with the generic worker pattern
func setupHighPerformanceLlogging() logging.Logger {
    // Create multiple async outputs for different log types
    infoOutput := logging.NewAsyncOutput(infoFileWriter, 1000)
    errorOutput := logging.NewAsyncOutput(errorFileWriter, 500)
    
    // Routing handler using unified interface
    router := &RoutingHandler{
        infoHandler:  logging.NewWithHandler(slog.NewJSONHandler(infoOutput, nil)),
        errorHandler: logging.NewWithHandler(slog.NewJSONHandler(errorOutput, nil)),
    }
    
    logger := logging.NewWithHandler(router)
    
    // All methods available - optimal performance
    return logger
}
```

### Performance Benchmarks

#### Running Benchmarks

```bash
go test ./pkg/logging -bench=. -benchmem
```

#### Benchmark Results with Unified Architecture

```
BenchmarkUnifiedLogger_Info               3051234    393.2 ns/op   189 B/op   1 allocs/op
BenchmarkUnifiedLogger_InfoWithFields     2847592    422.1 ns/op   245 B/op   2 allocs/op
BenchmarkUnifiedLogger_FluentInterface    2756483    435.8 ns/op   198 B/op   2 allocs/op
BenchmarkUnifiedLogger_LevelCheck         50000000    23.4 ns/op     0 B/op   0 allocs/op
BenchmarkAsyncOutput_HighThroughput       5847291    205.3 ns/op    48 B/op   1 allocs/op
```

#### Key Performance Improvements

- **Unified Interface**: Eliminates type assertion overhead
- **Optimized Level Checking**: Fast path for disabled levels (0 allocations)
- **Backend Delegation**: Slog level checking when using slog handlers
- **Async Workers**: ~50% faster async processing with generic pattern
- **Memory Efficiency**: Reduced allocations through better field management

### Performance Tips

1. **Use Level Checks for Expensive Operations**:
   ```go
   if logger.IsLevelEnabled(logging.DebugLevel) {
       logger.Debug("Expensive: %v", computeExpensiveData())
   }
   ```

2. **Pre-create Loggers with Static Fields**:
   ```go
   serviceLogger := logger.WithFields(staticServiceInfo)
   // Reuse serviceLogger instead of creating fields each time
   ```

3. **Use Async Processing for High-Throughput**:
   ```go
   asyncOutput := logging.NewAsyncOutput(fileWriter, 1000)
   logger := logging.NewWithLoggerConfig(configWithAsyncOutput)
   ```

4. **Leverage Fluent Interface for Complex Logs**:
   ```go
   // More efficient than multiple WithField calls
   logger.Fluent().Error().
       Str("service", service).
       Int("user_id", userID).
       Err(err).
       Msg("Operation failed")
   ```

5. **Use Appropriate Backends**:
   ```go
   // Use slog backend for structured logging
   logger := logging.NewSlogJSONLogger(logging.InfoLevel)
   
   // Use standard backend for simple text logging
   logger := logging.NewWithLevel(logging.InfoLevel)
   ```

## Custom Handlers

### Creating Custom Handlers

Implement the `slog.Handler` interface:

```go
type CustomHandler struct {
    output io.Writer
}

func (h *CustomHandler) Enabled(ctx context.Context, level slog.Level) bool {
    return true
}

func (h *CustomHandler) Handle(ctx context.Context, record slog.Record) error {
    // Custom handling logic
    _, err := fmt.Fprintf(h.output, "%s: %s\n", record.Level, record.Message)
    return err
}

func (h *CustomHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
    // Return new handler with attributes
    return h
}

func (h *CustomHandler) WithGroup(name string) slog.Handler {
    // Return new handler with group
    return h
}
```

### Using Custom Handlers

```go
customHandler := &CustomHandler{output: os.Stdout}
logger := logging.NewWithHandler(customHandler)
```

## Design Principles

All advanced features follow SOLID principles:

### Single Responsibility
- Each middleware has one purpose
- Each handler type has a specific use case
- Each extractor handles one type of value

### Open/Closed
- Extensible through interfaces
- New middlewares can be added without modifying existing code
- Handlers can be composed in infinite ways

### Liskov Substitution
- All handlers implement `slog.Handler`
- All middlewares implement `HandlerMiddleware`
- All extractors implement `ContextExtractor`

### Interface Segregation
- Clean, minimal interfaces
- No unnecessary methods
- Easy to implement and test

### Dependency Injection
- Configuration via builders
- Dependencies injected via constructors
- No global state

## Best Practices

1. **Compose Don't Inherit**: Use handler composition over custom implementations
2. **Use Builders**: Leverage HandlerBuilder for complex configurations
3. **Extract Reusable Logic**: Create custom extractors for app-specific context values
4. **Test Middleware**: Each middleware should be independently testable
5. **Monitor Performance**: Use benchmarks to measure impact of middlewares
6. **Document Custom Handlers**: Clearly document behavior and requirements

## Examples

See the `examples/` directory for complete working examples:
- `examples/slog/` - Basic slog integration
- `examples/custom-handlers/` - Advanced handler composition

## Further Reading

- [Slog Integration Guide](SLOG_INTEGRATION.md)
- [Main README](../README.md)
- [Contributing Guidelines](../CONTRIBUTING.md)

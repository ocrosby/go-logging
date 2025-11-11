# Advanced Features Guide

This document covers the advanced features added to the go-logging package.

## Table of Contents

1. [Context Value Extraction](#context-value-extraction)
2. [Handler Middleware](#handler-middleware)
3. [Handler Composition](#handler-composition)
4. [Performance Benchmarks](#performance-benchmarks)
5. [Custom Handlers](#custom-handlers)

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

## Performance Benchmarks

### Running Benchmarks

```bash
go test ./pkg/logging -bench=. -benchmem
```

### Benchmark Results

Comparison between StandardLogger and SlogLogger:

```
BenchmarkStandardLogger_Info         2748112    432.2 ns/op   418 B/op   3 allocs/op
BenchmarkSlogLogger_Info             2835986    427.2 ns/op   189 B/op   1 allocs/op
```

Key findings:
- Slog logger is slightly faster
- Slog logger uses ~55% less memory
- Slog logger has fewer allocations (1 vs 3)

### Performance Tips

1. **Use Level Checks**: Avoid expensive operations when not needed
   ```go
   if logger.IsLevelEnabled(logging.DebugLevel) {
       logger.Debug("Expensive: %v", computeExpensiveData())
   }
   ```

2. **Buffer Handlers**: Reduce I/O with buffering
   ```go
   handler := logging.NewBufferedHandler(baseHandler, 100)
   ```

3. **Sampling**: Reduce log volume in high-throughput scenarios
   ```go
   middleware := logging.SamplingMiddleware(10)
   ```

4. **Async Logging**: Non-blocking for performance-critical paths
   ```go
   handler := logging.NewAsyncHandler(baseHandler, 1000)
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

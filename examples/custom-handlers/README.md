# Custom Handlers Example

This example demonstrates advanced handler customization features in the go-logging package.

## ⚠️ For Most Users - Start Simple!

You probably don't need these advanced features! Try the simple approach first:

```go
import "github.com/ocrosby/go-logging/pkg/logging"

// Simple text logging
logger := logging.NewSimple()
logger.Info("Hello world")

// Simple JSON logging  
logger := logging.NewEasyJSON()
logger.Info("Hello world")

// Progressive configuration
logger := logging.NewEasyBuilder().
    Debug().
    JSON().
    Field("service", "my-app").
    Build()
```

**Only use the advanced features below if the simple functions don't meet your needs.**

## Features Demonstrated

### 1. Multi-Handler
Log to multiple destinations simultaneously:
```go
multiHandler := logging.NewMultiHandler(stdoutHandler, fileHandler)
logger := logging.NewWithHandler(multiHandler)
```

### 2. Handler Builder
Fluent interface for building complex handlers:
```go
handler := logging.NewHandlerBuilder(baseHandler).
    WithTimestamp().
    WithTraceContext().
    WithStaticFields(fields).
    Build()
```

### 3. Conditional Handler
Filter logs based on custom conditions:
```go
conditionalHandler := logging.NewConditionalHandler(handler, func(ctx context.Context, record slog.Record) bool {
    return record.Level >= slog.LevelWarn
})
```

### 4. Buffered Handler
Buffer log entries for batch processing:
```go
bufferedHandler := logging.NewBufferedHandler(handler, 5)
// ... log messages ...
bufferedHandler.Flush(ctx)
```

### 5. Middleware Chain
Compose multiple middlewares:
```go
middlewareHandler := logging.NewMiddlewareHandler(
    handler,
    logging.TimestampMiddleware(),
    logging.StaticFieldsMiddleware(fields),
    logging.ContextExtractorMiddleware(extractor),
)
```

## Available Middlewares

- **TimestampMiddleware()** - Add timestamps to log records
- **ContextExtractorMiddleware()** - Extract context values
- **LevelFilterMiddleware()** - Filter by minimum level
- **SamplingMiddleware()** - Sample logs at specified rate
- **CallerMiddleware()** - Add caller information
- **StaticFieldsMiddleware()** - Add static fields
- **RedactionMiddleware()** - Redact sensitive data
- **MetricsMiddleware()** - Record logging metrics

## Handler Composition Utilities

- **MultiHandler** - Log to multiple handlers
- **ConditionalHandler** - Conditional logging
- **BufferedHandler** - Buffered batch logging
- **AsyncHandler** - Asynchronous logging
- **RotatingHandler** - Round-robin handler rotation

## Running the Example

```bash
cd examples/custom-handlers
go run main.go
```

## Design Principles

All handlers and middlewares follow SOLID principles:
- **Single Responsibility**: Each handler/middleware has one purpose
- **Open/Closed**: Extensible through composition
- **Liskov Substitution**: All implement slog.Handler interface
- **Interface Segregation**: Clean, minimal interfaces
- **Dependency Injection**: Configuration via builders

## Performance Considerations

- **Buffered Handler**: Reduces I/O operations
- **Async Handler**: Non-blocking logging
- **Sampling Middleware**: Reduces log volume
- **Conditional Handler**: Filters before processing

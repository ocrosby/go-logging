# go-logging

A configurable, production-ready logging library for Go with support for structured logging, request tracing, and fluent interfaces.

## Features

- **Multiple Log Levels**: TRACE, DEBUG, INFO, WARN, ERROR, CRITICAL
- **Multiple Output Formats**: Text and JSON
- **Fluent Interface**: Chain methods for expressive logging
- **Request Tracing**: Built-in support for trace IDs, request IDs, and correlation IDs
- **Structured Logging**: Add contextual fields to log entries
- **HTTP Middleware**: Automatic request tracing and logging
- **Sensitive Data Redaction**: Built-in patterns to redact API keys and other sensitive data
- **Thread-Safe**: Concurrent logging with proper synchronization
- **Environment Configuration**: Configure from environment variables
- **Dependency Injection**: Design follows SOLID principles for easy testing

## Installation

```bash
go get github.com/ocrosby/go-logging
```

## Quick Start

### Basic Logging

```go
package main

import "github.com/ocrosby/go-logging/pkg/logging"

func main() {
    logger := logging.NewWithLevel(logging.InfoLevel)
    
    logger.Info("Application started")
    logger.Warn("This is a warning")
    logger.Error("This is an error")
}
```

### Fluent Interface

```go
logger := logging.NewStandardLogger(
    logging.NewConfig().
        WithLevel(logging.DebugLevel).
        WithJSONFormat().
        Build(),
)

logger.Fluent().Info().
    Str("service", "api").
    Int("user_id", 12345).
    Msg("User logged in")
```

### Request Tracing

```go
ctx := logging.NewContextWithTrace()
ctx = logging.WithRequestID(ctx, "req-123")

logger.Fluent().Info().
    Ctx(ctx).
    Str("operation", "fetch_user").
    Msg("Processing request")
```

### HTTP Middleware

```go
func main() {
    logger := logging.NewJSONLogger(logging.InfoLevel)
    
    mux := http.NewServeMux()
    mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
        ctx := r.Context()
        
        logger.Fluent().Info().
            Ctx(ctx).
            Msg("Handling request")
        
        fmt.Fprintf(w, "Hello!")
    })
    
    handler := logging.TracingMiddleware(logger)(mux)
    http.ListenAndServe(":8080", handler)
}
```

## Configuration

### Builder Pattern

```go
config := logging.NewConfig().
    WithLevel(logging.DebugLevel).
    WithJSONFormat().
    WithOutput(os.Stdout).
    IncludeFile(true).
    IncludeTime(true).
    AddRedactPattern(`password=\w+`).
    WithStaticField("service", "my-app").
    WithStaticField("version", "1.0.0").
    Build()

logger := logging.NewStandardLogger(config)
```

### Environment Variables

```go
// Reads LOG_LEVEL and LOG_FORMAT from environment
logger := logging.NewFromEnvironment()
```

Supported environment variables:
- `LOG_LEVEL`: TRACE, DEBUG, INFO, WARN, ERROR, CRITICAL
- `LOG_FORMAT`: json (text is default)

## Usage Patterns

### Traditional Logging

```go
logger.Info("User %s logged in", username)
logger.Error("Failed to connect: %v", err)
```

### Fluent Interface with Fields

```go
logger.Fluent().Error().
    Err(err).
    Str("host", "db.example.com").
    Int("port", 5432).
    Msg("Database connection failed")
```

### Structured Context

```go
logger = logger.WithFields(map[string]interface{}{
    "service": "api",
    "version": "1.0.0",
    "env": "production",
})

logger.Info("Message with static fields")
```

### Request Tracing with Context

```go
ctx := logging.WithTraceID(r.Context(), "trace-123")
ctx = logging.WithRequestID(ctx, "req-456")

logger.InfoContext(ctx, "Processing request")
```

## HTTP Headers

The middleware automatically handles these headers:
- `X-Trace-ID`: Trace identifier (generated if not present)
- `X-Request-ID`: Request identifier
- `X-Correlation-ID`: Correlation identifier

## Sensitive Data Redaction

Built-in redaction for API keys:

```go
url := "https://api.example.com?apiKey=secret123456"
redacted := logging.RedactedURL(url)
// Output: https://api.example.com?apiKey=secret1...<REDACTED>
```

Add custom redaction patterns:

```go
config := logging.NewConfig().
    AddRedactPattern(`password=\w+`).
    AddRedactPattern(`token=\w+`).
    Build()
```

## Log Levels

```go
const (
    TraceLevel    // Most verbose
    DebugLevel
    InfoLevel     // Default
    WarnLevel
    ErrorLevel
    CriticalLevel // Least verbose
)
```

Check if a level is enabled:

```go
if logger.IsLevelEnabled(logging.DebugLevel) {
    // Perform expensive debug logging
}
```

## Testing

The logger interface makes it easy to mock in tests:

```go
type mockLogger struct {
    entries []string
}

func (m *mockLogger) Info(msg string, args ...interface{}) {
    m.entries = append(m.entries, fmt.Sprintf(msg, args...))
}
```

## Examples

See the `examples/` directory for complete examples:
- `examples/basic/` - Basic logging usage
- `examples/fluent/` - Fluent interface examples
- `examples/http-server/` - HTTP middleware and tracing

## Design Principles

This library follows SOLID principles:
- **Single Responsibility**: Each component has a focused purpose
- **Open/Closed**: Extensible through interfaces
- **Liskov Substitution**: Logger interface is easily mockable
- **Interface Segregation**: Clean, minimal interfaces
- **Dependency Injection**: Configuration via builder pattern

## License

MIT License - see LICENSE file for details

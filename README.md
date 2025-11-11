# go-logging

[![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.19-blue)](https://golang.org/dl/)
[![Go Report Card](https://goreportcard.com/badge/github.com/ocrosby/go-logging)](https://goreportcard.com/report/github.com/ocrosby/go-logging)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A configurable, production-ready logging library for Go with support for structured logging, request tracing, and fluent interfaces.

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Usage](#usage)
  - [Basic Logging](#basic-logging)
  - [Fluent Interface](#fluent-interface)
  - [Request Tracing](#request-tracing)
  - [HTTP Middleware](#http-middleware)
- [Configuration](#configuration)
- [API Reference](#api-reference)
- [Examples](#examples)
- [Testing](#testing)
- [Design Principles](#design-principles)
- [Contributing](#contributing)
- [Changelog](#changelog)
- [Roadmap](#roadmap)
- [Support](#support)
- [Acknowledgments](#acknowledgments)
- [License](#license)

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
- **Zero Dependencies**: Only uses standard library (except for UUID generation)

## Installation

```bash
go get github.com/ocrosby/go-logging
```

### Requirements

- Go 1.19 or higher

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

## Usage

### Configuration

#### Builder Pattern

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

#### Environment Variables

```go
// Reads LOG_LEVEL and LOG_FORMAT from environment
logger := logging.NewFromEnvironment()
```

Supported environment variables:
- `LOG_LEVEL`: TRACE, DEBUG, INFO, WARN, ERROR, CRITICAL
- `LOG_FORMAT`: json (text is default)

### Usage Patterns

#### Traditional Logging

```go
logger.Info("User %s logged in", username)
logger.Error("Failed to connect: %v", err)
```

#### Fluent Interface with Fields

```go
logger.Fluent().Error().
    Err(err).
    Str("host", "db.example.com").
    Int("port", 5432).
    Msg("Database connection failed")
```

#### Structured Context

```go
logger = logger.WithFields(map[string]interface{}{
    "service": "api",
    "version": "1.0.0",
    "env": "production",
})

logger.Info("Message with static fields")
```

#### Request Tracing with Context

```go
ctx := logging.WithTraceID(r.Context(), "trace-123")
ctx = logging.WithRequestID(ctx, "req-456")

logger.InfoContext(ctx, "Processing request")
```

### HTTP Headers

The middleware automatically handles these headers:
- `X-Trace-ID`: Trace identifier (generated if not present)
- `X-Request-ID`: Request identifier
- `X-Correlation-ID`: Correlation identifier

### Sensitive Data Redaction

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

## API Reference

### Log Levels

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

### Logger Interface

```go
type Logger interface {
    Trace(msg string, args ...interface{})
    Debug(msg string, args ...interface{})
    Info(msg string, args ...interface{})
    Warn(msg string, args ...interface{})
    Error(msg string, args ...interface{})
    Critical(msg string, args ...interface{})
    
    WithField(key string, value interface{}) Logger
    WithFields(fields map[string]interface{}) Logger
    
    Fluent() FluentLogger
    
    IsLevelEnabled(level Level) bool
    SetLevel(level Level)
    GetLevel() Level
}
```

### Fluent Interface

```go
logger.Fluent().Info().
    Str(key, value string).
    Int(key string, value int).
    Int64(key string, value int64).
    Bool(key string, value bool).
    Err(err error).
    TraceID(id string).
    Ctx(ctx context.Context).
    Field(key string, value interface{}).
    Fields(fields map[string]interface{}).
    Msg(msg string)
```

### Factory Functions

```go
// Create logger with specific level
logger := logging.NewWithLevel(logging.InfoLevel)

// Create logger from environment variables
logger := logging.NewFromEnvironment()

// Create JSON logger
logger := logging.NewJSONLogger(logging.InfoLevel)

// Create text logger
logger := logging.NewTextLogger(logging.InfoLevel)

// Create with custom configuration
config := logging.NewConfig().Build()
logger := logging.NewStandardLogger(config)
```

### Context Functions

```go
// Create new trace ID
traceID := logging.NewTraceID()

// Add IDs to context
ctx = logging.WithTraceID(ctx, traceID)
ctx = logging.WithRequestID(ctx, requestID)
ctx = logging.WithCorrelationID(ctx, correlationID)

// Retrieve IDs from context
traceID, ok := logging.GetTraceID(ctx)
requestID, ok := logging.GetRequestID(ctx)
correlationID, ok := logging.GetCorrelationID(ctx)

// Create context with new trace ID
ctx := logging.NewContextWithTrace()
```

### Middleware Functions

```go
// Add tracing middleware
handler := logging.TracingMiddleware(logger)(yourHandler)

// Add request logger middleware
handler := logging.RequestLogger(logger, "User-Agent", "X-Custom-Header")(yourHandler)
```

## Examples

See the `examples/` directory for complete working examples:
- [`examples/basic/`](examples/basic/) - Basic logging usage
- [`examples/fluent/`](examples/fluent/) - Fluent interface examples
- [`examples/http-server/`](examples/http-server/) - HTTP middleware and tracing

## Testing

Run the test suite:

```bash
go test ./pkg/logging/... -v
```

Run tests with coverage:

```bash
go test ./pkg/logging/... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Mocking in Tests

The logger interface makes it easy to mock in tests:

```go
type mockLogger struct {
    entries []string
}

func (m *mockLogger) Info(msg string, args ...interface{}) {
    m.entries = append(m.entries, fmt.Sprintf(msg, args...))
}

// Implement other Logger interface methods...
```

## Design Principles

This library follows SOLID principles:
- **Single Responsibility**: Each component has a focused purpose
- **Open/Closed**: Extensible through interfaces
- **Liskov Substitution**: Logger interface is easily mockable
- **Interface Segregation**: Clean, minimal interfaces
- **Dependency Injection**: Configuration via builder pattern

### Architecture

```
pkg/logging/
├── logger.go           # Core Logger interface
├── level.go            # Log level definitions
├── config.go           # Configuration with builder pattern
├── standard_logger.go  # Standard logger implementation
├── fluent.go           # Fluent interface implementation
├── factory.go          # Factory functions
├── trace.go            # Request tracing utilities
├── middleware.go       # HTTP middleware
├── redactor.go         # Sensitive data redaction
└── http.go             # HTTP logging utilities
```

## Contributing

Contributions are welcome! Please read our [Contributing Guidelines](CONTRIBUTING.md) before submitting a PR.

### Development Setup

1. Fork the repository
2. Clone your fork: `git clone https://github.com/yourusername/go-logging.git`
3. Create a branch: `git checkout -b feature/my-feature`
4. Make your changes
5. Run tests: `go test ./...`
6. Commit your changes: `git commit -am 'Add new feature'`
7. Push to the branch: `git push origin feature/my-feature`
8. Submit a pull request

### Code Style

- Follow standard Go conventions and idioms
- Run `go fmt` before committing
- Ensure all tests pass
- Add tests for new functionality
- Update documentation as needed

## Changelog

See [CHANGELOG.md](CHANGELOG.md) for a list of changes in each release.

## Roadmap

### v1.1.0 (Planned)
- [ ] Add syslog support
- [ ] Add file rotation support
- [ ] Add performance benchmarks
- [ ] Add OpenTelemetry integration

### v1.2.0 (Planned)
- [ ] Add custom formatter support
- [ ] Add log sampling for high-volume scenarios
- [ ] Add async logging option
- [ ] Add log filtering capabilities

### Future
- [ ] Add metric collection
- [ ] Add log aggregation helpers
- [ ] Add cloud provider integrations (CloudWatch, Stackdriver, etc.)

## Support

### Getting Help

- **Documentation**: Read the full documentation in this README
- **Examples**: Check the `examples/` directory
- **Issues**: [GitHub Issues](https://github.com/ocrosby/go-logging/issues)
- **Discussions**: [GitHub Discussions](https://github.com/ocrosby/go-logging/discussions)

### Reporting Issues

If you find a bug or have a feature request, please create an issue on GitHub with:
- Clear description of the problem
- Steps to reproduce (for bugs)
- Expected vs actual behavior
- Go version and operating system
- Relevant code snippets or logs

### Security

If you discover a security vulnerability, please email security@example.com instead of creating a public issue.

## Acknowledgments

This library was inspired by logging patterns from:
- [Zerolog](https://github.com/rs/zerolog) - Fluent interface design
- [Zap](https://github.com/uber-go/zap) - Performance-focused logging
- [Logrus](https://github.com/sirupsen/logrus) - Structured logging approach

Special thanks to all contributors who help improve this library.

## Authors

- **Omar Crosby** - *Initial work* - [@ocrosby](https://github.com/ocrosby)

See also the list of [contributors](https://github.com/ocrosby/go-logging/contributors) who participated in this project.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

**Made with ❤️ for the Go community**

If you find this project useful, please consider giving it a ⭐️ on GitHub!

# go-logging

[![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.19-blue)](https://golang.org/dl/)
[![Go Report Card](https://goreportcard.com/badge/github.com/ocrosby/go-logging)](https://goreportcard.com/report/github.com/ocrosby/go-logging)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A modern, unified logging library for Go that seamlessly combines traditional logging with structured logging, featuring built-in slog integration, request tracing, fluent interfaces, and advanced async processing capabilities.

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

### 🎯 **Core Logging**
- **Unified Interface**: Single Logger interface with all level methods (Trace, Debug, Info, Warn, Error, Critical)
- **Dual Backend Support**: Seamlessly switch between standard Go logging and slog
- **Multiple Output Formats**: Text and JSON with customizable formatting
- **Thread-Safe**: Full concurrency support with proper synchronization

### 🔧 **Advanced Architecture**  
- **Consolidated Configuration**: Unified config system with backward compatibility
- **Async Processing**: Built-in async workers for high-performance logging
- **Handler Composition**: Advanced handler middleware and composition patterns
- **Level Dispatch**: Unified level method dispatch with automatic backend delegation

### 🌐 **Slog Integration**
- **Native slog Support**: Built on Go's standard `log/slog` for maximum compatibility
- **Custom Handler Support**: Use any `slog.Handler` (zerolog, zap, custom implementations)
- **Dynamic Level Control**: Runtime level changes with proper slog handler delegation
- **Attribute Preservation**: Full slog attribute and context support

### 📊 **Request Tracing**
- **Built-in Trace Support**: Trace IDs, request IDs, and correlation IDs
- **Context Propagation**: Automatic context field extraction and injection
- **HTTP Middleware**: Pre-built middleware for automatic request tracing
- **Header Integration**: Support for standard tracing headers

### 🎨 **Developer Experience**
- **Fluent Interface**: Expressive method chaining for readable logs
- **Structured Fields**: Type-safe field attachment with validation
- **Environment Config**: Automatic configuration from environment variables
- **Testing Support**: Mock-friendly interfaces with generated mocks

### 🛡️ **Production Ready**
- **Sensitive Data Redaction**: Built-in patterns for API keys, passwords, tokens
- **Error Handling**: Graceful degradation and error recovery
- **Performance Optimized**: Benchmarked and optimized for high-throughput scenarios
- **Memory Efficient**: Smart field copying and minimal allocations

## Installation

```bash
go get github.com/ocrosby/go-logging
```

### Requirements

- Go 1.19 or higher
- [Task](https://taskfile.dev) (optional, for using the Taskfile)

### Setup

This project uses [Task](https://taskfile.dev) for task automation. To set up your development environment:

1. Install Task (if not already installed):
   ```bash
   # macOS
   brew install go-task/tap/go-task
   
   # Linux
   sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d -b ~/.local/bin
   
   # Windows
   choco install go-task
   ```

2. Install development tools:
   ```bash
   task install-tools
   ```

3. Download dependencies:
   ```bash
   task deps
   ```

4. Build the project:
   ```bash
   task build
   ```

5. Run tests:
   ```bash
   task test
   ```

6. View all available tasks:
   ```bash
   task
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
// Create logger with fluent capabilities built-in
logger := logging.NewWithLevel(logging.DebugLevel)

// All loggers now support fluent interface directly
logger.Fluent().Info().
    Str("service", "api").
    Int("user_id", 12345).
    Msg("User logged in")

// Chain with context and error handling
logger.Fluent().Error().
    Err(err).
    Str("operation", "database_query").
    Int("retry_count", 3).
    Msg("Query failed after retries")
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

### Slog Integration

```go
import (
    "log/slog"
    "github.com/ocrosby/go-logging/pkg/logging"
)

func main() {
    // Use slog with default handler
    logger := logging.NewSlogTextLogger(logging.InfoLevel)
    logger.Info("Using slog backend")
    
    // Use custom slog handler
    handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
        Level: slog.LevelDebug,
        AddSource: true,
    })
    logger = logging.NewWithHandler(handler)
    logger.Debug("Custom handler with source location")
}
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

The Logger interface now provides a **unified, comprehensive API**:

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

    // Level-specific methods (all built-in)
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

    // Fluent interface (always available)
    Fluent() FluentLogger

    // Configuration (runtime changes)
    SetLevel(level Level)
    GetLevel() Level
}
```

**Key Improvements:**
- ✅ **All methods in one interface** - No more type assertions needed
- ✅ **Context support built-in** - Every level has a context variant  
- ✅ **Fluent interface included** - Available on all logger instances
- ✅ **Runtime configuration** - Change levels dynamically

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

// Create slog-based logger
logger := logging.NewSlogTextLogger(logging.InfoLevel)
logger := logging.NewSlogJSONLogger(logging.DebugLevel)

// Create with custom slog handler
handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
    Level: slog.LevelDebug,
})
logger := logging.NewWithHandler(handler)

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
- [`examples/slog/`](examples/slog/) - Slog integration with custom handlers
- [`examples/new-architecture/`](examples/new-architecture/) - Unified architecture showcase
- [`examples/custom-handlers/`](examples/custom-handlers/) - Advanced handler patterns

## Documentation

### 📚 **Comprehensive Guides**
- **[Architecture Guide](docs/ARCHITECTURE.md)** - Deep dive into the unified architecture and design decisions
- **[Examples Guide](docs/EXAMPLES.md)** - Comprehensive examples for all use cases and patterns
- **[Migration Guide](docs/MIGRATION.md)** - Smooth migration from older versions with zero breaking changes
- **[API Reference](docs/API_REFERENCE.md)** - Complete API documentation with all interfaces and functions

### 🔧 **Advanced Topics**  
- **[Advanced Features](docs/ADVANCED_FEATURES.md)** - Async processing, handler composition, custom middleware
- **[Slog Integration](docs/SLOG_INTEGRATION.md)** - Complete guide to slog backend integration

### 📈 **Project Information**
- **[Improvements Summary](docs/IMPROVEMENTS_SUMMARY.md)** - Overview of architectural improvements and benefits

## Testing

### Using Task

Run the test suite:

```bash
task test
```

Run tests with coverage:

```bash
task test-coverage
```

Run tests with race detector:

```bash
task test-race
```

Run benchmarks:

```bash
task bench
```

### Using Go directly

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

### Unified Architecture

The library now uses a **consolidated, modern architecture**:

```
pkg/logging/
├── logger.go              # Unified Logger interface (all methods included)
├── unified_logger.go      # Single implementation supporting both backends
├── level.go              # Log level definitions and parsing
├── level_dispatcher.go   # Unified level method dispatch system
├── config.go            # Legacy config (backward compatible)
├── config_new.go        # Modern structured configuration system
├── async_worker.go      # Generic async processing patterns
├── factory.go           # Enhanced factory functions
├── providers.go         # Dependency injection providers
├── trace.go             # Request tracing with context support
├── context_extractor.go # Context field extraction utilities
├── fluent.go            # Fluent interface implementation
├── middleware.go        # HTTP middleware for tracing
├── handler_interfaces.go # Unified handler interface system
├── handler_composition.go # Handler composition and middleware
├── handler_middleware.go # Handler middleware patterns
├── redactor.go          # Sensitive data redaction
├── registry.go          # Handler registry system
├── outputs.go           # Output implementations with async support
└── http.go              # HTTP logging utilities
```

**Key Architectural Improvements:**
- 🏗️ **Unified Logger**: Single implementation handles both standard and slog backends
- ⚡ **Async Workers**: Generic async processing with proper shutdown handling  
- 🔧 **Handler System**: Comprehensive handler composition and middleware
- 📊 **Structured Config**: Separated concerns (Core/Formatter/Output)
- 🎯 **Level Dispatch**: Centralized level method routing logic

## Contributing

Contributions are welcome! Please read our [Contributing Guidelines](CONTRIBUTING.md) before submitting a PR.

### Development Setup

1. Fork the repository
2. Clone your fork: `git clone https://github.com/yourusername/go-logging.git`
3. Set up development environment: `task install-tools && task deps`
4. Create a branch: `git checkout -b feature/my-feature`
5. Make your changes
6. Format and lint code: `task fmt && task lint`
7. Run tests: `task test`
8. Commit your changes: `git commit -am 'Add new feature'`
9. Push to the branch: `git push origin feature/my-feature`
10. Submit a pull request

### Code Style

- Follow standard Go conventions and idioms
- Format code with `task fmt` before committing
- Run linter with `task lint` to check code quality
- Ensure all tests pass with `task test`
- Add tests for new functionality
- Update documentation as needed
- Run `task check` to verify format, lint, tests, and build all pass

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

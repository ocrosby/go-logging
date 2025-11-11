# Slog Integration Guide

This document explains how the go-logging package integrates with Go's standard `log/slog` library.

## Overview

The go-logging package now wraps around Go's standard `log/slog` library, providing:
- Full backward compatibility with existing code
- Ability to use any `slog.Handler` implementation
- Easy integration with third-party logging libraries like zerolog and zap
- Support for custom log levels (TRACE and CRITICAL)

## Architecture

### Two Logger Implementations

The package now provides two logger implementations:

1. **StandardLogger** (`standard_logger.go`) - Original implementation
2. **SlogLogger** (`slog_logger.go`) - New slog-based implementation

Both implement the same `Logger` interface, ensuring backward compatibility.

### Level Mapping

Custom levels are mapped to slog levels:

```go
const (
    LevelTrace    = slog.Level(-8)  // Custom level more verbose than Debug
    LevelCritical = slog.Level(12)  // Custom level more severe than Error
)
```

Standard levels map directly:
- `DebugLevel` → `slog.LevelDebug`
- `InfoLevel` → `slog.LevelInfo`
- `WarnLevel` → `slog.LevelWarn`
- `ErrorLevel` → `slog.LevelError`

## Usage

### Using Slog Backend

#### Basic Usage

```go
// Create slog-based text logger
logger := logging.NewSlogTextLogger(logging.InfoLevel)
logger.Info("Using slog backend")

// Create slog-based JSON logger
jsonLogger := logging.NewSlogJSONLogger(logging.DebugLevel)
jsonLogger.Debug("JSON output with slog")
```

#### Custom Handler

```go
import "log/slog"

// Create custom handler
handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
    Level:     slog.LevelDebug,
    AddSource: true,
})

// Use custom handler
logger := logging.NewWithHandler(handler)
logger.Info("Using custom slog handler")
```

### Integration with Third-Party Libraries

#### Zerolog

```go
import (
    "github.com/rs/zerolog"
    "github.com/samber/slog-zerolog/v2"
)

zerologLogger := zerolog.New(os.Stdout)
handler := slogzerolog.Option{Logger: &zerologLogger}.NewZerologHandler()
logger := logging.NewWithHandler(handler)
```

#### Zap

```go
import (
    "go.uber.org/zap"
    "go.uber.org/zap/exp/zapslog"
)

zapLogger, _ := zap.NewProduction()
handler := zapslog.NewHandler(zapLogger.Core(), nil)
logger := logging.NewWithHandler(handler)
```

### Configuration

The `Config` struct supports slog configuration:

```go
config := logging.NewConfig().
    WithLevel(logging.InfoLevel).
    WithHandler(customHandler).  // Provide custom handler
    UseSlog(true).                // Enable slog backend
    Build()

logger := logging.New(func(b *logging.ConfigBuilder) {
    b.config = config
})
```

## Benefits

### Flexibility

By wrapping slog, you can:
- Use any logging backend that provides a `slog.Handler`
- Switch between different logging implementations without code changes
- Leverage the performance optimizations of various libraries

### Standard Library Integration

- Uses Go's standard structured logging approach
- Future-proof as slog evolves
- No need to learn multiple logging APIs

### Backward Compatibility

- Existing code continues to work without changes
- Can gradually migrate to slog-based loggers
- Factory functions maintain the same signatures

## Implementation Details

### SlogLogger Structure

```go
type slogLogger struct {
    mu            sync.RWMutex
    slog          *slog.Logger
    level         Level
    fields        []slog.Attr
    redactorChain RedactorChainInterface
}
```

### Key Methods

#### Log Method
```go
func (sl *slogLogger) log(ctx context.Context, level Level, msg string, args ...interface{})
```
- Converts custom levels to slog levels
- Applies redaction
- Adds context fields (request_id, etc.)
- Calls slog.LogAttrs()

#### WithField/WithFields
```go
func (sl *slogLogger) WithField(key string, value interface{}) Logger
```
- Creates new logger instance with additional attributes
- Uses slog.Attr for efficient attribute handling
- Maintains immutability pattern

### Provider Integration

The `ProvideLogger` function automatically selects the appropriate implementation:

```go
func ProvideLogger(config *Config, redactorChain RedactorChainInterface) Logger {
    if config.UseSlog {
        return NewSlogLoggerFromConfig(config, redactorChain)
    }
    return NewStandardLogger(config, redactorChain)
}
```

## Testing

All existing tests pass with the new implementation:
- Level handling tests
- Fluent interface tests
- Context propagation tests
- Redaction tests
- HTTP middleware tests

## Migration Guide

### From StandardLogger to SlogLogger

No code changes required! Simply use the new factory functions:

```go
// Before
logger := logging.NewWithLevel(logging.InfoLevel)

// After (using slog)
logger := logging.NewSlogTextLogger(logging.InfoLevel)
```

### Using Custom Handlers

```go
// Create your custom handler
handler := createYourCustomHandler()

// Use it with go-logging
logger := logging.NewWithHandler(handler)

// Everything else works the same
logger.Info("Message")
logger.WithField("key", "value").Info("Structured log")
logger.Fluent().Info().Str("user", "john").Msg("Fluent log")
```

## Performance Considerations

### Slog Benefits
- Efficient attribute handling
- Lazy evaluation support
- Optimized JSON encoding
- Better memory allocation patterns

### Custom Levels
- Custom levels (TRACE, CRITICAL) work seamlessly
- Handler implementations may need to handle these levels
- Most handlers will treat them as Debug/Error levels if not explicitly handled

## Future Enhancements

Potential improvements for slog integration:
- [ ] Context value extraction helpers
- [ ] Handler middleware support
- [ ] Performance benchmarks comparing implementations
- [ ] Additional handler examples (loki, datadog, etc.)
- [ ] Handler composition utilities

## Conclusion

The slog integration provides:
- **Flexibility**: Use any slog-compatible handler
- **Performance**: Leverage slog's optimizations
- **Compatibility**: Existing code works unchanged
- **Future-proof**: Built on Go's standard library

This approach follows the **Adapter Pattern**, wrapping slog while maintaining the existing Logger interface for backward compatibility and ease of use.

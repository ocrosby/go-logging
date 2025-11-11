# Slog Integration Example

This example demonstrates how to use the go-logging package with Go's standard `log/slog` library as the backend. This allows you to leverage slog's flexibility and use any third-party logging solution that implements the `slog.Handler` interface.

## Features

- **Slog Backend**: Uses Go's standard `log/slog` as the underlying logging mechanism
- **Custom Handlers**: Support for any `slog.Handler` implementation
- **Third-Party Integration**: Easy integration with popular logging libraries like zerolog, zap, etc.
- **Backward Compatible**: Maintains the same API as the existing logging package

## Usage

### Basic Slog Logger

```go
logger := logging.NewSlogTextLogger(logging.InfoLevel)
logger.Info("Application started with slog backend")
```

### JSON Format

```go
jsonLogger := logging.NewSlogJSONLogger(logging.DebugLevel)
jsonLogger.Debug("JSON formatted log with slog")
```

### Custom Handler

You can provide any custom `slog.Handler`:

```go
customHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
    Level:     slog.LevelDebug,
    AddSource: true,
})
customLogger := logging.NewWithHandler(customHandler)
customLogger.Info("Using custom slog handler")
```

### With Fields

```go
logger = logger.WithField("service", "slog-example")
logger.Info("Logger with static field")

logger = logger.WithFields(map[string]interface{}{
    "version": "2.0.0",
    "env":     "production",
})
logger.Info("Logger with multiple fields")
```

### Fluent Interface

```go
logger.Fluent().Info().
    Str("user", "john_doe").
    Int("attempts", 3).
    Msg("Login attempt")
```

## Integration with Third-Party Libraries

### Using with zerolog

```go
import (
    "github.com/rs/zerolog"
    "github.com/samber/slog-zerolog/v2"
)

zerologLogger := zerolog.New(os.Stdout)
handler := slogzerolog.Option{Logger: &zerologLogger}.NewZerologHandler()
logger := logging.NewWithHandler(handler)
```

### Using with zap

```go
import (
    "go.uber.org/zap"
    "go.uber.org/zap/exp/zapslog"
)

zapLogger, _ := zap.NewProduction()
handler := zapslog.NewHandler(zapLogger.Core(), nil)
logger := logging.NewWithHandler(handler)
```

## Level Mapping

The library maps its custom levels to slog levels:

- `TraceLevel` → Custom slog.Level(-8)
- `DebugLevel` → slog.LevelDebug
- `InfoLevel` → slog.LevelInfo
- `WarnLevel` → slog.LevelWarn
- `ErrorLevel` → slog.LevelError
- `CriticalLevel` → Custom slog.Level(12)

## Running the Example

```bash
cd examples/slog
go run main.go
```

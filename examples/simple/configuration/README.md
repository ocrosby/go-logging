# Configuration Example

Learn different ways to configure the go-logging library.

## Run the Example

```bash
go run main.go
```

## What This Example Shows

- **Builder pattern**: Use `NewEasyBuilder()` for flexible configuration
- **Static fields**: Add fields that appear in every log message
- **Environment configuration**: Use `NewFromEnvSimple()` with ENV vars
- **Different log levels**: Debug, Info, Warn, Error, Critical shortcuts

## Environment Variables

Try running with different environment variables:

```bash
# Debug level with JSON format
LOG_LEVEL=debug LOG_FORMAT=json go run main.go

# Info level with file information
LOG_LEVEL=info LOG_INCLUDE_FILE=true go run main.go
```

## Key Patterns

```go
// Builder with static fields
logger := logging.NewEasyBuilder().
    Level(logging.InfoLevel).
    JSON().
    Field("service", "my-app").
    Fields(map[string]any{
        "version": "1.0.0",
        "env": "production",
    }).
    Build()
```
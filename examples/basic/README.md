# Basic Logging Example

This example demonstrates the fundamental concepts and simplest ways to get started with the go-logging library. It's the perfect starting point for new users.

## What This Example Shows

### 1. **Simple Logging** - Zero Configuration
```go
logger := logging.NewSimple()
logger.Info("Application started")
```
The absolute simplest way to get started - just one line and you're logging!

### 2. **JSON Logging** - Easy Structured Output
```go
jsonLogger := logging.NewEasyJSON()
jsonLogger.Info("This will be formatted as JSON")
```
Perfect for production environments that need structured logs.

### 3. **Level Control** - Debug Mode
```go
debugLogger := logging.NewEasyJSONWithLevel(logging.DebugLevel)
debugLogger.Debug("This debug message will now appear")
```
Show how to control what log levels are displayed.

### 4. **Progressive Configuration** - Builder Pattern
```go
logger := logging.NewEasyBuilder().
    Level(logging.InfoLevel).
    JSON().
    WithFile().
    Field("service", "basic-example").
    Build()
```
Demonstrates the fluent builder pattern for more complex configurations.

### 5. **Environment Configuration** - Runtime Setup
```go
envLogger := logging.NewFromEnvSimple()
```
Shows how to configure logging from environment variables.

### 6. **Context Logging** - Request Tracing
```go
ctx := logging.NewContextWithTrace()
logger.InfoContext(ctx, "Message with trace ID from context")
```
Basic usage of context-aware logging with automatic trace ID generation.

## Expected Output

When you run this example, you'll see different output formats:

**Simple Logging (Text):**
```
[INFO] Application started
[WARN] This is a warning
[ERROR] This is an error
```

**JSON Logging:**
```json
{"level":"INFO","message":"This will be formatted as JSON","timestamp":"2025-11-12T..."}
```

**Debug Logging:**
```json
{"level":"DEBUG","message":"This debug message will now appear","timestamp":"2025-11-12T..."}
```

**Complex Configuration with Static Fields:**
```json
{"service":"basic-example","version":"1.0.0","env":"development","level":"INFO","message":"Logger with multiple configured options","timestamp":"2025-11-12T..."}
```

## Running the Example

```bash
cd examples/basic
go run main.go
```

### With Environment Variables

Try setting environment variables to see how they affect the output:

```bash
# Set log level and format
export LOG_LEVEL=debug
export LOG_FORMAT=json

cd examples/basic
go run main.go
```

## Key Learning Points

1. **Start Simple**: You only need one line to get started
2. **Progressive Complexity**: Add features as you need them
3. **Flexible Configuration**: Multiple ways to configure (code, environment, builder)
4. **Production Ready**: JSON output for structured logging
5. **Context Aware**: Built-in support for request tracing

## Next Steps

After running this example, try:

1. **[`simple/quick-start/`](../simple/quick-start/)** - Even simpler zero-config approach
2. **[`simple/structured/`](../simple/structured/)** - More structured logging examples  
3. **[`yaml-config/`](../yaml-config/)** - Most powerful configuration method
4. **[`fluent/`](../fluent/)** - Expressive fluent interface

## Best Practices Demonstrated

- **Start with defaults**: Use `NewSimple()` for immediate results
- **JSON for production**: Use `NewEasyJSON()` for structured output
- **Environment configuration**: Use `NewFromEnvSimple()` for runtime control
- **Static fields**: Add service context with `.Field()` and `.Fields()`
- **Context logging**: Use `InfoContext()` for request tracing
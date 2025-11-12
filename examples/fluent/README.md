# Fluent Interface Example

This example demonstrates the fluent interface for expressive, method-chained logging. The fluent interface provides a readable way to build complex log entries with multiple fields.

## What This Example Shows

### 1. **Basic Fluent Logging** - Method Chaining
```go
logger.Fluent().Info().
    Str("service", "fluent-example").
    Str("version", "1.0.0").
    Msg("Application started with fluent interface")
```
Build log entries by chaining field methods together for maximum readability.

### 2. **Multiple Data Types** - Type-Safe Fields
```go
logger.Fluent().Debug().
    Int("user_id", 12345).
    Str("username", "john_doe").
    Bool("active", true).
    Msg("User details")
```
Add different data types with type-specific methods: `Str()`, `Int()`, `Bool()`, etc.

### 3. **Error Logging** - Built-in Error Support
```go
err := errors.New("connection timeout")
logger.Fluent().Error().
    Err(err).
    Str("host", "db.example.com").
    Int("port", 5432).
    Msg("Database connection failed")
```
Use `Err(error)` to automatically include error details in a structured way.

### 4. **Context Integration** - Request Tracing
```go
ctx := logging.NewContextWithTrace()
ctx = logging.WithRequestID(ctx, "req-456")

logger.Fluent().Info().
    Ctx(ctx).
    Str("operation", "fetch_user").
    Msgf("Processing request for user %s", "john_doe")
```
Include context information like trace IDs and request IDs automatically.

## Expected Output

When you run this example, you'll see structured JSON output:

**Application Start:**
```json
{
  "level": "INFO",
  "message": "Application started with fluent interface",
  "service": "fluent-example",
  "version": "1.0.0",
  "timestamp": "2025-11-12T..."
}
```

**User Details:**
```json
{
  "level": "DEBUG", 
  "message": "User details",
  "user_id": 12345,
  "username": "john_doe",
  "active": true,
  "timestamp": "2025-11-12T..."
}
```

**Error with Context:**
```json
{
  "level": "ERROR",
  "message": "Database connection failed", 
  "error": "connection timeout",
  "host": "db.example.com",
  "port": 5432,
  "timestamp": "2025-11-12T..."
}
```

**With Request Context:**
```json
{
  "level": "INFO",
  "message": "Processing request for user john_doe",
  "operation": "fetch_user",
  "trace_id": "trace-abc123",
  "request_id": "req-456",
  "timestamp": "2025-11-12T..."
}
```

## Available Fluent Methods

### Field Methods
- `Str(key, value)` - String fields
- `Int(key, value)` - Integer fields  
- `Int64(key, value)` - 64-bit integer fields
- `Bool(key, value)` - Boolean fields
- `Err(error)` - Error fields (automatically uses "error" as key)
- `Field(key, value)` - Generic interface{} fields
- `Fields(map[string]interface{})` - Multiple fields at once

### Context Methods
- `Ctx(context)` - Add context fields (trace_id, request_id, etc.)
- `TraceID(id)` - Manually set trace ID

### Output Methods
- `Msg(message)` - Simple message
- `Msgf(format, args...)` - Formatted message (printf-style)

## Running the Example

```bash
cd examples/fluent
go run main.go
```

## When to Use Fluent Interface

The fluent interface is perfect when:

- **Complex Log Entries**: You need many fields in a single log entry
- **Type Safety**: You want compile-time checking of field types
- **Readability**: Method chaining makes the code self-documenting
- **Consistency**: All log entries follow the same structured pattern

## Comparison with Direct Logging

**Fluent Interface:**
```go
logger.Fluent().Error().
    Err(err).
    Str("host", "db.example.com").
    Int("port", 5432).
    Msg("Database connection failed")
```

**Direct Structured Logging:**
```go
logger.Error("Database connection failed",
    "error", err.Error(),
    "host", "db.example.com", 
    "port", 5432,
)
```

Both approaches produce the same output - choose based on your preference for readability vs. brevity.

## Next Steps

After running this example, try:

1. **[`simple/structured/`](../simple/structured/)** - Direct structured logging approach
2. **[`http-server/`](../http-server/)** - Fluent interface in web applications
3. **[`yaml-config/`](../yaml-config/)** - Configure fluent loggers with YAML

## Best Practices Demonstrated

- **Method Chaining**: Build complex log entries step by step
- **Type Safety**: Use specific methods for different data types
- **Error Handling**: Use `Err()` for consistent error logging
- **Context Integration**: Include request tracing with `Ctx()`
- **Readable Code**: Self-documenting log entry construction
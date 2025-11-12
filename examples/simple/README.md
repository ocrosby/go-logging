# Simple Logging Examples

This directory contains comprehensive, easy-to-understand examples that demonstrate the key features of the go-logging library. Each example is self-contained and focuses on a specific aspect of logging.

## Quick Start

If you're new to the library, start with these examples in order:

1. **[quick-start](./quick-start/)** - Get up and running in seconds
2. **[configuration](./configuration/)** - Learn different ways to configure loggers
3. **[structured](./structured/)** - Add structured data to your logs

## All Examples

### 🚀 [Quick Start](./quick-start/)
**Perfect for beginners**
- One-line logger creation
- Different output formats (text vs JSON)  
- Basic log levels

```bash
cd quick-start && go run main.go
```

### ⚙️ [Configuration](./configuration/)
**Learn configuration options**
- Builder pattern configuration
- Static fields for context
- Environment-based configuration
- File output with metadata

```bash
cd configuration && go run main.go
```

### 📊 [Structured Logging](./structured/)
**Add structure to your logs**
- Key-value pairs in log messages
- Service context with static fields
- Complex nested data structures
- Payment processing example

```bash
cd structured && go run main.go
```

### 🔗 [Context-Aware Logging](./context-logging/)
**Request tracing and context propagation**
- Automatic trace ID generation
- Request and correlation IDs
- Context propagation through function calls
- User and session tracking

```bash
cd context-logging && go run main.go
```

### ❌ [Error Handling](./error-handling/)
**Best practices for logging errors**
- Error vs warning distinction
- Structured error context
- Panic recovery
- Different error scenarios (network, validation, business logic)

```bash
cd error-handling && go run main.go
```

### 🌐 [HTTP Middleware](./middleware/)
**Web server integration**
- Request logging middleware
- Authentication logging
- CORS handling
- Response time tracking
- Status code monitoring

```bash
cd middleware && go run main.go
# Then visit http://localhost:8080/health
```

### ⚡ [Async Logging](./async/)
**High-performance logging**
- Asynchronous log processing
- High-throughput scenarios
- Graceful shutdown
- Queue management

```bash
cd async && go run main.go
```

## Common Patterns

### Simple Logger Creation
```go
// Text output to console
logger := logging.NewSimple()

// JSON output to console  
logger := logging.NewEasyJSON()

// Custom configuration
logger := logging.NewEasyBuilder().
    Level(logging.InfoLevel).
    JSON().
    Field("service", "my-app").
    Build()
```

### Context-Aware Logging
```go
// Create context with trace ID
ctx := logging.NewContextWithTrace()
ctx = logging.WithRequestID(ctx, "req_123")

// Log with context
logger.InfoContext(ctx, "Processing request", 
    "user_id", "user_456",
    "action", "update_profile",
)
```

### Structured Data
```go
// Simple key-value pairs
logger.Info("User login",
    "user_id", 12345,
    "email", "user@example.com", 
    "success", true,
)

// Complex nested data
logger.Info("Payment processed",
    "transaction_id", "txn_123",
    "amount", 99.99,
    "metadata", map[string]any{
        "gateway": "stripe",
        "currency": "USD",
    },
)
```

## Environment Configuration

Many examples support configuration via environment variables:

```bash
# Set log level
export LOG_LEVEL=debug

# Set output format  
export LOG_FORMAT=json

# Include file information
export LOG_INCLUDE_FILE=true

# Then run any example
cd quick-start && go run main.go
```

## Best Practices Demonstrated

### ✅ Do This
- Use structured logging with key-value pairs
- Include context (request IDs, trace IDs, user IDs)
- Log at appropriate levels (Debug < Info < Warn < Error)
- Use async logging for high-throughput applications
- Configure loggers once and reuse them
- Handle errors gracefully with proper context

### ❌ Avoid This
- String concatenation in log messages
- Logging sensitive data (passwords, tokens)
- Over-logging (too verbose in production)
- Blocking operations in logging code
- Inconsistent field naming across your application

## Production Considerations

### Performance
- Use async logging for high-throughput applications
- Consider buffered outputs for file logging
- Use appropriate queue sizes for async processing

### Security  
- Never log sensitive information
- Use redaction patterns for PII
- Consider log rotation for file outputs

### Monitoring
- Use structured JSON logs for log aggregation systems
- Include consistent fields across your application
- Use correlation IDs for request tracing

## Next Steps

After trying these examples:

1. **Explore Advanced Features**: Check out the `/examples/custom-handlers/` directory
2. **Integration**: Look at `/examples/slog/` for standard library integration  
3. **Dependency Injection**: See `/examples/di/` for wire integration
4. **Production Setup**: Review `/examples/http-server/` for real-world usage

## Need Help?

- 📖 [Full Documentation](../../docs/)
- 🏗️ [Architecture Guide](../../docs/ARCHITECTURE.md)
- 📚 [API Reference](../../docs/API_REFERENCE.md)
- 🚀 [Advanced Features](../../docs/ADVANCED_FEATURES.md)

Happy logging! 🪵
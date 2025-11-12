# HTTP Server Example

This example demonstrates how to integrate logging into a HTTP server with automatic request tracing, middleware, and context-aware logging throughout the request lifecycle.

## What This Example Shows

### 1. **HTTP Server Setup** - Production-Ready Logging
```go
config := logging.NewConfig().
    WithLevel(logging.InfoLevel).
    WithJSONFormat().
    Build()
logger := logging.NewStandardLogger(config, redactorChain)
```
Sets up structured JSON logging perfect for HTTP server applications.

### 2. **Request Handlers** - Context-Aware Logging
```go
mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    
    logger.Fluent().Info().
        Ctx(ctx).
        Str("handler", "hello").
        Msg("Handling hello request")
})
```
Each handler logs with context, automatically including trace and request IDs.

### 3. **Tracing Middleware** - Automatic Request Tracking
```go
handler := logging.TracingMiddleware(logger)(mux)
```
Wraps the entire server with tracing middleware that:
- Generates unique trace IDs for each request
- Extracts existing trace IDs from headers
- Adds request IDs and correlation IDs
- Logs request start and completion

### 4. **Multiple Endpoints** - Consistent Logging Pattern
```go
mux.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
    logger.Fluent().Info().
        Ctx(ctx).
        Str("handler", "user").
        Str("method", r.Method).
        Msg("Handling user request")
})
```
Shows how to maintain consistent logging across different endpoints.

## Expected Output

When you start the server and make requests, you'll see structured logs:

**Server Startup:**
```json
{
  "level": "INFO",
  "message": "Starting server on :8080",
  "timestamp": "2025-11-12T..."
}
```

**Request Handling (with automatic trace ID):**
```json
{
  "level": "INFO",
  "message": "Handling hello request",
  "handler": "hello",
  "trace_id": "trace-abc123",
  "request_id": "req-456",
  "timestamp": "2025-11-12T..."
}
```

**With Method Information:**
```json
{
  "level": "INFO", 
  "message": "Handling user request",
  "handler": "user",
  "method": "GET",
  "trace_id": "trace-def789",
  "request_id": "req-789",
  "timestamp": "2025-11-12T..."
}
```

## Running the Example

```bash
cd examples/http-server
go run main.go
```

The server will start on port 8080. In another terminal, test the endpoints:

```bash
# Test the hello endpoint
curl http://localhost:8080/hello

# Test the user endpoint  
curl http://localhost:8080/user

# Test with existing trace ID
curl -H "X-Trace-ID: custom-trace-123" http://localhost:8080/hello
```

## Key Features Demonstrated

### Automatic Request Tracing
- **Trace ID Generation**: Each request gets a unique trace ID
- **Header Support**: Respects existing `X-Trace-ID`, `X-Request-ID`, `X-Correlation-ID` headers
- **Context Propagation**: Trace information flows through the entire request

### Structured Logging
- **JSON Output**: Perfect for log aggregation systems
- **Consistent Fields**: Every log entry includes trace information
- **Handler Identification**: Easy to see which handler processed the request

### Production Ready
- **Error Handling**: Server startup failures are logged as critical
- **Middleware Pattern**: Clean separation of concerns
- **Context Usage**: Proper context handling throughout the request lifecycle

## HTTP Headers Supported

The tracing middleware automatically handles these headers:

| Header | Purpose | Example |
|--------|---------|---------|
| `X-Trace-ID` | Unique trace identifier | `trace-abc123` |
| `X-Request-ID` | Individual request identifier | `req-456` |  
| `X-Correlation-ID` | Cross-service correlation | `corr-789` |

If these headers aren't present, the middleware generates them automatically.

## Integration with Log Aggregation

This example produces logs perfect for:
- **ELK Stack** (Elasticsearch, Logstash, Kibana)
- **Fluentd/Fluent Bit**
- **Grafana Loki**
- **Datadog, New Relic, etc.**

The structured JSON format with consistent trace IDs makes it easy to:
- Track requests across services
- Debug issues by trace ID
- Create dashboards and alerts
- Monitor application performance

## Best Practices Demonstrated

1. **Structured Logging**: Use JSON format for HTTP servers
2. **Context Propagation**: Always use `r.Context()` in handlers
3. **Middleware Pattern**: Wrap your entire router with tracing middleware
4. **Consistent Fields**: Include handler and method information
5. **Error Logging**: Log server startup failures as critical events

## Next Steps

After running this example, try:

1. **[`simple/middleware/`](../simple/middleware/)** - Simpler middleware patterns
2. **[`simple/context-logging/`](../simple/context-logging/)** - More context examples
3. **[`yaml-config/`](../yaml-config/)** - Configure HTTP server logging via YAML

## Common Use Cases

This pattern works great for:
- **REST APIs** - Track API endpoint usage
- **Web Applications** - Monitor user interactions
- **Microservices** - Trace requests across service boundaries
- **Load Balanced Services** - Debug issues across multiple instances
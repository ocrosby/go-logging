# Examples Guide

This guide provides comprehensive examples of using the go-logging library in various scenarios.

## Table of Contents

- [Basic Usage](#basic-usage)
- [Configuration Examples](#configuration-examples)
- [Slog Integration](#slog-integration)
- [Fluent Interface](#fluent-interface)
- [Context and Tracing](#context-and-tracing)
- [HTTP Middleware](#http-middleware)
- [Advanced Patterns](#advanced-patterns)
- [Performance Examples](#performance-examples)
- [Testing Examples](#testing-examples)

## Basic Usage

### Simple Logging

```go
package main

import "github.com/ocrosby/go-logging/pkg/logging"

func main() {
    // Create logger with default settings
    logger := logging.NewWithLevel(logging.InfoLevel)
    
    // Use level-specific methods
    logger.Info("Application started")
    logger.Warn("This is a warning message")
    logger.Error("Error occurred: %v", someError)
    
    // Or use the generic Log method
    logger.Log(logging.InfoLevel, "Generic log message")
}
```

### With Fields

```go
func main() {
    logger := logging.NewWithLevel(logging.InfoLevel)
    
    // Attach static fields (immutable pattern)
    serviceLogger := logger.WithFields(map[string]interface{}{
        "service": "user-service",
        "version": "1.0.0",
        "environment": "production",
    })
    
    // This logger will include the fields in all messages
    serviceLogger.Info("Service started")
    
    // Add more fields
    userLogger := serviceLogger.WithField("user_id", 12345)
    userLogger.Info("User logged in")
    
    // Original logger unchanged
    logger.Info("This message has no extra fields")
}
```

## Configuration Examples

### Environment-Based Configuration

```go
func main() {
    // Set environment variables:
    // export LOG_LEVEL=DEBUG
    // export LOG_FORMAT=json
    
    logger := logging.NewFromEnvironment()
    logger.Debug("Debug message will show if LOG_LEVEL=DEBUG")
}
```

### Advanced Configuration

```go
func main() {
    // Using the new structured configuration system
    config := logging.NewLoggerConfig().
        WithCore(
            logging.NewCoreConfig().
                WithLevel(logging.DebugLevel).
                WithStaticField("service", "api-gateway").
                WithStaticField("datacenter", "us-east-1").
                Build(),
        ).
        WithFormatter(
            logging.NewFormatterConfig().
                WithJSONFormat().
                IncludeFile(true).
                IncludeTime(true).
                UseShortFile(true).
                AddRedactPattern(`password=\w+`).
                AddRedactPattern(`token=[\w-]+`).
                Build(),
        ).
        WithOutput(
            logging.NewOutputConfig().
                WithWriter(os.Stdout).
                Build(),
        ).
        Build()
    
    logger := logging.NewWithLoggerConfig(config)
    logger.Info("Configured logger ready")
}
```

### Legacy Configuration (Backward Compatible)

```go
func main() {
    // Old configuration API still works
    config := logging.NewConfig().
        WithLevel(logging.InfoLevel).
        WithJSONFormat().
        IncludeFile(true).
        AddRedactPattern(`apikey=\w+`).
        WithStaticFields(map[string]interface{}{
            "app": "my-service",
        }).
        Build()
    
    redactorChain := logging.NewRedactorChain()
    logger := logging.NewStandardLogger(config, redactorChain)
    logger.Info("Legacy config works fine")
}
```

## Slog Integration

### Basic Slog Usage

```go
func main() {
    // Create loggers with slog backend
    textLogger := logging.NewSlogTextLogger(logging.DebugLevel)
    jsonLogger := logging.NewSlogJSONLogger(logging.InfoLevel)
    
    textLogger.Info("Text formatted message")
    jsonLogger.Info("JSON formatted message")
}
```

### Custom Slog Handlers

```go
func main() {
    // Use any slog handler
    handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
        Level:     slog.LevelDebug,
        AddSource: true,
        ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
            // Custom attribute transformation
            if a.Key == slog.TimeKey {
                a.Value = slog.StringValue(time.Now().Format(time.RFC3339))
            }
            return a
        },
    })
    
    logger := logging.NewWithHandler(handler)
    logger.Debug("Message with custom handler")
}
```

### Third-Party Handler Integration

```go
import (
    "github.com/rs/zerolog"
    slogzerolog "github.com/samber/slog-zerolog"
)

func main() {
    // Use zerolog as slog handler
    zerologLogger := zerolog.New(os.Stdout).With().Timestamp().Logger()
    handler := slogzerolog.Option{Logger: &zerologLogger}.NewZerologHandler()
    
    logger := logging.NewWithHandler(handler)
    logger.Info("Using zerolog through slog")
}
```

## Fluent Interface

### Basic Fluent Logging

```go
func main() {
    logger := logging.NewWithLevel(logging.InfoLevel)
    
    // All loggers have fluent interface built-in
    logger.Fluent().Info().
        Str("service", "payment-service").
        Int("amount", 1000).
        Str("currency", "USD").
        Msg("Payment processed")
}
```

### Complex Fluent Examples

```go
func processOrder(logger logging.Logger, order *Order) error {
    // Fluent logging with error handling
    if err := validateOrder(order); err != nil {
        logger.Fluent().Error().
            Err(err).
            Str("order_id", order.ID).
            Int("user_id", order.UserID).
            Msg("Order validation failed")
        return err
    }
    
    // Success case with multiple data types
    logger.Fluent().Info().
        Str("order_id", order.ID).
        Int("user_id", order.UserID).
        Float64("total", order.Total).
        Bool("express", order.ExpressShipping).
        Int64("timestamp", time.Now().Unix()).
        Fields(map[string]interface{}{
            "items_count": len(order.Items),
            "payment_method": order.PaymentMethod,
        }).
        Msg("Order processed successfully")
    
    return nil
}
```

### Fluent with Context

```go
func handleRequest(ctx context.Context, logger logging.Logger) {
    // Extract trace info from context automatically
    logger.Fluent().Info().
        Ctx(ctx).  // Automatically includes trace_id, request_id, etc.
        Str("handler", "user_profile").
        Str("method", "GET").
        Msg("Processing request")
}
```

## Context and Tracing

### Manual Context Management

```go
func main() {
    logger := logging.NewWithLevel(logging.InfoLevel)
    
    // Create context with trace information
    ctx := context.Background()
    ctx = logging.WithTraceID(ctx, logging.NewTraceID())
    ctx = logging.WithRequestID(ctx, "req-12345")
    ctx = logging.WithCorrelationID(ctx, "corr-67890")
    
    // Context information automatically included
    logger.InfoContext(ctx, "Request processing started")
    
    // Retrieve context info
    if traceID, ok := logging.GetTraceID(ctx); ok {
        fmt.Printf("Trace ID: %s\n", traceID)
    }
}
```

### Automatic Context Creation

```go
func handleHTTPRequest(w http.ResponseWriter, r *http.Request) {
    // Create context with automatic trace ID generation
    ctx := logging.NewContextWithTrace()
    ctx = logging.WithRequestID(ctx, extractRequestID(r))
    
    logger := logging.NewWithLevel(logging.InfoLevel)
    
    // All subsequent logs will include trace information
    logger.InfoContext(ctx, "Request started")
    
    // Pass context through the call chain
    processRequest(ctx, logger, r)
}

func processRequest(ctx context.Context, logger logging.Logger, r *http.Request) {
    logger.InfoContext(ctx, "Processing in business logic")
    // Trace ID automatically propagated
}
```

## HTTP Middleware

### Basic Tracing Middleware

```go
func main() {
    logger := logging.NewJSONLogger(logging.InfoLevel)
    
    mux := http.NewServeMux()
    mux.HandleFunc("/api/users", handleUsers)
    mux.HandleFunc("/api/orders", handleOrders)
    
    // Add tracing middleware
    handler := logging.TracingMiddleware(logger)(mux)
    
    http.ListenAndServe(":8080", handler)
}

func handleUsers(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    logger := logging.NewWithLevel(logging.InfoLevel)
    
    // Context automatically contains trace information
    logger.InfoContext(ctx, "Fetching users")
    
    // Your business logic here
    w.WriteHeader(http.StatusOK)
}
```

### Request Logger Middleware

```go
func main() {
    logger := logging.NewWithLevel(logging.InfoLevel)
    
    mux := http.NewServeMux()
    mux.HandleFunc("/health", healthCheck)
    
    // Log request details with specific headers
    handler := logging.RequestLogger(logger, "User-Agent", "X-API-Key", "Authorization")(mux)
    
    http.ListenAndServe(":8080", handler)
}
```

### Combined Middleware Stack

```go
func main() {
    logger := logging.NewJSONLogger(logging.InfoLevel)
    
    mux := http.NewServeMux()
    mux.HandleFunc("/api/data", dataHandler)
    
    // Combine multiple middleware
    handler := logging.TracingMiddleware(logger)(
        logging.RequestLogger(logger, "User-Agent")(
            loggingMiddleware(logger)(mux),
        ),
    )
    
    http.ListenAndServe(":8080", handler)
}

func loggingMiddleware(logger logging.Logger) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            start := time.Now()
            
            logger.Fluent().Info().
                Str("method", r.Method).
                Str("path", r.URL.Path).
                Str("remote_addr", r.RemoteAddr).
                Msg("Request started")
            
            next.ServeHTTP(w, r)
            
            logger.Fluent().Info().
                Str("method", r.Method).
                Str("path", r.URL.Path).
                Int64("duration_ms", time.Since(start).Milliseconds()).
                Msg("Request completed")
        })
    }
}
```

## Advanced Patterns

### Conditional Logging

```go
func expensiveOperation(logger logging.Logger) {
    // Check if debug is enabled before expensive operations
    if logger.IsLevelEnabled(logging.DebugLevel) {
        data := performExpensiveDataCollection()
        logger.Fluent().Debug().
            Field("data", data).
            Msg("Debug data collected")
    }
    
    // Level-specific methods automatically check levels
    logger.Info("Operation completed")  // No pre-check needed
}
```

### Dynamic Level Changes

```go
func main() {
    logger := logging.NewWithLevel(logging.InfoLevel)
    
    logger.Info("Starting with INFO level")
    logger.Debug("This won't show")
    
    // Change level at runtime
    logger.SetLevel(logging.DebugLevel)
    logger.Debug("Now debug messages show")
    
    // Check current level
    currentLevel := logger.GetLevel()
    fmt.Printf("Current level: %v\n", currentLevel)
}
```

### Async Processing

```go
func main() {
    // Create async output for high-performance logging
    fileOutput := &logging.FileOutput{
        Filename: "app.log",
    }
    
    asyncOutput := logging.NewAsyncOutput(fileOutput, 1000) // 1000 item queue
    defer asyncOutput.Close()
    
    config := logging.NewLoggerConfig().
        WithOutput(
            logging.NewOutputConfig().
                WithWriter(asyncOutput).
                Build(),
        ).
        Build()
    
    logger := logging.NewWithLoggerConfig(config)
    
    // High-throughput logging
    for i := 0; i < 10000; i++ {
        logger.Fluent().Info().
            Int("iteration", i).
            Msg("Processing item")
    }
    
    // Ensure all logs are written before exit
    asyncOutput.Stop()
}
```

## Performance Examples

### High-Throughput Logging

```go
func benchmarkLogging() {
    logger := logging.NewSlogJSONLogger(logging.InfoLevel)
    
    // Pre-create logger with static fields to avoid allocation
    serviceLogger := logger.WithFields(map[string]interface{}{
        "service": "high-throughput-service",
        "version": "2.1.0",
    })
    
    start := time.Now()
    
    for i := 0; i < 100000; i++ {
        // Efficient logging without field allocation
        serviceLogger.Info("Processing item %d", i)
    }
    
    duration := time.Since(start)
    fmt.Printf("Logged 100k messages in %v\n", duration)
}
```

### Memory-Efficient Field Usage

```go
func efficientFieldLogging(logger logging.Logger) {
    // Good: Pre-create logger with common fields
    userLogger := logger.WithFields(map[string]interface{}{
        "user_id": 12345,
        "session": "sess-abc",
    })
    
    // Efficient: Reuse the logger
    userLogger.Info("User action 1")
    userLogger.Info("User action 2")
    
    // Less efficient: Creating fields each time
    // logger.WithField("user_id", 12345).Info("User action")
    // logger.WithField("user_id", 12345).Info("Another action")
}
```

## Testing Examples

### Using Mocks

```go
func TestServiceWithMocking(t *testing.T) {
    // Create mock logger
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    
    mockLogger := mocks.NewMockLogger(ctrl)
    
    // Set expectations
    mockLogger.EXPECT().
        WithField("user_id", 123).
        Return(mockLogger)
    
    mockLogger.EXPECT().
        Info("User processed successfully")
    
    // Test your service
    service := &UserService{Logger: mockLogger}
    service.ProcessUser(123)
}
```

### Test Logger Implementation

```go
type TestLogger struct {
    Entries []TestLogEntry
}

type TestLogEntry struct {
    Level   logging.Level
    Message string
    Fields  map[string]interface{}
}

func (t *TestLogger) Info(msg string, args ...interface{}) {
    t.Entries = append(t.Entries, TestLogEntry{
        Level:   logging.InfoLevel,
        Message: fmt.Sprintf(msg, args...),
        Fields:  make(map[string]interface{}),
    })
}

// Implement other Logger interface methods...

func TestBusinessLogic(t *testing.T) {
    testLogger := &TestLogger{}
    
    // Run your business logic
    processOrder(testLogger, &Order{ID: "123"})
    
    // Assert logging behavior
    if len(testLogger.Entries) != 1 {
        t.Errorf("Expected 1 log entry, got %d", len(testLogger.Entries))
    }
    
    entry := testLogger.Entries[0]
    if !strings.Contains(entry.Message, "Order processed") {
        t.Errorf("Expected log message about order processing")
    }
}
```

### Integration Testing

```go
func TestHTTPEndpointLogging(t *testing.T) {
    // Capture logs in buffer
    var logBuffer bytes.Buffer
    
    config := logging.NewLoggerConfig().
        WithFormatter(
            logging.NewFormatterConfig().
                WithJSONFormat().
                Build(),
        ).
        WithOutput(
            logging.NewOutputConfig().
                WithWriter(&logBuffer).
                Build(),
        ).
        Build()
    
    logger := logging.NewWithLoggerConfig(config)
    
    // Create test server with logging
    handler := logging.TracingMiddleware(logger)(
        http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            logger.InfoContext(r.Context(), "Handling test request")
            w.WriteHeader(http.StatusOK)
        }),
    )
    
    server := httptest.NewServer(handler)
    defer server.Close()
    
    // Make request
    resp, err := http.Get(server.URL + "/test")
    require.NoError(t, err)
    require.Equal(t, http.StatusOK, resp.StatusCode)
    
    // Verify logs
    logOutput := logBuffer.String()
    assert.Contains(t, logOutput, "Handling test request")
    assert.Contains(t, logOutput, "trace_id")
    
    // Parse JSON logs for detailed assertions
    lines := strings.Split(strings.TrimSpace(logOutput), "\n")
    for _, line := range lines {
        var logEntry map[string]interface{}
        err := json.Unmarshal([]byte(line), &logEntry)
        require.NoError(t, err)
        
        // Verify log structure
        assert.Contains(t, logEntry, "timestamp")
        assert.Contains(t, logEntry, "level")
        assert.Contains(t, logEntry, "message")
    }
}
```

---

These examples demonstrate the flexibility and power of the go-logging library across various use cases. The unified interface makes it easy to adopt consistent logging patterns throughout your application while maintaining high performance and extensive customization capabilities.
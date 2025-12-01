# Pluggable Backend Examples

This example demonstrates the **pluggable architecture** of go-logging by showing how to use different logging backends (zerolog, zap, slog) while maintaining the same Logger interface throughout your application.

## What is a Pluggable Backend?

The go-logging library uses Go's standard `log/slog` as an abstraction layer. Any `slog.Handler` implementation can be plugged in as the backend, allowing you to:

1. **Swap backends without changing application code** - Your application uses the `logging.Logger` interface, which remains constant regardless of the backend
2. **Leverage existing ecosystems** - Use high-performance loggers like zap or zerolog
3. **Maintain consistency** - Same API for all backends (Trace, Debug, Info, Warn, Error, Critical)
4. **Mix and match** - Use different backends in different parts of your application

## Architecture

```
Your Application Code
        ↓
logging.Logger Interface (consistent API)
        ↓
log/slog.Handler (abstraction layer)
        ↓
Backend Implementation (zerolog, zap, or slog)
```

## Examples in This Directory

### 1. Standard slog Backend
Shows using Go's built-in `log/slog` handlers (JSON and Text).

### 2. Zerolog Backend
Demonstrates using zerolog as the backend via the slog bridge.

### 3. Zap Backend
Demonstrates using Uber's zap logger as the backend via the slog bridge.

### 4. Common Log Format (CLF) Backend
Demonstrates a **custom pluggable backend** implementing the NCSA Common Log Format standard used by web servers (Apache, Nginx). This shows how you can create your own `slog.Handler` implementation for any custom format.

## Running the Examples

```bash
# Run all examples
go run main.go

# Or test individual functions by modifying main.go
```

## Key Takeaways

1. **Same Interface**: All three backends use `logging.Logger` interface
2. **One Line Change**: Switch backends by changing which handler you pass to `logging.NewWithHandler()`
3. **No Code Changes**: Your application code using `logger.Info()`, `logger.Error()`, etc. stays identical
4. **Performance Benefits**: Use high-performance backends (zap, zerolog) in production while keeping flexibility

## When to Use Each Backend

| Backend | Best For | Performance | Features |
|---------|----------|-------------|----------|
| **slog** | Simple apps, Go 1.21+ standard | Good | Built-in, no dependencies |
| **zerolog** | High-performance JSON logging | Excellent | Zero allocation, fastest |
| **zap** | Structured logging with type safety | Excellent | Strong typing, field validation |
| **CLF (custom)** | Web server access logs | Excellent | Optimized with pooling, buffering, O(1) lookups |

## Implementation Details

### How It Works

1. Each backend provides a `slog.Handler` implementation
2. `logging.NewWithHandler(handler)` wraps the handler in `unifiedLogger`
3. The `unifiedLogger` delegates all log calls to the slog handler
4. Your application code only depends on `logging.Logger` interface

### Key Files in the Library

- `pkg/logging/logger.go` - Logger interface definition
- `pkg/logging/unified_logger.go` - Implementation that wraps slog handlers
- `pkg/logging/factory.go` - `NewWithHandler()` factory function

### Creating Your Own Custom Backend

The Common Log Format handler (`clf_handler.go`) demonstrates how to create a custom backend:

1. **Implement `slog.Handler` interface**:
   ```go
   type CommonLogFormatHandler struct {
       writer io.Writer
       attrs  []slog.Attr
       groups []string
   }
   
   func (h *CommonLogFormatHandler) Enabled(ctx context.Context, level slog.Level) bool
   func (h *CommonLogFormatHandler) Handle(ctx context.Context, record slog.Record) error
   func (h *CommonLogFormatHandler) WithAttrs(attrs []slog.Attr) slog.Handler
   func (h *CommonLogFormatHandler) WithGroup(name string) slog.Handler
   ```

2. **Format the output** in your `Handle()` method:
   ```go
   // CLF Format: host ident authuser [timestamp] "request" status bytes
   logLine := fmt.Sprintf("%s %s %s [%s] \"%s\" %s %s\n",
       host, ident, authuser, timestamp, requestLine, status, bytes)
   ```

3. **Use it** with the logging library:
   ```go
   handler := NewCommonLogFormatHandler(os.Stdout)
   logger := logging.NewWithHandler(handler)
   logger.Info("Request logged")  // Outputs in CLF format
   ```

This pattern works for any custom format: CSV, XML, Protocol Buffers, database logging, etc.

### Performance Optimizations in CLF Handler

The CLF handler demonstrates high-performance logging techniques:

**Benchmark Results** (Apple M3 Pro):
```
BenchmarkCLFHandler-12              5043409    233.7 ns/op      0 B/op    0 allocs/op
BenchmarkCLFHandlerParallel-12      6535747    184.9 ns/op      0 B/op    0 allocs/op
BenchmarkCLFHandlerWithAttrs-12     2994426    396.3 ns/op    464 B/op    3 allocs/op
```

**Key Optimizations:**
1. **Byte slice pool** - Zero-allocation writes by reusing byte slices
2. **Lock-free attribute reads** - Immutable attrMap eliminates read locks
3. **Timestamp caching** - Caches formatted timestamp per second
4. **Map-based attribute storage** - O(1) lookups instead of O(n) slice iteration
5. **Buffered writer (64KB)** - Reduces syscalls by batching writes
6. **Direct byte appends** - No string concatenation or `fmt.Sprintf`

**Performance Comparison:**
- **Initial version**: ~2500 ns/op, 1500+ B/op, 20+ allocs/op
- **First optimization pass**: ~468 ns/op, 368 B/op, 6 allocs/op (5x faster)
- **Final optimization pass**: ~234 ns/op, 0 B/op, 0 allocs/op (10x faster)
- **Total improvement**: **10x faster, 100% zero-allocation, eliminates all allocations**

These techniques can be applied to any custom handler implementation.

## Further Reading

- [Go slog documentation](https://pkg.go.dev/log/slog)
- [Zerolog slog handler](https://github.com/rs/zerolog)
- [Zap slog handler](https://github.com/uber-go/zap/tree/master/exp/zapslog)

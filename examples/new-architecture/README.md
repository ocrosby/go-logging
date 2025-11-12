# New Architecture Example

This example demonstrates the unified logging architecture, showcasing how all logging interfaces are consolidated into a single, comprehensive `Logger` interface while maintaining full backward compatibility.

## What This Example Shows

### 1. **Unified Logger Interface** - Everything in One Place
```go
logger := logging.NewWithLevel(logging.InfoLevel)

// Core logging methods
logger.Log(logging.InfoLevel, "This uses the core Logger interface")
logger.LogContext(context.Background(), logging.InfoLevel, "With context")

// Level-specific methods (built-in)
logger.Info("Convenience Info method")
logger.Debug("Debug message")

// Fluent interface (built-in)
logger.Fluent().Info().
    Str("service", "example").
    Msg("Fluent interface")
```

**Key Point**: No more type assertions needed! Everything is available on the single `Logger` interface.

### 2. **New Configuration Structure** - Modular Design
```go
config := logging.NewLoggerConfig().
    WithLevel(logging.InfoLevel).
    WithJSONFormat().
    Build()

logger := logging.NewWithLoggerConfig(config)
```
The new configuration system separates concerns into Core, Formatter, and Output configurations.

### 3. **Direct Component Usage** - Fine-Grained Control
```go
// Use formatters directly
formatterConfig := logging.NewFormatterConfig().WithJSONFormat().Build()
formatter := logging.NewJSONFormatter(formatterConfig)

// Format entries manually
entry := logging.LogEntry{
    Level:   logging.InfoLevel,
    Message: "Direct formatter usage",
    Fields:  map[string]interface{}{"component": "demo"},
}

data, err := formatter.Format(entry)
```
Access individual components for maximum control.

### 4. **Handler Registry System** - Extensible Architecture
```go
handlers := logging.ListHandlers()
for _, name := range handlers {
    logger.Info("Available handler: %s", name)
}
```
Discover and use registered handlers dynamically.

## Expected Output

When you run this example, you'll see various output formats:

**Unified Interface Usage:**
```
[INFO] This uses the core Logger interface
[INFO] This includes context  
[INFO] This uses the convenience Info method
{"level":"INFO","message":"This uses the fluent interface","service":"example","version":1,"timestamp":"2025-11-12T..."}
```

**JSON Configuration:**
```json
{"level":"INFO","message":"This will be formatted as JSON","timestamp":"2025-11-12T..."}
{"level":"INFO","message":"With structured fields","user_id":123,"timestamp":"2025-11-12T..."}
```

**Direct Formatter Usage:**
```json
{"level":"INFO","message":"Direct formatter usage","component":"formatter-demo"}
```

**Handler Registry:**
```
[INFO] Available handler: text
[INFO] Available handler: json
```

## Architecture Benefits

### 1. **No More Type Assertions**

**Old Way (Required Type Assertions):**
```go
// Had to check for specific interfaces
if fluentLogger, ok := logger.(FluentCapable); ok {
    fluentLogger.Fluent().Info().Msg("message")
}

if levelLogger, ok := logger.(LevelLogger); ok {
    levelLogger.Info("message")
}
```

**New Way (Everything Built-In):**
```go
// Everything available directly
logger.Fluent().Info().Msg("message")
logger.Info("message")
```

### 2. **Unified Configuration**

**Structured Configuration System:**
- `CoreConfig` - Level, static fields
- `FormatterConfig` - Format, file info, redaction
- `OutputConfig` - Destination configuration
- `LoggerConfig` - Combines all configs

### 3. **Component Architecture**

```
Logger (Interface)
├── UnifiedLogger (Implementation)
├── Formatters (JSON, Text, Console)
├── Outputs (Writer, File, Multi, Async)
├── Handlers (Registry system)
└── Middleware (Composable patterns)
```

## Key Architectural Improvements

### 1. **Single Implementation**
- `UnifiedLogger` handles both standard and slog backends
- Automatic backend selection based on configuration
- Consistent behavior across all usage patterns

### 2. **Handler System** 
- Registry-based handler discovery
- Composition and middleware patterns
- Extensible architecture for custom handlers

### 3. **Level Dispatch**
- Centralized level method routing
- Consistent level filtering
- Performance optimized dispatching

### 4. **Async Processing**
- Generic async worker patterns
- Proper shutdown handling
- High-throughput scenarios support

## Running the Example

```bash
cd examples/new-architecture
go run main.go
```

## Comparing Old vs New Architecture

### Logger Creation

**Old:**
```go
config := logging.NewConfig().WithLevel(logging.InfoLevel).Build()
logger := logging.NewStandardLogger(config, redactorChain)

// Type assertion needed for different interfaces
fluentLogger := logger.(logging.FluentCapable)
```

**New:**
```go
logger := logging.NewWithLevel(logging.InfoLevel)

// Everything available directly
logger.Info("message")
logger.Fluent().Info().Msg("message")
```

### Configuration

**Old:**
```go
// Single monolithic config
config := logging.NewConfig().
    WithLevel(logging.InfoLevel).
    WithJSONFormat().
    IncludeFile(true).
    Build()
```

**New:**
```go
// Modular, composable configuration
config := logging.NewLoggerConfig().
    WithCore(logging.NewCoreConfig().WithLevel(logging.InfoLevel).Build()).
    WithFormatter(logging.NewFormatterConfig().WithJSONFormat().Build()).
    WithOutput(logging.NewOutputConfig().Build()).
    Build()
```

## Best Practices Demonstrated

1. **Use Unified Interface**: Take advantage of the single Logger interface
2. **Modular Configuration**: Use the new structured configuration system
3. **Direct Access**: Access components directly when you need fine-grained control
4. **Handler Registry**: Leverage the extensible handler system
5. **Backward Compatibility**: Old patterns still work, migrate at your pace

## Migration Path

This example shows you can:

1. **Start Simple**: Use existing patterns (they still work)
2. **Gradual Migration**: Move to new configuration when convenient  
3. **Advanced Usage**: Use direct component access for special needs
4. **Full Power**: Leverage the complete unified architecture

## Next Steps

After exploring this example:

1. **[`basic/`](../basic/)** - See how the unified interface simplifies basic usage
2. **[`custom-handlers/`](../custom-handlers/)** - Explore the advanced handler system
3. **[`yaml-config/`](../yaml-config/)** - Use YAML with the new architecture
4. **[Architecture Guide](../../docs/ARCHITECTURE.md)** - Deep dive into the design

## When to Use This Pattern

Perfect for:
- **Understanding the library**: See how everything fits together
- **Advanced customization**: Need direct component access
- **Migration planning**: Understanding new vs old patterns
- **Architecture exploration**: Learning the unified design principles
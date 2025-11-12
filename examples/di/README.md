# Dependency Injection Example

This example demonstrates how to integrate the go-logging library with Google Wire for dependency injection, enabling clean architecture patterns and testable code.

## What This Example Shows

### 1. **Wire Integration** - Clean Dependency Injection
```go
// wire.go
func InitializeLogger() logging.Logger {
    wire.Build(logging.DefaultSet)
    return nil
}
```
Uses Google Wire to automatically generate dependency injection code for logging.

### 2. **Generated Dependencies** - Automatic Wiring
The `wire_gen.go` file (automatically generated) contains the actual implementation that wires together:
- Logger configuration
- Redactor chains  
- Output destinations
- All required dependencies

### 3. **Clean Main Function** - Simple Usage
```go
func main() {
    logger := InitializeLogger()
    
    logger.Info("Application started using Wire DI")
    
    logger.Fluent().Info().
        Str("service", "di-example").
        Str("version", "1.0.0").
        Msg("Demonstrating dependency injection")
}
```
The main function is clean and focused, with all dependency complexity hidden.

### 4. **Logger Usage Patterns** - Production Examples
```go
// Context logging
ctx := context.Background()
logger.InfoContext(ctx, "Logging with context")

// Scoped logging with static fields
userLogger := logger.WithField("user_id", 12345)
userLogger.Info("User-specific logging")
```
Shows common patterns for using injected loggers.

## Expected Output

When you run this example:

```
[INFO] Application started using Wire DI
{"level":"INFO","message":"Demonstrating dependency injection","service":"di-example","version":"1.0.0","timestamp":"2025-11-12T..."}
[INFO] Logging with context
{"level":"INFO","message":"User-specific logging","user_id":12345,"timestamp":"2025-11-12T..."}
```

## File Structure

```
di/
├── main.go          # Application entry point
├── wire.go          # Wire dependency definitions  
└── wire_gen.go      # Generated dependency injection code (auto-generated)
```

## Running the Example

```bash
cd examples/di
go run main.go
```

### Prerequisites

This example uses Google Wire for code generation. The generated code (`wire_gen.go`) is already included, but if you want to regenerate it:

```bash
# Install wire (if not already installed)
go install github.com/google/wire/cmd/wire@latest

# Generate dependency injection code
cd examples/di
wire
```

## Benefits of Dependency Injection

### 1. **Testability**
```go
// In tests, you can inject mock loggers
func TestMyFunction(t *testing.T) {
    mockLogger := &MockLogger{}
    myFunction(mockLogger) // Easy to test
}
```

### 2. **Configuration Flexibility**
Wire makes it easy to swap configurations without changing application code:
- Development vs Production loggers
- Different output destinations
- Various logging levels

### 3. **Clean Architecture**
- Dependencies are explicit and managed centrally
- Business logic doesn't depend on concrete logging implementations
- Easy to modify logging behavior without touching application code

## Integration with Larger Applications

In larger applications, you would typically:

### 1. **Define Provider Sets**
```go
// providers.go
var ApplicationSet = wire.NewSet(
    logging.DefaultSet,
    DatabaseSet,
    HTTPServerSet,
    // ... other provider sets
)
```

### 2. **Inject into Services**
```go
type UserService struct {
    logger logging.Logger
    db     Database
}

func NewUserService(logger logging.Logger, db Database) *UserService {
    return &UserService{logger: logger, db: db}
}
```

### 3. **Wire Everything Together**
```go
// wire.go
func InitializeApplication() (*Application, error) {
    wire.Build(ApplicationSet)
    return nil, nil
}
```

## Available Provider Sets

The go-logging library provides these Wire provider sets:

| Provider Set | What It Provides | Use Case |
|-------------|------------------|----------|
| `logging.DefaultSet` | Default logger with standard configuration | Most applications |
| `logging.ProviderSet` | Core logging providers | Custom configurations |

## Testing with Dependency Injection

Wire makes testing easier by allowing mock injection:

```go
// test_wire.go (test build tag)
//go:build wireinject && test

func InitializeTestLogger() logging.Logger {
    wire.Build(
        MockLoggerSet,  // Use mocks instead of real logger
    )
    return nil
}
```

## Next Steps

After running this example:

1. **Learn Wire**: Check out [Google Wire documentation](https://github.com/google/wire)
2. **Explore Providers**: Look at `pkg/logging/providers.go` for available provider sets
3. **Integration**: Try [`new-architecture/`](../new-architecture/) for advanced architecture patterns
4. **Testing**: Create your own wire providers for testing scenarios

## Best Practices Demonstrated

1. **Separation of Concerns**: Dependency injection logic is separate from business logic
2. **Testability**: Dependencies can be easily mocked for testing
3. **Configuration Management**: Easy to change logging behavior centrally
4. **Code Generation**: Let Wire handle the complex dependency wiring
5. **Clean Interfaces**: Application code depends on interfaces, not concrete types

## When to Use Dependency Injection

Consider using DI with Wire when:

- **Large Applications**: Multiple services need loggers
- **Testing**: You need to mock dependencies frequently
- **Multiple Environments**: Different configurations for dev/prod
- **Clean Architecture**: You want explicit dependency management
- **Team Development**: Multiple developers need consistent patterns
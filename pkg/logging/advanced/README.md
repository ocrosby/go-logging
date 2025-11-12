# Advanced Logging Features

This package contains advanced logging features for users who need more control over their logging setup.

## Features

- **Custom Handlers**: Multi-handler, conditional, buffered, async, and rotating handlers
- **Handler Builders**: Fluent API for building complex handler configurations  
- **Middleware**: Custom middleware for handler chains
- **Handler Composition**: Tools for combining and composing handlers
- **Advanced Configuration**: Fine-grained control over logging behavior

## When to Use

Use these advanced features when:
- You need custom output destinations
- You want to filter logs based on complex conditions
- You need buffering or async processing
- You want to build middleware chains
- Simple factory functions don't meet your needs

## Getting Started

Most users should start with the simple factory functions in the main logging package:

```go
import "github.com/ocrosby/go-logging/pkg/logging"

// Simple case - this is usually what you want
logger := logging.NewSimple()
logger.Info("Hello world")

// Or with JSON formatting
logger := logging.NewEasyJSON()
logger.Info("Hello world")
```

Only import this advanced package if you need the extra power:

```go
import "github.com/ocrosby/go-logging/pkg/logging/advanced"
```

## Examples

See the `examples/custom-handlers/` directory for detailed examples of using these advanced features.
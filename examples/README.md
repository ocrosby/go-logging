# Go-Logging Examples

This directory contains a comprehensive collection of examples demonstrating various features and use cases of the go-logging library. The examples are organized by complexity and use case to help you find the right pattern for your needs.

## 🚀 Getting Started (New Users Start Here!)

### [`basic/`](basic/) - Basic Logging Usage
**Perfect for first-time users**
- Demonstrates the simplest way to get started with logging
- Shows simple factory functions: `NewSimple()`, `NewEasyJSON()` 
- Covers basic structured logging and different log levels
- **Best for**: First-time users, quick prototypes, simple applications

### [`simple/`](simple/) - Simplified Interface Examples
A collection of examples showing the simplified logging interface:

#### [`simple/quick-start/`](simple/quick-start/) - Zero Configuration Quick Start
- Absolute minimum setup required
- One-line logger creation with sensible defaults
- **Best for**: Getting started immediately, proof-of-concepts

#### [`simple/structured/`](simple/structured/) - Simple Structured Logging  
- Direct key-value pair logging without builders
- JSON formatting for structured data
- Static fields for service context
- **Best for**: Microservices, structured data logging, production apps

#### [`simple/configuration/`](simple/configuration/) - Configuration Examples
- Builder pattern with level shortcuts (`.Debug()`, `.Info()`, etc.)
- Environment-based configuration
- File output configuration
- **Best for**: Applications needing moderate configuration complexity

#### [`simple/context-logging/`](simple/context-logging/) - Context and Tracing
- Request tracing with trace IDs
- Context propagation examples
- **Best for**: Web applications, microservices with tracing

#### [`simple/error-handling/`](simple/error-handling/) - Error Handling Patterns
- Best practices for error logging
- Structured error information
- **Best for**: Production error handling and debugging

#### [`simple/middleware/`](simple/middleware/) - HTTP Middleware
- HTTP request/response logging
- Automatic trace ID generation
- **Best for**: Web servers, API services

#### [`simple/async/`](simple/async/) - Asynchronous Logging
- High-performance async logging patterns
- Non-blocking log operations
- **Best for**: High-throughput applications

## 🎛️ Advanced Configuration

### [`yaml-config/`](yaml-config/) - YAML-Based Configuration
**Most powerful configuration method**
- Complete YAML schema demonstration
- Preset configurations (development, production, debug, minimal, structured)
- File output with automatic directory creation
- Environment-based configuration loading
- Security redaction patterns
- **Best for**: Complex applications, multiple environments, enterprise setups

### [`fluent/`](fluent/) - Fluent Interface
- Method chaining for expressive logging
- Flexible field attachment
- Different data types (strings, integers, errors)
- **Best for**: Complex log entries, readable code, method chaining fans

## 🌐 Integration Examples

### [`http-server/`](http-server/) - HTTP Server Integration
- Complete HTTP server with logging middleware
- Request/response logging
- Automatic trace ID generation and propagation
- Context-aware logging throughout request lifecycle
- **Best for**: Web applications, REST APIs, HTTP services

### [`slog/`](slog/) - Go's slog Integration
- Using Go's standard `log/slog` as the backend
- Custom slog handlers (JSON, text)
- Third-party handler integration (zerolog, zap examples)
- Performance optimizations with slog
- **Best for**: Go 1.21+ applications, performance-critical logging, slog ecosystem integration

### [`di/`](di/) - Dependency Injection
- Google Wire integration for dependency injection
- Logger as an injected dependency
- Clean architecture patterns
- **Best for**: Large applications, clean architecture, testable code

## 🔧 Advanced Features

### [`custom-handlers/`](custom-handlers/) - Custom Handler Development
- Creating custom log handlers
- Handler composition and middleware
- Multi-output logging
- Buffered and async handlers
- **Best for**: Custom output destinations, complex routing, performance optimization

### [`new-architecture/`](new-architecture/) - Unified Architecture Showcase
- Demonstrates the new unified logger architecture
- Backend switching (standard vs slog)
- Handler composition examples
- **Best for**: Understanding the library architecture, advanced customization

## 📋 Quick Reference Guide

### For Different Use Cases:

| Use Case | Recommended Example | Key Features |
|----------|-------------------|--------------|
| **First time using the library** | [`basic/`](basic/) | Simple setup, multiple patterns |
| **Simple web application** | [`simple/quick-start/`](simple/quick-start/) + [`simple/middleware/`](simple/middleware/) | Zero config + HTTP middleware |
| **Microservice** | [`yaml-config/`](yaml-config/) | Production preset, structured logging |
| **High-performance application** | [`slog/`](slog/) + [`simple/async/`](simple/async/) | slog backend + async processing |
| **Complex enterprise application** | [`yaml-config/`](yaml-config/) + [`di/`](di/) | YAML config + dependency injection |
| **API with tracing** | [`simple/context-logging/`](simple/context-logging/) | Request tracing, context propagation |
| **Custom logging requirements** | [`custom-handlers/`](custom-handlers/) | Custom handlers, advanced routing |

### For Different Complexity Levels:

| Level | Examples | When to Use |
|-------|----------|-------------|
| **Beginner** | `basic/`, `simple/quick-start/` | Learning the library, simple apps |
| **Intermediate** | `simple/*`, `fluent/`, `yaml-config/` | Production apps, structured logging |
| **Advanced** | `custom-handlers/`, `di/`, `new-architecture/` | Custom requirements, large systems |

## 🏃‍♂️ Running the Examples

Each example is a standalone Go program. To run any example:

```bash
# Run the basic example
cd examples/basic
go run main.go

# Run the YAML configuration example
cd examples/yaml-config
go run main.go

# Run any simple example
cd examples/simple/structured
go run main.go
```

### Prerequisites

- Go 1.19 or higher
- Internet connection (for downloading dependencies)

Some examples may require additional setup:
- **YAML examples**: Config files are provided
- **HTTP examples**: Will start a local server
- **DI examples**: Uses code generation (already included)

## 📖 Learning Path

### Recommended progression for new users:

1. **Start**: [`basic/`](basic/) - Learn the fundamentals
2. **Explore**: [`simple/quick-start/`](simple/quick-start/) - See the simplest approach
3. **Structure**: [`simple/structured/`](simple/structured/) - Add structured logging
4. **Configure**: [`yaml-config/`](yaml-config/) - Use powerful YAML configuration
5. **Integrate**: [`http-server/`](http-server/) or [`simple/middleware/`](simple/middleware/) - Add to web apps
6. **Optimize**: [`slog/`](slog/) - Use high-performance backend
7. **Customize**: [`custom-handlers/`](custom-handlers/) - Advanced customization

## 🎯 Best Practices Demonstrated

These examples showcase logging best practices:

- **Security**: Sensitive data redaction patterns
- **Performance**: Async processing, slog backend usage
- **Maintainability**: YAML configuration, dependency injection
- **Observability**: Request tracing, structured logging
- **Flexibility**: Multiple output formats and destinations
- **Testing**: Mock-friendly interfaces and dependency injection

## 🆘 Need Help?

- **Start simple**: Begin with [`basic/`](basic/) or [`simple/quick-start/`](simple/quick-start/)
- **Check the main README**: [../README.md](../README.md) for comprehensive documentation
- **YAML configuration**: See [yaml-config/README.md](yaml-config/README.md) for detailed YAML guide
- **Advanced features**: Check [../docs/](../docs/) for detailed guides

Each example includes its own README.md or comprehensive comments explaining the specific features demonstrated.

## 🔄 Example Status

All examples are:
- ✅ **Tested**: Run in CI/CD pipeline
- ✅ **Up-to-date**: Use latest library features
- ✅ **Documented**: Include explanatory comments
- ✅ **Self-contained**: Can be run independently

Choose the example that best matches your use case and start coding! 🚀
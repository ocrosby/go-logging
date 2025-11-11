# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Planned
- Syslog support
- File rotation support
- Performance benchmarks
- OpenTelemetry integration

## [1.0.0] - 2025-01-11

### Added
- Initial release of go-logging
- Core logging functionality with multiple log levels (TRACE, DEBUG, INFO, WARN, ERROR, CRITICAL)
- Fluent interface for expressive logging
- Request tracing support with trace IDs, request IDs, and correlation IDs
- Structured logging with JSON and text output formats
- HTTP middleware for automatic request tracing
- Sensitive data redaction with built-in API key patterns
- Builder pattern for configuration
- Environment variable configuration support
- Thread-safe concurrent logging
- Comprehensive test suite with high coverage
- Three example applications demonstrating usage patterns
- Full API documentation

### Features

#### Logger Interface
- `Trace()`, `Debug()`, `Info()`, `Warn()`, `Error()`, `Critical()` methods
- Context-aware logging with `*Context()` methods
- `WithField()` and `WithFields()` for structured logging
- `Fluent()` for fluent interface access
- `IsLevelEnabled()` for conditional logging
- `SetLevel()` and `GetLevel()` for dynamic level changes

#### Fluent Interface
- Chainable methods: `Str()`, `Int()`, `Int64()`, `Bool()`, `Err()`
- Context support with `Ctx()`
- Trace ID support with `TraceID()`
- Field management with `Field()` and `Fields()`
- Message output with `Msg()` and `Msgf()`

#### Configuration
- Builder pattern with `NewConfig()`
- Support for log levels, output formats, and destinations
- File and time inclusion options
- Custom redaction patterns
- Static field support
- Environment variable loading

#### Request Tracing
- `NewTraceID()` for generating unique trace identifiers
- Context utilities: `WithTraceID()`, `WithRequestID()`, `WithCorrelationID()`
- Context retrievers: `GetTraceID()`, `GetRequestID()`, `GetCorrelationID()`
- `NewContextWithTrace()` for quick context creation

#### HTTP Support
- `TracingMiddleware()` for automatic request tracing
- `RequestLogger()` for request-level logging
- HTTP header handling (X-Trace-ID, X-Request-ID, X-Correlation-ID)
- Request and response logging utilities

#### Data Protection
- Built-in API key redaction
- Custom regex pattern support
- `RedactedURL()` utility function
- Redactor chain for multiple patterns

### Testing
- Comprehensive unit tests for all components
- Test coverage for level management, fluent interface, tracing, and HTTP utilities
- Mock-friendly interface design

### Documentation
- Complete README with usage examples
- API reference documentation
- Three example applications
- Contributing guidelines
- MIT License

[Unreleased]: https://github.com/ocrosby/go-logging/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/ocrosby/go-logging/releases/tag/v1.0.0

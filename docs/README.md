# Documentation Index

Welcome to the comprehensive documentation for the go-logging library! This directory contains detailed guides, references, and examples to help you get the most out of the unified logging architecture.

## 📖 **Documentation Structure**

### **Getting Started**
- **[Main README](../README.md)** - Quick start, features overview, and basic usage
- **[Migration Guide](MIGRATION.md)** - Upgrade from older versions with zero breaking changes

### **Core Documentation**
- **[Architecture Guide](ARCHITECTURE.md)** - Deep dive into unified architecture, design patterns, and implementation details
- **[API Reference](API_REFERENCE.md)** - Complete API documentation with all interfaces, functions, and examples
- **[Examples Guide](EXAMPLES.md)** - Comprehensive examples covering all use cases and advanced patterns

### **Advanced Topics**
- **[Advanced Features](ADVANCED_FEATURES.md)** - Async processing, handler composition, middleware, and performance optimization
- **[Slog Integration](SLOG_INTEGRATION.md)** - Complete guide to slog backend integration and custom handlers

### **Project Information**
- **[Improvements Summary](IMPROVEMENTS_SUMMARY.md)** - Overview of architectural improvements and consolidation benefits

## 🚀 **Quick Navigation**

### **New to go-logging?**
1. Start with the [Main README](../README.md) for basic usage
2. Check out [Examples Guide](EXAMPLES.md) for practical patterns  
3. Review [Architecture Guide](ARCHITECTURE.md) to understand the design

### **Upgrading from older version?**
1. Read the [Migration Guide](MIGRATION.md) - your code keeps working!
2. Learn about new features in [Improvements Summary](IMPROVEMENTS_SUMMARY.md)
3. Explore [Advanced Features](ADVANCED_FEATURES.md) for new capabilities

### **Looking for specific information?**
- **API details**: [API Reference](API_REFERENCE.md)
- **Code examples**: [Examples Guide](EXAMPLES.md)  
- **slog integration**: [Slog Integration](SLOG_INTEGRATION.md)
- **Performance tips**: [Advanced Features](ADVANCED_FEATURES.md#performance-optimization)
- **Testing strategies**: [Examples Guide](EXAMPLES.md#testing-examples)

### **Want to contribute?**
- [Contributing Guidelines](../CONTRIBUTING.md)
- [Architecture Guide](ARCHITECTURE.md) - Understand the design principles
- [Advanced Features](ADVANCED_FEATURES.md) - See extension points

## ✨ **Key Features Covered**

### **Unified Interface**
- Single Logger interface with all methods included
- No more type assertions or interface checking
- Consistent API across all logger types

### **Dual Backend Support**  
- Transparent switching between standard Go logging and slog
- Automatic backend selection based on configuration
- Full slog handler compatibility

### **Advanced Async Processing**
- Generic AsyncWorker pattern for high-performance logging
- Non-blocking async outputs and handlers
- Proper graceful shutdown handling

### **Handler System**
- Comprehensive handler composition and middleware
- Built-in handlers for common patterns (multi-output, buffering, routing)
- Easy custom handler creation

### **Performance Optimization**
- Memory-efficient field management
- Optimized level checking with backend delegation
- High-throughput async patterns

### **Developer Experience**
- Fluent interface available on all loggers
- Rich context support with automatic propagation
- Comprehensive testing utilities and mocks

## 📋 **Documentation Standards**

All documentation in this directory follows these principles:

- **Complete Examples**: Every feature has working code examples
- **Real-world Scenarios**: Examples reflect actual use cases
- **Performance Aware**: Performance implications are clearly stated
- **Migration Friendly**: Backward compatibility is preserved and documented
- **Testing Focused**: Testing strategies are included for all patterns

## 🔗 **External Resources**

- **[Go slog Documentation](https://pkg.go.dev/log/slog)** - Official slog package docs
- **[Go Context Package](https://pkg.go.dev/context)** - Context propagation patterns
- **[Structured Logging Best Practices](https://engineering.grab.com/structured-logging)** - Industry best practices

## 📞 **Need Help?**

- **Issues**: [GitHub Issues](https://github.com/ocrosby/go-logging/issues)
- **Discussions**: [GitHub Discussions](https://github.com/ocrosby/go-logging/discussions)
- **Examples**: Check the [`examples/`](../examples/) directory

---

**Happy Logging!** 🪵

The go-logging library provides a powerful, unified interface for all your logging needs while maintaining the simplicity and performance Go developers expect.
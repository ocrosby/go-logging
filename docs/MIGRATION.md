# Migration Guide

This guide helps you migrate to the latest version of go-logging with its unified architecture and consolidated interfaces.

## Table of Contents

- [Overview](#overview)
- [What Changed](#what-changed)
- [Backward Compatibility](#backward-compatibility)
- [Migration Steps](#migration-steps)
- [Interface Changes](#interface-changes)
- [Configuration Updates](#configuration-updates)
- [Common Migration Patterns](#common-migration-patterns)
- [Breaking Changes](#breaking-changes)
- [FAQ](#faq)

## Overview

The latest version of go-logging introduces a **unified architecture** that consolidates multiple interfaces into a single, comprehensive Logger interface. This change significantly simplifies the API while maintaining full backward compatibility.

### Key Benefits of Migrating

✅ **Simplified API**: No more type assertions - all methods available on every logger  
✅ **Better Performance**: Unified dispatch system and async processing improvements  
✅ **Enhanced slog Integration**: Better handler support and level delegation  
✅ **Cleaner Code**: Eliminate interface checking and type casting  
✅ **Future-Proof**: Built on modern architecture patterns  

## What Changed

### Before: Multiple Interfaces
```go
// Old way - multiple interfaces to check
var logger logging.Logger = getLogger()

// Need type assertions for different capabilities
if ll, ok := logger.(logging.LevelLogger); ok {
    ll.Info("Using level logger")
}

if cl, ok := logger.(logging.ContextLogger); ok {
    cl.InfoContext(ctx, "Using context logger") 
}

if fl, ok := logger.(logging.FluentCapable); ok {
    fl.Fluent().Info().Msg("Using fluent interface")
}

if configurable, ok := logger.(logging.ConfigurableLogger); ok {
    configurable.SetLevel(logging.DebugLevel)
}
```

### After: Unified Interface
```go
// New way - everything available directly
var logger logging.Logger = getLogger()

// All methods available without type assertions
logger.Info("Direct level method")
logger.InfoContext(ctx, "Direct context method")
logger.Fluent().Info().Msg("Direct fluent interface")
logger.SetLevel(logging.DebugLevel)
```

## Backward Compatibility

**⚠️ Good News**: Your existing code continues to work unchanged! The library maintains full backward compatibility.

### Legacy APIs Still Work

```go
// These patterns continue to work exactly as before
config := logging.NewConfig().WithLevel(logging.InfoLevel).Build()
logger := logging.NewStandardLogger(config, redactorChain)
logger.Log(logging.InfoLevel, "Still works")

// Type assertions still work (but are no longer needed)
if ll, ok := logger.(logging.LevelLogger); ok {
    ll.Info("This still works")
}
```

### Gradual Migration

You can migrate gradually:
1. Update to the new version (existing code keeps working)
2. Remove type assertions where convenient
3. Adopt new patterns for new code
4. Eventually remove all legacy patterns

## Migration Steps

### Step 1: Update Dependencies

```bash
go get github.com/ocrosby/go-logging@latest
go mod tidy
```

### Step 2: Remove Type Assertions (Optional but Recommended)

**Before:**
```go
func logUserAction(logger logging.Logger, action string) {
    if ll, ok := logger.(logging.LevelLogger); ok {
        ll.Info("User action: %s", action)
    } else {
        logger.Log(logging.InfoLevel, "User action: %s", action)
    }
}
```

**After:**
```go
func logUserAction(logger logging.Logger, action string) {
    // Direct method call - works with any logger
    logger.Info("User action: %s", action)
}
```

### Step 3: Simplify Context Logging

**Before:**
```go
func logWithContext(logger logging.Logger, ctx context.Context, msg string) {
    if cl, ok := logger.(logging.ContextLogger); ok {
        cl.InfoContext(ctx, msg)
    } else {
        logger.LogContext(ctx, logging.InfoLevel, msg)
    }
}
```

**After:**
```go
func logWithContext(logger logging.Logger, ctx context.Context, msg string) {
    logger.InfoContext(ctx, msg)  // Always available
}
```

### Step 4: Simplify Fluent Interface Usage

**Before:**
```go
func logWithFluent(logger logging.Logger, userID int) {
    if fl, ok := logger.(logging.FluentCapable); ok {
        fl.Fluent().Info().Int("user_id", userID).Msg("User logged in")
    } else {
        logger.WithField("user_id", userID).Info("User logged in")
    }
}
```

**After:**
```go
func logWithFluent(logger logging.Logger, userID int) {
    logger.Fluent().Info().Int("user_id", userID).Msg("User logged in")
}
```

### Step 5: Update Configuration (Optional)

**Legacy (continues to work):**
```go
config := logging.NewConfig().
    WithLevel(logging.InfoLevel).
    WithJSONFormat().
    Build()

logger := logging.NewStandardLogger(config, redactorChain)
```

**Modern approach:**
```go
logger := logging.NewWithLevel(logging.InfoLevel)
// or
config := logging.NewLoggerConfig().
    WithCore(logging.NewCoreConfig().WithLevel(logging.InfoLevel).Build()).
    WithFormatter(logging.NewFormatterConfig().WithJSONFormat().Build()).
    Build()

logger := logging.NewWithLoggerConfig(config)
```

## Interface Changes

### Logger Interface Evolution

The `Logger` interface now includes all methods that were previously spread across multiple interfaces:

```go
type Logger interface {
    // Core methods (always existed)
    Log(level Level, msg string, args ...interface{})
    LogContext(ctx context.Context, level Level, msg string, args ...interface{})
    WithField(key string, value interface{}) Logger
    WithFields(fields map[string]interface{}) Logger
    IsLevelEnabled(level Level) bool

    // Level methods (previously in LevelLogger)
    Trace(msg string, args ...interface{})
    Debug(msg string, args ...interface{})
    Info(msg string, args ...interface{})
    Warn(msg string, args ...interface{})
    Error(msg string, args ...interface{})
    Critical(msg string, args ...interface{})

    // Context methods (previously in ContextLogger)  
    TraceContext(ctx context.Context, msg string, args ...interface{})
    DebugContext(ctx context.Context, msg string, args ...interface{})
    InfoContext(ctx context.Context, msg string, args ...interface{})
    WarnContext(ctx context.Context, msg string, args ...interface{})
    ErrorContext(ctx context.Context, msg string, args ...interface{})
    CriticalContext(ctx context.Context, msg string, args ...interface{})

    // Fluent interface (previously in FluentCapable)
    Fluent() FluentLogger

    // Configuration (previously in ConfigurableLogger)
    SetLevel(level Level)
    GetLevel() Level
}
```

### Removed Interfaces

These interfaces still exist for backward compatibility but are no longer needed:

- ✅ `LevelLogger` - Methods now in main `Logger` interface
- ✅ `ContextLogger` - Methods now in main `Logger` interface  
- ✅ `FluentCapable` - Method now in main `Logger` interface
- ✅ `ConfigurableLogger` - Methods now in main `Logger` interface

## Configuration Updates

### New Structured Configuration

The new configuration system separates concerns:

```go
// Modern structured approach
config := logging.NewLoggerConfig().
    WithCore(
        logging.NewCoreConfig().
            WithLevel(logging.DebugLevel).
            WithStaticField("service", "api").
            Build(),
    ).
    WithFormatter(
        logging.NewFormatterConfig().
            WithJSONFormat().
            IncludeFile(true).
            AddRedactPattern(`password=\w+`).
            Build(),
    ).
    WithOutput(
        logging.NewOutputConfig().
            WithWriter(os.Stdout).
            Build(),
    ).
    Build()
```

### Simplified Factory Functions

New factory functions reduce boilerplate:

```go
// Simple level-based creation
logger := logging.NewWithLevel(logging.InfoLevel)

// Slog integration
logger := logging.NewSlogJSONLogger(logging.DebugLevel)
logger := logging.NewWithHandler(customHandler)

// Environment-based
logger := logging.NewFromEnvironment()
```

## Common Migration Patterns

### Pattern 1: HTTP Middleware

**Before:**
```go
func loggingMiddleware(logger logging.Logger) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            if ll, ok := logger.(logging.LevelLogger); ok {
                ll.Info("Request started: %s %s", r.Method, r.URL.Path)
            }
            
            next.ServeHTTP(w, r)
            
            if cl, ok := logger.(logging.ContextLogger); ok {
                cl.InfoContext(r.Context(), "Request completed")
            }
        })
    }
}
```

**After:**
```go
func loggingMiddleware(logger logging.Logger) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            logger.Info("Request started: %s %s", r.Method, r.URL.Path)
            next.ServeHTTP(w, r)
            logger.InfoContext(r.Context(), "Request completed")
        })
    }
}
```

### Pattern 2: Service Layer Logging

**Before:**
```go
type UserService struct {
    logger logging.Logger
}

func (s *UserService) CreateUser(ctx context.Context, user *User) error {
    // Complex interface checking
    if cl, ok := s.logger.(logging.ContextLogger); ok {
        cl.InfoContext(ctx, "Creating user: %s", user.Email)
    }
    
    if err := s.validateUser(user); err != nil {
        if fl, ok := s.logger.(logging.FluentCapable); ok {
            fl.Fluent().Error().Err(err).Str("email", user.Email).Msg("Validation failed")
        }
        return err
    }
    
    return nil
}
```

**After:**
```go
type UserService struct {
    logger logging.Logger
}

func (s *UserService) CreateUser(ctx context.Context, user *User) error {
    s.logger.InfoContext(ctx, "Creating user: %s", user.Email)
    
    if err := s.validateUser(user); err != nil {
        s.logger.Fluent().Error().
            Err(err).
            Str("email", user.Email).
            Msg("Validation failed")
        return err
    }
    
    return nil
}
```

### Pattern 3: Dynamic Configuration

**Before:**
```go
func adjustLogLevel(logger logging.Logger, newLevel logging.Level) {
    if configurable, ok := logger.(logging.ConfigurableLogger); ok {
        currentLevel := configurable.GetLevel()
        fmt.Printf("Changing level from %v to %v\n", currentLevel, newLevel)
        configurable.SetLevel(newLevel)
    } else {
        fmt.Println("Logger doesn't support level changes")
    }
}
```

**After:**
```go
func adjustLogLevel(logger logging.Logger, newLevel logging.Level) {
    currentLevel := logger.GetLevel()
    fmt.Printf("Changing level from %v to %v\n", currentLevel, newLevel)
    logger.SetLevel(newLevel)
}
```

## Breaking Changes

### None! 

There are **no breaking changes** in this release. All existing code continues to work exactly as before.

### Deprecated (but still functional)

These patterns are deprecated but continue to work:

- Type assertions for logger capabilities
- Individual interface types (`LevelLogger`, etc.)
- Old configuration builder methods

## FAQ

### Q: Do I need to change my existing code?

**A: No**, your existing code will continue to work unchanged. The new unified interface is fully backward compatible.

### Q: Should I migrate to the new patterns?

**A: Recommended but not required**. The new patterns provide:
- Cleaner, simpler code
- Better performance
- Future-proofing
- Enhanced IDE support

### Q: What about performance?

**A: Performance is improved**. The unified interface eliminates type assertions and provides more efficient dispatch.

### Q: Can I mix old and new patterns?

**A: Yes**, you can gradually migrate. Old and new patterns can coexist in the same codebase.

### Q: Are there any runtime differences?

**A: Only improvements**. Better level checking, improved context handling, and enhanced slog integration.

### Q: How do I test the migration?

**A: Your existing tests should pass unchanged**. New code can use the simplified patterns.

### Q: What about third-party integrations?

**A: They continue to work**. Any code expecting the old interfaces will work exactly as before.

## Migration Checklist

- [ ] Update dependency to latest version
- [ ] Verify all tests pass (they should!)
- [ ] Identify areas with type assertions
- [ ] Remove unnecessary type assertions
- [ ] Simplify logger interface usage
- [ ] Consider adopting new configuration patterns
- [ ] Update documentation/examples
- [ ] Train team on new simplified patterns

## Getting Help

If you encounter issues during migration:

1. **Check the examples**: See `examples/` directory for patterns
2. **Review the tests**: Test files show both old and new patterns  
3. **Create an issue**: [GitHub Issues](https://github.com/ocrosby/go-logging/issues)
4. **Ask questions**: [GitHub Discussions](https://github.com/ocrosby/go-logging/discussions)

---

The migration to the unified interface represents a significant improvement in the library's usability while maintaining complete backward compatibility. Take your time with the migration and enjoy the cleaner, more powerful API!
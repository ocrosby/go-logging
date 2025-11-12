package logging

import (
	"log/slog"
	"os"
	"sync"
)

var (
	defaultLogger Logger
	once          sync.Once
)

func GetDefaultLogger() Logger {
	once.Do(func() {
		config := ProvideConfig()
		redactorChain := ProvideRedactorChain(config)
		defaultLogger = ProvideLogger(config, redactorChain)
	})
	return defaultLogger
}

func SetDefaultLogger(logger Logger) {
	defaultLogger = logger
}

func New(options ...func(*ConfigBuilder)) Logger {
	builder := NewConfig()

	for _, option := range options {
		option(builder)
	}

	config := builder.Build()
	redactorChain := ProvideRedactorChain(config)
	return ProvideLogger(config, redactorChain)
}

// NewWithLoggerConfig creates a new logger using the new configuration structure.
func NewWithLoggerConfig(config *LoggerConfig) Logger {
	if config == nil {
		config = NewLoggerConfig().Build()
	}
	redactorChain := ProvideRedactorChainFromLoggerConfig(config)
	return ProvideLoggerFromConfig(config, redactorChain)
}

func NewFromEnvironment() Logger {
	config := NewConfig().
		FromEnvironment().
		Build()
	redactorChain := ProvideRedactorChain(config)
	return ProvideLogger(config, redactorChain)
}

func NewWithLevel(level Level) Logger {
	config := NewConfig().
		WithLevel(level).
		Build()
	redactorChain := ProvideRedactorChain(config)
	return ProvideLogger(config, redactorChain)
}

func NewWithLevelString(level string) Logger {
	config := NewConfig().
		WithLevelString(level).
		Build()
	redactorChain := ProvideRedactorChain(config)
	return ProvideLogger(config, redactorChain)
}

func NewJSONLogger(level Level) Logger {
	config := NewConfig().
		WithLevel(level).
		WithJSONFormat().
		Build()
	redactorChain := ProvideRedactorChain(config)
	return ProvideLogger(config, redactorChain)
}

func NewTextLogger(level Level) Logger {
	config := NewConfig().
		WithLevel(level).
		WithTextFormat().
		Build()
	redactorChain := ProvideRedactorChain(config)
	return ProvideLogger(config, redactorChain)
}

func NewWithHandler(handler slog.Handler) Logger {
	config := NewConfig().
		WithHandler(handler).
		Build()
	redactorChain := ProvideRedactorChain(config)
	return ProvideLogger(config, redactorChain)
}

func NewSlogJSONLogger(level Level) Logger {
	config := NewConfig().
		WithLevel(level).
		WithJSONFormat().
		UseSlog(true).
		Build()
	redactorChain := ProvideRedactorChain(config)
	return ProvideLogger(config, redactorChain)
}

func NewSlogTextLogger(level Level) Logger {
	config := NewConfig().
		WithLevel(level).
		WithTextFormat().
		UseSlog(true).
		Build()
	redactorChain := ProvideRedactorChain(config)
	return ProvideLogger(config, redactorChain)
}

func Trace(msg string, args ...interface{}) {
	logger := GetDefaultLogger()
	logger.Trace(msg, args...)
}

func Debug(msg string, args ...interface{}) {
	logger := GetDefaultLogger()
	logger.Debug(msg, args...)
}

func Info(msg string, args ...interface{}) {
	logger := GetDefaultLogger()
	logger.Info(msg, args...)
}

func Warn(msg string, args ...interface{}) {
	logger := GetDefaultLogger()
	logger.Warn(msg, args...)
}

func Error(msg string, args ...interface{}) {
	logger := GetDefaultLogger()
	logger.Error(msg, args...)
}

func Critical(msg string, args ...interface{}) {
	logger := GetDefaultLogger()
	logger.Critical(msg, args...)
}

func T() Logger {
	return GetDefaultLogger()
}

func D() Logger {
	return GetDefaultLogger()
}

func I() Logger {
	return GetDefaultLogger()
}

func E() Logger {
	return GetDefaultLogger()
}

func IsDebugEnabled() bool {
	return GetDefaultLogger().IsLevelEnabled(DebugLevel)
}

func IsTraceEnabled() bool {
	return GetDefaultLogger().IsLevelEnabled(TraceLevel)
}

func MustGetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		Critical("Required environment variable not set: %s", key)
		os.Exit(1)
	}
	return value
}

// NewStandardLogger creates a standard logger with the specified config and redactor chain.
// Deprecated: Use New() or NewWithLoggerConfig() for new code.
func NewStandardLogger(config *Config, redactorChain RedactorChainInterface) Logger {
	return ProvideLogger(config, redactorChain)
}

// NewSlogLogger creates a slog-based logger with the specified config and redactor chain.
// Deprecated: Use New() or NewWithLoggerConfig() for new code.
func NewSlogLogger(config *Config, redactorChain RedactorChainInterface) Logger {
	return ProvideLogger(config, redactorChain)
}

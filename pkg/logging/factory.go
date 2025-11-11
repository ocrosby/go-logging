package logging

import (
	"os"
	"sync"
)

var (
	defaultLogger Logger
	once          sync.Once
)

func GetDefaultLogger() Logger {
	once.Do(func() {
		config := NewConfig().
			FromEnvironment().
			Build()
		defaultLogger = NewStandardLogger(config)
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

	return NewStandardLogger(builder.Build())
}

func NewFromEnvironment() Logger {
	config := NewConfig().
		FromEnvironment().
		Build()
	return NewStandardLogger(config)
}

func NewWithLevel(level Level) Logger {
	config := NewConfig().
		WithLevel(level).
		Build()
	return NewStandardLogger(config)
}

func NewWithLevelString(level string) Logger {
	config := NewConfig().
		WithLevelString(level).
		Build()
	return NewStandardLogger(config)
}

func NewJSONLogger(level Level) Logger {
	config := NewConfig().
		WithLevel(level).
		WithJSONFormat().
		Build()
	return NewStandardLogger(config)
}

func NewTextLogger(level Level) Logger {
	config := NewConfig().
		WithLevel(level).
		WithTextFormat().
		Build()
	return NewStandardLogger(config)
}

func Trace(msg string, args ...interface{}) {
	GetDefaultLogger().Trace(msg, args...)
}

func Debug(msg string, args ...interface{}) {
	GetDefaultLogger().Debug(msg, args...)
}

func Info(msg string, args ...interface{}) {
	GetDefaultLogger().Info(msg, args...)
}

func Warn(msg string, args ...interface{}) {
	GetDefaultLogger().Warn(msg, args...)
}

func Error(msg string, args ...interface{}) {
	GetDefaultLogger().Error(msg, args...)
}

func Critical(msg string, args ...interface{}) {
	GetDefaultLogger().Critical(msg, args...)
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

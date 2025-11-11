package logging

import "context"

type Logger interface {
	Trace(msg string, args ...interface{})
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Critical(msg string, args ...interface{})

	TraceContext(ctx context.Context, msg string, args ...interface{})
	DebugContext(ctx context.Context, msg string, args ...interface{})
	InfoContext(ctx context.Context, msg string, args ...interface{})
	WarnContext(ctx context.Context, msg string, args ...interface{})
	ErrorContext(ctx context.Context, msg string, args ...interface{})
	CriticalContext(ctx context.Context, msg string, args ...interface{})

	WithField(key string, value interface{}) Logger
	WithFields(fields map[string]interface{}) Logger

	IsLevelEnabled(level Level) bool

	SetLevel(level Level)
	GetLevel() Level

	Fluent() FluentLogger
}

type FluentLogger interface {
	Trace() *FluentEntry
	Debug() *FluentEntry
	Info() *FluentEntry
	Warn() *FluentEntry
	Error() *FluentEntry
	Critical() *FluentEntry
}

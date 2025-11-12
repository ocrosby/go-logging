package logging

import "context"

// LevelDispatcher provides a unified way to dispatch level-specific logging calls
// to the core Log and LogContext methods. This eliminates code duplication and
// provides a centralized place for level method dispatch logic.
type LevelDispatcher struct {
	logger Logger
}

// NewLevelDispatcher creates a new dispatcher for the given logger.
func NewLevelDispatcher(logger Logger) *LevelDispatcher {
	return &LevelDispatcher{logger: logger}
}

// DispatchTrace logs at TRACE level
func (d *LevelDispatcher) DispatchTrace(msg string, args ...interface{}) {
	d.logger.Log(TraceLevel, msg, args...)
}

// DispatchDebug logs at DEBUG level
func (d *LevelDispatcher) DispatchDebug(msg string, args ...interface{}) {
	d.logger.Log(DebugLevel, msg, args...)
}

// DispatchInfo logs at INFO level
func (d *LevelDispatcher) DispatchInfo(msg string, args ...interface{}) {
	d.logger.Log(InfoLevel, msg, args...)
}

// DispatchWarn logs at WARN level
func (d *LevelDispatcher) DispatchWarn(msg string, args ...interface{}) {
	d.logger.Log(WarnLevel, msg, args...)
}

// DispatchError logs at ERROR level
func (d *LevelDispatcher) DispatchError(msg string, args ...interface{}) {
	d.logger.Log(ErrorLevel, msg, args...)
}

// DispatchCritical logs at CRITICAL level
func (d *LevelDispatcher) DispatchCritical(msg string, args ...interface{}) {
	d.logger.Log(CriticalLevel, msg, args...)
}

// DispatchTraceContext logs at TRACE level with context
func (d *LevelDispatcher) DispatchTraceContext(ctx context.Context, msg string, args ...interface{}) {
	d.logger.LogContext(ctx, TraceLevel, msg, args...)
}

// DispatchDebugContext logs at DEBUG level with context
func (d *LevelDispatcher) DispatchDebugContext(ctx context.Context, msg string, args ...interface{}) {
	d.logger.LogContext(ctx, DebugLevel, msg, args...)
}

// DispatchInfoContext logs at INFO level with context
func (d *LevelDispatcher) DispatchInfoContext(ctx context.Context, msg string, args ...interface{}) {
	d.logger.LogContext(ctx, InfoLevel, msg, args...)
}

// DispatchWarnContext logs at WARN level with context
func (d *LevelDispatcher) DispatchWarnContext(ctx context.Context, msg string, args ...interface{}) {
	d.logger.LogContext(ctx, WarnLevel, msg, args...)
}

// DispatchErrorContext logs at ERROR level with context
func (d *LevelDispatcher) DispatchErrorContext(ctx context.Context, msg string, args ...interface{}) {
	d.logger.LogContext(ctx, ErrorLevel, msg, args...)
}

// DispatchCriticalContext logs at CRITICAL level with context
func (d *LevelDispatcher) DispatchCriticalContext(ctx context.Context, msg string, args ...interface{}) {
	d.logger.LogContext(ctx, CriticalLevel, msg, args...)
}

// LoggerLevelMethods provides default implementations for all level methods
// that delegate to Log/LogContext. This can be embedded in logger implementations.
type LoggerLevelMethods struct {
	dispatcher *LevelDispatcher
}

// InitLevelMethods initializes the level methods with the core logger
func (l *LoggerLevelMethods) InitLevelMethods(coreLogger Logger) {
	l.dispatcher = NewLevelDispatcher(coreLogger)
}

func (l *LoggerLevelMethods) Trace(msg string, args ...interface{}) {
	l.dispatcher.DispatchTrace(msg, args...)
}

func (l *LoggerLevelMethods) Debug(msg string, args ...interface{}) {
	l.dispatcher.DispatchDebug(msg, args...)
}

func (l *LoggerLevelMethods) Info(msg string, args ...interface{}) {
	l.dispatcher.DispatchInfo(msg, args...)
}

func (l *LoggerLevelMethods) Warn(msg string, args ...interface{}) {
	l.dispatcher.DispatchWarn(msg, args...)
}

func (l *LoggerLevelMethods) Error(msg string, args ...interface{}) {
	l.dispatcher.DispatchError(msg, args...)
}

func (l *LoggerLevelMethods) Critical(msg string, args ...interface{}) {
	l.dispatcher.DispatchCritical(msg, args...)
}

func (l *LoggerLevelMethods) TraceContext(ctx context.Context, msg string, args ...interface{}) {
	l.dispatcher.DispatchTraceContext(ctx, msg, args...)
}

func (l *LoggerLevelMethods) DebugContext(ctx context.Context, msg string, args ...interface{}) {
	l.dispatcher.DispatchDebugContext(ctx, msg, args...)
}

func (l *LoggerLevelMethods) InfoContext(ctx context.Context, msg string, args ...interface{}) {
	l.dispatcher.DispatchInfoContext(ctx, msg, args...)
}

func (l *LoggerLevelMethods) WarnContext(ctx context.Context, msg string, args ...interface{}) {
	l.dispatcher.DispatchWarnContext(ctx, msg, args...)
}

func (l *LoggerLevelMethods) ErrorContext(ctx context.Context, msg string, args ...interface{}) {
	l.dispatcher.DispatchErrorContext(ctx, msg, args...)
}

func (l *LoggerLevelMethods) CriticalContext(ctx context.Context, msg string, args ...interface{}) {
	l.dispatcher.DispatchCriticalContext(ctx, msg, args...)
}

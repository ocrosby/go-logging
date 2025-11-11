package logging

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
)

const (
	LevelTrace    = slog.Level(-8)
	LevelCritical = slog.Level(12)
)

type slogLogger struct {
	mu            sync.RWMutex
	slog          *slog.Logger
	level         Level
	fields        []slog.Attr
	redactorChain RedactorChainInterface
}

func NewSlogLogger(handler slog.Handler, redactorChain RedactorChainInterface) Logger {
	if handler == nil {
		handler = slog.Default().Handler()
	}
	if redactorChain == nil {
		redactorChain = NewRedactorChain()
	}

	return &slogLogger{
		slog:          slog.New(handler),
		level:         InfoLevel,
		fields:        make([]slog.Attr, 0),
		redactorChain: redactorChain,
	}
}

func (sl *slogLogger) levelToSlog(level Level) slog.Level {
	switch level {
	case TraceLevel:
		return LevelTrace
	case DebugLevel:
		return slog.LevelDebug
	case InfoLevel:
		return slog.LevelInfo
	case WarnLevel:
		return slog.LevelWarn
	case ErrorLevel:
		return slog.LevelError
	case CriticalLevel:
		return LevelCritical
	default:
		return slog.LevelInfo
	}
}

func (sl *slogLogger) log(ctx context.Context, level Level, msg string, args ...interface{}) {
	sl.mu.RLock()
	defer sl.mu.RUnlock()

	if !sl.isLevelEnabledInternal(level) {
		return
	}

	slogLevel := sl.levelToSlog(level)

	if len(args) > 0 {
		msg = sl.redactorChain.Redact(formatMessage(msg, args...))
	} else {
		msg = sl.redactorChain.Redact(msg)
	}

	if ctx == nil {
		ctx = context.Background()
	}

	logAttrs := make([]slog.Attr, 0, len(sl.fields)+1)
	logAttrs = append(logAttrs, sl.fields...)

	if reqID, ok := GetRequestID(ctx); ok {
		logAttrs = append(logAttrs, slog.String("request_id", reqID))
	}

	sl.slog.LogAttrs(ctx, slogLevel, msg, logAttrs...)
}

func (sl *slogLogger) isLevelEnabledInternal(level Level) bool {
	return level >= sl.level
}

func (sl *slogLogger) Trace(msg string, args ...interface{}) {
	sl.log(context.Background(), TraceLevel, msg, args...)
}

func (sl *slogLogger) Debug(msg string, args ...interface{}) {
	sl.log(context.Background(), DebugLevel, msg, args...)
}

func (sl *slogLogger) Info(msg string, args ...interface{}) {
	sl.log(context.Background(), InfoLevel, msg, args...)
}

func (sl *slogLogger) Warn(msg string, args ...interface{}) {
	sl.log(context.Background(), WarnLevel, msg, args...)
}

func (sl *slogLogger) Error(msg string, args ...interface{}) {
	sl.log(context.Background(), ErrorLevel, msg, args...)
}

func (sl *slogLogger) Critical(msg string, args ...interface{}) {
	sl.log(context.Background(), CriticalLevel, msg, args...)
}

func (sl *slogLogger) TraceContext(ctx context.Context, msg string, args ...interface{}) {
	sl.log(ctx, TraceLevel, msg, args...)
}

func (sl *slogLogger) DebugContext(ctx context.Context, msg string, args ...interface{}) {
	sl.log(ctx, DebugLevel, msg, args...)
}

func (sl *slogLogger) InfoContext(ctx context.Context, msg string, args ...interface{}) {
	sl.log(ctx, InfoLevel, msg, args...)
}

func (sl *slogLogger) WarnContext(ctx context.Context, msg string, args ...interface{}) {
	sl.log(ctx, WarnLevel, msg, args...)
}

func (sl *slogLogger) ErrorContext(ctx context.Context, msg string, args ...interface{}) {
	sl.log(ctx, ErrorLevel, msg, args...)
}

func (sl *slogLogger) CriticalContext(ctx context.Context, msg string, args ...interface{}) {
	sl.log(ctx, CriticalLevel, msg, args...)
}

func (sl *slogLogger) WithField(key string, value interface{}) Logger {
	sl.mu.RLock()
	defer sl.mu.RUnlock()

	newFields := make([]slog.Attr, len(sl.fields), len(sl.fields)+1)
	copy(newFields, sl.fields)
	newFields = append(newFields, slog.Any(key, value))

	return &slogLogger{
		slog:          sl.slog,
		level:         sl.level,
		fields:        newFields,
		redactorChain: sl.redactorChain,
	}
}

func (sl *slogLogger) WithFields(fields map[string]interface{}) Logger {
	sl.mu.RLock()
	defer sl.mu.RUnlock()

	newFields := make([]slog.Attr, len(sl.fields), len(sl.fields)+len(fields))
	copy(newFields, sl.fields)

	for k, v := range fields {
		newFields = append(newFields, slog.Any(k, v))
	}

	return &slogLogger{
		slog:          sl.slog,
		level:         sl.level,
		fields:        newFields,
		redactorChain: sl.redactorChain,
	}
}

func (sl *slogLogger) IsLevelEnabled(level Level) bool {
	sl.mu.RLock()
	defer sl.mu.RUnlock()
	return sl.isLevelEnabledInternal(level)
}

func (sl *slogLogger) SetLevel(level Level) {
	sl.mu.Lock()
	defer sl.mu.Unlock()
	sl.level = level
}

func (sl *slogLogger) GetLevel() Level {
	sl.mu.RLock()
	defer sl.mu.RUnlock()
	return sl.level
}

func (sl *slogLogger) Fluent() FluentLogger {
	return &fluentLoggerWrapper{logger: Logger(sl)}
}

func formatMessage(msg string, args ...interface{}) string {
	if len(args) == 0 {
		return msg
	}

	for i := 0; i < len(msg); i++ {
		if msg[i] == '%' && i+1 < len(msg) {
			return formatWithPercent(msg, args...)
		}
	}

	return msg
}

func formatWithPercent(msg string, args ...interface{}) string {
	return fmt.Sprintf(msg, args...)
}

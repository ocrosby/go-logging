package logging

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"runtime"
	"sync"
	"time"
)

type standardLogger struct {
	mu            sync.RWMutex
	config        *Config
	fields        map[string]interface{}
	textLoggers   map[Level]*log.Logger
	discard       *log.Logger
	redactorChain *RedactorChain
}

func NewStandardLogger(config *Config) Logger {
	if config == nil {
		config = NewConfig().Build()
	}

	sl := &standardLogger{
		config:        config,
		fields:        make(map[string]interface{}),
		textLoggers:   make(map[Level]*log.Logger),
		discard:       log.New(io.Discard, "", 0),
		redactorChain: NewRedactorChain(config.RedactPatterns...),
	}

	if config.Format == TextFormat {
		sl.initTextLoggers()
	}

	return sl
}

func (sl *standardLogger) initTextLoggers() {
	flags := 0
	if sl.config.IncludeTime {
		flags |= log.LstdFlags | log.Lmsgprefix
	}
	if sl.config.IncludeFile {
		if sl.config.UseShortFile {
			flags |= log.Lshortfile
		} else {
			flags |= log.Llongfile
		}
	}

	for level := range levelNames {
		prefix := fmt.Sprintf("[%s] ", level)
		sl.textLoggers[level] = log.New(sl.config.Output, prefix, flags)
	}
}

func (sl *standardLogger) log(level Level, msg string, args ...interface{}) {
	sl.mu.RLock()
	defer sl.mu.RUnlock()

	if !sl.isLevelEnabledInternal(level) {
		return
	}

	message := fmt.Sprintf(msg, args...)
	message = sl.redactorChain.Redact(message)

	if sl.config.Format == JSONFormat {
		sl.logJSON(level, message, nil)
	} else {
		sl.logText(level, message)
	}
}

func (sl *standardLogger) logContext(ctx context.Context, level Level, msg string, args ...interface{}) {
	sl.mu.RLock()
	defer sl.mu.RUnlock()

	if !sl.isLevelEnabledInternal(level) {
		return
	}

	message := fmt.Sprintf(msg, args...)
	message = sl.redactorChain.Redact(message)

	if sl.config.Format == JSONFormat {
		sl.logJSON(level, message, ctx)
	} else {
		sl.logText(level, message)
	}
}

func (sl *standardLogger) logText(level Level, message string) {
	logger := sl.textLoggers[level]
	if logger == nil {
		logger = sl.discard
	}
	logger.Output(3, message)
}

func (sl *standardLogger) logJSON(level Level, message string, ctx context.Context) {
	entry := make(map[string]interface{})

	if sl.config.IncludeTime {
		entry["timestamp"] = time.Now().UTC().Format(time.RFC3339)
	}

	entry["level"] = level.String()
	entry["message"] = message

	if sl.config.IncludeFile {
		if _, file, line, ok := runtime.Caller(3); ok {
			if sl.config.UseShortFile {
				short := file
				for i := len(file) - 1; i > 0; i-- {
					if file[i] == '/' {
						short = file[i+1:]
						break
					}
				}
				file = short
			}
			entry["file"] = fmt.Sprintf("%s:%d", file, line)
		}
	}

	for k, v := range sl.config.StaticFields {
		entry[k] = v
	}

	for k, v := range sl.fields {
		entry[k] = v
	}

	if ctx != nil {
		if reqID, ok := ctx.Value("request_id").(string); ok {
			entry["request_id"] = reqID
		}
	}

	jsonBytes, err := json.Marshal(entry)
	if err != nil {
		return
	}

	fmt.Fprintln(sl.config.Output, string(jsonBytes))
}

func (sl *standardLogger) isLevelEnabledInternal(level Level) bool {
	return level >= sl.config.Level
}

func (sl *standardLogger) Trace(msg string, args ...interface{}) {
	sl.log(TraceLevel, msg, args...)
}

func (sl *standardLogger) Debug(msg string, args ...interface{}) {
	sl.log(DebugLevel, msg, args...)
}

func (sl *standardLogger) Info(msg string, args ...interface{}) {
	sl.log(InfoLevel, msg, args...)
}

func (sl *standardLogger) Warn(msg string, args ...interface{}) {
	sl.log(WarnLevel, msg, args...)
}

func (sl *standardLogger) Error(msg string, args ...interface{}) {
	sl.log(ErrorLevel, msg, args...)
}

func (sl *standardLogger) Critical(msg string, args ...interface{}) {
	sl.log(CriticalLevel, msg, args...)
}

func (sl *standardLogger) Fluent() FluentLogger {
	return &fluentLoggerWrapper{logger: sl}
}

func (sl *standardLogger) TraceContext(ctx context.Context, msg string, args ...interface{}) {
	sl.logContext(ctx, TraceLevel, msg, args...)
}

func (sl *standardLogger) DebugContext(ctx context.Context, msg string, args ...interface{}) {
	sl.logContext(ctx, DebugLevel, msg, args...)
}

func (sl *standardLogger) InfoContext(ctx context.Context, msg string, args ...interface{}) {
	sl.logContext(ctx, InfoLevel, msg, args...)
}

func (sl *standardLogger) WarnContext(ctx context.Context, msg string, args ...interface{}) {
	sl.logContext(ctx, WarnLevel, msg, args...)
}

func (sl *standardLogger) ErrorContext(ctx context.Context, msg string, args ...interface{}) {
	sl.logContext(ctx, ErrorLevel, msg, args...)
}

func (sl *standardLogger) CriticalContext(ctx context.Context, msg string, args ...interface{}) {
	sl.logContext(ctx, CriticalLevel, msg, args...)
}

func (sl *standardLogger) WithField(key string, value interface{}) Logger {
	sl.mu.RLock()
	newFields := make(map[string]interface{}, len(sl.fields)+1)
	for k, v := range sl.fields {
		newFields[k] = v
	}
	sl.mu.RUnlock()

	newFields[key] = value

	return &standardLogger{
		config:        sl.config,
		fields:        newFields,
		textLoggers:   sl.textLoggers,
		discard:       sl.discard,
		redactorChain: sl.redactorChain,
	}
}

func (sl *standardLogger) WithFields(fields map[string]interface{}) Logger {
	sl.mu.RLock()
	newFields := make(map[string]interface{}, len(sl.fields)+len(fields))
	for k, v := range sl.fields {
		newFields[k] = v
	}
	sl.mu.RUnlock()

	for k, v := range fields {
		newFields[k] = v
	}

	return &standardLogger{
		config:        sl.config,
		fields:        newFields,
		textLoggers:   sl.textLoggers,
		discard:       sl.discard,
		redactorChain: sl.redactorChain,
	}
}

func (sl *standardLogger) IsLevelEnabled(level Level) bool {
	sl.mu.RLock()
	defer sl.mu.RUnlock()
	return sl.isLevelEnabledInternal(level)
}

func (sl *standardLogger) SetLevel(level Level) {
	sl.mu.Lock()
	defer sl.mu.Unlock()
	sl.config.Level = level
}

func (sl *standardLogger) GetLevel() Level {
	sl.mu.RLock()
	defer sl.mu.RUnlock()
	return sl.config.Level
}

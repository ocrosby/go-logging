package logging

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"runtime"
	"sync"
	"time"

	"github.com/ocrosby/go-logging/pkg/logging/internal"
)

// unifiedLogger is a single implementation that provides all logger interfaces
// and can adapt between standard Go logging and slog backends.
type unifiedLogger struct {
	mu            sync.RWMutex
	config        *LoggerConfig
	fields        map[string]interface{}
	textLoggers   map[Level]*log.Logger
	slogLogger    *slog.Logger
	discard       *log.Logger
	redactorChain RedactorChainInterface
}

// NewUnifiedLogger creates a new unified logger implementation.
func NewUnifiedLogger(config *LoggerConfig, redactorChain RedactorChainInterface) Logger {
	if config == nil {
		config = NewLoggerConfig().Build()
	}
	if redactorChain == nil {
		redactorChain = NewRedactorChain()
	}

	ul := &unifiedLogger{
		config:        config,
		fields:        make(map[string]interface{}),
		textLoggers:   make(map[Level]*log.Logger),
		discard:       log.New(io.Discard, "", 0),
		redactorChain: redactorChain,
	}

	// Initialize based on configuration
	if config.UseSlog {
		handler := config.Handler
		if handler == nil {
			if config.Formatter.Format == JSONFormat {
				handler = slog.NewJSONHandler(config.Output.Writer, &slog.HandlerOptions{
					Level: ul.levelToSlog(config.Core.Level),
				})
			} else {
				handler = slog.NewTextHandler(config.Output.Writer, &slog.HandlerOptions{
					Level: ul.levelToSlog(config.Core.Level),
				})
			}
		}
		ul.slogLogger = slog.New(handler)
	} else if config.Formatter.Format == TextFormat {
		ul.initTextLoggers()
	}

	return ul
}

func (ul *unifiedLogger) initTextLoggers() {
	flags := 0
	if ul.config.Formatter.IncludeTime {
		flags |= log.LstdFlags | log.Lmsgprefix
	}
	if ul.config.Formatter.IncludeFile {
		if ul.config.Formatter.UseShortFile {
			flags |= log.Lshortfile
		} else {
			flags |= log.Llongfile
		}
	}

	for level := range levelNames {
		prefix := fmt.Sprintf("[%s] ", level)
		ul.textLoggers[level] = log.New(ul.config.Output.Writer, prefix, flags)
	}
}

func (ul *unifiedLogger) levelToSlog(level Level) slog.Level {
	levelMap := map[Level]slog.Level{
		TraceLevel:    slog.Level(-8), // Custom trace level
		DebugLevel:    slog.LevelDebug,
		InfoLevel:     slog.LevelInfo,
		WarnLevel:     slog.LevelWarn,
		ErrorLevel:    slog.LevelError,
		CriticalLevel: slog.Level(12), // Custom critical level
	}

	if slogLevel, ok := levelMap[level]; ok {
		return slogLevel
	}
	return slog.LevelInfo
}

// Core Logger interface implementation
func (ul *unifiedLogger) Log(level Level, msg string, args ...interface{}) {
	ul.LogContext(context.Background(), level, msg, args...)
}

func (ul *unifiedLogger) LogContext(ctx context.Context, level Level, msg string, args ...interface{}) {
	ul.mu.RLock()
	defer ul.mu.RUnlock()

	if !ul.isLevelEnabledInternal(level) {
		return
	}

	message := fmt.Sprintf(msg, args...)
	message = ul.redactorChain.Redact(message)

	if ul.config.UseSlog {
		ul.logSlog(ctx, level, message)
	} else if ul.config.Formatter.Format == JSONFormat {
		ul.logJSON(level, message, ctx)
	} else {
		ul.logText(level, message)
	}
}

func (ul *unifiedLogger) WithField(key string, value interface{}) Logger {
	ul.mu.RLock()
	newFields := make(map[string]interface{}, len(ul.fields)+1)
	for k, v := range ul.fields {
		newFields[k] = v
	}
	ul.mu.RUnlock()

	newFields[key] = value

	return &unifiedLogger{
		config:        ul.config,
		fields:        newFields,
		textLoggers:   ul.textLoggers,
		slogLogger:    ul.slogLogger,
		discard:       ul.discard,
		redactorChain: ul.redactorChain,
	}
}

func (ul *unifiedLogger) WithFields(fields map[string]interface{}) Logger {
	ul.mu.RLock()
	newFields := make(map[string]interface{}, len(ul.fields)+len(fields))
	for k, v := range ul.fields {
		newFields[k] = v
	}
	ul.mu.RUnlock()

	for k, v := range fields {
		newFields[k] = v
	}

	return &unifiedLogger{
		config:        ul.config,
		fields:        newFields,
		textLoggers:   ul.textLoggers,
		slogLogger:    ul.slogLogger,
		discard:       ul.discard,
		redactorChain: ul.redactorChain,
	}
}

func (ul *unifiedLogger) IsLevelEnabled(level Level) bool {
	ul.mu.RLock()
	defer ul.mu.RUnlock()
	return ul.isLevelEnabledInternal(level)
}

func (ul *unifiedLogger) isLevelEnabledInternal(level Level) bool {
	// When using slog, delegate to the slog handler for level checking
	if ul.config.UseSlog && ul.slogLogger != nil {
		return ul.slogLogger.Enabled(context.Background(), ul.levelToSlog(level))
	}
	// For standard logging, use config level
	return level >= ul.config.Core.Level
}

// LevelLogger interface implementation
func (ul *unifiedLogger) Trace(msg string, args ...interface{}) {
	ul.Log(TraceLevel, msg, args...)
}

func (ul *unifiedLogger) Debug(msg string, args ...interface{}) {
	ul.Log(DebugLevel, msg, args...)
}

func (ul *unifiedLogger) Info(msg string, args ...interface{}) {
	ul.Log(InfoLevel, msg, args...)
}

func (ul *unifiedLogger) Warn(msg string, args ...interface{}) {
	ul.Log(WarnLevel, msg, args...)
}

func (ul *unifiedLogger) Error(msg string, args ...interface{}) {
	ul.Log(ErrorLevel, msg, args...)
}

func (ul *unifiedLogger) Critical(msg string, args ...interface{}) {
	ul.Log(CriticalLevel, msg, args...)
}

// ContextLogger interface implementation
func (ul *unifiedLogger) TraceContext(ctx context.Context, msg string, args ...interface{}) {
	ul.LogContext(ctx, TraceLevel, msg, args...)
}

func (ul *unifiedLogger) DebugContext(ctx context.Context, msg string, args ...interface{}) {
	ul.LogContext(ctx, DebugLevel, msg, args...)
}

func (ul *unifiedLogger) InfoContext(ctx context.Context, msg string, args ...interface{}) {
	ul.LogContext(ctx, InfoLevel, msg, args...)
}

func (ul *unifiedLogger) WarnContext(ctx context.Context, msg string, args ...interface{}) {
	ul.LogContext(ctx, WarnLevel, msg, args...)
}

func (ul *unifiedLogger) ErrorContext(ctx context.Context, msg string, args ...interface{}) {
	ul.LogContext(ctx, ErrorLevel, msg, args...)
}

func (ul *unifiedLogger) CriticalContext(ctx context.Context, msg string, args ...interface{}) {
	ul.LogContext(ctx, CriticalLevel, msg, args...)
}

// ConfigurableLogger interface implementation
func (ul *unifiedLogger) SetLevel(level Level) {
	ul.mu.Lock()
	defer ul.mu.Unlock()
	ul.config.Core.Level = level
}

func (ul *unifiedLogger) GetLevel() Level {
	ul.mu.RLock()
	defer ul.mu.RUnlock()
	return ul.config.Core.Level
}

// FluentCapable interface implementation
func (ul *unifiedLogger) Fluent() FluentLogger {
	return &fluentLoggerWrapper{logger: ul}
}

// Internal logging methods
func (ul *unifiedLogger) logSlog(ctx context.Context, level Level, message string) {
	if ul.slogLogger == nil {
		return
	}

	slogLevel := ul.levelToSlog(level)
	logAttrs := make([]slog.Attr, 0, len(ul.fields)+len(ul.config.Core.StaticFields)+1)

	// Add static fields
	for k, v := range ul.config.Core.StaticFields {
		logAttrs = append(logAttrs, slog.Any(k, v))
	}

	// Add instance fields
	for k, v := range ul.fields {
		logAttrs = append(logAttrs, slog.Any(k, v))
	}

	// Add context fields using the correct context keys
	if requestID, ok := GetRequestID(ctx); ok && requestID != "" {
		logAttrs = append(logAttrs, slog.String("request_id", requestID))
	}
	if traceID, ok := GetTraceID(ctx); ok && traceID != "" {
		logAttrs = append(logAttrs, slog.String("trace_id", traceID))
	}
	if correlationID, ok := GetCorrelationID(ctx); ok && correlationID != "" {
		logAttrs = append(logAttrs, slog.String("correlation_id", correlationID))
	}

	ul.slogLogger.LogAttrs(ctx, slogLevel, message, logAttrs...)
}

func (ul *unifiedLogger) logText(level Level, message string) {
	logger := ul.textLoggers[level]
	if logger == nil {
		logger = ul.discard
	}
	_ = logger.Output(3, message)
}

func (ul *unifiedLogger) logJSON(level Level, message string, ctx context.Context) {
	entry := ul.createBaseEntry(level, message)
	ul.addFileInfo(entry)
	ul.addStaticFields(entry)
	ul.addInstanceFields(entry)
	ul.addContextFields(entry, ctx)
	ul.writeJSON(entry)
}

func (ul *unifiedLogger) createBaseEntry(level Level, message string) map[string]interface{} {
	entry := make(map[string]interface{})

	if ul.config.Formatter.IncludeTime {
		entry["timestamp"] = time.Now().UTC().Format(time.RFC3339)
	}

	entry["level"] = level.String()
	entry["message"] = message

	return entry
}

func (ul *unifiedLogger) addFileInfo(entry map[string]interface{}) {
	if !ul.config.Formatter.IncludeFile {
		return
	}

	if _, file, line, ok := runtime.Caller(4); ok {
		entry["file"] = ul.formatFilename(file, line)
	}
}

func (ul *unifiedLogger) formatFilename(file string, line int) string {
	return internal.FormatFilename(file, line, ul.config.Formatter.UseShortFile)
}

func (ul *unifiedLogger) addStaticFields(entry map[string]interface{}) {
	for k, v := range ul.config.Core.StaticFields {
		entry[k] = v
	}
}

func (ul *unifiedLogger) addInstanceFields(entry map[string]interface{}) {
	for k, v := range ul.fields {
		entry[k] = v
	}
}

func (ul *unifiedLogger) addContextFields(entry map[string]interface{}, ctx context.Context) {
	if requestID, ok := GetRequestID(ctx); ok && requestID != "" {
		entry["request_id"] = requestID
	}
	if traceID, ok := GetTraceID(ctx); ok && traceID != "" {
		entry["trace_id"] = traceID
	}
	if correlationID, ok := GetCorrelationID(ctx); ok && correlationID != "" {
		entry["correlation_id"] = correlationID
	}
}

func (ul *unifiedLogger) writeJSON(entry map[string]interface{}) {
	jsonBytes, err := json.Marshal(entry)
	if err != nil {
		return
	}

	fmt.Fprintln(ul.config.Output.Writer, string(jsonBytes))
}

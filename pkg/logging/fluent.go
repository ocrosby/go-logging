package logging

import "context"

type fluentLoggerWrapper struct {
	logger *standardLogger
}

func (w *fluentLoggerWrapper) Trace() *FluentEntry {
	return &FluentEntry{
		logger: w.logger,
		level:  TraceLevel,
		fields: make(map[string]interface{}),
	}
}

func (w *fluentLoggerWrapper) Debug() *FluentEntry {
	return &FluentEntry{
		logger: w.logger,
		level:  DebugLevel,
		fields: make(map[string]interface{}),
	}
}

func (w *fluentLoggerWrapper) Info() *FluentEntry {
	return &FluentEntry{
		logger: w.logger,
		level:  InfoLevel,
		fields: make(map[string]interface{}),
	}
}

func (w *fluentLoggerWrapper) Warn() *FluentEntry {
	return &FluentEntry{
		logger: w.logger,
		level:  WarnLevel,
		fields: make(map[string]interface{}),
	}
}

func (w *fluentLoggerWrapper) Error() *FluentEntry {
	return &FluentEntry{
		logger: w.logger,
		level:  ErrorLevel,
		fields: make(map[string]interface{}),
	}
}

func (w *fluentLoggerWrapper) Critical() *FluentEntry {
	return &FluentEntry{
		logger: w.logger,
		level:  CriticalLevel,
		fields: make(map[string]interface{}),
	}
}

type FluentEntry struct {
	logger  Logger
	level   Level
	fields  map[string]interface{}
	ctx     context.Context
	traceID string
}

func (e *FluentEntry) Field(key string, value interface{}) *FluentEntry {
	e.fields[key] = value
	return e
}

func (e *FluentEntry) Fields(fields map[string]interface{}) *FluentEntry {
	for k, v := range fields {
		e.fields[k] = v
	}
	return e
}

func (e *FluentEntry) Str(key, value string) *FluentEntry {
	e.fields[key] = value
	return e
}

func (e *FluentEntry) Int(key string, value int) *FluentEntry {
	e.fields[key] = value
	return e
}

func (e *FluentEntry) Int64(key string, value int64) *FluentEntry {
	e.fields[key] = value
	return e
}

func (e *FluentEntry) Bool(key string, value bool) *FluentEntry {
	e.fields[key] = value
	return e
}

func (e *FluentEntry) Err(err error) *FluentEntry {
	if err != nil {
		e.fields["error"] = err.Error()
	}
	return e
}

func (e *FluentEntry) TraceID(id string) *FluentEntry {
	e.traceID = id
	e.fields["trace_id"] = id
	return e
}

func (e *FluentEntry) Ctx(ctx context.Context) *FluentEntry {
	e.ctx = ctx

	if traceID, ok := GetTraceID(ctx); ok {
		e.TraceID(traceID)
	}
	if requestID, ok := GetRequestID(ctx); ok {
		e.fields["request_id"] = requestID
	}
	if correlationID, ok := GetCorrelationID(ctx); ok {
		e.fields["correlation_id"] = correlationID
	}

	return e
}

func (e *FluentEntry) Msg(msg string) {
	logger := e.logger.WithFields(e.fields)

	if e.ctx != nil {
		e.logWithContext(logger, msg)
	} else {
		e.logDirect(logger, msg)
	}
}

func (e *FluentEntry) Msgf(format string, args ...interface{}) {
	logger := e.logger.WithFields(e.fields)

	if e.ctx != nil {
		e.logWithContextf(logger, format, args...)
	} else {
		e.logDirectf(logger, format, args...)
	}
}

func (e *FluentEntry) logWithContext(logger Logger, msg string) {
	switch e.level {
	case TraceLevel:
		logger.TraceContext(e.ctx, msg)
	case DebugLevel:
		logger.DebugContext(e.ctx, msg)
	case InfoLevel:
		logger.InfoContext(e.ctx, msg)
	case WarnLevel:
		logger.WarnContext(e.ctx, msg)
	case ErrorLevel:
		logger.ErrorContext(e.ctx, msg)
	case CriticalLevel:
		logger.CriticalContext(e.ctx, msg)
	}
}

func (e *FluentEntry) logWithContextf(logger Logger, format string, args ...interface{}) {
	switch e.level {
	case TraceLevel:
		logger.TraceContext(e.ctx, format, args...)
	case DebugLevel:
		logger.DebugContext(e.ctx, format, args...)
	case InfoLevel:
		logger.InfoContext(e.ctx, format, args...)
	case WarnLevel:
		logger.WarnContext(e.ctx, format, args...)
	case ErrorLevel:
		logger.ErrorContext(e.ctx, format, args...)
	case CriticalLevel:
		logger.CriticalContext(e.ctx, format, args...)
	}
}

func (e *FluentEntry) logDirect(logger Logger, msg string) {
	switch e.level {
	case TraceLevel:
		logger.Trace(msg)
	case DebugLevel:
		logger.Debug(msg)
	case InfoLevel:
		logger.Info(msg)
	case WarnLevel:
		logger.Warn(msg)
	case ErrorLevel:
		logger.Error(msg)
	case CriticalLevel:
		logger.Critical(msg)
	}
}

func (e *FluentEntry) logDirectf(logger Logger, format string, args ...interface{}) {
	switch e.level {
	case TraceLevel:
		logger.Trace(format, args...)
	case DebugLevel:
		logger.Debug(format, args...)
	case InfoLevel:
		logger.Info(format, args...)
	case WarnLevel:
		logger.Warn(format, args...)
	case ErrorLevel:
		logger.Error(format, args...)
	case CriticalLevel:
		logger.Critical(format, args...)
	}
}

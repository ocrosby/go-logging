package logging

import "context"

type fluentLoggerWrapper struct {
	logger Logger
}

func (w *fluentLoggerWrapper) createEntry(level Level) *FluentEntry {
	return &FluentEntry{
		logger: w.logger,
		level:  level,
		fields: make(map[string]interface{}),
	}
}

func (w *fluentLoggerWrapper) Trace() *FluentEntry {
	return w.createEntry(TraceLevel)
}

func (w *fluentLoggerWrapper) Debug() *FluentEntry {
	return w.createEntry(DebugLevel)
}

func (w *fluentLoggerWrapper) Info() *FluentEntry {
	return w.createEntry(InfoLevel)
}

func (w *fluentLoggerWrapper) Warn() *FluentEntry {
	return w.createEntry(WarnLevel)
}

func (w *fluentLoggerWrapper) Error() *FluentEntry {
	return w.createEntry(ErrorLevel)
}

func (w *fluentLoggerWrapper) Critical() *FluentEntry {
	return w.createEntry(CriticalLevel)
}

// FluentEntry represents a fluent logging entry that can be configured
// with structured fields before being output. Methods can be chained
// until Msg() or Msgf() is called to output the log entry.
//
// Example:
//
//	logger.Fluent().Info().
//		Str("service", "api").
//		Int("port", 8080).
//		Bool("ssl", true).
//		Msg("Server started")
type FluentEntry struct {
	logger  Logger
	level   Level
	fields  map[string]interface{}
	ctx     context.Context
	traceID string
}

// Field adds a key-value pair to the log entry and returns the entry for chaining.
func (e *FluentEntry) Field(key string, value interface{}) *FluentEntry {
	e.fields[key] = value
	return e
}

// Fields adds multiple key-value pairs to the log entry and returns the entry for chaining.
func (e *FluentEntry) Fields(fields map[string]interface{}) *FluentEntry {
	for k, v := range fields {
		e.fields[k] = v
	}
	return e
}

// Str adds a string field to the log entry and returns the entry for chaining.
func (e *FluentEntry) Str(key, value string) *FluentEntry {
	e.fields[key] = value
	return e
}

// Int adds an integer field to the log entry and returns the entry for chaining.
func (e *FluentEntry) Int(key string, value int) *FluentEntry {
	e.fields[key] = value
	return e
}

// Int64 adds a 64-bit integer field to the log entry and returns the entry for chaining.
func (e *FluentEntry) Int64(key string, value int64) *FluentEntry {
	e.fields[key] = value
	return e
}

// Bool adds a boolean field to the log entry and returns the entry for chaining.
func (e *FluentEntry) Bool(key string, value bool) *FluentEntry {
	e.fields[key] = value
	return e
}

// Err adds an error field to the log entry and returns the entry for chaining.
// If err is nil, no field is added.
func (e *FluentEntry) Err(err error) *FluentEntry {
	if err != nil {
		e.fields["error"] = err.Error()
	}
	return e
}

// TraceID adds a trace identifier to the log entry and returns the entry for chaining.
// The trace ID will appear as "trace_id" in the output.
func (e *FluentEntry) TraceID(id string) *FluentEntry {
	e.traceID = id
	e.fields["trace_id"] = id
	return e
}

// Ctx adds context information to the log entry and returns the entry for chaining.
// Automatically extracts trace_id, request_id, and correlation_id from the context if present.
//
// Example:
//
//	ctx := logging.WithTraceID(context.Background(), "trace-123")
//	logger.Fluent().Info().Ctx(ctx).Msg("Processing request")
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

// Msg outputs the log entry with the specified message.
// This is the terminal method that actually writes the log.
func (e *FluentEntry) Msg(msg string) {
	logger := e.logger.WithFields(e.fields)
	e.dispatch(logger, msg, nil)
}

// Msgf outputs the log entry with a formatted message.
// This is the terminal method that actually writes the log.
//
// Example:
//
//	logger.Fluent().Info().
//		Str("user", username).
//		Msgf("User %s logged in at %s", username, time.Now())
func (e *FluentEntry) Msgf(format string, args ...interface{}) {
	logger := e.logger.WithFields(e.fields)
	e.dispatch(logger, format, args)
}

type levelMethod func(string, ...interface{})
type contextLevelMethod func(context.Context, string, ...interface{})

var levelMethodMap = map[Level]func(Logger) levelMethod{
	TraceLevel: func(l Logger) levelMethod {
		return l.Trace
	},
	DebugLevel: func(l Logger) levelMethod {
		return l.Debug
	},
	InfoLevel: func(l Logger) levelMethod {
		return l.Info
	},
	WarnLevel: func(l Logger) levelMethod {
		return l.Warn
	},
	ErrorLevel: func(l Logger) levelMethod {
		return l.Error
	},
	CriticalLevel: func(l Logger) levelMethod {
		return l.Critical
	},
}

var contextLevelMethodMap = map[Level]func(Logger) contextLevelMethod{
	TraceLevel: func(l Logger) contextLevelMethod {
		return l.TraceContext
	},
	DebugLevel: func(l Logger) contextLevelMethod {
		return l.DebugContext
	},
	InfoLevel: func(l Logger) contextLevelMethod {
		return l.InfoContext
	},
	WarnLevel: func(l Logger) contextLevelMethod {
		return l.WarnContext
	},
	ErrorLevel: func(l Logger) contextLevelMethod {
		return l.ErrorContext
	},
	CriticalLevel: func(l Logger) contextLevelMethod {
		return l.CriticalContext
	},
}

func (e *FluentEntry) dispatch(logger Logger, format string, args []interface{}) {
	if e.ctx != nil {
		if methodGetter, ok := contextLevelMethodMap[e.level]; ok {
			method := methodGetter(logger)
			if args == nil {
				method(e.ctx, format)
			} else {
				method(e.ctx, format, args...)
			}
		}
	} else {
		if methodGetter, ok := levelMethodMap[e.level]; ok {
			method := methodGetter(logger)
			if args == nil {
				method(format)
			} else {
				method(format, args...)
			}
		}
	}
}

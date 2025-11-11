package logging

import (
	"context"
	"log/slog"
	"runtime"
	"time"
)

type HandlerMiddleware interface {
	Handle(ctx context.Context, record slog.Record, next HandlerFunc) error
}

type HandlerFunc func(context.Context, slog.Record) error

type handlerMiddlewareFunc func(context.Context, slog.Record, HandlerFunc) error

func (f handlerMiddlewareFunc) Handle(ctx context.Context, record slog.Record, next HandlerFunc) error {
	return f(ctx, record, next)
}

type MiddlewareHandler struct {
	handler     slog.Handler
	middlewares []HandlerMiddleware
}

func NewMiddlewareHandler(handler slog.Handler, middlewares ...HandlerMiddleware) *MiddlewareHandler {
	return &MiddlewareHandler{
		handler:     handler,
		middlewares: middlewares,
	}
}

func (h *MiddlewareHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

func (h *MiddlewareHandler) Handle(ctx context.Context, record slog.Record) error {
	chain := h.buildChain()
	return chain(ctx, record)
}

func (h *MiddlewareHandler) buildChain() HandlerFunc {
	final := func(ctx context.Context, record slog.Record) error {
		return h.handler.Handle(ctx, record)
	}

	for i := len(h.middlewares) - 1; i >= 0; i-- {
		middleware := h.middlewares[i]
		next := final
		final = func(ctx context.Context, record slog.Record) error {
			return middleware.Handle(ctx, record, next)
		}
	}

	return final
}

func (h *MiddlewareHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &MiddlewareHandler{
		handler:     h.handler.WithAttrs(attrs),
		middlewares: h.middlewares,
	}
}

func (h *MiddlewareHandler) WithGroup(name string) slog.Handler {
	return &MiddlewareHandler{
		handler:     h.handler.WithGroup(name),
		middlewares: h.middlewares,
	}
}

func TimestampMiddleware() HandlerMiddleware {
	return handlerMiddlewareFunc(func(ctx context.Context, record slog.Record, next HandlerFunc) error {
		if record.Time.IsZero() {
			record.Time = time.Now()
		}
		return next(ctx, record)
	})
}

func ContextExtractorMiddleware(extractor ContextExtractor) HandlerMiddleware {
	return handlerMiddlewareFunc(func(ctx context.Context, record slog.Record, next HandlerFunc) error {
		attrs := extractor.Extract(ctx)
		if len(attrs) > 0 {
			record.AddAttrs(attrs...)
		}
		return next(ctx, record)
	})
}

func LevelFilterMiddleware(minLevel slog.Level) HandlerMiddleware {
	return handlerMiddlewareFunc(func(ctx context.Context, record slog.Record, next HandlerFunc) error {
		if record.Level < minLevel {
			return nil
		}
		return next(ctx, record)
	})
}

func SamplingMiddleware(rate int) HandlerMiddleware {
	counter := 0
	return handlerMiddlewareFunc(func(ctx context.Context, record slog.Record, next HandlerFunc) error {
		counter++
		if counter%rate != 0 {
			return nil
		}
		return next(ctx, record)
	})
}

func CallerMiddleware(skip int) HandlerMiddleware {
	return handlerMiddlewareFunc(func(ctx context.Context, record slog.Record, next HandlerFunc) error {
		if record.PC == 0 {
			record.PC = getPCs(skip + 3)
		}
		return next(ctx, record)
	})
}

func getPCs(skip int) uintptr {
	var pcs [1]uintptr
	n := runtime.Callers(skip+1, pcs[:])
	if n > 0 {
		return pcs[0]
	}
	return 0
}

func StaticFieldsMiddleware(fields map[string]interface{}) HandlerMiddleware {
	attrs := make([]slog.Attr, 0, len(fields))
	for k, v := range fields {
		attrs = append(attrs, slog.Any(k, v))
	}

	return handlerMiddlewareFunc(func(ctx context.Context, record slog.Record, next HandlerFunc) error {
		record.AddAttrs(attrs...)
		return next(ctx, record)
	})
}

func RedactionMiddleware(redactor Redactor) HandlerMiddleware {
	return handlerMiddlewareFunc(func(ctx context.Context, record slog.Record, next HandlerFunc) error {
		record.Message = redactor.Redact(record.Message)
		return next(ctx, record)
	})
}

type LoggingMiddleware struct {
	onBefore func(context.Context, slog.Record)
	onAfter  func(context.Context, slog.Record, error)
}

func NewLoggingMiddleware(onBefore func(context.Context, slog.Record), onAfter func(context.Context, slog.Record, error)) HandlerMiddleware {
	return &LoggingMiddleware{
		onBefore: onBefore,
		onAfter:  onAfter,
	}
}

func (m *LoggingMiddleware) Handle(ctx context.Context, record slog.Record, next HandlerFunc) error {
	if m.onBefore != nil {
		m.onBefore(ctx, record)
	}

	err := next(ctx, record)

	if m.onAfter != nil {
		m.onAfter(ctx, record, err)
	}

	return err
}

func MetricsMiddleware(recordMetric func(level slog.Level)) HandlerMiddleware {
	return handlerMiddlewareFunc(func(ctx context.Context, record slog.Record, next HandlerFunc) error {
		if recordMetric != nil {
			recordMetric(record.Level)
		}
		return next(ctx, record)
	})
}

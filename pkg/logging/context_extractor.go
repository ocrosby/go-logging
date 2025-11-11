package logging

import (
	"context"
	"log/slog"
)

type ContextExtractor interface {
	Extract(ctx context.Context) []slog.Attr
}

type contextExtractorFunc func(context.Context) []slog.Attr

func (f contextExtractorFunc) Extract(ctx context.Context) []slog.Attr {
	return f(ctx)
}

type CompositeContextExtractor struct {
	extractors []ContextExtractor
}

func NewCompositeContextExtractor(extractors ...ContextExtractor) *CompositeContextExtractor {
	return &CompositeContextExtractor{
		extractors: extractors,
	}
}

func (c *CompositeContextExtractor) Extract(ctx context.Context) []slog.Attr {
	var attrs []slog.Attr
	for _, extractor := range c.extractors {
		attrs = append(attrs, extractor.Extract(ctx)...)
	}
	return attrs
}

func (c *CompositeContextExtractor) Add(extractor ContextExtractor) {
	c.extractors = append(c.extractors, extractor)
}

func TraceContextExtractor() ContextExtractor {
	return contextExtractorFunc(func(ctx context.Context) []slog.Attr {
		var attrs []slog.Attr

		if traceID, ok := GetTraceID(ctx); ok {
			attrs = append(attrs, slog.String("trace_id", traceID))
		}

		if requestID, ok := GetRequestID(ctx); ok {
			attrs = append(attrs, slog.String("request_id", requestID))
		}

		if correlationID, ok := GetCorrelationID(ctx); ok {
			attrs = append(attrs, slog.String("correlation_id", correlationID))
		}

		return attrs
	})
}

type ContextKey string

func CustomContextExtractor(key ContextKey, attrName string) ContextExtractor {
	return contextExtractorFunc(func(ctx context.Context) []slog.Attr {
		if value := ctx.Value(key); value != nil {
			return []slog.Attr{slog.Any(attrName, value)}
		}
		return nil
	})
}

func StringContextExtractor(key ContextKey, attrName string) ContextExtractor {
	return contextExtractorFunc(func(ctx context.Context) []slog.Attr {
		if value := ctx.Value(key); value != nil {
			if str, ok := value.(string); ok {
				return []slog.Attr{slog.String(attrName, str)}
			}
		}
		return nil
	})
}

func IntContextExtractor(key ContextKey, attrName string) ContextExtractor {
	return contextExtractorFunc(func(ctx context.Context) []slog.Attr {
		if value := ctx.Value(key); value != nil {
			if i, ok := value.(int); ok {
				return []slog.Attr{slog.Int(attrName, i)}
			}
		}
		return nil
	})
}

func Int64ContextExtractor(key ContextKey, attrName string) ContextExtractor {
	return contextExtractorFunc(func(ctx context.Context) []slog.Attr {
		if value := ctx.Value(key); value != nil {
			if i, ok := value.(int64); ok {
				return []slog.Attr{slog.Int64(attrName, i)}
			}
		}
		return nil
	})
}

func BoolContextExtractor(key ContextKey, attrName string) ContextExtractor {
	return contextExtractorFunc(func(ctx context.Context) []slog.Attr {
		if value := ctx.Value(key); value != nil {
			if b, ok := value.(bool); ok {
				return []slog.Attr{slog.Bool(attrName, b)}
			}
		}
		return nil
	})
}

package internal

import (
	"context"
	"fmt"
	"regexp"
	"strings"
)

// ContextFields represents extracted context information
type ContextFields struct {
	TraceID       string
	RequestID     string
	CorrelationID string
}

// FormatFilename formats a filename and line number for logging output
func FormatFilename(file string, line int, useShort bool) string {
	if useShort {
		// Extract just the filename
		if idx := strings.LastIndex(file, "/"); idx >= 0 {
			file = file[idx+1:]
		}
	}
	return fmt.Sprintf("%s:%d", file, line)
}

// ApplyRedactionPatterns applies redaction patterns to a message
func ApplyRedactionPatterns(message string, patterns []*regexp.Regexp) string {
	result := message
	for _, pattern := range patterns {
		result = pattern.ReplaceAllString(result, "[REDACTED]")
	}
	return result
}

// ExtractContextFields extracts common context fields used throughout logging
func ExtractContextFields(ctx context.Context) ContextFields {
	if ctx == nil {
		return ContextFields{}
	}

	fields := ContextFields{}

	fields.TraceID = extractTraceID(ctx)
	fields.RequestID = extractRequestID(ctx)
	fields.CorrelationID = extractCorrelationID(ctx)

	return fields
}

func extractTraceID(ctx context.Context) string {
	if val := ctx.Value(contextKey("trace_id")); val != nil {
		if traceID, ok := val.(string); ok {
			return traceID
		}
	}
	return ""
}

func extractRequestID(ctx context.Context) string {
	if val := ctx.Value(contextKey("request_id")); val != nil {
		if requestID, ok := val.(string); ok {
			return requestID
		}
	}
	return ""
}

func extractCorrelationID(ctx context.Context) string {
	if val := ctx.Value(contextKey("correlation_id")); val != nil {
		if correlationID, ok := val.(string); ok {
			return correlationID
		}
	}
	return ""
}

// Context keys (duplicated to avoid circular import)
type contextKey string

// Remove unused constants - only keep the contextKey type

// AddContextFieldsToMap adds context fields to a map for JSON/structured output
func (cf ContextFields) AddToMap(data map[string]interface{}) {
	if cf.TraceID != "" {
		data["trace_id"] = cf.TraceID
	}
	if cf.RequestID != "" {
		data["request_id"] = cf.RequestID
	}
	if cf.CorrelationID != "" {
		data["correlation_id"] = cf.CorrelationID
	}
}

// FormatForText formats context fields for text output
func (cf ContextFields) FormatForText() string {
	var parts []string
	if cf.TraceID != "" {
		parts = append(parts, fmt.Sprintf("trace_id=%s", cf.TraceID))
	}
	if cf.RequestID != "" {
		parts = append(parts, fmt.Sprintf("request_id=%s", cf.RequestID))
	}
	if cf.CorrelationID != "" {
		parts = append(parts, fmt.Sprintf("correlation_id=%s", cf.CorrelationID))
	}

	if len(parts) > 0 {
		return fmt.Sprintf("[%s]", strings.Join(parts, " "))
	}
	return ""
}

// HasAnyFields returns true if any context fields are present
func (cf ContextFields) HasAnyFields() bool {
	return cf.TraceID != "" || cf.RequestID != "" || cf.CorrelationID != ""
}

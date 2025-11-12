package logging

import (
	"net/http"
	"time"
)

const (
	HeaderTraceID       = "X-Trace-ID"
	HeaderRequestID     = "X-Request-ID"
	HeaderCorrelationID = "X-Correlation-ID"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
	written    int64
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(b)
	rw.written += int64(n)
	return n, err
}

func TracingMiddleware(logger Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			ctx := r.Context()

			traceID := r.Header.Get(HeaderTraceID)
			if traceID == "" {
				traceID = NewTraceID()
			}
			ctx = WithTraceID(ctx, traceID)

			requestID := r.Header.Get(HeaderRequestID)
			if requestID != "" {
				ctx = WithRequestID(ctx, requestID)
			}

			correlationID := r.Header.Get(HeaderCorrelationID)
			if correlationID != "" {
				ctx = WithCorrelationID(ctx, correlationID)
			}

			w.Header().Set(HeaderTraceID, traceID)

			rw := &responseWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			logger.Fluent().Info().
				Ctx(ctx).
				Str("method", r.Method).
				Str("path", RedactedURL(r.URL.String())).
				Str("remote_addr", r.RemoteAddr).
				Str("user_agent", r.UserAgent()).
				Msg("Request started")

			next.ServeHTTP(rw, r.WithContext(ctx))

			duration := time.Since(start)

			logger.Fluent().Info().
				Ctx(ctx).
				Str("method", r.Method).
				Str("path", RedactedURL(r.URL.String())).
				Int("status", rw.statusCode).
				Int64("bytes", rw.written).
				Int64("duration_ms", duration.Milliseconds()).
				Msg("Request completed")
		})
	}
}

func RequestLogger(logger Logger, headers ...string) func(http.Handler) http.Handler {
	if len(headers) == 0 {
		headers = []string{"User-Agent"}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			traceID, _ := GetTraceID(ctx)

			entry := logger.Fluent().Info().
				TraceID(traceID).
				Str("method", r.Method).
				Str("path", RedactedURL(r.URL.String())).
				Str("headers", RequestHeaders(r, headers))
			entry.Msg("HTTP Request")

			next.ServeHTTP(w, r)
		})
	}
}

package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/ocrosby/go-logging/pkg/logging"
)

func main() {
	exampleMultiHandler()
	exampleHandlerBuilder()
	exampleConditionalHandler()
	exampleBufferedHandler()
	exampleMiddlewareChain()
}

func exampleMultiHandler() {
	stdoutHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	fileHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelWarn,
	})

	multiHandler := logging.NewMultiHandler(stdoutHandler, fileHandler)
	logger := logging.NewWithHandler(multiHandler)

	logger.Info("Info message - appears in stdout only")
	logger.Warn("Warning message - appears in both handlers")
	logger.Error("Error message - appears in both handlers")
}

func exampleHandlerBuilder() {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	builtHandler := logging.NewHandlerBuilder(handler).
		WithTimestamp().
		WithTraceContext().
		WithStaticFields(map[string]interface{}{
			"service": "custom-handler-example",
			"version": "1.0.0",
		}).
		Build()

	logger := logging.NewWithHandler(builtHandler)

	ctx := logging.WithTraceID(context.Background(), "trace-abc-123")
	logger.InfoContext(ctx, "Logging with handler builder")
}

func exampleConditionalHandler() {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	conditionalHandler := logging.NewConditionalHandler(handler, func(ctx context.Context, record slog.Record) bool {
		return record.Level >= slog.LevelWarn
	})

	logger := logging.NewWithHandler(conditionalHandler)

	logger.Debug("Debug - filtered out")
	logger.Info("Info - filtered out")
	logger.Warn("Warning - logged")
	logger.Error("Error - logged")
}

func exampleBufferedHandler() {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	bufferedHandler := logging.NewBufferedHandler(handler, 5)
	logger := logging.NewWithHandler(bufferedHandler)

	logger.Info("Message 1")
	logger.Info("Message 2")
	logger.Info("Message 3")
	logger.Info("Message 4")
	logger.Info("Message 5")

	_ = bufferedHandler.Flush(context.Background())
}

func exampleMiddlewareChain() {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	middlewareHandler := logging.NewMiddlewareHandler(
		handler,
		logging.TimestampMiddleware(),
		logging.StaticFieldsMiddleware(map[string]interface{}{
			"app": "example",
		}),
		logging.ContextExtractorMiddleware(logging.TraceContextExtractor()),
	)

	logger := logging.NewWithHandler(middlewareHandler)

	ctx := logging.WithTraceID(context.Background(), "trace-xyz-789")
	ctx = logging.WithRequestID(ctx, "req-123")

	logger.InfoContext(ctx, "Logging with middleware chain")
}

package logging

import (
	"bytes"
	"context"
	"log/slog"
	"testing"
)

func BenchmarkStandardLogger_Info(b *testing.B) {
	var buf bytes.Buffer
	config := NewConfig().
		WithLevel(InfoLevel).
		WithOutput(&buf).
		WithTextFormat().
		Build()
	logger := NewStandardLogger(config, NewRedactorChain())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("benchmark message")
	}
}

func BenchmarkSlogLogger_Info(b *testing.B) {
	var buf bytes.Buffer
	handler := slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	logger := NewSlogLogger(handler, NewRedactorChain())
	logger.SetLevel(InfoLevel)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("benchmark message")
	}
}

func BenchmarkStandardLogger_InfoWithFields(b *testing.B) {
	var buf bytes.Buffer
	config := NewConfig().
		WithLevel(InfoLevel).
		WithOutput(&buf).
		WithTextFormat().
		Build()
	logger := NewStandardLogger(config, NewRedactorChain()).
		WithField("service", "test").
		WithField("version", "1.0.0")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("benchmark message")
	}
}

func BenchmarkSlogLogger_InfoWithFields(b *testing.B) {
	var buf bytes.Buffer
	handler := slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	logger := NewSlogLogger(handler, NewRedactorChain()).
		WithField("service", "test").
		WithField("version", "1.0.0")
	logger.SetLevel(InfoLevel)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("benchmark message")
	}
}

func BenchmarkStandardLogger_InfoContext(b *testing.B) {
	var buf bytes.Buffer
	config := NewConfig().
		WithLevel(InfoLevel).
		WithOutput(&buf).
		WithTextFormat().
		Build()
	logger := NewStandardLogger(config, NewRedactorChain())
	ctx := WithTraceID(context.Background(), "trace-123")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.InfoContext(ctx, "benchmark message")
	}
}

func BenchmarkSlogLogger_InfoContext(b *testing.B) {
	var buf bytes.Buffer
	handler := slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	logger := NewSlogLogger(handler, NewRedactorChain())
	logger.SetLevel(InfoLevel)
	ctx := WithTraceID(context.Background(), "trace-123")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.InfoContext(ctx, "benchmark message")
	}
}

func BenchmarkStandardLogger_JSON(b *testing.B) {
	var buf bytes.Buffer
	config := NewConfig().
		WithLevel(InfoLevel).
		WithOutput(&buf).
		WithJSONFormat().
		Build()
	logger := NewStandardLogger(config, NewRedactorChain())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("benchmark message")
	}
}

func BenchmarkSlogLogger_JSON(b *testing.B) {
	var buf bytes.Buffer
	handler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	logger := NewSlogLogger(handler, NewRedactorChain())
	logger.SetLevel(InfoLevel)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("benchmark message")
	}
}

func BenchmarkStandardLogger_Fluent(b *testing.B) {
	var buf bytes.Buffer
	config := NewConfig().
		WithLevel(InfoLevel).
		WithOutput(&buf).
		WithJSONFormat().
		Build()
	logger := NewStandardLogger(config, NewRedactorChain())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Fluent().Info().
			Str("service", "test").
			Int("count", 42).
			Msg("benchmark message")
	}
}

func BenchmarkSlogLogger_Fluent(b *testing.B) {
	var buf bytes.Buffer
	handler := slog.NewJSONHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	logger := NewSlogLogger(handler, NewRedactorChain())
	logger.SetLevel(InfoLevel)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Fluent().Info().
			Str("service", "test").
			Int("count", 42).
			Msg("benchmark message")
	}
}

func BenchmarkStandardLogger_Formatting(b *testing.B) {
	var buf bytes.Buffer
	config := NewConfig().
		WithLevel(InfoLevel).
		WithOutput(&buf).
		WithTextFormat().
		Build()
	logger := NewStandardLogger(config, NewRedactorChain())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("user %s logged in with id %d", "john", 123)
	}
}

func BenchmarkSlogLogger_Formatting(b *testing.B) {
	var buf bytes.Buffer
	handler := slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	logger := NewSlogLogger(handler, NewRedactorChain())
	logger.SetLevel(InfoLevel)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("user %s logged in with id %d", "john", 123)
	}
}

func BenchmarkStandardLogger_LevelCheck(b *testing.B) {
	var buf bytes.Buffer
	config := NewConfig().
		WithLevel(WarnLevel).
		WithOutput(&buf).
		WithTextFormat().
		Build()
	logger := NewStandardLogger(config, NewRedactorChain())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if logger.IsLevelEnabled(DebugLevel) {
			logger.Debug("this should not execute")
		}
	}
}

func BenchmarkSlogLogger_LevelCheck(b *testing.B) {
	var buf bytes.Buffer
	handler := slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelWarn,
	})
	logger := NewSlogLogger(handler, NewRedactorChain())
	logger.SetLevel(WarnLevel)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if logger.IsLevelEnabled(DebugLevel) {
			logger.Debug("this should not execute")
		}
	}
}

func BenchmarkStandardLogger_WithFieldAllocation(b *testing.B) {
	var buf bytes.Buffer
	config := NewConfig().
		WithLevel(InfoLevel).
		WithOutput(&buf).
		WithTextFormat().
		Build()
	baseLogger := NewStandardLogger(config, NewRedactorChain())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger := baseLogger.WithField("request_id", "req-123")
		logger.Info("benchmark message")
	}
}

func BenchmarkSlogLogger_WithFieldAllocation(b *testing.B) {
	var buf bytes.Buffer
	handler := slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	baseLogger := NewSlogLogger(handler, NewRedactorChain())
	baseLogger.SetLevel(InfoLevel)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger := baseLogger.WithField("request_id", "req-123")
		logger.Info("benchmark message")
	}
}

func BenchmarkContextExtractor(b *testing.B) {
	extractor := TraceContextExtractor()
	ctx := context.Background()
	ctx = WithTraceID(ctx, "trace-123")
	ctx = WithRequestID(ctx, "req-456")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = extractor.Extract(ctx)
	}
}

func BenchmarkCompositeContextExtractor_Benchmark(b *testing.B) {
	const userKey ContextKey = "user"
	composite := NewCompositeContextExtractor(
		TraceContextExtractor(),
		StringContextExtractor(userKey, "username"),
	)
	ctx := context.Background()
	ctx = WithTraceID(ctx, "trace-123")
	ctx = context.WithValue(ctx, userKey, "john_doe")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = composite.Extract(ctx)
	}
}

func BenchmarkRedaction_StandardLogger(b *testing.B) {
	var buf bytes.Buffer
	config := NewConfig().
		WithLevel(InfoLevel).
		WithOutput(&buf).
		WithTextFormat().
		AddRedactPattern(`password=\w+`).
		Build()
	redactorChain := ProvideRedactorChain(config)
	logger := NewStandardLogger(config, redactorChain)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("user logged in with password=secret123")
	}
}

func BenchmarkRedaction_SlogLogger(b *testing.B) {
	var buf bytes.Buffer
	handler := slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	config := NewConfig().
		AddRedactPattern(`password=\w+`).
		Build()
	redactorChain := ProvideRedactorChain(config)
	logger := NewSlogLogger(handler, redactorChain)
	logger.SetLevel(InfoLevel)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("user logged in with password=secret123")
	}
}

func BenchmarkParallelLogging_StandardLogger(b *testing.B) {
	var buf bytes.Buffer
	config := NewConfig().
		WithLevel(InfoLevel).
		WithOutput(&buf).
		WithTextFormat().
		Build()
	logger := NewStandardLogger(config, NewRedactorChain())

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("parallel benchmark message")
		}
	})
}

func BenchmarkParallelLogging_SlogLogger(b *testing.B) {
	var buf bytes.Buffer
	handler := slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	logger := NewSlogLogger(handler, NewRedactorChain())
	logger.SetLevel(InfoLevel)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("parallel benchmark message")
		}
	})
}

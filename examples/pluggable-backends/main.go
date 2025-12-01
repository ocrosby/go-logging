package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/ocrosby/go-logging/pkg/logging"
	"github.com/rs/zerolog"
	"go.uber.org/zap"
	"go.uber.org/zap/exp/zapslog"
)

func main() {
	fmt.Println("=== Pluggable Backend Examples ===\n")
	fmt.Println("This demonstrates how you can use different logging backends")
	fmt.Println("(slog, zerolog, zap, custom) while maintaining the same Logger interface.\n")

	exampleStandardSlog()
	fmt.Println()

	exampleZerolog()
	fmt.Println()

	exampleZap()
	fmt.Println()

	DemoCommonLogFormat()
	fmt.Println()

	DemoCommonLogFormatWithLoggingLibrary()
	fmt.Println()

	demonstrateInterfaceConsistency()
}

func exampleStandardSlog() {
	fmt.Println("--- 1. Standard slog Backend ---")

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	logger := logging.NewWithHandler(handler)

	logger.Info("Using standard slog backend")
	logger.Debug("Debug message from slog")
	logger.Warn("Warning from slog")

	logger = logger.WithFields(map[string]interface{}{
		"backend": "slog",
		"version": "1.0",
	})
	logger.Info("With static fields")
}

func exampleZerolog() {
	fmt.Println("--- 2. Zerolog Backend (High Performance) ---")

	zerologLogger := zerolog.New(os.Stdout).
		With().
		Timestamp().
		Logger()

	handler := zerologSlogHandler{logger: zerologLogger}

	logger := logging.NewWithHandler(handler)

	logger.Info("Using zerolog backend")
	logger.Debug("Debug message from zerolog")
	logger.Warn("Warning from zerolog")

	logger = logger.WithFields(map[string]interface{}{
		"backend": "zerolog",
		"version": "1.0",
	})
	logger.Error("Error with zerolog backend")
}

func exampleZap() {
	fmt.Println("--- 3. Zap Backend (Uber's Logger) ---")

	zapLogger, _ := zap.NewProduction()
	defer zapLogger.Sync()

	handler := zapslog.NewHandler(zapLogger.Core(), zapslog.WithCallerSkip(3))

	logger := logging.NewWithHandler(handler)

	logger.Info("Using zap backend")
	logger.Debug("Debug message from zap")
	logger.Warn("Warning from zap")

	logger = logger.WithFields(map[string]interface{}{
		"backend": "zap",
		"version": "1.0",
	})
	logger.Error("Error with zap backend")
}

func demonstrateInterfaceConsistency() {
	fmt.Println("--- 4. Interface Consistency Demo ---")
	fmt.Println("Notice: All three loggers use the SAME interface\n")

	backends := []struct {
		name    string
		logger  logging.Logger
		backend string
	}{
		{
			name:    "slog",
			logger:  createSlogLogger(),
			backend: "Go standard slog",
		},
		{
			name:    "zerolog",
			logger:  createZerologLogger(),
			backend: "rs/zerolog",
		},
		{
			name:    "zap",
			logger:  createZapLogger(),
			backend: "uber-go/zap",
		},
	}

	ctx := logging.NewContextWithTrace()
	ctx = logging.WithRequestID(ctx, "req-12345")

	for _, b := range backends {
		fmt.Printf("\nUsing %s backend (%s):\n", b.name, b.backend)

		logger := b.logger.WithField("backend", b.name)

		logger.Info("Standard info log")

		logger.InfoContext(ctx, "Log with context")

		logger.Fluent().Info().
			Str("user", "john_doe").
			Int("attempt", 3).
			Msg("Fluent interface works too")
	}

	fmt.Println("\n✅ All backends use the same logging.Logger interface!")
	fmt.Println("✅ Application code stays identical regardless of backend!")
}

func createSlogLogger() logging.Logger {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	return logging.NewWithHandler(handler)
}

func createZerologLogger() logging.Logger {
	zerologLogger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	handler := zerologSlogHandler{logger: zerologLogger}
	return logging.NewWithHandler(handler)
}

func createZapLogger() logging.Logger {
	zapLogger, _ := zap.NewProduction()
	handler := zapslog.NewHandler(zapLogger.Core(), zapslog.WithCallerSkip(3))
	return logging.NewWithHandler(handler)
}

type zerologSlogHandler struct {
	logger zerolog.Logger
	attrs  []slog.Attr
	groups []string
}

func (h zerologSlogHandler) Enabled(_ context.Context, level slog.Level) bool {
	zerologLevel := h.slogLevelToZerolog(level)
	return h.logger.GetLevel() <= zerologLevel
}

func (h zerologSlogHandler) Handle(_ context.Context, record slog.Record) error {
	var event *zerolog.Event

	switch record.Level {
	case slog.LevelDebug:
		event = h.logger.Debug()
	case slog.LevelInfo:
		event = h.logger.Info()
	case slog.LevelWarn:
		event = h.logger.Warn()
	case slog.LevelError:
		event = h.logger.Error()
	default:
		event = h.logger.Info()
	}

	for _, attr := range h.attrs {
		h.addAttr(event, attr)
	}

	record.Attrs(func(attr slog.Attr) bool {
		h.addAttr(event, attr)
		return true
	})

	event.Msg(record.Message)
	return nil
}

func (h zerologSlogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newAttrs := make([]slog.Attr, len(h.attrs)+len(attrs))
	copy(newAttrs, h.attrs)
	copy(newAttrs[len(h.attrs):], attrs)

	return zerologSlogHandler{
		logger: h.logger,
		attrs:  newAttrs,
		groups: h.groups,
	}
}

func (h zerologSlogHandler) WithGroup(name string) slog.Handler {
	newGroups := make([]string, len(h.groups)+1)
	copy(newGroups, h.groups)
	newGroups[len(h.groups)] = name

	return zerologSlogHandler{
		logger: h.logger,
		attrs:  h.attrs,
		groups: newGroups,
	}
}

func (h zerologSlogHandler) addAttr(event *zerolog.Event, attr slog.Attr) {
	key := attr.Key
	if len(h.groups) > 0 {
		for _, group := range h.groups {
			key = group + "." + key
		}
	}

	switch attr.Value.Kind() {
	case slog.KindString:
		event.Str(key, attr.Value.String())
	case slog.KindInt64:
		event.Int64(key, attr.Value.Int64())
	case slog.KindUint64:
		event.Uint64(key, attr.Value.Uint64())
	case slog.KindFloat64:
		event.Float64(key, attr.Value.Float64())
	case slog.KindBool:
		event.Bool(key, attr.Value.Bool())
	case slog.KindDuration:
		event.Dur(key, attr.Value.Duration())
	case slog.KindTime:
		event.Time(key, attr.Value.Time())
	default:
		event.Interface(key, attr.Value.Any())
	}
}

func (h zerologSlogHandler) slogLevelToZerolog(level slog.Level) zerolog.Level {
	switch {
	case level < slog.LevelInfo:
		return zerolog.DebugLevel
	case level < slog.LevelWarn:
		return zerolog.InfoLevel
	case level < slog.LevelError:
		return zerolog.WarnLevel
	default:
		return zerolog.ErrorLevel
	}
}

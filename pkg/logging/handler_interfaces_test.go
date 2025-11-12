package logging

import (
	"bytes"
	"context"
	"log/slog"
	"testing"
	"time"
)

func TestNewBaseHandler(t *testing.T) {
	buf := &bytes.Buffer{}
	handler := slog.NewTextHandler(buf, nil)

	baseHandler := NewBaseHandler("test-handler", handler)

	if baseHandler == nil {
		t.Error("expected base handler to be created")
	}

	if baseHandler.Name() != "test-handler" {
		t.Errorf("expected name 'test-handler', got %s", baseHandler.Name())
	}
}

func TestBaseHandler_Enabled(t *testing.T) {
	buf := &bytes.Buffer{}
	handler := slog.NewTextHandler(buf, nil)
	baseHandler := NewBaseHandler("test-handler", handler)

	ctx := context.Background()

	// Should delegate to underlying handler
	enabled := baseHandler.Enabled(ctx, slog.LevelInfo)
	if !enabled {
		t.Error("expected handler to be enabled for info level")
	}
}

func TestBaseHandler_Handle(t *testing.T) {
	buf := &bytes.Buffer{}
	handler := slog.NewTextHandler(buf, nil)
	baseHandler := NewBaseHandler("test-handler", handler)

	ctx := context.Background()
	record := slog.NewRecord(time.Now(), slog.LevelInfo, "test message", 0)

	err := baseHandler.Handle(ctx, record)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	output := buf.String()
	if output == "" {
		t.Error("expected output from handler")
	}
}

func TestBaseHandler_WithAttrs(t *testing.T) {
	buf := &bytes.Buffer{}
	handler := slog.NewTextHandler(buf, nil)
	baseHandler := NewBaseHandler("test-handler", handler)

	attrs := []slog.Attr{slog.String("key", "value")}
	newHandler := baseHandler.WithAttrs(attrs)

	if newHandler == nil {
		t.Error("expected WithAttrs to return a handler")
	}
}

func TestBaseHandler_WithGroup(t *testing.T) {
	buf := &bytes.Buffer{}
	handler := slog.NewTextHandler(buf, nil)
	baseHandler := NewBaseHandler("test-handler", handler)

	newHandler := baseHandler.WithGroup("testgroup")

	if newHandler == nil {
		t.Error("expected WithGroup to return a handler")
	}
}

func TestBaseHandler_WithMiddleware(t *testing.T) {
	buf := &bytes.Buffer{}
	handler := slog.NewTextHandler(buf, nil)
	baseHandler := NewBaseHandler("test-handler", handler)

	middleware := handlerMiddlewareFunc(func(ctx context.Context, record slog.Record, next HandlerFunc) error {
		record.AddAttrs(slog.String("middleware", "test"))
		return next(ctx, record)
	})

	newHandler := baseHandler.WithMiddleware(middleware)

	if newHandler == nil {
		t.Error("expected WithMiddleware to return a handler")
	}
}

func TestBaseHandler_Create(t *testing.T) {
	buf := &bytes.Buffer{}
	handler := slog.NewTextHandler(buf, nil)
	baseHandler := NewBaseHandler("test-handler", handler)

	config := map[string]interface{}{"key": "value"}
	newHandler, err := baseHandler.Create(config)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if newHandler == nil {
		t.Error("expected Create to return a handler")
	}
}

func TestBaseHandler_ConfigType(t *testing.T) {
	buf := &bytes.Buffer{}
	handler := slog.NewTextHandler(buf, nil)
	baseHandler := NewBaseHandler("test-handler", handler)

	configType := baseHandler.ConfigType()

	// ConfigType may return nil for default implementation
	_ = configType // Just test that method exists and doesn't panic
}

func TestBaseHandler_Close(t *testing.T) {
	buf := &bytes.Buffer{}
	handler := slog.NewTextHandler(buf, nil)
	baseHandler := NewBaseHandler("test-handler", handler)

	err := baseHandler.Close()

	// Should not error for basic close
	if err != nil {
		t.Errorf("unexpected error on close: %v", err)
	}
}

func TestNewHandlerCompositor(t *testing.T) {
	compositor := NewHandlerCompositor()

	if compositor == nil {
		t.Error("expected compositor to be created")
	}
}

func TestHandlerCompositor_Add(t *testing.T) {
	buf := &bytes.Buffer{}
	handler1 := slog.NewTextHandler(buf, nil)
	handler2 := slog.NewJSONHandler(buf, nil)

	compositor := NewHandlerCompositor()
	compositor.Add(handler1)
	compositor.Add(handler2)

	// Should have added handlers (implementation detail)
}

func TestHandlerCompositor_Multi(t *testing.T) {
	buf := &bytes.Buffer{}
	handler1 := slog.NewTextHandler(buf, nil)
	handler2 := slog.NewJSONHandler(buf, nil)

	compositor := NewHandlerCompositor()
	compositor.Add(handler1)
	compositor.Add(handler2)

	multiHandler := compositor.Multi()

	if multiHandler == nil {
		t.Error("expected Multi to return a handler")
	}
}

func TestHandlerCompositor_Chain(t *testing.T) {
	buf := &bytes.Buffer{}
	handler := slog.NewTextHandler(buf, nil)

	middleware1 := handlerMiddlewareFunc(func(ctx context.Context, record slog.Record, next HandlerFunc) error {
		return next(ctx, record)
	})

	middleware2 := handlerMiddlewareFunc(func(ctx context.Context, record slog.Record, next HandlerFunc) error {
		record.AddAttrs(slog.String("chained", "true"))
		return next(ctx, record)
	})

	compositor := NewHandlerCompositor()
	compositor.Add(handler) // Add handler first

	chainedHandler := compositor.Chain(middleware1, middleware2)

	if chainedHandler == nil {
		t.Error("expected Chain to return a handler")
	}
}

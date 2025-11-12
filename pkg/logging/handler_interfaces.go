package logging

import (
	"context"
	"log/slog"
)

// UnifiedHandlerInterface provides a consolidated interface for all handler operations.
// This combines functionality from HandlerFactory, NamedHandlerFactory, and HandlerMiddleware.
type UnifiedHandlerInterface interface {
	slog.Handler

	// Factory methods
	Create(config interface{}) (slog.Handler, error)
	Name() string
	ConfigType() interface{}

	// Middleware support
	WithMiddleware(middleware ...HandlerMiddleware) slog.Handler

	// Lifecycle management
	Close() error
}

// BaseHandler provides a foundation for implementing UnifiedHandlerInterface
type BaseHandler struct {
	name       string
	handler    slog.Handler
	middleware []HandlerMiddleware
}

// NewBaseHandler creates a new base handler
func NewBaseHandler(name string, handler slog.Handler) *BaseHandler {
	return &BaseHandler{
		name:    name,
		handler: handler,
	}
}

func (h *BaseHandler) Name() string {
	return h.name
}

func (h *BaseHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

func (h *BaseHandler) Handle(ctx context.Context, record slog.Record) error {
	if len(h.middleware) == 0 {
		return h.handler.Handle(ctx, record)
	}

	// Apply middleware chain
	return h.executeMiddleware(ctx, record, 0)
}

func (h *BaseHandler) executeMiddleware(ctx context.Context, record slog.Record, index int) error {
	if index >= len(h.middleware) {
		return h.handler.Handle(ctx, record)
	}

	next := func(ctx context.Context, record slog.Record) error {
		return h.executeMiddleware(ctx, record, index+1)
	}

	return h.middleware[index].Handle(ctx, record, next)
}

func (h *BaseHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &BaseHandler{
		name:       h.name,
		handler:    h.handler.WithAttrs(attrs),
		middleware: h.middleware,
	}
}

func (h *BaseHandler) WithGroup(name string) slog.Handler {
	return &BaseHandler{
		name:       h.name,
		handler:    h.handler.WithGroup(name),
		middleware: h.middleware,
	}
}

func (h *BaseHandler) WithMiddleware(middleware ...HandlerMiddleware) slog.Handler {
	newMiddleware := make([]HandlerMiddleware, len(h.middleware)+len(middleware))
	copy(newMiddleware, h.middleware)
	copy(newMiddleware[len(h.middleware):], middleware)

	return &BaseHandler{
		name:       h.name,
		handler:    h.handler,
		middleware: newMiddleware,
	}
}

func (h *BaseHandler) Create(config interface{}) (slog.Handler, error) {
	// Default implementation - subclasses should override
	return h, nil
}

func (h *BaseHandler) ConfigType() interface{} {
	// Default implementation - subclasses should override
	return nil
}

func (h *BaseHandler) Close() error {
	// Default implementation - no cleanup needed
	return nil
}

// HandlerCompositor provides utilities for composing handlers
type HandlerCompositor struct {
	handlers []slog.Handler
}

// NewHandlerCompositor creates a new handler compositor
func NewHandlerCompositor() *HandlerCompositor {
	return &HandlerCompositor{}
}

// Add adds a handler to the composition
func (c *HandlerCompositor) Add(handler slog.Handler) *HandlerCompositor {
	c.handlers = append(c.handlers, handler)
	return c
}

// Multi creates a multi-handler that writes to all composed handlers
func (c *HandlerCompositor) Multi() slog.Handler {
	return NewMultiHandler(c.handlers...)
}

// Chain creates a handler chain with middleware
func (c *HandlerCompositor) Chain(middleware ...HandlerMiddleware) slog.Handler {
	if len(c.handlers) == 0 {
		return nil
	}

	base := c.handlers[0]
	if bh, ok := base.(*BaseHandler); ok {
		return bh.WithMiddleware(middleware...)
	}

	// Wrap in BaseHandler to add middleware support
	wrapped := NewBaseHandler("chained", base)
	return wrapped.WithMiddleware(middleware...)
}

// HandlerType represents different types of handlers for easier management
type HandlerType int

const (
	TextHandlerType HandlerType = iota
	JSONHandlerType
	CustomHandlerType
	AsyncHandlerType
	MultiHandlerType
)

// HandlerTypeInfo provides metadata about handler types
type HandlerTypeInfo struct {
	Type        HandlerType
	Name        string
	Description string
}

var HandlerTypes = []HandlerTypeInfo{
	{TextHandlerType, "text", "Text format handler"},
	{JSONHandlerType, "json", "JSON format handler"},
	{CustomHandlerType, "custom", "Custom handler implementation"},
	{AsyncHandlerType, "async", "Asynchronous handler wrapper"},
	{MultiHandlerType, "multi", "Multiple handlers composition"},
}

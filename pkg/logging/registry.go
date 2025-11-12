package logging

import (
	"fmt"
	"log/slog"
	"sync"
)

// NamedHandlerFactory creates slog handlers with the given configuration.
type NamedHandlerFactory interface {
	// Name returns the unique name of this handler factory.
	Name() string

	// Create creates a new handler instance with the given configuration.
	Create(config interface{}) (slog.Handler, error)

	// ConfigType returns the expected configuration type for this factory.
	ConfigType() interface{}
}

// HandlerRegistry manages registered handler factories.
type HandlerRegistry struct {
	mu        sync.RWMutex
	factories map[string]NamedHandlerFactory
}

var (
	defaultRegistry = &HandlerRegistry{
		factories: make(map[string]NamedHandlerFactory),
	}
)

// GetDefaultRegistry returns the default global handler registry.
func GetDefaultRegistry() *HandlerRegistry {
	return defaultRegistry
}

// RegisterHandler registers a handler factory with the default registry.
func RegisterHandler(factory NamedHandlerFactory) error {
	return defaultRegistry.RegisterHandler(factory)
}

// CreateHandler creates a handler using a registered factory from the default registry.
func CreateHandler(name string, config interface{}) (slog.Handler, error) {
	return defaultRegistry.CreateHandler(name, config)
}

// ListHandlers returns all registered handler names from the default registry.
func ListHandlers() []string {
	return defaultRegistry.ListHandlers()
}

// RegisterHandler registers a handler factory.
func (r *HandlerRegistry) RegisterHandler(factory NamedHandlerFactory) error {
	if factory == nil {
		return fmt.Errorf("handler factory cannot be nil")
	}

	name := factory.Name()
	if name == "" {
		return fmt.Errorf("handler factory name cannot be empty")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.factories[name]; exists {
		return fmt.Errorf("handler factory with name %q already registered", name)
	}

	r.factories[name] = factory
	return nil
}

// CreateHandler creates a handler using a registered factory.
func (r *HandlerRegistry) CreateHandler(name string, config interface{}) (slog.Handler, error) {
	r.mu.RLock()
	factory, exists := r.factories[name]
	r.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("no handler factory registered with name %q", name)
	}

	return factory.Create(config)
}

// GetFactory returns a registered handler factory by name.
func (r *HandlerRegistry) GetFactory(name string) (NamedHandlerFactory, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	factory, exists := r.factories[name]
	return factory, exists
}

// ListHandlers returns all registered handler names.
func (r *HandlerRegistry) ListHandlers() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.factories))
	for name := range r.factories {
		names = append(names, name)
	}
	return names
}

// UnregisterHandler removes a handler factory from the registry.
func (r *HandlerRegistry) UnregisterHandler(name string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.factories[name]; !exists {
		return false
	}

	delete(r.factories, name)
	return true
}

// Clear removes all registered handler factories.
func (r *HandlerRegistry) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.factories = make(map[string]NamedHandlerFactory)
}

// Built-in handler factories

// JSONHandlerFactory creates JSON handlers.
type JSONHandlerFactory struct{}

func (f *JSONHandlerFactory) Name() string {
	return jsonFormatString
}

func (f *JSONHandlerFactory) ConfigType() interface{} {
	return &OutputConfig{}
}

func (f *JSONHandlerFactory) Create(config interface{}) (slog.Handler, error) {
	outputConfig, ok := config.(*OutputConfig)
	if !ok {
		return nil, fmt.Errorf("expected *OutputConfig, got %T", config)
	}

	return slog.NewJSONHandler(outputConfig.Writer, &slog.HandlerOptions{
		Level: slog.LevelDebug, // Default level, can be overridden
	}), nil
}

// TextHandlerFactory creates text handlers.
type TextHandlerFactory struct{}

func (f *TextHandlerFactory) Name() string {
	return textFormatString
}

func (f *TextHandlerFactory) ConfigType() interface{} {
	return &OutputConfig{}
}

func (f *TextHandlerFactory) Create(config interface{}) (slog.Handler, error) {
	outputConfig, ok := config.(*OutputConfig)
	if !ok {
		return nil, fmt.Errorf("expected *OutputConfig, got %T", config)
	}

	return slog.NewTextHandler(outputConfig.Writer, &slog.HandlerOptions{
		Level: slog.LevelDebug, // Default level, can be overridden
	}), nil
}

// MultiHandlerFactory creates multi-output handlers.
type MultiHandlerFactory struct{}

func (f *MultiHandlerFactory) Name() string {
	return "multi"
}

type MultiHandlerConfig struct {
	Handlers []slog.Handler
}

func (f *MultiHandlerFactory) ConfigType() interface{} {
	return &MultiHandlerConfig{}
}

func (f *MultiHandlerFactory) Create(config interface{}) (slog.Handler, error) {
	multiConfig, ok := config.(*MultiHandlerConfig)
	if !ok {
		return nil, fmt.Errorf("expected *MultiHandlerConfig, got %T", config)
	}

	if len(multiConfig.Handlers) == 0 {
		return nil, fmt.Errorf("multi handler requires at least one handler")
	}

	return NewMultiHandler(multiConfig.Handlers...), nil
}

// Initialize default handler factories
func init() {
	_ = RegisterHandler(&JSONHandlerFactory{})
	_ = RegisterHandler(&TextHandlerFactory{})
	_ = RegisterHandler(&MultiHandlerFactory{})
}

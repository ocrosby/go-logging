package logging

import (
	"bytes"
	"log/slog"
	"testing"
)

// Mock handler factory for testing
type mockNamedHandlerFactory struct {
	handlerName string
}

func (m *mockNamedHandlerFactory) Name() string {
	return m.handlerName
}

func (m *mockNamedHandlerFactory) ConfigType() interface{} {
	return map[string]interface{}{}
}

func (m *mockNamedHandlerFactory) Create(config interface{}) (slog.Handler, error) {
	buf := &bytes.Buffer{}
	return slog.NewTextHandler(buf, nil), nil
}

func TestGetDefaultRegistry_Simple(t *testing.T) {
	registry := GetDefaultRegistry()
	if registry == nil {
		t.Fatal("expected default registry to be created")
	}

	// Should return same instance
	registry2 := GetDefaultRegistry()
	if registry != registry2 {
		t.Error("expected same registry instance")
	}
}

func TestRegisterHandler_Simple(t *testing.T) {
	factory := &mockNamedHandlerFactory{handlerName: "test-simple"}

	err := RegisterHandler(factory)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Should be able to create handler
	handler, err := CreateHandler("test-simple", map[string]interface{}{})
	if err != nil {
		t.Errorf("unexpected error creating handler: %v", err)
	}

	if handler == nil {
		t.Error("expected handler to be created")
	}
}

func TestCreateHandler_NotFound(t *testing.T) {
	// Try to create a non-existent handler
	_, err := CreateHandler("definitely-non-existent-handler", nil)
	if err == nil {
		t.Error("expected error for non-existent handler")
	}
}

func TestListHandlers_Simple(t *testing.T) {
	handlers := ListHandlers()
	if handlers == nil {
		t.Error("expected handlers list to be non-nil")
	}
}

func TestTextHandlerFactory_Simple(t *testing.T) {
	factory := &TextHandlerFactory{}

	if factory.Name() != "text" {
		t.Errorf("expected name 'text', got %s", factory.Name())
	}

	configType := factory.ConfigType()
	if configType == nil {
		t.Error("expected config type to be non-nil")
	}

	// Create handler with proper config
	config := NewOutputConfig().Build()
	handler, err := factory.Create(config)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if handler == nil {
		t.Error("expected handler to be created")
	}
}

func TestJSONHandlerFactory_Simple(t *testing.T) {
	factory := &JSONHandlerFactory{}

	if factory.Name() != "json" {
		t.Errorf("expected name 'json', got %s", factory.Name())
	}

	configType := factory.ConfigType()
	if configType == nil {
		t.Error("expected config type to be non-nil")
	}

	// Create handler with proper config
	config := NewOutputConfig().Build()
	handler, err := factory.Create(config)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if handler == nil {
		t.Error("expected handler to be created")
	}
}

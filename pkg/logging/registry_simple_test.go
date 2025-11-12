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

func TestUnregisterHandler(t *testing.T) {
	factory := &mockNamedHandlerFactory{handlerName: "test-unregister"}

	// Register handler first
	err := RegisterHandler(factory)
	if err != nil {
		t.Errorf("unexpected error registering handler: %v", err)
	}

	// Verify it exists
	_, err = CreateHandler("test-unregister", map[string]interface{}{})
	if err != nil {
		t.Errorf("handler should exist after registration: %v", err)
	}

	// Unregister it
	registry := GetDefaultRegistry()
	registry.UnregisterHandler("test-unregister")

	// Should not exist anymore
	_, err = CreateHandler("test-unregister", map[string]interface{}{})
	if err == nil {
		t.Error("handler should not exist after unregistration")
	}
}

func TestRegistryClear(t *testing.T) {
	factory1 := &mockNamedHandlerFactory{handlerName: "test-clear-1"}
	factory2 := &mockNamedHandlerFactory{handlerName: "test-clear-2"}

	// Register handlers
	_ = RegisterHandler(factory1)
	_ = RegisterHandler(factory2)

	// Verify they exist
	_, err1 := CreateHandler("test-clear-1", map[string]interface{}{})
	_, err2 := CreateHandler("test-clear-2", map[string]interface{}{})

	if err1 != nil || err2 != nil {
		t.Error("handlers should exist before clear")
	}

	// Clear registry
	registry := GetDefaultRegistry()
	registry.Clear()

	// Should not exist anymore
	_, err1 = CreateHandler("test-clear-1", map[string]interface{}{})
	_, err2 = CreateHandler("test-clear-2", map[string]interface{}{})

	if err1 == nil || err2 == nil {
		t.Error("handlers should not exist after clear")
	}
}

func TestGetFactory(t *testing.T) {
	factory := &mockNamedHandlerFactory{handlerName: "test-get-factory"}

	// Register handler
	_ = RegisterHandler(factory)

	// Get factory
	registry := GetDefaultRegistry()
	retrievedFactory, exists := registry.GetFactory("test-get-factory")

	if !exists {
		t.Error("expected factory to exist")
	}

	if retrievedFactory == nil {
		t.Error("expected factory to be retrieved")
	}

	if retrievedFactory.Name() != "test-get-factory" {
		t.Errorf("expected factory name 'test-get-factory', got %s", retrievedFactory.Name())
	}

	// Test non-existent factory
	nonExistentFactory, exists := registry.GetFactory("non-existent")
	if exists {
		t.Error("expected factory to not exist")
	}
	if nonExistentFactory != nil {
		t.Error("expected nil for non-existent factory")
	}
}

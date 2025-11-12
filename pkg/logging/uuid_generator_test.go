package logging

import (
	"regexp"
	"testing"
)

// MockUUIDGenerator is a test implementation that returns predictable UUIDs
type MockUUIDGenerator struct {
	counter int
	prefix  string
}

func NewMockUUIDGenerator(prefix string) *MockUUIDGenerator {
	return &MockUUIDGenerator{prefix: prefix}
}

func (m *MockUUIDGenerator) Generate() string {
	m.counter++
	return m.prefix + "-" + string(rune('0'+m.counter))
}

func TestDefaultUUIDGenerator_Generate(t *testing.T) {
	generator := NewDefaultUUIDGenerator()

	// Generate a UUID and verify it matches UUID v4 format
	uuid := generator.Generate()

	// UUID v4 regex pattern
	uuidPattern := `^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`
	matched, err := regexp.MatchString(uuidPattern, uuid)
	if err != nil {
		t.Fatalf("regex error: %v", err)
	}

	if !matched {
		t.Errorf("generated UUID %s does not match UUID v4 format", uuid)
	}

	// Verify uniqueness by generating multiple UUIDs
	uuids := make(map[string]bool)
	for i := 0; i < 1000; i++ {
		uuid := generator.Generate()
		if uuids[uuid] {
			t.Errorf("duplicate UUID generated: %s", uuid)
			break
		}
		uuids[uuid] = true
	}
}

func TestSetUUIDGenerator(t *testing.T) {
	// Save original generator
	original := GetUUIDGenerator()
	defer SetUUIDGenerator(original)

	// Test with mock generator
	mock := NewMockUUIDGenerator("test")
	SetUUIDGenerator(mock)

	current := GetUUIDGenerator()
	if current != mock {
		t.Error("SetUUIDGenerator did not set the generator correctly")
	}

	// Test that NewTraceID uses the new generator
	traceID := NewTraceID()
	expected := "test-1"
	if traceID != expected {
		t.Errorf("expected %s, got %s", expected, traceID)
	}
}

func TestSetUUIDGenerator_Nil(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic when setting nil generator")
		}
	}()

	SetUUIDGenerator(nil)
}

func TestNewTraceID_UsesDefaultGenerator(t *testing.T) {
	// Save original generator
	original := GetUUIDGenerator()
	defer SetUUIDGenerator(original)

	// Set mock generator
	mock := NewMockUUIDGenerator("trace")
	SetUUIDGenerator(mock)

	// Test that NewTraceID uses the injected generator
	id1 := NewTraceID()
	id2 := NewTraceID()

	if id1 != "trace-1" {
		t.Errorf("expected trace-1, got %s", id1)
	}
	if id2 != "trace-2" {
		t.Errorf("expected trace-2, got %s", id2)
	}
}

func TestUUIDGeneratorInterface(t *testing.T) {
	// Ensure DefaultUUIDGenerator implements UUIDGenerator interface
	var _ UUIDGenerator = &DefaultUUIDGenerator{}
	var _ UUIDGenerator = &MockUUIDGenerator{}
}

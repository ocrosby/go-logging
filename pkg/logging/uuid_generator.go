package logging

import (
	"crypto/rand"
	"fmt"
)

// UUIDGenerator defines the interface for generating unique identifiers.
// This allows users to provide custom UUID generation strategies.
type UUIDGenerator interface {
	// Generate creates a new UUID string
	Generate() string
}

// DefaultUUIDGenerator implements UUIDGenerator using crypto/rand with UUID v4 format.
// This is the default implementation that maintains backward compatibility.
type DefaultUUIDGenerator struct{}

// NewDefaultUUIDGenerator creates a new instance of DefaultUUIDGenerator.
func NewDefaultUUIDGenerator() *DefaultUUIDGenerator {
	return &DefaultUUIDGenerator{}
}

// Generate creates a new UUID v4 string using crypto/rand.
// This maintains the same behavior as the original NewTraceID function.
func (g *DefaultUUIDGenerator) Generate() string {
	var uuid [16]byte
	_, err := rand.Read(uuid[:])
	if err != nil {
		panic(err)
	}
	uuid[6] = (uuid[6] & 0x0f) | 0x40
	uuid[8] = (uuid[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:])
}

// defaultGenerator is the package-level default generator
var defaultGenerator UUIDGenerator = NewDefaultUUIDGenerator()

// SetUUIDGenerator sets the package-level UUID generator.
// This allows users to inject their own UUID generation strategy.
//
// Example:
//
//	logging.SetUUIDGenerator(&MyCustomUUIDGenerator{})
func SetUUIDGenerator(generator UUIDGenerator) {
	if generator == nil {
		panic("UUID generator cannot be nil")
	}
	defaultGenerator = generator
}

// GetUUIDGenerator returns the current package-level UUID generator.
func GetUUIDGenerator() UUIDGenerator {
	return defaultGenerator
}

package logging

import (
	"testing"
)

func TestLevelString(t *testing.T) {
	tests := []struct {
		level    Level
		expected string
	}{
		{TraceLevel, "TRACE"},
		{DebugLevel, "DEBUG"},
		{InfoLevel, "INFO"},
		{WarnLevel, "WARN"},
		{ErrorLevel, "ERROR"},
		{CriticalLevel, "CRITICAL"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if got := tt.level.String(); got != tt.expected {
				t.Errorf("Level.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestParseLevel(t *testing.T) {
	tests := []struct {
		input    string
		expected Level
		ok       bool
	}{
		{"TRACE", TraceLevel, true},
		{"trace", TraceLevel, true},
		{"DEBUG", DebugLevel, true},
		{"debug", DebugLevel, true},
		{"INFO", InfoLevel, true},
		{"info", InfoLevel, true},
		{"WARN", WarnLevel, true},
		{"warn", WarnLevel, true},
		{"ERROR", ErrorLevel, true},
		{"error", ErrorLevel, true},
		{"CRITICAL", CriticalLevel, true},
		{"critical", CriticalLevel, true},
		{"INVALID", TraceLevel, false},
		{"", TraceLevel, false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, ok := ParseLevel(tt.input)
			if ok != tt.ok {
				t.Errorf("ParseLevel(%q) ok = %v, want %v", tt.input, ok, tt.ok)
			}
			if ok && got != tt.expected {
				t.Errorf("ParseLevel(%q) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}

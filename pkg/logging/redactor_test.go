package logging

import (
	"regexp"
	"testing"
)

func TestRedactAPIKeys(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "URL with apiKey",
			input:    "https://api.example.com/v1/data?apiKey=abcd1234qwer",
			expected: "https://api.example.com/v1/data?apiKey=abcd123...<REDACTED>",
		},
		{
			name:     "URL with apikey (lowercase)",
			input:    "https://api.example.com/v1/data?apikey=abcd1234qwer",
			expected: "https://api.example.com/v1/data?apikey=abcd123...<REDACTED>",
		},
		{
			name:     "URL with apiKey and other params",
			input:    "https://api.example.com/v1/data?language=en&apiKey=abcd1234qwer&units=m",
			expected: "https://api.example.com/v1/data?language=en&apiKey=abcd123...<REDACTED>&units=m",
		},
		{
			name:     "URL without apiKey",
			input:    "https://api.example.com/v1/data?language=en",
			expected: "https://api.example.com/v1/data?language=en",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RedactAPIKeys(tt.input)
			if got != tt.expected {
				t.Errorf("RedactAPIKeys() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestRegexRedactor(t *testing.T) {
	pattern := regexp.MustCompile(`password=\w+`)
	redactor := NewRegexRedactor(pattern, "password=***")

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Redact password",
			input:    "login?username=user&password=secret123",
			expected: "login?username=user&password=***",
		},
		{
			name:     "No password",
			input:    "login?username=user",
			expected: "login?username=user",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := redactor.Redact(tt.input)
			if got != tt.expected {
				t.Errorf("RegexRedactor.Redact() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestRedactorChain(t *testing.T) {
	pattern1 := regexp.MustCompile(`password=\w+`)
	pattern2 := regexp.MustCompile(`token=\w+`)

	chain := NewRedactorChain()
	chain.AddRedactor(NewRegexRedactor(pattern1, "password=***"))
	chain.AddRedactor(NewRegexRedactor(pattern2, "token=***"))

	input := "url?password=secret&token=abc123"
	expected := "url?password=***&token=***"

	got := chain.Redact(input)
	if got != expected {
		t.Errorf("RedactorChain.Redact() = %v, want %v", got, expected)
	}
}

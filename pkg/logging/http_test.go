package logging

import (
	"net/http"
	"regexp"
	"strings"
	"testing"
)

func TestRedactedURL(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		expected string
	}{
		{
			name:     "Simple path",
			url:      "/",
			expected: "/",
		},
		{
			name:     "URL with apiKey",
			url:      "https://xyz/v1/geocode/1/2/indices/achePain/daypart/3day.json?apiKey=abcd1234qwer",
			expected: "https://xyz/v1/geocode/1/2/indices/achePain/daypart/3day.json?apiKey=abcd123...<REDACTED>",
		},
		{
			name:     "URL with apiKey and language param",
			url:      "https://xyz/v1/geocode/1/2/indices/achePain/daypart/3day.json?language=en&apiKey=abcd1234qwer",
			expected: "https://xyz/v1/geocode/1/2/indices/achePain/daypart/3day.json?language=en&apiKey=abcd123...<REDACTED>",
		},
		{
			name:     "URL without apiKey",
			url:      "https://xyz/v1/geocode/1/2/indices/achePain/daypart/3day.json?language=en",
			expected: "https://xyz/v1/geocode/1/2/indices/achePain/daypart/3day.json?language=en",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RedactedURL(tt.url)
			if got != tt.expected {
				t.Errorf("RedactedURL() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestRequestHeaders(t *testing.T) {
	tests := []struct {
		name           string
		headersToPrint []string
		reqHeaders     map[string][]string
		expected       string
	}{
		{
			name:           "Single header, single value",
			headersToPrint: []string{"X-Test-Header"},
			reqHeaders:     map[string][]string{"X-Test-Header": {"value1"}},
			expected:       "X-Test-Header: value1",
		},
		{
			name:           "Single header, multiple values",
			headersToPrint: []string{"X-Test-Header"},
			reqHeaders:     map[string][]string{"X-Test-Header": {"value1", "value2"}},
			expected:       "X-Test-Header: value1, value2",
		},
		{
			name:           "Multiple headers",
			headersToPrint: []string{"X-Test-Header", "X-Another-Header"},
			reqHeaders:     map[string][]string{"X-Test-Header": {"value1"}, "X-Another-Header": {"value2"}},
			expected:       "X-Test-Header: value1 | X-Another-Header: value2",
		},
		{
			name:           "No matching headers",
			headersToPrint: []string{"X-Nonexistent-Header"},
			reqHeaders:     map[string][]string{"X-Test-Header": {"value1"}},
			expected:       "",
		},
		{
			name:           "Empty headersToPrint",
			headersToPrint: []string{},
			reqHeaders:     map[string][]string{"X-Test-Header": {"value1"}},
			expected:       "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &http.Request{
				Header: make(http.Header),
			}

			for key, values := range tt.reqHeaders {
				req.Header[key] = values
			}

			got := RequestHeaders(req, tt.headersToPrint)
			if got != tt.expected {
				t.Errorf("RequestHeaders() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestGetDefaultHeaders(t *testing.T) {
	req := &http.Request{
		Header: make(http.Header),
	}
	req.Header["User-Agent"] = []string{"TestAgent/1.0"}

	expected := "User-Agent: TestAgent/1.0"
	got := GetDefaultHeaders(req)

	if got != expected {
		t.Errorf("GetDefaultHeaders() = %v, want %v", got, expected)
	}
}

func TestLogHTTPRequest(t *testing.T) {
	// Create a test logger
	logger := NewWithLevel(InfoLevel)

	// Create a test request
	req, err := http.NewRequest("POST", "https://api.example.com/users?id=123", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "test-client/1.0")

	// This should not panic
	LogHTTPRequest(logger, req, "Content-Type", "User-Agent")
}

func TestLogHTTPRequest_DefaultHeaders(t *testing.T) {
	logger := NewWithLevel(InfoLevel)

	req, err := http.NewRequest("GET", "https://example.com", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	// Should use default headers when none specified
	LogHTTPRequest(logger, req)
}

func TestLogHTTPResponse(t *testing.T) {
	logger := NewWithLevel(InfoLevel)

	// This should not panic
	LogHTTPResponse(logger, 200, "https://api.example.com/users")
}

func TestLogHTTPResponse_WithRedaction(t *testing.T) {
	logger := NewWithLevel(InfoLevel)

	// URL with API key should be redacted
	LogHTTPResponse(logger, 404, "https://api.example.com/users?apiKey=secret123")
}

func TestRedactedURL_WithCustomRedactor(t *testing.T) {
	// Create a custom redactor
	re := regexp.MustCompile(`password=\w+`)
	customRedactor := NewRegexRedactor(re, "[PASSWORD_REDACTED]")

	url := "https://example.com/login?user=john&password=secret123"
	result := RedactedURL(url, customRedactor)

	if !strings.Contains(result, "[PASSWORD_REDACTED]") {
		t.Errorf("expected custom redaction, got %s", result)
	}

	if strings.Contains(result, "secret123") {
		t.Error("expected password to be redacted")
	}
}

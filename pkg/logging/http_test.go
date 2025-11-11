package logging

import (
	"net/http"
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

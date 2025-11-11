package logging

import (
	"fmt"
	"net/http"
	"strings"
)

func RedactedURL(url string, redactors ...Redactor) string {
	result := url
	if len(redactors) == 0 {
		return RedactAPIKeys(url)
	}
	for _, redactor := range redactors {
		result = redactor.Redact(result)
	}
	return result
}

func RequestHeaders(r *http.Request, headersToPrint []string) string {
	var sb strings.Builder

	for i, headerName := range headersToPrint {
		if values, ok := r.Header[headerName]; ok {
			if i > 0 {
				sb.WriteString(" | ")
			}
			sb.WriteString(fmt.Sprintf("%s: %s", headerName, strings.Join(values, ", ")))
		}
	}

	return sb.String()
}

func GetDefaultHeaders(r *http.Request) string {
	return RequestHeaders(r, []string{"User-Agent"})
}

func LogHTTPRequest(logger Logger, r *http.Request, headers ...string) {
	if len(headers) == 0 {
		headers = []string{"User-Agent"}
	}

	url := RedactedURL(r.URL.String())
	headerStr := RequestHeaders(r, headers)

	logger.Info("HTTP Request: %s %s | Headers: %s", r.Method, url, headerStr)
}

func LogHTTPResponse(logger Logger, statusCode int, url string) {
	logger.Info("HTTP Response: %d | URL: %s", statusCode, RedactedURL(url))
}

package main

import (
	"bytes"
	"io"
	"log/slog"
	"testing"
	"time"
)

func BenchmarkCLFHandler(b *testing.B) {
	handler := NewCommonLogFormatHandler(io.Discard)

	handlerWithAttrs := handler.WithAttrs([]slog.Attr{
		slog.String("host", "192.168.1.100"),
		slog.String("ident", "-"),
		slog.String("authuser", "testuser"),
		slog.String("method", "GET"),
		slog.String("path", "/api/test"),
		slog.String("protocol", "HTTP/1.1"),
		slog.String("status", "200"),
		slog.String("bytes", "1234"),
	})

	logger := slog.New(handlerWithAttrs)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		logger.Info("Request processed")
	}
}

func BenchmarkCLFHandlerParallel(b *testing.B) {
	handler := NewCommonLogFormatHandler(io.Discard)

	handlerWithAttrs := handler.WithAttrs([]slog.Attr{
		slog.String("host", "192.168.1.100"),
		slog.String("ident", "-"),
		slog.String("authuser", "testuser"),
		slog.String("method", "GET"),
		slog.String("path", "/api/test"),
		slog.String("protocol", "HTTP/1.1"),
		slog.String("status", "200"),
		slog.String("bytes", "1234"),
	})

	logger := slog.New(handlerWithAttrs)

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("Request processed")
		}
	})
}

func BenchmarkCLFHandlerWithAttrs(b *testing.B) {
	handler := NewCommonLogFormatHandler(io.Discard)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		handlerWithAttrs := handler.WithAttrs([]slog.Attr{
			slog.String("host", "192.168.1.100"),
			slog.String("ident", "-"),
			slog.String("authuser", "testuser"),
			slog.String("method", "GET"),
			slog.String("path", "/api/test"),
			slog.String("protocol", "HTTP/1.1"),
			slog.String("status", "200"),
			slog.String("bytes", "1234"),
		})

		logger := slog.New(handlerWithAttrs)
		logger.Info("Request processed")
	}
}

func TestCLFHandlerOutput(t *testing.T) {
	var buf bytes.Buffer
	handler := NewCommonLogFormatHandler(&buf)

	handlerWithAttrs := handler.WithAttrs([]slog.Attr{
		slog.String("host", "192.168.1.100"),
		slog.String("ident", "-"),
		slog.String("authuser", "alice"),
		slog.String("method", "GET"),
		slog.String("path", "/index.html"),
		slog.String("protocol", "HTTP/1.1"),
		slog.String("status", "200"),
		slog.String("bytes", "9481"),
	})

	logger := slog.New(handlerWithAttrs)
	logger.Info("Request processed")

	handler.Flush()

	output := buf.String()

	expectedParts := []string{
		"192.168.1.100",
		"-",
		"alice",
		"GET /index.html HTTP/1.1",
		"200",
		"9481",
	}

	for _, part := range expectedParts {
		if !bytes.Contains([]byte(output), []byte(part)) {
			t.Errorf("Expected output to contain %q, got: %s", part, output)
		}
	}
}

func TestCLFHandlerFlush(t *testing.T) {
	var buf bytes.Buffer
	handler := NewCommonLogFormatHandler(&buf)

	handlerWithAttrs := handler.WithAttrs([]slog.Attr{
		slog.String("host", "127.0.0.1"),
		slog.String("ident", "-"),
		slog.String("authuser", "-"),
		slog.String("method", "POST"),
		slog.String("path", "/api/test"),
		slog.String("protocol", "HTTP/1.1"),
		slog.String("status", "201"),
		slog.String("bytes", "42"),
	})

	logger := slog.New(handlerWithAttrs)
	logger.Info("Request 1")
	logger.Info("Request 2")
	logger.Info("Request 3")

	if buf.Len() == 0 {
		err := handler.Flush()
		if err != nil {
			t.Errorf("Flush failed: %v", err)
		}
	}

	if buf.Len() == 0 {
		t.Error("Expected buffered output after flush")
	}
}

func TestCLFHandlerConcurrency(t *testing.T) {
	handler := NewCommonLogFormatHandler(io.Discard)

	handlerWithAttrs := handler.WithAttrs([]slog.Attr{
		slog.String("host", "192.168.1.100"),
		slog.String("ident", "-"),
		slog.String("authuser", "testuser"),
		slog.String("method", "GET"),
		slog.String("path", "/api/test"),
		slog.String("protocol", "HTTP/1.1"),
		slog.String("status", "200"),
		slog.String("bytes", "1234"),
	})

	logger := slog.New(handlerWithAttrs)

	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 100; j++ {
				logger.Info("Concurrent request")
			}
			done <- true
		}()
	}

	for i := 0; i < 10; i++ {
		<-done
	}
}

func TestCLFHandlerDefaults(t *testing.T) {
	var buf bytes.Buffer
	handler := NewCommonLogFormatHandler(&buf)

	logger := slog.New(handler)
	logger.Info("Request without attrs")

	handler.Flush()

	output := buf.String()

	expectedDefaults := []string{
		"-",
		"GET",
		"/",
		"HTTP/1.1",
		"200",
		"0",
	}

	for _, def := range expectedDefaults {
		if !bytes.Contains([]byte(output), []byte(def)) {
			t.Errorf("Expected output to contain default %q, got: %s", def, output)
		}
	}
}

func TestCLFHandlerTimestamp(t *testing.T) {
	var buf bytes.Buffer
	handler := NewCommonLogFormatHandler(&buf)

	logger := slog.New(handler)

	now := time.Now()
	logger.Info("Request with timestamp")

	handler.Flush()

	output := buf.String()

	expectedMonth := now.Format("Jan")
	expectedYear := now.Format("2006")

	if !bytes.Contains([]byte(output), []byte(expectedMonth)) {
		t.Errorf("Expected timestamp to contain month %q, got: %s", expectedMonth, output)
	}

	if !bytes.Contains([]byte(output), []byte(expectedYear)) {
		t.Errorf("Expected timestamp to contain year %q, got: %s", expectedYear, output)
	}
}

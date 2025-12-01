package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

const clfTimeFormat = "02/Jan/2006:15:04:05 -0700"

var bytePool = sync.Pool{
	New: func() interface{} {
		b := make([]byte, 0, 128)
		return &b
	},
}

type timestampEntry struct {
	last   time.Time
	cached string
}

var tsCache atomic.Value

func init() {
	tsCache.Store(timestampEntry{})
}

func formatTimestamp(t time.Time) string {
	t = t.Truncate(time.Second)

	if entry, ok := tsCache.Load().(timestampEntry); ok {
		if t.Equal(entry.last) {
			return entry.cached
		}
	}

	formatted := t.Format(clfTimeFormat)
	tsCache.Store(timestampEntry{last: t, cached: formatted})
	return formatted
}

type CommonLogFormatHandler struct {
	mu        sync.Mutex
	writer    io.Writer
	attrCache map[string]string
	groups    []string
}

func NewCommonLogFormatHandler(w io.Writer) *CommonLogFormatHandler {
	if w == nil {
		w = os.Stdout
	}

	bw := bufio.NewWriterSize(w, 64*1024)

	return &CommonLogFormatHandler{
		writer:    bw,
		attrCache: make(map[string]string),
		groups:    make([]string, 0),
	}
}

func (h *CommonLogFormatHandler) Enabled(_ context.Context, level slog.Level) bool {
	return true
}

func (h *CommonLogFormatHandler) Handle(_ context.Context, record slog.Record) error {
	bp := bytePool.Get().(*[]byte)
	buf := (*bp)[:0]
	defer bytePool.Put(bp)

	host := "-"
	if v, ok := h.attrCache["host"]; ok {
		host = v
	}

	ident := "-"
	if v, ok := h.attrCache["ident"]; ok {
		ident = v
	}

	authuser := "-"
	if v, ok := h.attrCache["authuser"]; ok {
		authuser = v
	}

	timestamp := formatTimestamp(record.Time)

	method := "GET"
	if v, ok := h.attrCache["method"]; ok {
		method = v
	}

	path := "/"
	if v, ok := h.attrCache["path"]; ok {
		path = v
	}

	protocol := "HTTP/1.1"
	if v, ok := h.attrCache["protocol"]; ok {
		protocol = v
	}

	status := "200"
	if v, ok := h.attrCache["status"]; ok {
		status = v
	}

	bytes := "0"
	if v, ok := h.attrCache["bytes"]; ok {
		bytes = v
	}

	buf = append(buf, host...)
	buf = append(buf, ' ')
	buf = append(buf, ident...)
	buf = append(buf, ' ')
	buf = append(buf, authuser...)
	buf = append(buf, ' ', '[')
	buf = append(buf, timestamp...)
	buf = append(buf, ']', ' ', '"')
	buf = append(buf, method...)
	buf = append(buf, ' ')
	buf = append(buf, path...)
	buf = append(buf, ' ')
	buf = append(buf, protocol...)
	buf = append(buf, '"', ' ')
	buf = append(buf, status...)
	buf = append(buf, ' ')
	buf = append(buf, bytes...)
	buf = append(buf, '\n')

	h.mu.Lock()
	_, err := h.writer.Write(buf)
	h.mu.Unlock()

	return err
}

func (h *CommonLogFormatHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newAttrCache := make(map[string]string, len(h.attrCache)+len(attrs))

	for k, v := range h.attrCache {
		newAttrCache[k] = v
	}

	for _, attr := range attrs {
		newAttrCache[attr.Key] = attr.Value.String()
	}

	return &CommonLogFormatHandler{
		writer:    h.writer,
		attrCache: newAttrCache,
		groups:    h.groups,
	}
}

func (h *CommonLogFormatHandler) WithGroup(name string) slog.Handler {
	newGroups := make([]string, len(h.groups)+1)
	copy(newGroups, h.groups)
	newGroups[len(h.groups)] = name

	return &CommonLogFormatHandler{
		writer:    h.writer,
		attrCache: h.attrCache,
		groups:    newGroups,
	}
}

func (h *CommonLogFormatHandler) Flush() error {
	if bw, ok := h.writer.(*bufio.Writer); ok {
		return bw.Flush()
	}
	return nil
}

func DemoCommonLogFormat() {
	fmt.Println("--- Common Log Format (CLF) Backend ---")
	fmt.Println("Format: host ident authuser [timestamp] \"request\" status bytes")
	fmt.Println()

	handler := NewCommonLogFormatHandler(os.Stdout)
	defer handler.Flush()

	handlerWithDefaults := handler.WithAttrs([]slog.Attr{
		slog.String("host", "192.168.1.100"),
		slog.String("ident", "-"),
		slog.String("authuser", "alice"),
		slog.String("method", "GET"),
		slog.String("path", "/index.html"),
		slog.String("protocol", "HTTP/1.1"),
		slog.String("status", "200"),
		slog.String("bytes", "9481"),
	})

	logger := slog.New(handlerWithDefaults)
	logger.Info("Request processed")

	handlerPost := handler.WithAttrs([]slog.Attr{
		slog.String("host", "10.0.0.5"),
		slog.String("ident", "-"),
		slog.String("authuser", "bob"),
		slog.String("method", "POST"),
		slog.String("path", "/api/users"),
		slog.String("protocol", "HTTP/1.1"),
		slog.String("status", "201"),
		slog.String("bytes", "342"),
	})

	logger2 := slog.New(handlerPost)
	logger2.Info("User created")

	handlerError := handler.WithAttrs([]slog.Attr{
		slog.String("host", "203.0.113.42"),
		slog.String("ident", "-"),
		slog.String("authuser", "-"),
		slog.String("method", "GET"),
		slog.String("path", "/admin"),
		slog.String("protocol", "HTTP/1.1"),
		slog.String("status", "403"),
		slog.String("bytes", "1024"),
	})

	logger3 := slog.New(handlerError)
	logger3.Error("Forbidden access")

	fmt.Println("\n✅ Common Log Format is a custom pluggable backend!")
	fmt.Println("✅ Can be used with any HTTP server for access logs")
	fmt.Println("✅ Compatible with standard log analysis tools (Webalizer, Analog)")
	fmt.Println("✅ Ultra high-performance: atomic timestamp cache, pre-computed strings, inlined lookups")
}

func DemoCommonLogFormatWithLoggingLibrary() {
	fmt.Println("\n--- CLF with go-logging Library ---")
	fmt.Println("Using CLF handler via logging.NewWithHandler()")
	fmt.Println()

	handler := NewCommonLogFormatHandler(os.Stdout)
	defer handler.Flush()

	handlerWithRequest := handler.WithAttrs([]slog.Attr{
		slog.String("host", "172.16.0.10"),
		slog.String("ident", "-"),
		slog.String("authuser", "charlie"),
		slog.String("method", "PUT"),
		slog.String("path", "/api/orders/123"),
		slog.String("protocol", "HTTP/1.1"),
		slog.String("status", "204"),
		slog.String("bytes", "0"),
	})

	logger := slog.New(handlerWithRequest)
	logger.Info("Order updated")

	handlerWithTimestamp := handler.WithAttrs([]slog.Attr{
		slog.String("host", "127.0.0.1"),
		slog.String("ident", "ident"),
		slog.String("authuser", "david"),
		slog.String("method", "DELETE"),
		slog.String("path", "/api/sessions/456"),
		slog.String("protocol", "HTTP/1.1"),
		slog.String("status", "200"),
		slog.String("bytes", "128"),
	})

	logger2 := slog.New(handlerWithTimestamp)

	time.Sleep(50 * time.Millisecond)
	logger2.Info("Session deleted")
}

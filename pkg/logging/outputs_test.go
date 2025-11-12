package logging

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

type closableBuffer struct {
	*bytes.Buffer
	closed bool
	mu     sync.Mutex
}

func (cb *closableBuffer) Close() error {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.closed = true
	return nil
}

func (cb *closableBuffer) IsClosed() bool {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	return cb.closed
}

func TestNewWriterOutput(t *testing.T) {
	buf := &bytes.Buffer{}
	output := NewWriterOutput(buf)

	if output == nil {
		t.Fatal("expected output to be created")
	}

	if output.writer != buf {
		t.Error("expected writer to match provided buffer")
	}
}

func TestWriterOutput_Write(t *testing.T) {
	buf := &bytes.Buffer{}
	output := NewWriterOutput(buf)

	testData := []byte("test log message\n")
	err := output.Write(testData)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if buf.String() != "test log message\n" {
		t.Errorf("expected 'test log message\\n', got '%s'", buf.String())
	}
}

func TestWriterOutput_Close(t *testing.T) {
	// Test with closable writer
	buf := &closableBuffer{Buffer: &bytes.Buffer{}}
	output := NewWriterOutput(buf)

	err := output.Close()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if !buf.IsClosed() {
		t.Error("expected buffer to be closed")
	}

	// Test with non-closable writer
	simpleBuf := &bytes.Buffer{}
	output2 := NewWriterOutput(simpleBuf)

	err = output2.Close()
	if err != nil {
		t.Errorf("unexpected error for non-closable writer: %v", err)
	}
}

func TestNewFileOutput(t *testing.T) {
	tempDir := t.TempDir()
	filename := filepath.Join(tempDir, "test.log")

	output, err := NewFileOutput(filename)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer output.Close()

	if output.filename != filename {
		t.Errorf("expected filename '%s', got '%s'", filename, output.filename)
	}
}

func TestNewFileOutput_CreateDirectory(t *testing.T) {
	tempDir := t.TempDir()
	filename := filepath.Join(tempDir, "subdir", "test.log")

	output, err := NewFileOutput(filename)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer output.Close()

	// Check that directory was created
	if _, err := os.Stat(filepath.Dir(filename)); os.IsNotExist(err) {
		t.Error("expected directory to be created")
	}
}

func TestFileOutput_Write(t *testing.T) {
	tempDir := t.TempDir()
	filename := filepath.Join(tempDir, "test.log")

	output, err := NewFileOutput(filename)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer output.Close()

	testData := []byte("test log message\n")
	err = output.Write(testData)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Read file contents
	content, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	if string(content) != "test log message\n" {
		t.Errorf("expected 'test log message\\n', got '%s'", string(content))
	}
}

func TestFileOutput_WriteAfterClose(t *testing.T) {
	tempDir := t.TempDir()
	filename := filepath.Join(tempDir, "test.log")

	output, err := NewFileOutput(filename)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Close the output
	err = output.Close()
	if err != nil {
		t.Errorf("unexpected error on close: %v", err)
	}

	// Try to write after close
	err = output.Write([]byte("should fail"))
	if err == nil {
		t.Error("expected error when writing to closed output")
	}
}

func TestFileOutput_Close(t *testing.T) {
	tempDir := t.TempDir()
	filename := filepath.Join(tempDir, "test.log")

	output, err := NewFileOutput(filename)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Write some data first
	testData := []byte("test data\n")
	err = output.Write(testData)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Close
	err = output.Close()
	if err != nil {
		t.Errorf("unexpected error on close: %v", err)
	}

	// Multiple closes should be safe
	err = output.Close()
	if err != nil {
		t.Errorf("unexpected error on second close: %v", err)
	}
}

func TestNewBufferedOutput(t *testing.T) {
	buf := &bytes.Buffer{}
	underlying := NewWriterOutput(buf)

	buffered := NewBufferedOutput(underlying, 1024, 100*time.Millisecond)
	defer buffered.Close()

	if buffered.output != underlying {
		t.Error("expected underlying output to match")
	}
}

func TestBufferedOutput_Write(t *testing.T) {
	buf := &bytes.Buffer{}
	underlying := NewWriterOutput(buf)

	buffered := NewBufferedOutput(underlying, 1024, 100*time.Millisecond)
	defer buffered.Close()

	testData := []byte("test message")
	err := buffered.Write(testData)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Force flush
	err = buffered.Flush()
	if err != nil {
		t.Errorf("unexpected error on flush: %v", err)
	}

	if !strings.Contains(buf.String(), "test message") {
		t.Error("expected message to be written after flush")
	}
}

func TestBufferedOutput_AutoFlush(t *testing.T) {
	buf := &bytes.Buffer{}
	underlying := NewWriterOutput(buf)

	// Short flush interval
	buffered := NewBufferedOutput(underlying, 1024, 10*time.Millisecond)
	defer buffered.Close()

	testData := []byte("auto flush test")
	err := buffered.Write(testData)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Wait for auto flush
	time.Sleep(50 * time.Millisecond)

	if !strings.Contains(buf.String(), "auto flush test") {
		t.Error("expected message to be auto-flushed")
	}
}

func TestNewMultiOutput(t *testing.T) {
	buf1 := &bytes.Buffer{}
	buf2 := &bytes.Buffer{}
	out1 := NewWriterOutput(buf1)
	out2 := NewWriterOutput(buf2)

	multi := NewMultiOutput(out1, out2)

	testData := []byte("test message")
	err := multi.Write(testData)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if buf1.String() != "test message" {
		t.Error("expected message in first buffer")
	}

	if buf2.String() != "test message" {
		t.Error("expected message in second buffer")
	}
}

func TestMultiOutput_AddRemoveOutput(t *testing.T) {
	buf1 := &bytes.Buffer{}
	buf2 := &bytes.Buffer{}
	out1 := NewWriterOutput(buf1)
	out2 := NewWriterOutput(buf2)

	multi := NewMultiOutput()
	multi.AddOutput(out1)
	multi.AddOutput(out2)

	// Write to both
	testData := []byte("test message")
	err := multi.Write(testData)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Remove one output
	multi.RemoveOutput(out1)

	// Write again
	testData2 := []byte(" second message")
	err = multi.Write(testData2)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if buf1.String() != "test message" {
		t.Error("expected first buffer to have only first message after removal")
	}

	if !strings.Contains(buf2.String(), "second message") {
		t.Error("expected second buffer to have both messages")
	}
}

func TestMultiOutput_Close(t *testing.T) {
	buf1 := &closableBuffer{Buffer: &bytes.Buffer{}}
	buf2 := &closableBuffer{Buffer: &bytes.Buffer{}}
	out1 := NewWriterOutput(buf1)
	out2 := NewWriterOutput(buf2)

	multi := NewMultiOutput(out1, out2)

	err := multi.Close()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if !buf1.IsClosed() {
		t.Error("expected first buffer to be closed")
	}

	if !buf2.IsClosed() {
		t.Error("expected second buffer to be closed")
	}
}

func TestNewAsyncOutput(t *testing.T) {
	buf := &bytes.Buffer{}
	underlying := NewWriterOutput(buf)

	async := NewAsyncOutput(underlying, 10)
	defer async.Close()

	if async == nil {
		t.Fatal("expected async output to be created")
	}
}

func TestAsyncOutput_Write(t *testing.T) {
	buf := &bytes.Buffer{}
	underlying := NewWriterOutput(buf)

	async := NewAsyncOutput(underlying, 10)
	defer async.Close()

	testData := []byte("async test message")
	err := async.Write(testData)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Wait a bit for async processing
	time.Sleep(10 * time.Millisecond)

	if !strings.Contains(buf.String(), "async test message") {
		t.Error("expected message to be written asynchronously")
	}
}

func TestAsyncOutput_Stop(t *testing.T) {
	buf := &bytes.Buffer{}
	underlying := NewWriterOutput(buf)

	async := NewAsyncOutput(underlying, 10)

	// Write some data
	testData := []byte("stop test")
	err := async.Write(testData)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Stop
	err = async.Stop()
	if err != nil {
		t.Errorf("unexpected error on stop: %v", err)
	}

	// Should still have processed the data
	if !strings.Contains(buf.String(), "stop test") {
		t.Error("expected message to be written before stop")
	}
}

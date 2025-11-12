package logging

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// WriterOutput wraps an io.Writer to implement the Output interface.
type WriterOutput struct {
	writer io.Writer
	mu     sync.Mutex
}

// NewWriterOutput creates a new WriterOutput.
func NewWriterOutput(w io.Writer) *WriterOutput {
	return &WriterOutput{writer: w}
}

// Write writes data to the underlying writer.
func (o *WriterOutput) Write(data []byte) error {
	o.mu.Lock()
	defer o.mu.Unlock()
	_, err := o.writer.Write(data)
	return err
}

// Close closes the output if the underlying writer implements io.Closer.
func (o *WriterOutput) Close() error {
	if closer, ok := o.writer.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}

// FileOutput writes log entries to a file.
type FileOutput struct {
	filename string
	file     *os.File
	mu       sync.Mutex
}

// NewFileOutput creates a new FileOutput that writes to the specified file.
func NewFileOutput(filename string) (*FileOutput, error) {
	// Create directory if it doesn't exist
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	return &FileOutput{
		filename: filename,
		file:     file,
	}, nil
}

// Write writes data to the file.
func (o *FileOutput) Write(data []byte) error {
	o.mu.Lock()
	defer o.mu.Unlock()

	if o.file == nil {
		return fmt.Errorf("file output is closed")
	}

	_, err := o.file.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write to log file: %w", err)
	}

	// Sync to ensure data is written to disk
	return o.file.Sync()
}

// Close closes the file.
func (o *FileOutput) Close() error {
	o.mu.Lock()
	defer o.mu.Unlock()

	if o.file != nil {
		err := o.file.Close()
		o.file = nil
		return err
	}
	return nil
}

// BufferedOutput buffers writes and flushes them periodically or when full.
type BufferedOutput struct {
	output        Output
	buffer        *bufio.Writer
	bufferSize    int
	flushTimer    *time.Timer
	flushInterval time.Duration
	mu            sync.Mutex
	closed        bool
}

// NewBufferedOutput creates a new BufferedOutput with the specified buffer size and flush interval.
func NewBufferedOutput(output Output, bufferSize int, flushInterval time.Duration) *BufferedOutput {
	bo := &BufferedOutput{
		output:        output,
		bufferSize:    bufferSize,
		flushInterval: flushInterval,
	}

	// Create a buffer that writes to our output
	bo.buffer = bufio.NewWriterSize(&outputWriter{output: output}, bufferSize)

	// Start periodic flush timer
	if flushInterval > 0 {
		bo.flushTimer = time.AfterFunc(flushInterval, bo.periodicFlush)
	}

	return bo
}

// outputWriter is a helper to make Output compatible with io.Writer
type outputWriter struct {
	output Output
}

func (ow *outputWriter) Write(p []byte) (n int, err error) {
	err = ow.output.Write(p)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

// Write writes data to the buffer.
func (bo *BufferedOutput) Write(data []byte) error {
	bo.mu.Lock()
	defer bo.mu.Unlock()

	if bo.closed {
		return fmt.Errorf("buffered output is closed")
	}

	_, err := bo.buffer.Write(data)
	return err
}

// Flush forces all buffered data to be written to the underlying output.
func (bo *BufferedOutput) Flush() error {
	bo.mu.Lock()
	defer bo.mu.Unlock()

	if bo.closed {
		return fmt.Errorf("buffered output is closed")
	}

	return bo.buffer.Flush()
}

// periodicFlush is called by the timer to flush buffered data.
func (bo *BufferedOutput) periodicFlush() {
	_ = bo.Flush()

	// Restart the timer
	bo.mu.Lock()
	if !bo.closed && bo.flushInterval > 0 {
		bo.flushTimer = time.AfterFunc(bo.flushInterval, bo.periodicFlush)
	}
	bo.mu.Unlock()
}

// Close flushes any remaining data and closes the underlying output.
func (bo *BufferedOutput) Close() error {
	bo.mu.Lock()
	defer bo.mu.Unlock()

	if bo.closed {
		return nil
	}

	bo.closed = true

	// Stop the flush timer
	if bo.flushTimer != nil {
		bo.flushTimer.Stop()
	}

	// Flush any remaining data
	if err := bo.buffer.Flush(); err != nil {
		return err
	}

	// Close the underlying output
	return bo.output.Close()
}

// MultiOutput writes to multiple outputs simultaneously.
type MultiOutput struct {
	outputs []Output
	mu      sync.RWMutex
}

// NewMultiOutput creates a new MultiOutput that writes to multiple outputs.
func NewMultiOutput(outputs ...Output) *MultiOutput {
	return &MultiOutput{outputs: outputs}
}

// Write writes data to all outputs.
func (mo *MultiOutput) Write(data []byte) error {
	mo.mu.RLock()
	defer mo.mu.RUnlock()

	var firstErr error
	for _, output := range mo.outputs {
		if err := output.Write(data); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}

// Close closes all outputs.
func (mo *MultiOutput) Close() error {
	mo.mu.Lock()
	defer mo.mu.Unlock()

	var firstErr error
	for _, output := range mo.outputs {
		if err := output.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}

// AddOutput adds a new output to the multi-output.
func (mo *MultiOutput) AddOutput(output Output) {
	mo.mu.Lock()
	defer mo.mu.Unlock()
	mo.outputs = append(mo.outputs, output)
}

// RemoveOutput removes an output from the multi-output.
func (mo *MultiOutput) RemoveOutput(output Output) {
	mo.mu.Lock()
	defer mo.mu.Unlock()

	for i, o := range mo.outputs {
		if o == output {
			mo.outputs = append(mo.outputs[:i], mo.outputs[i+1:]...)
			return
		}
	}
}

// AsyncOutput processes writes asynchronously in a background goroutine.
type AsyncOutput struct {
	output Output
	worker *AsyncWorker[[]byte]
}

// NewAsyncOutput creates a new AsyncOutput with the specified queue size.
func NewAsyncOutput(output Output, queueSize int) *AsyncOutput {
	ao := &AsyncOutput{output: output}

	ao.worker = NewAsyncWorker(AsyncWorkerConfig[[]byte]{
		QueueSize: queueSize,
		Processor: func(data []byte) error {
			return ao.output.Write(data)
		},
	})

	return ao
}

// Write queues data for asynchronous writing.
func (ao *AsyncOutput) Write(data []byte) error {
	if ao.worker.IsClosed() {
		return fmt.Errorf("async output is closed")
	}

	// Make a copy of the data since it might be modified by the caller
	dataCopy := make([]byte, len(data))
	copy(dataCopy, data)

	if ao.worker.Submit(dataCopy) {
		return nil
	}
	// Queue is full, write synchronously as fallback
	return ao.output.Write(data)
}

// Stop gracefully shuts down the async processing.
func (ao *AsyncOutput) Stop() error {
	return ao.worker.Stop()
}

// Close stops async processing and closes the underlying output.
func (ao *AsyncOutput) Close() error {
	if err := ao.Stop(); err != nil {
		return err
	}
	return ao.output.Close()
}

// RotatingFileOutput writes to files with automatic rotation based on size or time.
type RotatingFileOutput struct {
	pattern     string        // File pattern with placeholders
	maxSize     int64         // Maximum size in bytes
	maxAge      time.Duration // Maximum age
	current     *os.File
	currentSize int64
	mu          sync.Mutex
}

// NewRotatingFileOutput creates a new rotating file output.
func NewRotatingFileOutput(pattern string, maxSize int64, maxAge time.Duration) *RotatingFileOutput {
	return &RotatingFileOutput{
		pattern: pattern,
		maxSize: maxSize,
		maxAge:  maxAge,
	}
}

// Write writes data to the current file, rotating if necessary.
func (rfo *RotatingFileOutput) Write(data []byte) error {
	rfo.mu.Lock()
	defer rfo.mu.Unlock()

	// Check if we need to rotate
	if rfo.shouldRotate(int64(len(data))) {
		if err := rfo.rotate(); err != nil {
			return fmt.Errorf("failed to rotate log file: %w", err)
		}
	}

	// Ensure we have a current file
	if rfo.current == nil {
		if err := rfo.openNew(); err != nil {
			return fmt.Errorf("failed to open log file: %w", err)
		}
	}

	// Write the data
	n, err := rfo.current.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write to log file: %w", err)
	}

	rfo.currentSize += int64(n)

	// Sync to ensure data is written
	return rfo.current.Sync()
}

// shouldRotate determines if the log file should be rotated.
func (rfo *RotatingFileOutput) shouldRotate(dataSize int64) bool {
	if rfo.current == nil {
		return false
	}

	// Check size limit
	if rfo.maxSize > 0 && rfo.currentSize+dataSize > rfo.maxSize {
		return true
	}

	// Check age limit (simplified - would need to track file creation time)
	if rfo.maxAge > 0 {
		if stat, err := rfo.current.Stat(); err == nil {
			if time.Since(stat.ModTime()) > rfo.maxAge {
				return true
			}
		}
	}

	return false
}

// rotate closes the current file and prepares for a new one.
func (rfo *RotatingFileOutput) rotate() error {
	if rfo.current != nil {
		if err := rfo.current.Close(); err != nil {
			return err
		}
		rfo.current = nil
		rfo.currentSize = 0
	}
	return nil
}

// openNew opens a new log file.
func (rfo *RotatingFileOutput) openNew() error {
	// Generate filename from pattern (simplified)
	filename := rfo.generateFilename()

	// Create directory if it doesn't exist
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}

	rfo.current = file

	// Get current size if file already exists
	if stat, err := file.Stat(); err == nil {
		rfo.currentSize = stat.Size()
	}

	return nil
}

// generateFilename generates a filename from the pattern.
func (rfo *RotatingFileOutput) generateFilename() string {
	// Simple implementation - replace placeholders with current timestamp
	now := time.Now()
	filename := rfo.pattern

	// Replace common placeholders
	filename = fmt.Sprintf(filename, now.Format("2006-01-02-15-04-05"))

	return filename
}

// Close closes the current file.
func (rfo *RotatingFileOutput) Close() error {
	rfo.mu.Lock()
	defer rfo.mu.Unlock()

	if rfo.current != nil {
		err := rfo.current.Close()
		rfo.current = nil
		rfo.currentSize = 0
		return err
	}
	return nil
}

package logging

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestNewAsyncWorker(t *testing.T) {
	processed := make([]string, 0)
	var mu sync.Mutex

	processor := func(item string) error {
		mu.Lock()
		processed = append(processed, item)
		mu.Unlock()
		return nil
	}

	config := AsyncWorkerConfig[string]{
		QueueSize: 10,
		Processor: processor,
	}

	worker := NewAsyncWorker(config)
	defer func() { _ = worker.Stop() }()

	if worker == nil {
		t.Fatal("expected worker to be created")
	}

	if worker.QueueCapacity() != 10 {
		t.Errorf("expected queue capacity 10, got %d", worker.QueueCapacity())
	}
}

func TestNewAsyncWorker_PanicOnNilProcessor(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic when processor is nil")
		}
	}()

	config := AsyncWorkerConfig[string]{
		QueueSize: 10,
		Processor: nil,
	}

	NewAsyncWorker(config)
}

func TestAsyncWorker_Submit(t *testing.T) {
	var processed []string
	var mu sync.Mutex

	processor := func(item string) error {
		mu.Lock()
		processed = append(processed, item)
		mu.Unlock()
		return nil
	}

	config := AsyncWorkerConfig[string]{
		QueueSize: 2,
		Processor: processor,
	}

	worker := NewAsyncWorker(config)
	defer func() { _ = worker.Stop() }()

	// Test successful submit
	if !worker.Submit("item1") {
		t.Error("expected submit to succeed")
	}

	if !worker.Submit("item2") {
		t.Error("expected submit to succeed")
	}

	// Test queue full - may or may not fail depending on timing
	// so we don't test this specific case as it's timing dependent

	// Wait for processing
	time.Sleep(10 * time.Millisecond)

	mu.Lock()
	if len(processed) < 1 {
		t.Error("expected at least one item to be processed")
	}
	mu.Unlock()
}

func TestAsyncWorker_SubmitBlocking(t *testing.T) {
	var processed []string
	var mu sync.Mutex

	processor := func(item string) error {
		mu.Lock()
		processed = append(processed, item)
		mu.Unlock()
		time.Sleep(5 * time.Millisecond) // Slow processing
		return nil
	}

	config := AsyncWorkerConfig[string]{
		QueueSize: 1,
		Processor: processor,
	}

	worker := NewAsyncWorker(config)
	defer func() { _ = worker.Stop() }()

	// Submit first item
	if !worker.SubmitBlocking("item1") {
		t.Error("expected blocking submit to succeed")
	}

	// Submit second item (should block briefly but succeed)
	done := make(chan bool)
	go func() {
		result := worker.SubmitBlocking("item2")
		done <- result
	}()

	select {
	case result := <-done:
		if !result {
			t.Error("expected blocking submit to succeed")
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("blocking submit took too long")
	}

	// Wait for processing
	time.Sleep(20 * time.Millisecond)

	mu.Lock()
	if len(processed) != 2 {
		t.Errorf("expected 2 items processed, got %d", len(processed))
	}
	mu.Unlock()
}

func TestAsyncWorker_Stop(t *testing.T) {
	var processed int32

	processor := func(item string) error {
		atomic.AddInt32(&processed, 1)
		return nil
	}

	config := AsyncWorkerConfig[string]{
		QueueSize: 10,
		Processor: processor,
	}

	worker := NewAsyncWorker(config)

	// Submit some items
	worker.Submit("item1")
	worker.Submit("item2")
	worker.Submit("item3")

	// Stop worker
	err := worker.Stop()
	if err != nil {
		t.Errorf("expected no error on stop, got %v", err)
	}

	// Verify worker is closed
	if !worker.IsClosed() {
		t.Error("expected worker to be closed after stop")
	}

	// Verify submissions fail after stop
	if worker.Submit("item4") {
		t.Error("expected submit to fail after stop")
	}

	if worker.SubmitBlocking("item5") {
		t.Error("expected blocking submit to fail after stop")
	}

	// Multiple stops should be safe
	err = worker.Stop()
	if err != nil {
		t.Errorf("expected no error on second stop, got %v", err)
	}
}

func TestAsyncWorker_WithShutdownCallback(t *testing.T) {
	var shutdownCalled bool

	processor := func(item string) error {
		return nil
	}

	onShutdown := func() error {
		shutdownCalled = true
		return nil
	}

	config := AsyncWorkerConfig[string]{
		QueueSize:  5,
		Processor:  processor,
		OnShutdown: onShutdown,
	}

	worker := NewAsyncWorker(config)
	worker.Submit("item1")

	err := worker.Stop()
	if err != nil {
		t.Errorf("expected no error on stop, got %v", err)
	}

	if !shutdownCalled {
		t.Error("expected shutdown callback to be called")
	}
}

func TestAsyncWorker_QueueMethods(t *testing.T) {
	processor := func(item string) error {
		time.Sleep(50 * time.Millisecond) // Slow processing
		return nil
	}

	config := AsyncWorkerConfig[string]{
		QueueSize: 5,
		Processor: processor,
	}

	worker := NewAsyncWorker(config)
	defer func() { _ = worker.Stop() }()

	// Test initial state
	if worker.QueueSize() != 0 {
		t.Errorf("expected queue size 0, got %d", worker.QueueSize())
	}

	if worker.QueueCapacity() != 5 {
		t.Errorf("expected queue capacity 5, got %d", worker.QueueCapacity())
	}

	// Add items
	worker.Submit("item1")
	worker.Submit("item2")

	// Check queue size (might be 0, 1, or 2 depending on processing speed)
	queueSize := worker.QueueSize()
	if queueSize < 0 || queueSize > 2 {
		t.Errorf("expected queue size 0-2, got %d", queueSize)
	}
}

func TestAsyncWorker_DrainOnShutdown(t *testing.T) {
	var processed []string
	var mu sync.Mutex

	processor := func(item string) error {
		mu.Lock()
		processed = append(processed, item)
		mu.Unlock()
		return nil
	}

	config := AsyncWorkerConfig[string]{
		QueueSize: 10,
		Processor: processor,
	}

	worker := NewAsyncWorker(config)

	// Fill queue
	for i := 0; i < 5; i++ {
		worker.Submit("item")
	}

	// Stop immediately (should drain queue)
	_ = worker.Stop()

	mu.Lock()
	processedCount := len(processed)
	mu.Unlock()

	if processedCount != 5 {
		t.Errorf("expected 5 items processed during drain, got %d", processedCount)
	}
}

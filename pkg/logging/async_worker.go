package logging

import (
	"sync"
)

// AsyncWorker provides a unified pattern for asynchronous processing with proper shutdown.
// This extracts the common pattern used in AsyncOutput and AsyncHandler.
type AsyncWorker[T any] struct {
	queue      chan T
	done       chan struct{}
	wg         sync.WaitGroup
	closed     bool
	mu         sync.Mutex
	processor  func(T) error
	onShutdown func() error
}

// AsyncWorkerConfig configures an AsyncWorker
type AsyncWorkerConfig[T any] struct {
	QueueSize  int
	Processor  func(T) error
	OnShutdown func() error // Optional cleanup on shutdown
}

// NewAsyncWorker creates a new async worker with the specified configuration
func NewAsyncWorker[T any](config AsyncWorkerConfig[T]) *AsyncWorker[T] {
	if config.Processor == nil {
		panic("processor function is required")
	}

	worker := &AsyncWorker[T]{
		queue:      make(chan T, config.QueueSize),
		done:       make(chan struct{}),
		processor:  config.Processor,
		onShutdown: config.OnShutdown,
	}

	worker.wg.Add(1)
	go worker.run()

	return worker
}

// Submit adds an item to the processing queue
func (w *AsyncWorker[T]) Submit(item T) bool {
	w.mu.Lock()
	closed := w.closed
	w.mu.Unlock()

	if closed {
		return false
	}

	select {
	case w.queue <- item:
		return true
	default:
		return false // Queue is full
	}
}

// SubmitBlocking adds an item to the queue, blocking if queue is full
func (w *AsyncWorker[T]) SubmitBlocking(item T) bool {
	w.mu.Lock()
	closed := w.closed
	w.mu.Unlock()

	if closed {
		return false
	}

	w.queue <- item
	return true
}

// run is the main worker loop
func (w *AsyncWorker[T]) run() {
	defer w.wg.Done()

	for {
		select {
		case item := <-w.queue:
			_ = w.processor(item)
		case <-w.done:
			// Drain remaining items
			for {
				select {
				case item := <-w.queue:
					_ = w.processor(item)
				default:
					if w.onShutdown != nil {
						_ = w.onShutdown()
					}
					return
				}
			}
		}
	}
}

// Stop gracefully shuts down the worker
func (w *AsyncWorker[T]) Stop() error {
	w.mu.Lock()
	if w.closed {
		w.mu.Unlock()
		return nil
	}
	w.closed = true
	w.mu.Unlock()

	close(w.done)
	w.wg.Wait()
	return nil
}

// IsClosed returns whether the worker is closed
func (w *AsyncWorker[T]) IsClosed() bool {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.closed
}

// QueueSize returns the current number of items in the queue
func (w *AsyncWorker[T]) QueueSize() int {
	return len(w.queue)
}

// QueueCapacity returns the maximum capacity of the queue
func (w *AsyncWorker[T]) QueueCapacity() int {
	return cap(w.queue)
}

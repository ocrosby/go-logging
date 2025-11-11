package logging

import (
	"context"
	"log/slog"
	"sync"
)

type HandlerFactory interface {
	Create(output interface{}) slog.Handler
}

type handlerFactoryFunc func(interface{}) slog.Handler

func (f handlerFactoryFunc) Create(output interface{}) slog.Handler {
	return f(output)
}

type MultiHandler struct {
	handlers []slog.Handler
	mu       sync.RWMutex
}

func NewMultiHandler(handlers ...slog.Handler) *MultiHandler {
	return &MultiHandler{
		handlers: handlers,
	}
}

func (h *MultiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, handler := range h.handlers {
		if handler.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

func (h *MultiHandler) Handle(ctx context.Context, record slog.Record) error {
	h.mu.RLock()
	defer h.mu.RUnlock()

	var firstErr error
	for _, handler := range h.handlers {
		if handler.Enabled(ctx, record.Level) {
			if err := handler.Handle(ctx, record.Clone()); err != nil && firstErr == nil {
				firstErr = err
			}
		}
	}
	return firstErr
}

func (h *MultiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	h.mu.RLock()
	defer h.mu.RUnlock()

	newHandlers := make([]slog.Handler, len(h.handlers))
	for i, handler := range h.handlers {
		newHandlers[i] = handler.WithAttrs(attrs)
	}
	return &MultiHandler{handlers: newHandlers}
}

func (h *MultiHandler) WithGroup(name string) slog.Handler {
	h.mu.RLock()
	defer h.mu.RUnlock()

	newHandlers := make([]slog.Handler, len(h.handlers))
	for i, handler := range h.handlers {
		newHandlers[i] = handler.WithGroup(name)
	}
	return &MultiHandler{handlers: newHandlers}
}

func (h *MultiHandler) AddHandler(handler slog.Handler) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.handlers = append(h.handlers, handler)
}

func (h *MultiHandler) RemoveHandler(handler slog.Handler) {
	h.mu.Lock()
	defer h.mu.Unlock()

	for i, hdlr := range h.handlers {
		if hdlr == handler {
			h.handlers = append(h.handlers[:i], h.handlers[i+1:]...)
			return
		}
	}
}

type ConditionalHandler struct {
	handler   slog.Handler
	condition func(context.Context, slog.Record) bool
}

func NewConditionalHandler(handler slog.Handler, condition func(context.Context, slog.Record) bool) *ConditionalHandler {
	return &ConditionalHandler{
		handler:   handler,
		condition: condition,
	}
}

func (h *ConditionalHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

func (h *ConditionalHandler) Handle(ctx context.Context, record slog.Record) error {
	if h.condition != nil && !h.condition(ctx, record) {
		return nil
	}
	return h.handler.Handle(ctx, record)
}

func (h *ConditionalHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &ConditionalHandler{
		handler:   h.handler.WithAttrs(attrs),
		condition: h.condition,
	}
}

func (h *ConditionalHandler) WithGroup(name string) slog.Handler {
	return &ConditionalHandler{
		handler:   h.handler.WithGroup(name),
		condition: h.condition,
	}
}

type BufferedHandler struct {
	handler slog.Handler
	buffer  []slog.Record
	maxSize int
	mu      sync.Mutex
	flushFn func([]slog.Record) error
}

func NewBufferedHandler(handler slog.Handler, maxSize int) *BufferedHandler {
	return &BufferedHandler{
		handler: handler,
		buffer:  make([]slog.Record, 0, maxSize),
		maxSize: maxSize,
	}
}

func (h *BufferedHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

func (h *BufferedHandler) Handle(ctx context.Context, record slog.Record) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.buffer = append(h.buffer, record)

	if len(h.buffer) >= h.maxSize {
		return h.flushInternal(ctx)
	}

	return nil
}

func (h *BufferedHandler) Flush(ctx context.Context) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.flushInternal(ctx)
}

func (h *BufferedHandler) flushInternal(ctx context.Context) error {
	if len(h.buffer) == 0 {
		return nil
	}

	var firstErr error
	for _, record := range h.buffer {
		if err := h.handler.Handle(ctx, record); err != nil && firstErr == nil {
			firstErr = err
		}
	}

	h.buffer = h.buffer[:0]
	return firstErr
}

func (h *BufferedHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &BufferedHandler{
		handler: h.handler.WithAttrs(attrs),
		buffer:  make([]slog.Record, 0, h.maxSize),
		maxSize: h.maxSize,
		flushFn: h.flushFn,
	}
}

func (h *BufferedHandler) WithGroup(name string) slog.Handler {
	return &BufferedHandler{
		handler: h.handler.WithGroup(name),
		buffer:  make([]slog.Record, 0, h.maxSize),
		maxSize: h.maxSize,
		flushFn: h.flushFn,
	}
}

type AsyncHandler struct {
	handler slog.Handler
	queue   chan slog.Record
	done    chan struct{}
	wg      sync.WaitGroup
}

func NewAsyncHandler(handler slog.Handler, queueSize int) *AsyncHandler {
	ah := &AsyncHandler{
		handler: handler,
		queue:   make(chan slog.Record, queueSize),
		done:    make(chan struct{}),
	}

	ah.wg.Add(1)
	go ah.worker()

	return ah
}

func (h *AsyncHandler) worker() {
	defer h.wg.Done()

	for {
		select {
		case record := <-h.queue:
			_ = h.handler.Handle(context.Background(), record)
		case <-h.done:
			for record := range h.queue {
				_ = h.handler.Handle(context.Background(), record)
			}
			return
		}
	}
}

func (h *AsyncHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

func (h *AsyncHandler) Handle(ctx context.Context, record slog.Record) error {
	select {
	case h.queue <- record:
		return nil
	default:
		return h.handler.Handle(ctx, record)
	}
}

func (h *AsyncHandler) Close() {
	close(h.done)
	h.wg.Wait()
	close(h.queue)
}

func (h *AsyncHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return NewAsyncHandler(h.handler.WithAttrs(attrs), cap(h.queue))
}

func (h *AsyncHandler) WithGroup(name string) slog.Handler {
	return NewAsyncHandler(h.handler.WithGroup(name), cap(h.queue))
}

type RotatingHandler struct {
	handlers []slog.Handler
	current  int
	mu       sync.Mutex
}

func NewRotatingHandler(handlers ...slog.Handler) *RotatingHandler {
	return &RotatingHandler{
		handlers: handlers,
		current:  0,
	}
}

func (h *RotatingHandler) Enabled(ctx context.Context, level slog.Level) bool {
	h.mu.Lock()
	defer h.mu.Unlock()

	if len(h.handlers) == 0 {
		return false
	}
	return h.handlers[h.current].Enabled(ctx, level)
}

func (h *RotatingHandler) Handle(ctx context.Context, record slog.Record) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if len(h.handlers) == 0 {
		return nil
	}

	err := h.handlers[h.current].Handle(ctx, record)
	h.current = (h.current + 1) % len(h.handlers)
	return err
}

func (h *RotatingHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	h.mu.Lock()
	defer h.mu.Unlock()

	newHandlers := make([]slog.Handler, len(h.handlers))
	for i, handler := range h.handlers {
		newHandlers[i] = handler.WithAttrs(attrs)
	}
	return &RotatingHandler{handlers: newHandlers}
}

func (h *RotatingHandler) WithGroup(name string) slog.Handler {
	h.mu.Lock()
	defer h.mu.Unlock()

	newHandlers := make([]slog.Handler, len(h.handlers))
	for i, handler := range h.handlers {
		newHandlers[i] = handler.WithGroup(name)
	}
	return &RotatingHandler{handlers: newHandlers}
}

type HandlerBuilder struct {
	handler     slog.Handler
	middlewares []HandlerMiddleware
	extractors  []ContextExtractor
}

func NewHandlerBuilder(handler slog.Handler) *HandlerBuilder {
	return &HandlerBuilder{
		handler:     handler,
		middlewares: make([]HandlerMiddleware, 0),
		extractors:  make([]ContextExtractor, 0),
	}
}

func (b *HandlerBuilder) WithMiddleware(middleware HandlerMiddleware) *HandlerBuilder {
	b.middlewares = append(b.middlewares, middleware)
	return b
}

func (b *HandlerBuilder) WithTimestamp() *HandlerBuilder {
	return b.WithMiddleware(TimestampMiddleware())
}

func (b *HandlerBuilder) WithStaticFields(fields map[string]interface{}) *HandlerBuilder {
	return b.WithMiddleware(StaticFieldsMiddleware(fields))
}

func (b *HandlerBuilder) WithContextExtractor(extractor ContextExtractor) *HandlerBuilder {
	b.extractors = append(b.extractors, extractor)
	return b
}

func (b *HandlerBuilder) WithTraceContext() *HandlerBuilder {
	return b.WithContextExtractor(TraceContextExtractor())
}

func (b *HandlerBuilder) WithLevelFilter(minLevel slog.Level) *HandlerBuilder {
	return b.WithMiddleware(LevelFilterMiddleware(minLevel))
}

func (b *HandlerBuilder) WithRedaction(redactor Redactor) *HandlerBuilder {
	return b.WithMiddleware(RedactionMiddleware(redactor))
}

func (b *HandlerBuilder) WithSampling(rate int) *HandlerBuilder {
	return b.WithMiddleware(SamplingMiddleware(rate))
}

func (b *HandlerBuilder) Build() slog.Handler {
	if len(b.extractors) > 0 {
		composite := NewCompositeContextExtractor(b.extractors...)
		b.middlewares = append(b.middlewares, ContextExtractorMiddleware(composite))
	}

	if len(b.middlewares) == 0 {
		return b.handler
	}

	return NewMiddlewareHandler(b.handler, b.middlewares...)
}

func MultiHandlerBuilder(handlers ...slog.Handler) *HandlerBuilder {
	return NewHandlerBuilder(NewMultiHandler(handlers...))
}

func ConditionalHandlerBuilder(handler slog.Handler, condition func(context.Context, slog.Record) bool) *HandlerBuilder {
	return NewHandlerBuilder(NewConditionalHandler(handler, condition))
}

func BufferedHandlerBuilder(handler slog.Handler, maxSize int) *HandlerBuilder {
	return NewHandlerBuilder(NewBufferedHandler(handler, maxSize))
}

func AsyncHandlerBuilder(handler slog.Handler, queueSize int) *HandlerBuilder {
	return NewHandlerBuilder(NewAsyncHandler(handler, queueSize))
}

package events

import (
	"context"
	"log/slog"
	"sync"
)

type InProcessBus struct {
	handlers map[string][]Handler
	mu       sync.RWMutex
	logger   *slog.Logger
}

func NewInProcessBus(logger *slog.Logger) *InProcessBus {
	return &InProcessBus{
		handlers: make(map[string][]Handler),
		logger:   logger,
	}
}

func (b *InProcessBus) Publish(ctx context.Context, topic string, event Event) error {
	b.mu.RLock()
	handlers := make([]Handler, len(b.handlers[topic]))
	copy(handlers, b.handlers[topic])
	b.mu.RUnlock()

	var wg sync.WaitGroup
	for _, handler := range handlers {
		wg.Add(1)
		go func(h Handler) {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					b.logger.Error("event handler panic",
						"topic", topic, "event_type", event.Type, "panic", r)
				}
			}()
			if err := h(ctx, event); err != nil {
				b.logger.Error("event handler error",
					"topic", topic, "event_type", event.Type, "error", err)
			}
		}(handler)
	}
	wg.Wait()

	b.logger.Info("event published",
		"topic", topic, "event_type", event.Type, "event_id", event.ID)

	return nil
}

func (b *InProcessBus) Subscribe(topic string, handler Handler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[topic] = append(b.handlers[topic], handler)
}

func (b *InProcessBus) Close() error {
	return nil
}

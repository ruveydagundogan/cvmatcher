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
	handlers := b.handlers[topic]
	b.mu.RUnlock()

	for _, handler := range handlers {
		go func(h Handler) {
			if err := h(ctx, event); err != nil {
				b.logger.Error("event handler error",
					"topic", topic,
					"event_type", event.Type,
					"error", err,
				)
			}
		}(handler)
	}

	b.logger.Info("event published",
		"topic", topic,
		"event_type", event.Type,
		"event_id", event.ID,
	)

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

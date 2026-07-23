package events

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const (
	TopicLLMScoring = "cvmatcher.scoring"
	TopicIAM        = "cvmatcher.iam"
	TopicAudit      = "cvmatcher.audit"
)

type Event struct {
	ID        string      `json:"id"`
	Type      string      `json:"type"`
	Source    string      `json:"source"`
	Data      interface{} `json:"data"`
	Timestamp time.Time   `json:"timestamp"`
}

type Handler func(ctx context.Context, event Event) error

type EventBus interface {
	Publish(ctx context.Context, topic string, event Event) error
	Subscribe(topic string, handler Handler)
	Close() error
}

func NewEvent(eventType, source string, data interface{}) Event {
	return Event{
		ID:        uuid.New().String(),
		Type:      eventType,
		Source:    source,
		Data:      data,
		Timestamp: time.Now().UTC(),
	}
}

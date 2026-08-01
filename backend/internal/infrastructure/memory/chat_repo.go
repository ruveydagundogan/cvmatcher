package memory

import (
	"context"
	"sync"
	"time"

	chatmodel "github.com/ruveydagundogan/cvmatcher/backend/internal/domain/chat/model"
)

type InMemoryChatRepository struct {
	mu    sync.RWMutex
	convs map[string]*chatmodel.Conversation
	msgs  map[string][]*chatmodel.Message
}

func NewInMemoryChatRepository() *InMemoryChatRepository {
	return &InMemoryChatRepository{
		convs: make(map[string]*chatmodel.Conversation),
		msgs:  make(map[string][]*chatmodel.Message),
	}
}

func (r *InMemoryChatRepository) ListByUser(ctx context.Context, userID string) ([]*chatmodel.Conversation, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*chatmodel.Conversation, 0)
	for _, c := range r.convs {
		if c.UserID == userID {
			result = append(result, c)
		}
	}
	return result, nil
}

func (r *InMemoryChatRepository) FindByID(ctx context.Context, id string) (*chatmodel.Conversation, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	c, ok := r.convs[id]
	if !ok {
		return nil, nil
	}
	return c, nil
}

func (r *InMemoryChatRepository) SaveConversation(ctx context.Context, conv *chatmodel.Conversation) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.convs[conv.ID] = conv
	return nil
}

func (r *InMemoryChatRepository) DeleteConversation(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.convs, id)
	delete(r.msgs, id)
	return nil
}

func (r *InMemoryChatRepository) TouchConversation(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if c, ok := r.convs[id]; ok {
		c.UpdatedAt = time.Now().UTC()
	}
	return nil
}

func (r *InMemoryChatRepository) ListMessages(ctx context.Context, conversationID string) ([]*chatmodel.Message, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.msgs[conversationID], nil
}

func (r *InMemoryChatRepository) SaveMessage(ctx context.Context, msg *chatmodel.Message) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.msgs[msg.ConversationID] = append(r.msgs[msg.ConversationID], msg)
	return nil
}

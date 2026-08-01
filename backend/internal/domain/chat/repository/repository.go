package repository

import (
	"context"

	chatmodel "github.com/ruveydagundogan/cvmatcher/backend/internal/domain/chat/model"
)

type ChatRepository interface {
	// Conversations
	ListByUser(ctx context.Context, userID string) ([]*chatmodel.Conversation, error)
	FindByID(ctx context.Context, id string) (*chatmodel.Conversation, error)
	SaveConversation(ctx context.Context, conv *chatmodel.Conversation) error
	DeleteConversation(ctx context.Context, id string) error
	TouchConversation(ctx context.Context, id string) error

	// Messages
	ListMessages(ctx context.Context, conversationID string) ([]*chatmodel.Message, error)
	SaveMessage(ctx context.Context, msg *chatmodel.Message) error
}

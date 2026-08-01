package model

import (
	"time"

	"github.com/google/uuid"
)

type Conversation struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Title     string    `json:"title"`
	CVID      string    `json:"cv_id,omitempty"`
	JDID      string    `json:"jd_id,omitempty"`
	MatchID   string    `json:"match_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Message struct {
	ID             string    `json:"id"`
	ConversationID string    `json:"conversation_id"`
	Role           string    `json:"role"` // user | assistant
	Content        string    `json:"content"`
	TokenCount     int       `json:"token_count"`
	CreatedAt      time.Time `json:"created_at"`
}

func NewConversation(userID, title, cvID, jdID, matchID string) *Conversation {
	if title == "" {
		title = "New Chat"
	}
	return &Conversation{
		ID:        uuid.New().String(),
		UserID:    userID,
		Title:     title,
		CVID:      cvID,
		JDID:      jdID,
		MatchID:   matchID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
}

func NewMessage(conversationID, role, content string, tokenCount int) *Message {
	return &Message{
		ID:             uuid.New().String(),
		ConversationID: conversationID,
		Role:           role,
		Content:        content,
		TokenCount:     tokenCount,
		CreatedAt:      time.Now().UTC(),
	}
}

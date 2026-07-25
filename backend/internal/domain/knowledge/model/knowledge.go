package model

import "time"

type KnowledgeEntry struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	Tags        []string  `json:"tags,omitempty"`
	Category    string    `json:"category,omitempty"`
	Source      string    `json:"source,omitempty"`
	Embedding   []float64 `json:"-"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type KnowledgeSearchResult struct {
	Entry   KnowledgeEntry `json:"entry"`
	Score   float64       `json:"score"`
}

func NewKnowledgeEntry(userID, title, content string) *KnowledgeEntry {
	return &KnowledgeEntry{
		UserID:    userID,
		Title:     title,
		Content:   content,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
}

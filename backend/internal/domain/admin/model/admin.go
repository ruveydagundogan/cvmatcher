package model

import "time"

type Adapter struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	FilePath    string    `json:"file_path"`
	Active      bool      `json:"active"`
	ModelName   string    `json:"model_name"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type SystemPrompt struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Content   string    `json:"content"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type LLMSettings struct {
	MaxTokens      int     `json:"max_tokens"`
	Temperature    float64 `json:"temperature"`
	TopP           float64 `json:"top_p"`
	ContextLength  int     `json:"context_length"`
	ModelName      string  `json:"model_name"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type LogEntry struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Query       string    `json:"query"`
	Response    string    `json:"response"`
	Model       string    `json:"model"`
	Adapter     string    `json:"adapter,omitempty"`
	DurationMs  int64     `json:"duration_ms"`
	TokenCount  int       `json:"token_count"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

func NewAdapter(name, description, filePath, modelName string) *Adapter {
	return &Adapter{
		Name:        name,
		Description: description,
		FilePath:    filePath,
		Active:      false,
		ModelName:   modelName,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}
}

func NewSystemPrompt(name, content string) *SystemPrompt {
	return &SystemPrompt{
		Name:      name,
		Content:   content,
		Active:    false,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
}

func DefaultSettings() *LLMSettings {
	return &LLMSettings{
		MaxTokens:     2048,
		Temperature:   0.7,
		TopP:          0.9,
		ContextLength: 4096,
		ModelName:     "qwen2.5:1.5b-instruct",
		UpdatedAt:     time.Now().UTC(),
	}
}

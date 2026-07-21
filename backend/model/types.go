package model

// ScoreRequest represents the incoming request from the frontend.
type ScoreRequest struct {
	Prompt   string `json:"prompt" binding:"required"`
	Response string `json:"response" binding:"required"`
}

// ScoreResponse represents the response sent back to the frontend.
type ScoreResponse struct {
	Score int `json:"score"`
}

// HistoryItem represents one scored LLM interaction.
type HistoryItem struct {
	Prompt   string `json:"prompt"`
	Response string `json:"response"`
	Score    int    `json:"score"`
}
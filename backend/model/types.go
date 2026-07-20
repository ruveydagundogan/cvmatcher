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

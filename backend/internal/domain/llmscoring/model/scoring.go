package model

import (
	"time"

	"github.com/google/uuid"
)

type ScoreRequest struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Prompt      string    `json:"prompt"`
	Response    string    `json:"response"`
	Score       int       `json:"score"`
	Model       string    `json:"model"`
	InferenceMs int64     `json:"inference_ms"`
	WordCount   int       `json:"word_count"`
	CharCount   int       `json:"char_count"`
	Category    string    `json:"category"`
	CreatedAt   time.Time `json:"created_at"`
}

func NewScoreRequest(userID, prompt, response, model string, inferenceMs int64) *ScoreRequest {
	wordCount := countWords(response)
	return &ScoreRequest{
		ID:          uuid.New().String(),
		UserID:      userID,
		Prompt:      prompt,
		Response:    response,
		Model:       model,
		InferenceMs: inferenceMs,
		WordCount:   wordCount,
		CharCount:   len(response),
		CreatedAt:   time.Now().UTC(),
	}
}

func countWords(s string) int {
	count := 0
	inWord := false
	for i := 0; i < len(s); i++ {
		if s[i] == ' ' || s[i] == '\n' || s[i] == '\t' {
			if inWord {
				count++
				inWord = false
			}
		} else {
			inWord = true
		}
	}
	if inWord {
		count++
	}
	return count
}

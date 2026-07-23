package repository

import (
	"context"

	"github.com/ruveydagundogan/cvmatcher/backend/internal/domain/llmscoring/model"
)

type ScoringRepository interface {
	Save(ctx context.Context, score *model.ScoreRequest) error
	FindByID(ctx context.Context, id string) (*model.ScoreRequest, error)
	FindByUserID(ctx context.Context, userID string, offset, limit int) ([]*model.ScoreRequest, int, error)
	DeleteByUserID(ctx context.Context, userID string) error
	GetStats(ctx context.Context, userID string) (*ScoringStats, error)
}

type ScoringStats struct {
	TotalRequests int     `json:"total_requests"`
	AverageScore  float64 `json:"average_score"`
	TotalWords    int     `json:"total_words"`
	TotalChars    int     `json:"total_chars"`
}

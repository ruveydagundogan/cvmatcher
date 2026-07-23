package repository

import (
	"context"

	"github.com/ruveydagundogan/llm-decision-score/backend/internal/domain/matching/model"
)

type MatchingRepository interface {
	Save(ctx context.Context, match *model.MatchResult) error
	FindByID(ctx context.Context, id string) (*model.MatchResult, error)
	FindByCVAndJD(ctx context.Context, cvID, jdID string) (*model.MatchResult, error)
	FindByUserID(ctx context.Context, userID string, offset, limit int) ([]*model.MatchResult, int, error)
	GetDashboardStats(ctx context.Context, userID string) (*model.DashboardStats, error)
}

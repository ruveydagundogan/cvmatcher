package repository

import (
	"context"

	"github.com/ruveydagundogan/cvmatcher/backend/internal/domain/cv/model"
)

type CVRepository interface {
	Save(ctx context.Context, cv *model.CV) error
	FindByID(ctx context.Context, id string) (*model.CV, error)
	FindByUserID(ctx context.Context, userID string, offset, limit int) ([]*model.CV, int, error)
	Update(ctx context.Context, cv *model.CV) error
	Delete(ctx context.Context, id string) error
}

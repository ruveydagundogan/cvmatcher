package repository

import (
	"context"

	"github.com/ruveydagundogan/cvmatcher/backend/internal/domain/knowledge/model"
)

type KnowledgeRepository interface {
	Save(ctx context.Context, entry *model.KnowledgeEntry) error
	FindByID(ctx context.Context, id string) (*model.KnowledgeEntry, error)
	FindByUserID(ctx context.Context, userID string, offset, limit int) ([]*model.KnowledgeEntry, int, error)
	Search(ctx context.Context, query string, tags []string, limit int) ([]*model.KnowledgeSearchResult, error)
	Update(ctx context.Context, entry *model.KnowledgeEntry) error
	Delete(ctx context.Context, id string) error
	ListCategories(ctx context.Context) ([]string, error)
}

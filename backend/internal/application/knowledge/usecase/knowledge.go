package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"
	knowledge "github.com/ruveydagundogan/cvmatcher/backend/internal/domain/knowledge/model"
	"github.com/ruveydagundogan/cvmatcher/backend/internal/domain/knowledge/repository"
	apperrors "github.com/ruveydagundogan/cvmatcher/backend/internal/shared/errors"
)

type KnowledgeUseCase struct {
	repo repository.KnowledgeRepository
	log  *slog.Logger
}

func NewKnowledgeUseCase(repo repository.KnowledgeRepository, log *slog.Logger) *KnowledgeUseCase {
	return &KnowledgeUseCase{repo: repo, log: log}
}

func (uc *KnowledgeUseCase) Create(ctx context.Context, userID, title, content, category, source string, tags []string) (*knowledge.KnowledgeEntry, error) {
	entry := &knowledge.KnowledgeEntry{
		ID:        uuid.New().String(),
		UserID:    userID,
		Title:     title,
		Content:   content,
		Tags:      tags,
		Category:  category,
		Source:    source,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	if err := uc.repo.Save(ctx, entry); err != nil {
		return nil, fmt.Errorf("save knowledge: %w", err)
	}
	return entry, nil
}

func (uc *KnowledgeUseCase) GetByID(ctx context.Context, id string) (*knowledge.KnowledgeEntry, error) {
	return uc.repo.FindByID(ctx, id)
}

func (uc *KnowledgeUseCase) List(ctx context.Context, userID string, offset, limit int) ([]*knowledge.KnowledgeEntry, int, error) {
	return uc.repo.FindByUserID(ctx, userID, offset, limit)
}

func (uc *KnowledgeUseCase) Search(ctx context.Context, query string, tags []string, limit int) ([]*knowledge.KnowledgeSearchResult, error) {
	if query == "" && len(tags) == 0 {
		return nil, apperrors.Validation("query or tags required")
	}
	if limit <= 0 {
		limit = 10
	}
	return uc.repo.Search(ctx, query, tags, limit)
}

func (uc *KnowledgeUseCase) Update(ctx context.Context, entry *knowledge.KnowledgeEntry) error {
	entry.UpdatedAt = time.Now().UTC()
	return uc.repo.Update(ctx, entry)
}

func (uc *KnowledgeUseCase) Delete(ctx context.Context, id string) error {
	return uc.repo.Delete(ctx, id)
}

func (uc *KnowledgeUseCase) ListCategories(ctx context.Context) ([]string, error) {
	return uc.repo.ListCategories(ctx)
}

func (uc *KnowledgeUseCase) IndexContent(ctx context.Context, userID, content string) error {
	lines := strings.Split(content, "\n")
	var title string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			title = line
			if len(title) > 100 {
				title = title[:100]
			}
			break
		}
	}
	if title == "" {
		title = fmt.Sprintf("Knowledge %s", time.Now().Format("2006-01-02"))
	}

	entry := &knowledge.KnowledgeEntry{
		ID:        uuid.New().String(),
		UserID:    userID,
		Title:     title,
		Content:   content,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	return uc.repo.Save(ctx, entry)
}

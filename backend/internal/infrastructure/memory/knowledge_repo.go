package memory

import (
	"context"
	"strings"
	"sync"

	"github.com/ruveydagundogan/cvmatcher/backend/internal/domain/knowledge/model"
)

type KnowledgeRepository struct {
	mu      sync.RWMutex
	entries map[string]*model.KnowledgeEntry
}

func NewKnowledgeRepository() *KnowledgeRepository {
	return &KnowledgeRepository{
		entries: make(map[string]*model.KnowledgeEntry),
	}
}

func (r *KnowledgeRepository) Save(ctx context.Context, entry *model.KnowledgeEntry) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.entries[entry.ID] = entry
	return nil
}

func (r *KnowledgeRepository) FindByID(ctx context.Context, id string) (*model.KnowledgeEntry, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	entry, ok := r.entries[id]
	if !ok {
		return nil, errNotFound
	}
	return entry, nil
}

func (r *KnowledgeRepository) FindByUserID(ctx context.Context, userID string, offset, limit int) ([]*model.KnowledgeEntry, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*model.KnowledgeEntry
	for _, e := range r.entries {
		if e.UserID == userID {
			result = append(result, e)
		}
	}
	total := len(result)
	if offset >= total {
		return []*model.KnowledgeEntry{}, total, nil
	}
	end := offset + limit
	if end > total {
		end = total
	}
	return result[offset:end], total, nil
}

func (r *KnowledgeRepository) Search(ctx context.Context, query string, tags []string, limit int) ([]*model.KnowledgeSearchResult, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	lower := strings.ToLower(query)
	var results []*model.KnowledgeSearchResult
	for _, e := range r.entries {
		content := strings.ToLower(e.Title + " " + e.Content)
		if strings.Contains(content, lower) {
			results = append(results, &model.KnowledgeSearchResult{Entry: *e, Score: 1.0})
		}
	}
	if len(results) > limit {
		results = results[:limit]
	}
	return results, nil
}

func (r *KnowledgeRepository) Update(ctx context.Context, entry *model.KnowledgeEntry) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.entries[entry.ID] = entry
	return nil
}

func (r *KnowledgeRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.entries, id)
	return nil
}

func (r *KnowledgeRepository) ListCategories(ctx context.Context) ([]string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	seen := make(map[string]bool)
	for _, e := range r.entries {
		if e.Category != "" {
			seen[e.Category] = true
		}
	}
	var cats []string
	for c := range seen {
		cats = append(cats, c)
	}
	return cats, nil
}

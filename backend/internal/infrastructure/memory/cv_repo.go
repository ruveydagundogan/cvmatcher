package memory

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"github.com/ruveydagundogan/llm-decision-score/backend/internal/domain/cv/model"
)

type InMemoryCVRepo struct {
	mu   sync.RWMutex
	cvs  map[string]*model.CV
}

func NewInMemoryCVRepo() *InMemoryCVRepo {
	return &InMemoryCVRepo{cvs: make(map[string]*model.CV)}
}

func (r *InMemoryCVRepo) Save(ctx context.Context, cv *model.CV) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.cvs[cv.ID] = cv
	return nil
}

func (r *InMemoryCVRepo) FindByID(ctx context.Context, id string) (*model.CV, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	cv, ok := r.cvs[id]
	if !ok {
		return nil, fmt.Errorf("cv not found")
	}
	return cv, nil
}

func (r *InMemoryCVRepo) FindByUserID(ctx context.Context, userID string, offset, limit int) ([]*model.CV, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var cvs []*model.CV
	for _, cv := range r.cvs {
		if cv.UserID == userID {
			cvs = append(cvs, cv)
		}
	}
	sort.Slice(cvs, func(i, j int) bool {
		return cvs[i].CreatedAt.After(cvs[j].CreatedAt)
	})

	total := len(cvs)
	if offset >= total {
		return []*model.CV{}, total, nil
	}
	end := offset + limit
	if end > total {
		end = total
	}
	return cvs[offset:end], total, nil
}

func (r *InMemoryCVRepo) Update(ctx context.Context, cv *model.CV) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.cvs[cv.ID] = cv
	return nil
}

func (r *InMemoryCVRepo) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.cvs, id)
	return nil
}

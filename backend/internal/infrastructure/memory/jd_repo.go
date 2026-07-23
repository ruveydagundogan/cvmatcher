package memory

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"github.com/ruveydagundogan/llm-decision-score/backend/internal/domain/jobdescription/model"
)

type InMemoryJDRepo struct {
	mu  sync.RWMutex
	jds map[string]*model.JobDescription
}

func NewInMemoryJDRepo() *InMemoryJDRepo {
	return &InMemoryJDRepo{jds: make(map[string]*model.JobDescription)}
}

func (r *InMemoryJDRepo) Save(ctx context.Context, jd *model.JobDescription) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.jds[jd.ID] = jd
	return nil
}

func (r *InMemoryJDRepo) FindByID(ctx context.Context, id string) (*model.JobDescription, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	jd, ok := r.jds[id]
	if !ok {
		return nil, fmt.Errorf("jd not found")
	}
	return jd, nil
}

func (r *InMemoryJDRepo) FindByUserID(ctx context.Context, userID string, offset, limit int) ([]*model.JobDescription, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var jds []*model.JobDescription
	for _, jd := range r.jds {
		if jd.UserID == userID {
			jds = append(jds, jd)
		}
	}
	sort.Slice(jds, func(i, j int) bool {
		return jds[i].CreatedAt.After(jds[j].CreatedAt)
	})

	total := len(jds)
	if offset >= total {
		return []*model.JobDescription{}, total, nil
	}
	end := offset + limit
	if end > total {
		end = total
	}
	return jds[offset:end], total, nil
}

func (r *InMemoryJDRepo) Update(ctx context.Context, jd *model.JobDescription) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.jds[jd.ID] = jd
	return nil
}

func (r *InMemoryJDRepo) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.jds, id)
	return nil
}

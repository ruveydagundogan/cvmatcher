package memory

import (
	"context"
	"fmt"
	"sync"

	"github.com/ruveydagundogan/cvmatcher/backend/internal/domain/iam/model"
)

type InMemoryUserRepo struct {
	mu    sync.RWMutex
	users map[string]*model.User
}

func NewInMemoryUserRepo() *InMemoryUserRepo {
	return &InMemoryUserRepo{users: make(map[string]*model.User)}
}

func (r *InMemoryUserRepo) Create(_ context.Context, user *model.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.users[user.ID] = user
	return nil
}

func (r *InMemoryUserRepo) FindByID(_ context.Context, id string) (*model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if u, ok := r.users[id]; ok {
		return u, nil
	}
	return nil, fmt.Errorf("not found")
}

func (r *InMemoryUserRepo) FindByEmail(_ context.Context, email string) (*model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, u := range r.users {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, fmt.Errorf("not found")
}

func (r *InMemoryUserRepo) Update(_ context.Context, user *model.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.users[user.ID] = user
	return nil
}

func (r *InMemoryUserRepo) Delete(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.users, id)
	return nil
}

func (r *InMemoryUserRepo) List(_ context.Context, offset, limit int) ([]*model.User, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	all := make([]*model.User, 0, len(r.users))
	for _, u := range r.users {
		all = append(all, u)
	}
	total := len(all)
	if offset >= total {
		return nil, total, nil
	}
	end := offset + limit
	if end > total {
		end = total
	}
	return all[offset:end], total, nil
}

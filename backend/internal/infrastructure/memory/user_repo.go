package memory

import (
	"context"
	"sync"
	"time"

	"github.com/ruveydagundogan/cvmatcher/backend/internal/domain/iam/model"
)

type InMemoryUserRepo struct {
	mu    sync.RWMutex
	users map[string]*model.User
}

func NewInMemoryUserRepo() *InMemoryUserRepo {
	repo := &InMemoryUserRepo{users: make(map[string]*model.User)}
	repo.seed()
	return repo
}

func (r *InMemoryUserRepo) seed() {
	now := time.Now()
	r.users["a0000000-0000-0000-0000-000000000001"] = &model.User{
		ID:           "a0000000-0000-0000-0000-000000000001",
		FirstName:    "Admin",
		LastName:     "User",
		Email:        "admin@cvmatcher.com",
		PasswordHash: "$2a$10$AnorUp.ommMnkp1ZUwQ/f.QkWDgirOkanHjMzEMXDhmeeWf2Uq.x2",
		Status:       "active",
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	r.users["a0000000-0000-0000-0000-000000000002"] = &model.User{
		ID:           "a0000000-0000-0000-0000-000000000002",
		FirstName:    "Regular",
		LastName:     "User",
		Email:        "user@cvmatcher.com",
		PasswordHash: "$2a$10$AnorUp.ommMnkp1ZUwQ/f.QkWDgirOkanHjMzEMXDhmeeWf2Uq.x2",
		Status:       "active",
		CreatedAt:    now,
		UpdatedAt:    now,
	}
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
	return nil, nil
}

func (r *InMemoryUserRepo) FindByEmail(_ context.Context, email string) (*model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, u := range r.users {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, nil
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

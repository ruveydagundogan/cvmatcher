package memory

import (
	"context"
	"fmt"
	"sync"

	"github.com/ruveydagundogan/llm-decision-score/backend/internal/domain/iam/model"
)

type InMemoryRoleRepo struct {
	mu    sync.RWMutex
	roles map[string]*model.Role
}

func NewInMemoryRoleRepo() *InMemoryRoleRepo {
	r := &InMemoryRoleRepo{roles: make(map[string]*model.Role)}
	r.roles["user"] = &model.Role{
		ID:          "role-user",
		Name:        "user",
		Permissions: []string{"score:write", "history:read", "history:write", "profile:read", "profile:write"},
	}
	r.roles["admin"] = &model.Role{
		ID:          "role-admin",
		Name:        "admin",
		Permissions: []string{"*"},
	}
	return r
}

func (r *InMemoryRoleRepo) FindByID(_ context.Context, id string) (*model.Role, error) {
	for _, role := range r.roles {
		if role.ID == id {
			return role, nil
		}
	}
	return nil, fmt.Errorf("role not found")
}

func (r *InMemoryRoleRepo) FindByName(_ context.Context, name string) (*model.Role, error) {
	if role, ok := r.roles[name]; ok {
		return role, nil
	}
	return nil, fmt.Errorf("role not found")
}

func (r *InMemoryRoleRepo) AssignToUser(_ context.Context, _, _ string) error {
	return nil
}

func (r *InMemoryRoleRepo) GetUserRole(_ context.Context, _ string) (*model.UserRole, error) {
	return &model.UserRole{RoleName: "user"}, nil
}

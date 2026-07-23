package repository

import (
	"context"

	"github.com/ruveydagundogan/cvmatcher/backend/internal/domain/iam/model"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	FindByID(ctx context.Context, id string) (*model.User, error)
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, offset, limit int) ([]*model.User, int, error)
}

type RoleRepository interface {
	FindByID(ctx context.Context, id string) (*model.Role, error)
	FindByName(ctx context.Context, name string) (*model.Role, error)
	AssignToUser(ctx context.Context, userID, roleID string) error
	GetUserRole(ctx context.Context, userID string) (*model.UserRole, error)
}

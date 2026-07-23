package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ruveydagundogan/cvmatcher/backend/internal/domain/iam/model"
)

type RoleRepository struct {
	pool *pgxpool.Pool
}

func NewRoleRepository(pool *pgxpool.Pool) *RoleRepository {
	return &RoleRepository{pool: pool}
}

func (r *RoleRepository) FindByID(ctx context.Context, id string) (*model.Role, error) {
	role := &model.Role{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, name, permissions FROM roles WHERE id = $1`, id,
	).Scan(&role.ID, &role.Name, &role.Permissions)
	if err != nil {
		return nil, fmt.Errorf("role not found: %w", err)
	}
	return role, nil
}

func (r *RoleRepository) FindByName(ctx context.Context, name string) (*model.Role, error) {
	role := &model.Role{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, name, permissions FROM roles WHERE name = $1`, name,
	).Scan(&role.ID, &role.Name, &role.Permissions)
	if err != nil {
		return nil, fmt.Errorf("role not found: %w", err)
	}
	return role, nil
}

func (r *RoleRepository) AssignToUser(ctx context.Context, userID, roleID string) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO user_roles (user_id, role_id, created_at) VALUES ($1, $2, NOW())
		 ON CONFLICT (user_id) DO UPDATE SET role_id = $2`,
		userID, roleID,
	)
	return err
}

func (r *RoleRepository) GetUserRole(ctx context.Context, userID string) (*model.UserRole, error) {
	ur := &model.UserRole{}
	err := r.pool.QueryRow(ctx,
		`SELECT ur.user_id, ur.role_id, r.name
		 FROM user_roles ur JOIN roles r ON ur.role_id = r.id
		 WHERE ur.user_id = $1`, userID,
	).Scan(&ur.UserID, &ur.RoleID, &ur.RoleName)
	if err != nil {
		return nil, fmt.Errorf("user role not found: %w", err)
	}
	return ur, nil
}

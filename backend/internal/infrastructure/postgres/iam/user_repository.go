package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ruveydagundogan/llm-decision-score/backend/internal/domain/iam/model"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO users (id, email, password_hash, first_name, last_name, status, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		user.ID, user.Email, user.PasswordHash, user.FirstName, user.LastName, user.Status, user.CreatedAt, user.UpdatedAt,
	)
	return err
}

func (r *UserRepository) FindByID(ctx context.Context, id string) (*model.User, error) {
	user := &model.User{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, email, password_hash, first_name, last_name, status, created_at, updated_at
		 FROM users WHERE id = $1`, id,
	).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.FirstName, &user.LastName, &user.Status, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return user, nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	user := &model.User{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, email, password_hash, first_name, last_name, status, created_at, updated_at
		 FROM users WHERE email = $1`, email,
	).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.FirstName, &user.LastName, &user.Status, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return user, nil
}

func (r *UserRepository) Update(ctx context.Context, user *model.User) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE users SET first_name = $1, last_name = $2, status = $3, updated_at = NOW()
		 WHERE id = $4`,
		user.FirstName, user.LastName, user.Status, user.ID,
	)
	return err
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM users WHERE id = $1`, id)
	return err
}

func (r *UserRepository) List(ctx context.Context, offset, limit int) ([]*model.User, int, error) {
	var total int
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM users`).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.pool.Query(ctx,
		`SELECT id, email, password_hash, first_name, last_name, status, created_at, updated_at
		 FROM users ORDER BY created_at DESC LIMIT $1 OFFSET $2`, limit, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	users := make([]*model.User, 0, min(limit, 64))
	for rows.Next() {
		user := &model.User{}
		if err := rows.Scan(&user.ID, &user.Email, &user.PasswordHash, &user.FirstName, &user.LastName, &user.Status, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, 0, err
		}
		users = append(users, user)
	}

	return users, total, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

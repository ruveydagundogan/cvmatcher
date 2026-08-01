package admin

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"

	adminmodel "github.com/ruveydagundogan/cvmatcher/backend/internal/domain/admin/model"
)

type AdminRepository struct {
	pool *pgxpool.Pool
}

func NewAdminRepository(pool *pgxpool.Pool) *AdminRepository {
	return &AdminRepository{pool: pool}
}

func (r *AdminRepository) ListAdapters(ctx context.Context) ([]*adminmodel.Adapter, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, name, description, file_path, active, model_name, created_at, updated_at
		 FROM admin_adapters ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	adapters := make([]*adminmodel.Adapter, 0)
	for rows.Next() {
		a := &adminmodel.Adapter{}
		if err := rows.Scan(&a.ID, &a.Name, &a.Description, &a.FilePath, &a.Active, &a.ModelName, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, err
		}
		adapters = append(adapters, a)
	}
	return adapters, rows.Err()
}

func (r *AdminRepository) SaveAdapter(ctx context.Context, adapter *adminmodel.Adapter) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO admin_adapters (id, name, description, file_path, active, model_name, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
		 ON CONFLICT (id) DO UPDATE SET
		   name=EXCLUDED.name, description=EXCLUDED.description, file_path=EXCLUDED.file_path,
		   active=EXCLUDED.active, model_name=EXCLUDED.model_name, updated_at=NOW()`,
		adapter.ID, adapter.Name, adapter.Description, adapter.FilePath, adapter.Active, adapter.ModelName, adapter.CreatedAt, adapter.UpdatedAt,
	)
	return err
}

func (r *AdminRepository) UpdateAdapter(ctx context.Context, adapter *adminmodel.Adapter) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE admin_adapters SET name=$1, description=$2, file_path=$3, active=$4, model_name=$5, updated_at=NOW() WHERE id=$6`,
		adapter.Name, adapter.Description, adapter.FilePath, adapter.Active, adapter.ModelName, adapter.ID,
	)
	return err
}

func (r *AdminRepository) DeleteAdapter(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM admin_adapters WHERE id=$1`, id)
	return err
}

func (r *AdminRepository) ListPrompts(ctx context.Context) ([]*adminmodel.SystemPrompt, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, name, content, active, created_at, updated_at
		 FROM system_prompts ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	prompts := make([]*adminmodel.SystemPrompt, 0)
	for rows.Next() {
		p := &adminmodel.SystemPrompt{}
		if err := rows.Scan(&p.ID, &p.Name, &p.Content, &p.Active, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		prompts = append(prompts, p)
	}
	return prompts, rows.Err()
}

func (r *AdminRepository) SavePrompt(ctx context.Context, prompt *adminmodel.SystemPrompt) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO system_prompts (id, name, content, active, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6)
		 ON CONFLICT (id) DO UPDATE SET
		   name=EXCLUDED.name, content=EXCLUDED.content, active=EXCLUDED.active, updated_at=NOW()`,
		prompt.ID, prompt.Name, prompt.Content, prompt.Active, prompt.CreatedAt, prompt.UpdatedAt,
	)
	return err
}

func (r *AdminRepository) UpdatePrompt(ctx context.Context, prompt *adminmodel.SystemPrompt) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE system_prompts SET name=$1, content=$2, active=$3, updated_at=NOW() WHERE id=$4`,
		prompt.Name, prompt.Content, prompt.Active, prompt.ID,
	)
	return err
}

func (r *AdminRepository) DeletePrompt(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM system_prompts WHERE id=$1`, id)
	return err
}

func (r *AdminRepository) GetActivePrompt(ctx context.Context) (*adminmodel.SystemPrompt, error) {
	p := &adminmodel.SystemPrompt{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, name, content, active, created_at, updated_at FROM system_prompts WHERE active=true LIMIT 1`,
	).Scan(&p.ID, &p.Name, &p.Content, &p.Active, &p.CreatedAt, &p.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (r *AdminRepository) GetSettings(ctx context.Context) (*adminmodel.LLMSettings, error) {
	s := &adminmodel.LLMSettings{}
	err := r.pool.QueryRow(ctx,
		`SELECT max_tokens, temperature, top_p, context_length, model_name, updated_at
		 FROM llm_settings ORDER BY id LIMIT 1`,
	).Scan(&s.MaxTokens, &s.Temperature, &s.TopP, &s.ContextLength, &s.ModelName, &s.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (r *AdminRepository) SaveSettings(ctx context.Context, settings *adminmodel.LLMSettings) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE llm_settings SET max_tokens=$1, temperature=$2, top_p=$3, context_length=$4, model_name=$5, updated_at=NOW()`,
		settings.MaxTokens, settings.Temperature, settings.TopP, settings.ContextLength, settings.ModelName,
	)
	if err != nil {
		return err
	}
	// Ensure at least one row exists
	_, err = r.pool.Exec(ctx,
		`INSERT INTO llm_settings (max_tokens, temperature, top_p, context_length, model_name)
		 SELECT $1,$2,$3,$4,$5 WHERE NOT EXISTS (SELECT 1 FROM llm_settings)`,
		settings.MaxTokens, settings.Temperature, settings.TopP, settings.ContextLength, settings.ModelName,
	)
	return err
}

func (r *AdminRepository) SaveLog(ctx context.Context, entry *adminmodel.LogEntry) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO query_logs (id, user_id, query, response, model, adapter, duration_ms, token_count, status, created_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`,
		entry.ID, entry.UserID, entry.Query, entry.Response, entry.Model, entry.Adapter, entry.DurationMs, entry.TokenCount, entry.Status, entry.CreatedAt,
	)
	return err
}

func (r *AdminRepository) ListLogs(ctx context.Context, offset, limit int) ([]*adminmodel.LogEntry, int, error) {
	var total int
	err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM query_logs`).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.pool.Query(ctx,
		`SELECT id, user_id, query, response, model, adapter, duration_ms, token_count, status, created_at
		 FROM query_logs ORDER BY created_at DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	logs := make([]*adminmodel.LogEntry, 0)
	for rows.Next() {
		e := &adminmodel.LogEntry{}
		if err := rows.Scan(&e.ID, &e.UserID, &e.Query, &e.Response, &e.Model, &e.Adapter, &e.DurationMs, &e.TokenCount, &e.Status, &e.CreatedAt); err != nil {
			return nil, 0, err
		}
		logs = append(logs, e)
	}
	return logs, total, rows.Err()
}

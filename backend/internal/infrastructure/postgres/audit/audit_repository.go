package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ruveydagundogan/llm-decision-score/backend/internal/domain/audit/model"
)

type AuditRepository struct {
	pool *pgxpool.Pool
}

func NewAuditRepository(pool *pgxpool.Pool) *AuditRepository {
	return &AuditRepository{pool: pool}
}

func (r *AuditRepository) Save(ctx context.Context, log *model.AuditLog) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO audit_logs (id, user_id, action, resource, resource_id, details, ip_address, user_agent, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		log.ID, log.UserID, log.Action, log.Resource, log.ResourceID,
		log.Details, log.IPAddress, log.UserAgent, log.CreatedAt,
	)
	return err
}

func (r *AuditRepository) FindByUserID(ctx context.Context, userID string, offset, limit int) ([]*model.AuditLog, int, error) {
	var total int
	err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM audit_logs WHERE user_id = $1`, userID,
	).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.pool.Query(ctx,
		`SELECT id, user_id, action, resource, resource_id, details, ip_address, user_agent, created_at
		 FROM audit_logs WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
		userID, limit, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var logs []*model.AuditLog
	for rows.Next() {
		log := &model.AuditLog{}
		if err := rows.Scan(&log.ID, &log.UserID, &log.Action, &log.Resource, &log.ResourceID,
			&log.Details, &log.IPAddress, &log.UserAgent, &log.CreatedAt); err != nil {
			return nil, 0, err
		}
		logs = append(logs, log)
	}

	return logs, total, nil
}

func (r *AuditRepository) FindByResource(ctx context.Context, resource, resourceID string) ([]*model.AuditLog, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, user_id, action, resource, resource_id, details, ip_address, user_agent, created_at
		 FROM audit_logs WHERE resource = $1 AND resource_id = $2 ORDER BY created_at DESC`,
		resource, resourceID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []*model.AuditLog
	for rows.Next() {
		log := &model.AuditLog{}
		if err := rows.Scan(&log.ID, &log.UserID, &log.Action, &log.Resource, &log.ResourceID,
			&log.Details, &log.IPAddress, &log.UserAgent, &log.CreatedAt); err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}

	return logs, nil
}

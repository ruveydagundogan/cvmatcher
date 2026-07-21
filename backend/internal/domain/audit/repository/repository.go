package repository

import (
	"context"

	"github.com/ruveydagundogan/llm-decision-score/backend/internal/domain/audit/model"
)

type AuditRepository interface {
	Save(ctx context.Context, log *model.AuditLog) error
	FindByUserID(ctx context.Context, userID string, offset, limit int) ([]*model.AuditLog, int, error)
	FindByResource(ctx context.Context, resource, resourceID string) ([]*model.AuditLog, error)
}

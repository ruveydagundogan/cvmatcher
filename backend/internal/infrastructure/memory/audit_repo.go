package memory

import (
	"context"
	"sync"

	"github.com/ruveydagundogan/cvmatcher/backend/internal/domain/audit/model"
)

const maxInMemoryAuditLogs = 50000

type InMemoryAuditRepo struct {
	mu   sync.RWMutex
	logs []*model.AuditLog
}

func NewInMemoryAuditRepo() *InMemoryAuditRepo {
	return &InMemoryAuditRepo{}
}

func (r *InMemoryAuditRepo) Save(_ context.Context, log *model.AuditLog) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if len(r.logs) >= maxInMemoryAuditLogs {
		cutoff := maxInMemoryAuditLogs / 10
		r.logs = r.logs[cutoff:]
	}
	r.logs = append(r.logs, log)
	return nil
}

func (r *InMemoryAuditRepo) FindByUserID(_ context.Context, userID string, offset, limit int) ([]*model.AuditLog, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	filtered := make([]*model.AuditLog, 0, 32)
	for _, l := range r.logs {
		if l.UserID == userID {
			filtered = append(filtered, l)
		}
	}

	total := len(filtered)
	if offset >= total {
		return nil, total, nil
	}
	end := offset + limit
	if end > total {
		end = total
	}
	return filtered[offset:end], total, nil
}

func (r *InMemoryAuditRepo) FindByResource(_ context.Context, resource, resourceID string) ([]*model.AuditLog, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	filtered := make([]*model.AuditLog, 0, 16)
	for _, l := range r.logs {
		if l.Resource == resource && l.ResourceID == resourceID {
			filtered = append(filtered, l)
		}
	}
	return filtered, nil
}

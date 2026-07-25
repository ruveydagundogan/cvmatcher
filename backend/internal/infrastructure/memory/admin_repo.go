package memory

import (
	"context"
	"errors"
	"sync"

	"github.com/ruveydagundogan/cvmatcher/backend/internal/domain/admin/model"
)

var errNotFound = errors.New("not found")

type AdminRepository struct {
	mu       sync.RWMutex
	adapters map[string]*model.Adapter
	prompts  map[string]*model.SystemPrompt
	settings *model.LLMSettings
	logs     []*model.LogEntry
}

func NewAdminRepository() *AdminRepository {
	return &AdminRepository{
		adapters: make(map[string]*model.Adapter),
		prompts:  make(map[string]*model.SystemPrompt),
		settings: model.DefaultSettings(),
	}
}

func (r *AdminRepository) ListAdapters(ctx context.Context) ([]*model.Adapter, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*model.Adapter, 0)
	for _, a := range r.adapters {
		result = append(result, a)
	}
	return result, nil
}

func (r *AdminRepository) SaveAdapter(ctx context.Context, adapter *model.Adapter) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.adapters[adapter.ID] = adapter
	return nil
}

func (r *AdminRepository) UpdateAdapter(ctx context.Context, adapter *model.Adapter) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.adapters[adapter.ID] = adapter
	return nil
}

func (r *AdminRepository) DeleteAdapter(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.adapters, id)
	return nil
}

func (r *AdminRepository) ListPrompts(ctx context.Context) ([]*model.SystemPrompt, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*model.SystemPrompt, 0)
	for _, p := range r.prompts {
		result = append(result, p)
	}
	return result, nil
}

func (r *AdminRepository) SavePrompt(ctx context.Context, prompt *model.SystemPrompt) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.prompts[prompt.ID] = prompt
	return nil
}

func (r *AdminRepository) UpdatePrompt(ctx context.Context, prompt *model.SystemPrompt) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.prompts[prompt.ID] = prompt
	return nil
}

func (r *AdminRepository) DeletePrompt(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.prompts, id)
	return nil
}

func (r *AdminRepository) GetActivePrompt(ctx context.Context) (*model.SystemPrompt, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, p := range r.prompts {
		if p.Active {
			return p, nil
		}
	}
	return nil, errNotFound
}

func (r *AdminRepository) GetSettings(ctx context.Context) (*model.LLMSettings, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if r.settings == nil {
		return nil, errNotFound
	}
	return r.settings, nil
}

func (r *AdminRepository) SaveSettings(ctx context.Context, settings *model.LLMSettings) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.settings = settings
	return nil
}

func (r *AdminRepository) SaveLog(ctx context.Context, entry *model.LogEntry) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.logs = append(r.logs, entry)
	return nil
}

func (r *AdminRepository) ListLogs(ctx context.Context, offset, limit int) ([]*model.LogEntry, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	total := len(r.logs)
	if offset >= total {
		return []*model.LogEntry{}, total, nil
	}
	end := offset + limit
	if end > total {
		end = total
	}
	return r.logs[offset:end], total, nil
}

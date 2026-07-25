package repository

import (
	"context"

	"github.com/ruveydagundogan/cvmatcher/backend/internal/domain/admin/model"
)

type AdminRepository interface {
	// Adapters
	ListAdapters(ctx context.Context) ([]*model.Adapter, error)
	SaveAdapter(ctx context.Context, adapter *model.Adapter) error
	UpdateAdapter(ctx context.Context, adapter *model.Adapter) error
	DeleteAdapter(ctx context.Context, id string) error

	// System Prompts
	ListPrompts(ctx context.Context) ([]*model.SystemPrompt, error)
	SavePrompt(ctx context.Context, prompt *model.SystemPrompt) error
	UpdatePrompt(ctx context.Context, prompt *model.SystemPrompt) error
	DeletePrompt(ctx context.Context, id string) error
	GetActivePrompt(ctx context.Context) (*model.SystemPrompt, error)

	// Settings
	GetSettings(ctx context.Context) (*model.LLMSettings, error)
	SaveSettings(ctx context.Context, settings *model.LLMSettings) error

	// Logs
	SaveLog(ctx context.Context, log *model.LogEntry) error
	ListLogs(ctx context.Context, offset, limit int) ([]*model.LogEntry, int, error)
}

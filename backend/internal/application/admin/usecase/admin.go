package usecase

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/ruveydagundogan/cvmatcher/backend/internal/domain/admin/model"
	"github.com/ruveydagundogan/cvmatcher/backend/internal/domain/admin/repository"
	apperrors "github.com/ruveydagundogan/cvmatcher/backend/internal/shared/errors"
)

type AdminUseCase struct {
	repo repository.AdminRepository
	log  *slog.Logger
}

func NewAdminUseCase(repo repository.AdminRepository, log *slog.Logger) *AdminUseCase {
	return &AdminUseCase{repo: repo, log: log}
}

func (uc *AdminUseCase) ListAdapters(ctx context.Context) ([]*model.Adapter, error) {
	return uc.repo.ListAdapters(ctx)
}

func (uc *AdminUseCase) CreateAdapter(ctx context.Context, name, description, filePath, modelName string) (*model.Adapter, error) {
	adapter := model.NewAdapter(name, description, filePath, modelName)
	adapter.ID = uuid.New().String()
	if err := uc.repo.SaveAdapter(ctx, adapter); err != nil {
		return nil, apperrors.Internal("failed to save adapter", err)
	}
	return adapter, nil
}

func (uc *AdminUseCase) UpdateAdapter(ctx context.Context, adapter *model.Adapter) error {
	adapter.UpdatedAt = time.Now().UTC()
	return uc.repo.UpdateAdapter(ctx, adapter)
}

func (uc *AdminUseCase) DeleteAdapter(ctx context.Context, id string) error {
	return uc.repo.DeleteAdapter(ctx, id)
}

func (uc *AdminUseCase) ListPrompts(ctx context.Context) ([]*model.SystemPrompt, error) {
	return uc.repo.ListPrompts(ctx)
}

func (uc *AdminUseCase) CreatePrompt(ctx context.Context, name, content string) (*model.SystemPrompt, error) {
	prompt := model.NewSystemPrompt(name, content)
	prompt.ID = uuid.New().String()
	if err := uc.repo.SavePrompt(ctx, prompt); err != nil {
		return nil, apperrors.Internal("failed to save prompt", err)
	}
	return prompt, nil
}

func (uc *AdminUseCase) UpdatePrompt(ctx context.Context, prompt *model.SystemPrompt) error {
	prompt.UpdatedAt = time.Now().UTC()
	return uc.repo.UpdatePrompt(ctx, prompt)
}

func (uc *AdminUseCase) DeletePrompt(ctx context.Context, id string) error {
	return uc.repo.DeletePrompt(ctx, id)
}

func (uc *AdminUseCase) ActivatePrompt(ctx context.Context, id string) (*model.SystemPrompt, error) {
	prompts, err := uc.repo.ListPrompts(ctx)
	if err != nil {
		return nil, err
	}
	var target *model.SystemPrompt
	for _, p := range prompts {
		if p.ID == id {
			target = p
		}
		p.Active = false
		uc.repo.UpdatePrompt(ctx, p)
	}
	if target == nil {
		return nil, apperrors.NotFound("prompt not found")
	}
	target.Active = true
	target.UpdatedAt = time.Now().UTC()
	if err := uc.repo.UpdatePrompt(ctx, target); err != nil {
		return nil, err
	}
	return target, nil
}

func (uc *AdminUseCase) GetSettings(ctx context.Context) (*model.LLMSettings, error) {
	settings, err := uc.repo.GetSettings(ctx)
	if err != nil || settings == nil {
		return model.DefaultSettings(), nil
	}
	return settings, nil
}

func (uc *AdminUseCase) SaveSettings(ctx context.Context, settings *model.LLMSettings) error {
	settings.UpdatedAt = time.Now().UTC()
	return uc.repo.SaveSettings(ctx, settings)
}

func (uc *AdminUseCase) ListLogs(ctx context.Context, offset, limit int) ([]*model.LogEntry, int, error) {
	return uc.repo.ListLogs(ctx, offset, limit)
}

func (uc *AdminUseCase) LogQuery(ctx context.Context, userID, query, responseText, modelName, adapter string, durationMs int64, tokenCount int, status string) error {
	entry := &model.LogEntry{
		ID:         uuid.New().String(),
		UserID:     userID,
		Query:      query,
		Response:   responseText,
		Model:      modelName,
		Adapter:    adapter,
		DurationMs: durationMs,
		TokenCount: tokenCount,
		Status:     status,
		CreatedAt:  time.Now().UTC(),
	}
	if saveErr := uc.repo.SaveLog(ctx, entry); saveErr != nil {
		uc.log.Error("failed to save log", "error", saveErr)
	}
	return nil
}

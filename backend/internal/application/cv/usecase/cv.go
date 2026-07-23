package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"time"

	cvmodel "github.com/ruveydagundogan/llm-decision-score/backend/internal/domain/cv/model"
	cvrepo "github.com/ruveydagundogan/llm-decision-score/backend/internal/domain/cv/repository"
	"github.com/ruveydagundogan/llm-decision-score/backend/internal/infrastructure/llm"
)

type CVUseCase struct {
	repo      cvrepo.CVRepository
	llmClient *llm.Client
	log       *slog.Logger
}

func NewCVUseCase(repo cvrepo.CVRepository, llmClient *llm.Client, log *slog.Logger) *CVUseCase {
	return &CVUseCase{repo: repo, llmClient: llmClient, log: log}
}

func (uc *CVUseCase) Create(ctx context.Context, userID, title, content string) (*cvmodel.CV, error) {
	cv := cvmodel.NewCV(userID, title, content)
	cv.Status = "pending"
	cv.CreatedAt = time.Now().UTC()
	cv.UpdatedAt = cv.CreatedAt
	if err := uc.repo.Save(ctx, cv); err != nil {
		return nil, fmt.Errorf("save cv: %w", err)
	}
	return cv, nil
}

func (uc *CVUseCase) GetByID(ctx context.Context, id string) (*cvmodel.CV, error) {
	return uc.repo.FindByID(ctx, id)
}

func (uc *CVUseCase) ListByUser(ctx context.Context, userID string, offset, limit int) ([]*cvmodel.CV, int, error) {
	return uc.repo.FindByUserID(ctx, userID, offset, limit)
}

func (uc *CVUseCase) Delete(ctx context.Context, id string) error {
	return uc.repo.Delete(ctx, id)
}

func (uc *CVUseCase) ParseWithLLM(ctx context.Context, cvID string) (*cvmodel.CV, error) {
	cv, err := uc.repo.FindByID(ctx, cvID)
	if err != nil {
		return nil, fmt.Errorf("find cv: %w", err)
	}
	if cv.Content == "" {
		return nil, fmt.Errorf("cv content is empty")
	}

	prompt := fmt.Sprintf(llm.CVParsePrompt, cv.Content)
	response, err := uc.llmClient.ChatCompletion(ctx, []llm.ChatMessage{
		{Role: "user", Content: prompt},
	}, 1024)
	if err != nil {
		cv.Status = "failed"
		uc.repo.Update(ctx, cv)
		return nil, fmt.Errorf("llm parse failed: %w", err)
	}

	response = cleanJSON(response)
	var parsed struct {
		Skills     []string             `json:"skills"`
		Experience []cvmodel.ParsedExperience `json:"experience"`
		Education  []cvmodel.ParsedEducation  `json:"education"`
		Summary    string               `json:"summary"`
	}
	if err := json.Unmarshal([]byte(response), &parsed); err != nil {
		cv.Status = "failed"
		uc.repo.Update(ctx, cv)
		return nil, fmt.Errorf("parse llm response: %w", err)
	}

	cv.ParsedSkills = parsed.Skills
	cv.ParsedExperience = parsed.Experience
	cv.ParsedEducation = parsed.Education
	cv.ParsedSummary = parsed.Summary
	cv.Status = "completed"
	cv.UpdatedAt = time.Now().UTC()
	if err := uc.repo.Update(ctx, cv); err != nil {
		return nil, fmt.Errorf("update cv: %w", err)
	}
	return cv, nil
}

func cleanJSON(s string) string {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "```json")
	s = strings.TrimPrefix(s, "```")
	s = strings.TrimSuffix(s, "```")
	s = strings.TrimSpace(s)
	return s
}

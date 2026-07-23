package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"time"

	jdmodel "github.com/ruveydagundogan/llm-decision-score/backend/internal/domain/jobdescription/model"
	jdrepo "github.com/ruveydagundogan/llm-decision-score/backend/internal/domain/jobdescription/repository"
	"github.com/ruveydagundogan/llm-decision-score/backend/internal/infrastructure/llm"
)

type JDUseCase struct {
	repo      jdrepo.JobDescriptionRepository
	llmClient *llm.Client
	log       *slog.Logger
}

func NewJDUseCase(repo jdrepo.JobDescriptionRepository, llmClient *llm.Client, log *slog.Logger) *JDUseCase {
	return &JDUseCase{repo: repo, llmClient: llmClient, log: log}
}

func (uc *JDUseCase) Create(ctx context.Context, userID, title, content string) (*jdmodel.JobDescription, error) {
	jd := jdmodel.NewJobDescription(userID, title, content)
	jd.CreatedAt = time.Now().UTC()
	jd.UpdatedAt = jd.CreatedAt
	if err := uc.repo.Save(ctx, jd); err != nil {
		return nil, fmt.Errorf("save jd: %w", err)
	}
	return jd, nil
}

func (uc *JDUseCase) GetByID(ctx context.Context, id string) (*jdmodel.JobDescription, error) {
	return uc.repo.FindByID(ctx, id)
}

func (uc *JDUseCase) ListByUser(ctx context.Context, userID string, offset, limit int) ([]*jdmodel.JobDescription, int, error) {
	return uc.repo.FindByUserID(ctx, userID, offset, limit)
}

func (uc *JDUseCase) Update(ctx context.Context, jd *jdmodel.JobDescription) error {
	jd.UpdatedAt = time.Now().UTC()
	return uc.repo.Update(ctx, jd)
}

func (uc *JDUseCase) Delete(ctx context.Context, id string) error {
	return uc.repo.Delete(ctx, id)
}

func (uc *JDUseCase) AnalyzeWithLLM(ctx context.Context, jdID string) (*jdmodel.JobDescription, error) {
	jd, err := uc.repo.FindByID(ctx, jdID)
	if err != nil {
		return nil, fmt.Errorf("find jd: %w", err)
	}
	if jd.Content == "" {
		return nil, fmt.Errorf("jd content is empty")
	}

	prompt := fmt.Sprintf(llm.JDAnalyzePrompt, jd.Content)
	response, err := uc.llmClient.ChatCompletion(ctx, []llm.ChatMessage{
		{Role: "user", Content: prompt},
	}, 1024)
	if err != nil {
		return nil, fmt.Errorf("llm analyze failed: %w", err)
	}

	response = cleanJSON(response)
	var parsed struct {
		RequiredSkills  []string `json:"required_skills"`
		PreferredSkills []string `json:"preferred_skills"`
		ExperienceLevel string   `json:"experience_level"`
		EmploymentType  string   `json:"employment_type"`
	}
	if err := json.Unmarshal([]byte(response), &parsed); err != nil {
		return nil, fmt.Errorf("parse llm response: %w", err)
	}

	jd.RequiredSkills = parsed.RequiredSkills
	jd.PreferredSkills = parsed.PreferredSkills
	jd.ExperienceLevel = parsed.ExperienceLevel
	jd.EmploymentType = parsed.EmploymentType
	jd.UpdatedAt = time.Now().UTC()
	if err := uc.repo.Update(ctx, jd); err != nil {
		return nil, fmt.Errorf("update jd: %w", err)
	}
	return jd, nil
}

func cleanJSON(s string) string {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "```json")
	s = strings.TrimPrefix(s, "```")
	s = strings.TrimSuffix(s, "```")
	s = strings.TrimSpace(s)
	return s
}

package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"regexp"
	"strings"
	"time"

	cvmodel "github.com/ruveydagundogan/cvmatcher/backend/internal/domain/cv/model"
	cvrepo "github.com/ruveydagundogan/cvmatcher/backend/internal/domain/cv/repository"
	"github.com/ruveydagundogan/cvmatcher/backend/internal/infrastructure/llm"
	apperrors "github.com/ruveydagundogan/cvmatcher/backend/internal/shared/errors"
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
		return nil, apperrors.NotFound("CV not found")
	}
	if cv.Content == "" {
		return nil, apperrors.Validation("CV content is empty")
	}

	prompt := fmt.Sprintf(llm.CVParsePrompt, cv.Content)
	response, err := uc.llmClient.ChatCompletion(ctx, []llm.ChatMessage{
		{Role: "user", Content: prompt},
	}, 1024)
	if err != nil {
		uc.log.Warn("llm parse failed, using fallback", "error", err)
		return uc.fallbackParse(ctx, cv)
	}

	response = cleanJSON(response)
	var parsed struct {
		Skills     []string                 `json:"skills"`
		Experience []cvmodel.ParsedExperience `json:"experience"`
		Education  []cvmodel.ParsedEducation  `json:"education"`
		Summary    string                   `json:"summary"`
	}
	if err := json.Unmarshal([]byte(response), &parsed); err != nil {
		uc.log.Warn("llm response parse failed, using fallback", "error", err)
		return uc.fallbackParse(ctx, cv)
	}

	cv.ParsedSkills = parsed.Skills
	cv.ParsedExperience = parsed.Experience
	cv.ParsedEducation = parsed.Education
	cv.ParsedSummary = parsed.Summary
	cv.Status = "completed"
	cv.UpdatedAt = time.Now().UTC()
	if err := uc.repo.Update(ctx, cv); err != nil {
		return nil, apperrors.Internal("failed to update CV", err)
	}
	return cv, nil
}

func (uc *CVUseCase) fallbackParse(ctx context.Context, cv *cvmodel.CV) (*cvmodel.CV, error) {
	content := cv.Content

	email := extractEmail(content)
	lines := strings.Split(content, "\n")

	var skills []string
	knownSkills := []string{
		"go", "golang", "python", "java", "javascript", "typescript", "rust", "c++", "c#",
		"react", "vue", "angular", "next.js", "node.js", "express",
		"postgresql", "mysql", "mongodb", "redis", "sql", "database",
		"docker", "kubernetes", "aws", "gcp", "azure", "linux", "git",
		"rest api", "graphql", "microservices", "ci/cd", "agile",
		"machine learning", "ai", "data science", "deep learning",
	}

	lower := strings.ToLower(content)
	for _, skill := range knownSkills {
		if strings.Contains(lower, skill) {
			skills = append(skills, skill)
		}
	}

	var summary string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.Contains(strings.ToUpper(trimmed), "ÖZET") ||
		   strings.Contains(strings.ToUpper(trimmed), "SUMMARY") ||
		   strings.Contains(strings.ToUpper(trimmed), "PROFILE") ||
		   strings.Contains(strings.ToUpper(trimmed), "ABOUT") {
			for j := 1; j < len(lines) && j < 5; j++ {
				if j+1 < len(lines) {
					next := strings.TrimSpace(lines[j+1])
					if next != "" && !strings.Contains(strings.ToUpper(next), "DENEY") && !strings.Contains(strings.ToUpper(next), "EXPERIENCE") {
						summary += next + " "
					}
				}
			}
			break
		}
	}

	if summary == "" && len(lines) > 0 {
		summary = strings.TrimSpace(lines[0])
	}

	cv.ParsedSkills = uniqueStrings(skills)
	cv.ParsedSummary = strings.TrimSpace(summary)
	if email != "" {
		cv.ParsedSummary = "Email: " + email + ". " + cv.ParsedSummary
	}

	uc.log.Info("fallback parse completed", "cv_id", cv.ID, "skills_count", len(skills))
	cv.Status = "completed"
	cv.UpdatedAt = time.Now().UTC()
	if err := uc.repo.Update(ctx, cv); err != nil {
		return nil, apperrors.Internal("failed to update CV", err)
	}
	return cv, nil
}

func extractEmail(s string) string {
	re := regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)
	matches := re.FindString(s)
	return matches
}

func extractPhone(s string) string {
	re := regexp.MustCompile(`(\+?\d{1,3}[\s-]?)?\(?\d{3}\)?[\s-]?\d{3}[\s-]?\d{2,4}[\s-]?\d{2,4}`)
	matches := re.FindString(s)
	return matches
}

func uniqueStrings(s []string) []string {
	seen := make(map[string]bool)
	var result []string
	for _, v := range s {
		v = strings.TrimSpace(v)
		if v != "" && !seen[v] {
			seen[v] = true
			result = append(result, v)
		}
	}
	return result
}

func cleanJSON(s string) string {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "```json")
	s = strings.TrimPrefix(s, "```")
	s = strings.TrimSuffix(s, "```")
	s = strings.TrimSpace(s)
	return s
}

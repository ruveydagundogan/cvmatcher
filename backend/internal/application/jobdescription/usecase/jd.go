package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"time"

	jdmodel "github.com/ruveydagundogan/cvmatcher/backend/internal/domain/jobdescription/model"
	jdrepo "github.com/ruveydagundogan/cvmatcher/backend/internal/domain/jobdescription/repository"
	"github.com/ruveydagundogan/cvmatcher/backend/internal/infrastructure/llm"
	apperrors "github.com/ruveydagundogan/cvmatcher/backend/internal/shared/errors"
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
		return nil, apperrors.Internal("failed to save job description", err)
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
		return nil, apperrors.NotFound("job description not found")
	}
	if jd.Content == "" {
		return nil, apperrors.Validation("job description content is empty")
	}

	prompt := fmt.Sprintf(llm.JDAnalyzePrompt, jd.Content)
	response, err := uc.llmClient.ChatCompletionWithModel(ctx, llm.ModelCVParse, []llm.ChatMessage{
		{Role: "user", Content: prompt},
	}, 1024)
	if err != nil {
		uc.log.Warn("llm analyze failed, using fallback", "error", err)
		return uc.fallbackAnalyze(ctx, jd)
	}

	response = cleanJSON(response)
	var parsed struct {
		RequiredSkills  []string `json:"required_skills"`
		PreferredSkills []string `json:"preferred_skills"`
		ExperienceLevel string   `json:"experience_level"`
		EmploymentType  string   `json:"employment_type"`
	}
	if err := json.Unmarshal([]byte(response), &parsed); err != nil {
		uc.log.Warn("llm analyze response parse failed, using fallback", "error", err)
		return uc.fallbackAnalyze(ctx, jd)
	}

	jd.RequiredSkills = parsed.RequiredSkills
	jd.PreferredSkills = parsed.PreferredSkills
	jd.ExperienceLevel = parsed.ExperienceLevel
	jd.EmploymentType = parsed.EmploymentType
	jd.UpdatedAt = time.Now().UTC()
	if err := uc.repo.Update(ctx, jd); err != nil {
		return nil, apperrors.Internal("failed to update job description", err)
	}
	return jd, nil
}

func (uc *JDUseCase) fallbackAnalyze(ctx context.Context, jd *jdmodel.JobDescription) (*jdmodel.JobDescription, error) {
	content := jd.Content
	lower := strings.ToLower(content)

	var requiredSkills []string
	keywordSkills := []string{
		"go", "golang", "python", "java", "javascript", "typescript", "rust", "c++",
		"react", "vue", "angular", "next.js", "node.js",
		"postgresql", "mysql", "mongodb", "redis", "sql",
		"docker", "kubernetes", "aws", "gcp", "azure", "linux", "git",
		"rest api", "graphql", "microservices",
		"machine learning", "ai", "data science",
	}
	for _, skill := range keywordSkills {
		if strings.Contains(lower, skill) {
			requiredSkills = append(requiredSkills, skill)
		}
	}

	expLevel := "mid"
	if strings.Contains(lower, "senior") || strings.Contains(lower, "kıdemli") {
		expLevel = "senior"
	} else if strings.Contains(lower, "junior") || strings.Contains(lower, "entry") || strings.Contains(lower, "yeni") {
		expLevel = "junior"
	}

	empType := "full-time"
	if strings.Contains(lower, "part-time") || strings.Contains(lower, "yarı zamanlı") {
		empType = "part-time"
	} else if strings.Contains(lower, "contract") || strings.Contains(lower, "sözleşmeli") {
		empType = "contract"
	} else if strings.Contains(lower, "freelance") {
		empType = "freelance"
	}

	jd.RequiredSkills = uniqueStrings(requiredSkills)
	jd.ExperienceLevel = expLevel
	jd.EmploymentType = empType
	jd.UpdatedAt = time.Now().UTC()

	uc.log.Info("fallback analyze completed", "jd_id", jd.ID, "skills_count", len(requiredSkills))
	if err := uc.repo.Update(ctx, jd); err != nil {
		return nil, apperrors.Internal("failed to update job description", err)
	}
	return jd, nil
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

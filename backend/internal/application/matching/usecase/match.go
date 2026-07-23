package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"time"

	cvmodel "github.com/ruveydagundogan/cvmatcher/backend/internal/domain/cv/model"
	cvrepo "github.com/ruveydagundogan/cvmatcher/backend/internal/domain/cv/repository"
	jdmodel "github.com/ruveydagundogan/cvmatcher/backend/internal/domain/jobdescription/model"
	jdrepo "github.com/ruveydagundogan/cvmatcher/backend/internal/domain/jobdescription/repository"
	matchmodel "github.com/ruveydagundogan/cvmatcher/backend/internal/domain/matching/model"
	matchrepo "github.com/ruveydagundogan/cvmatcher/backend/internal/domain/matching/repository"
	"github.com/ruveydagundogan/cvmatcher/backend/internal/infrastructure/llm"
)

type MatchUseCase struct {
	cvRepo      cvrepo.CVRepository
	jdRepo      jdrepo.JobDescriptionRepository
	matchRepo   matchrepo.MatchingRepository
	llmClient   *llm.Client
	log         *slog.Logger
}

func NewMatchUseCase(
	cvRepo cvrepo.CVRepository,
	jdRepo jdrepo.JobDescriptionRepository,
	matchRepo matchrepo.MatchingRepository,
	llmClient *llm.Client,
	log *slog.Logger,
) *MatchUseCase {
	return &MatchUseCase{
		cvRepo:    cvRepo,
		jdRepo:    jdRepo,
		matchRepo: matchRepo,
		llmClient: llmClient,
		log:       log,
	}
}

func (uc *MatchUseCase) RunMatch(ctx context.Context, userID, cvID, jdID string) (*matchmodel.MatchResult, error) {
	cv, err := uc.cvRepo.FindByID(ctx, cvID)
	if err != nil {
		return nil, fmt.Errorf("cv not found: %w", err)
	}
	if cv.UserID != userID {
		return nil, fmt.Errorf("cv does not belong to user")
	}

	jd, err := uc.jdRepo.FindByID(ctx, jdID)
	if err != nil {
		return nil, fmt.Errorf("jd not found: %w", err)
	}
	if jd.UserID != userID {
		return nil, fmt.Errorf("jd does not belong to user")
	}

	existing, err := uc.matchRepo.FindByCVAndJD(ctx, cvID, jdID)
	if err == nil && existing != nil {
		return existing, nil
	}

	if cv.Status != "completed" {
		if err := uc.parseCV(ctx, cv); err != nil {
			return nil, err
		}
	}
	if len(jd.RequiredSkills) == 0 {
		if err := uc.analyzeJD(ctx, jd); err != nil {
			return nil, err
		}
	}

	return uc.runLLMMatch(ctx, userID, cv, jd)
}

func (uc *MatchUseCase) parseCV(ctx context.Context, cv *cvmodel.CV) error {
	prompt := fmt.Sprintf(llm.CVParsePrompt, cv.Content)
	response, err := uc.llmClient.ChatCompletion(ctx, []llm.ChatMessage{
		{Role: "user", Content: prompt},
	}, 1024)
	if err != nil {
		return fmt.Errorf("llm parse cv: %w", err)
	}

	cleaned := cleanJSON(response)
	var parsed struct {
		Skills     []string                 `json:"skills"`
		Experience []cvmodel.ParsedExperience `json:"experience"`
		Education  []cvmodel.ParsedEducation  `json:"education"`
		Summary    string                   `json:"summary"`
	}
	if err := json.Unmarshal([]byte(cleaned), &parsed); err != nil {
		return fmt.Errorf("parse cv response: %w", err)
	}

	cv.ParsedSkills = parsed.Skills
	cv.ParsedExperience = parsed.Experience
	cv.ParsedEducation = parsed.Education
	cv.ParsedSummary = parsed.Summary
	cv.Status = "completed"
	return uc.cvRepo.Update(ctx, cv)
}

func (uc *MatchUseCase) analyzeJD(ctx context.Context, jd *jdmodel.JobDescription) error {
	prompt := fmt.Sprintf(llm.JDAnalyzePrompt, jd.Content)
	response, err := uc.llmClient.ChatCompletion(ctx, []llm.ChatMessage{
		{Role: "user", Content: prompt},
	}, 1024)
	if err != nil {
		return fmt.Errorf("llm analyze jd: %w", err)
	}

	cleaned := cleanJSON(response)
	var parsed struct {
		RequiredSkills  []string `json:"required_skills"`
		PreferredSkills []string `json:"preferred_skills"`
		ExperienceLevel string   `json:"experience_level"`
		EmploymentType  string   `json:"employment_type"`
	}
	if err := json.Unmarshal([]byte(cleaned), &parsed); err != nil {
		return fmt.Errorf("parse jd response: %w", err)
	}

	jd.RequiredSkills = parsed.RequiredSkills
	jd.PreferredSkills = parsed.PreferredSkills
	jd.ExperienceLevel = parsed.ExperienceLevel
	jd.EmploymentType = parsed.EmploymentType
	return uc.jdRepo.Update(ctx, jd)
}

func (uc *MatchUseCase) runLLMMatch(ctx context.Context, userID string, cv *cvmodel.CV, jd *jdmodel.JobDescription) (*matchmodel.MatchResult, error) {
	expJSON, _ := json.Marshal(cv.ParsedExperience)
	eduJSON, _ := json.Marshal(cv.ParsedEducation)

	skillsStr := strings.Join(cv.ParsedSkills, ", ")
	reqSkillsStr := strings.Join(jd.RequiredSkills, ", ")
	prefSkillsStr := strings.Join(jd.PreferredSkills, ", ")

	prompt := fmt.Sprintf(llm.CVJDMatchPrompt,
		skillsStr, string(expJSON), string(eduJSON), cv.ParsedSummary,
		reqSkillsStr, prefSkillsStr, jd.ExperienceLevel, jd.Content,
	)

	response, err := uc.llmClient.ChatCompletion(ctx, []llm.ChatMessage{
		{Role: "user", Content: prompt},
	}, 1024)
	if err != nil {
		return nil, fmt.Errorf("llm match: %w", err)
	}

	cleaned := cleanJSON(response)
	var match matchmodel.MatchResult
	if err := json.Unmarshal([]byte(cleaned), &match); err != nil {
		return nil, fmt.Errorf("parse match response: %w", err)
	}

	result := matchmodel.NewMatchResult(userID, cv.ID, jd.ID)
	result.OverallScore = match.OverallScore
	result.SkillMatchScore = match.SkillMatchScore
	result.ExperienceScore = match.ExperienceScore
	result.EducationScore = match.EducationScore
	result.LLMAnalysis = match.LLMAnalysis
	result.MatchedSkills = match.MatchedSkills
	result.MissingSkills = match.MissingSkills
	result.CreatedAt = time.Now().UTC()
	result.CVTitle = cv.Title
	result.JDTitle = jd.Title

	if err := uc.matchRepo.Save(ctx, result); err != nil {
		return nil, fmt.Errorf("save match: %w", err)
	}
	return result, nil
}

func (uc *MatchUseCase) GetByID(ctx context.Context, id string) (*matchmodel.MatchResult, error) {
	return uc.matchRepo.FindByID(ctx, id)
}

func (uc *MatchUseCase) ListByUser(ctx context.Context, userID string, offset, limit int) ([]*matchmodel.MatchResult, int, error) {
	return uc.matchRepo.FindByUserID(ctx, userID, offset, limit)
}

func (uc *MatchUseCase) GetDashboardStats(ctx context.Context, userID string) (*matchmodel.DashboardStats, error) {
	return uc.matchRepo.GetDashboardStats(ctx, userID)
}

func cleanJSON(s string) string {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "```json")
	s = strings.TrimPrefix(s, "```")
	s = strings.TrimSuffix(s, "```")
	s = strings.TrimSpace(s)
	return s
}

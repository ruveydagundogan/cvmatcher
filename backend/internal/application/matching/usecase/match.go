package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"math"
	"regexp"
	"strconv"
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
		normalizeResult(existing)
		return existing, nil
	}

	if cv.Status != "completed" {
		if err := uc.parseCV(ctx, cv); err != nil {
			uc.log.Warn("cv parse failed, proceeding with empty skills", "error", err)
		}
	}
	if len(jd.RequiredSkills) == 0 {
		if err := uc.analyzeJD(ctx, jd); err != nil {
			uc.log.Warn("jd analyze failed, proceeding with empty skills", "error", err)
		}
	}

	return uc.runMatch(ctx, userID, cv, jd)
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

func (uc *MatchUseCase) runMatch(ctx context.Context, userID string, cv *cvmodel.CV, jd *jdmodel.JobDescription) (*matchmodel.MatchResult, error) {
	if len(cv.ParsedSkills) == 0 {
		cv.ParsedSkills = extractSkillsFromContent(cv.Content)
	}
	if len(jd.RequiredSkills) == 0 {
		jd.RequiredSkills = extractSkillsFromContent(jd.Content)
	}

	response, err := uc.llmMatch(ctx, cv, jd)
	if err != nil {
		uc.log.Warn("llm match failed, using fallback scoring", "error", err)
		return uc.fallbackMatch(ctx, userID, cv, jd)
	}

	cleaned := cleanJSON(response)
	var parsed struct {
		OverallScore    float64  `json:"overall_score"`
		SkillMatchScore float64  `json:"skill_match_score"`
		ExperienceScore float64  `json:"experience_score"`
		EducationScore  float64  `json:"education_score"`
		MatchedSkills   []string `json:"matched_skills"`
		MissingSkills   []string `json:"missing_skills"`
		Analysis        string   `json:"analysis"`
	}

	if err := json.Unmarshal([]byte(cleaned), &parsed); err != nil {
		uc.log.Warn("llm match response parse failed, using fallback scoring", "error", err)
		return uc.fallbackMatch(ctx, userID, cv, jd)
	}

	result := matchmodel.NewMatchResult(userID, cv.ID, jd.ID)
	result.OverallScore = parsed.OverallScore
	result.SkillMatchScore = parsed.SkillMatchScore
	result.ExperienceScore = parsed.ExperienceScore
	result.EducationScore = parsed.EducationScore
	result.LLMAnalysis = parsed.Analysis
	result.MatchedSkills = parsed.MatchedSkills
	result.MissingSkills = parsed.MissingSkills
	result.CreatedAt = time.Now().UTC()
	result.CVTitle = cv.Title
	result.JDTitle = jd.Title

	deduplicateSkills(&result.MatchedSkills, &result.MissingSkills)

	validateScores(result)

	if err := uc.matchRepo.Save(ctx, result); err != nil {
		return nil, fmt.Errorf("save match: %w", err)
	}
	return result, nil
}

func validateScores(r *matchmodel.MatchResult) {
	r.OverallScore = norm(r.OverallScore)
	r.SkillMatchScore = norm(r.SkillMatchScore)
	r.ExperienceScore = norm(r.ExperienceScore)
	r.EducationScore = norm(r.EducationScore)

	if r.OverallScore > 1.0 {
		r.OverallScore = 1.0
	}
	if r.SkillMatchScore > 1.0 {
		r.SkillMatchScore = 1.0
	}
	if r.ExperienceScore > 1.0 {
		r.ExperienceScore = 1.0
	}
	if r.EducationScore > 1.0 {
		r.EducationScore = 1.0
	}

	if len(r.MatchedSkills) > 0 && r.OverallScore < 0.1 {
		minScore := float64(len(r.MatchedSkills)) * 0.15
		if minScore > 0.9 {
			minScore = 0.9
		}
		if r.OverallScore < minScore {
			r.OverallScore = minScore
		}
	}
}

func deduplicateSkills(matched, missing *[]string) {
	matchedSet := make(map[string]bool)
	for _, s := range *matched {
		matchedSet[strings.ToLower(strings.TrimSpace(s))] = true
	}
	var filtered []string
	for _, s := range *missing {
		if !matchedSet[strings.ToLower(strings.TrimSpace(s))] {
			filtered = append(filtered, s)
		}
	}
	*missing = filtered
}

func (uc *MatchUseCase) llmMatch(ctx context.Context, cv *cvmodel.CV, jd *jdmodel.JobDescription) (string, error) {
	expJSON, _ := json.Marshal(cv.ParsedExperience)
	eduJSON, _ := json.Marshal(cv.ParsedEducation)

	skillsStr := strings.Join(cv.ParsedSkills, ", ")
	reqSkillsStr := strings.Join(jd.RequiredSkills, ", ")
	prefSkillsStr := strings.Join(jd.PreferredSkills, ", ")

	cvContent := cv.Content
	if cvContent == "" {
		cvContent = skillsStr
	}

	prompt := fmt.Sprintf(llm.CVJDMatchPrompt,
		cvContent, skillsStr, string(expJSON), string(eduJSON), cv.ParsedSummary,
		jd.Content,
		reqSkillsStr, prefSkillsStr, jd.ExperienceLevel,
	)

	uc.log.Info("llm match request",
		"cv_id", cv.ID, "jd_id", jd.ID,
		"cv_skills_len", len(cv.ParsedSkills),
		"jd_skills_len", len(jd.RequiredSkills),
	)

	resp, err := uc.llmClient.ChatCompletion(ctx, []llm.ChatMessage{
		{Role: "user", Content: prompt},
	}, 2048)
	if err != nil {
		return "", err
	}

	uc.log.Info("llm match response", "response_preview", truncateStr(resp, 200))
	return resp, nil
}

func (uc *MatchUseCase) fallbackMatch(ctx context.Context, userID string, cv *cvmodel.CV, jd *jdmodel.JobDescription) (*matchmodel.MatchResult, error) {
	cvContent := cv.Content
	if cvContent == "" {
		cvContent = strings.Join(cv.ParsedSkills, ", ")
	}
	jdContent := jd.Content
	if jdContent == "" {
		jdContent = strings.Join(jd.RequiredSkills, ", ")
	}

	extractedCVSkills := extractSkillsFromContent(cvContent)
	extractedJDSkills := extractSkillsFromContent(jdContent)

	cvYears := extractYearsOfExperience(cvContent)
	jdYears := extractYearsOfExperience(jdContent)

	cvEdu := extractEducationLevel(cvContent)
	jdEdu := extractEducationLevel(jdContent)

	cvSkillSet := make(map[string]bool)
	for _, s := range extractedCVSkills {
		cvSkillSet[s] = true
	}

	var matchedSkills, missingSkills []string
	for _, s := range extractedJDSkills {
		if cvSkillSet[s] {
			matchedSkills = append(matchedSkills, s)
		} else {
			missingSkills = append(missingSkills, s)
		}
	}

	total := len(extractedJDSkills)
	matched := len(matchedSkills)

	var skillScore float64
	if total > 0 {
		skillScore = math.Round(float64(matched)/float64(total)*100) / 100
	}

	var expScore float64 = 0.5
	if cvYears > 0 && jdYears > 0 {
		if cvYears >= jdYears {
			expScore = math.Min(1.0, 0.5+float64(cvYears-jdYears)*0.05)
		} else {
			expScore = math.Max(0.2, 0.5-float64(jdYears-cvYears)*0.1)
		}
	}

	var eduScore float64 = 0.5
	eduLevels := map[string]int{"unknown": 0, "bachelor": 1, "master": 2, "phd": 3}
	cvEduLevel := eduLevels[cvEdu]
	jdEduLevel := eduLevels[jdEdu]
	if cvEduLevel >= jdEduLevel {
		eduScore = 1.0
	} else if jdEduLevel > 0 {
		eduScore = float64(cvEduLevel) / float64(jdEduLevel)
	}

	overallScore := skillScore*0.5 + expScore*0.25 + eduScore*0.25

	if overallScore > 1.0 {
		overallScore = 1.0
	}
	if overallScore < 0.1 && matched > 0 {
		overallScore = math.Min(0.3, float64(matched)*0.1)
	}

	result := matchmodel.NewMatchResult(userID, cv.ID, jd.ID)
	result.OverallScore = overallScore
	result.SkillMatchScore = skillScore
	result.ExperienceScore = math.Round(expScore*100) / 100
	result.EducationScore = math.Round(eduScore*100) / 100
	result.MatchedSkills = matchedSkills
	result.MissingSkills = missingSkills

	deduplicateSkills(&result.MatchedSkills, &result.MissingSkills)

	analysis := fmt.Sprintf("Fallback analysis: %d of %d required skills matched. ", matched, total)
	if overallScore >= 0.7 {
		analysis += "Good overall fit."
	} else if overallScore >= 0.4 {
		analysis += "Moderate fit - consider reviewing missing skills."
	} else {
		analysis += "Low match - significant skill gaps identified."
	}
	result.LLMAnalysis = analysis
	result.CreatedAt = time.Now().UTC()
	result.CVTitle = cv.Title
	result.JDTitle = jd.Title

	validateScores(result)

	uc.log.Info("fallback match completed",
		"cv_id", cv.ID, "jd_id", jd.ID,
		"matched", matched, "total", total,
		"cv_years", cvYears, "jd_years", jdYears,
		"cv_edu", cvEdu, "jd_edu", jdEdu,
		"score", overallScore,
	)

	if err := uc.matchRepo.Save(ctx, result); err != nil {
		return nil, fmt.Errorf("save match: %w", err)
	}
	return result, nil
}

func (uc *MatchUseCase) GetByID(ctx context.Context, id string) (*matchmodel.MatchResult, error) {
	m, err := uc.matchRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	normalizeResult(m)
	return m, nil
}

func (uc *MatchUseCase) ListByUser(ctx context.Context, userID string, offset, limit int) ([]*matchmodel.MatchResult, int, error) {
	matches, total, err := uc.matchRepo.FindByUserID(ctx, userID, offset, limit)
	if err != nil {
		return nil, 0, err
	}
	for _, m := range matches {
		normalizeResult(m)
	}
	return matches, total, nil
}

func (uc *MatchUseCase) GetDashboardStats(ctx context.Context, userID string) (*matchmodel.DashboardStats, error) {
	s, err := uc.matchRepo.GetDashboardStats(ctx, userID)
	if err != nil {
		return nil, err
	}
	s.AverageScore = norm(s.AverageScore)
	for _, m := range s.RecentMatches {
		normalizeResult(m)
	}
	return s, nil
}

func (uc *MatchUseCase) Delete(ctx context.Context, id string) error {
	return uc.matchRepo.Delete(ctx, id)
}

func normalizeResult(m *matchmodel.MatchResult) {
	m.OverallScore = norm(m.OverallScore)
	m.SkillMatchScore = norm(m.SkillMatchScore)
	m.ExperienceScore = norm(m.ExperienceScore)
	m.EducationScore = norm(m.EducationScore)
	validateScores(m)
}

func extractSkillsFromContent(content string) []string {
	knownSkills := []string{
		"go", "golang", "python", "java", "javascript", "typescript", "rust", "c++", "c#",
		"ruby", "scala", "kotlin", "swift", "php", "perl", "bash", "shell", "elixir",
		"react", "vue", "angular", "next.js", "svelte", "jquery", "html", "css", "sass",
		"node.js", "express", "django", "flask", "fastapi", "spring boot", "rails",
		"postgresql", "mysql", "mongodb", "redis", "sql", "database", "sqlite",
		"oracle", "mariadb", "cassandra", "dynamodb", "elasticsearch", "bigquery",
		"docker", "kubernetes", "k8s", "aws", "gcp", "azure", "linux", "git",
		"terraform", "ansible", "helm", "jenkins", "gitlab ci", "github actions", "ci/cd",
		"rest api", "graphql", "grpc", "microservices", "event-driven", "kafka",
		"rabbitmq", "nginx", "apache", "prometheus", "grafana", "datadog",
		"machine learning", "deep learning", "ai", "data science", "nlp", "computer vision",
		"tensorflow", "pytorch", "scikit-learn", "pandas", "numpy", "jupyter", "spark",
		"agile", "scrum", "jira", "confluence", "leadership", "mentoring",
		"system design", "architecture", "distributed systems", "cloud", "devops",
		"react native", "flutter", "android", "ios", "mobile",
		"blockchain", "solidity", "web3", "smart contracts",
		"cybersecurity", "penetration testing", "security", "oauth", "jwt",
		"testing", "unit test", "integration test", "tdd", "cypress", "jest",
		"nosql", "sql", "etl", "data pipeline", "data warehouse",
	}
	lower := strings.ToLower(content)
	var skills []string
	for _, s := range knownSkills {
		if strings.Contains(lower, s) {
			skills = append(skills, s)
		}
	}
	return uniqueStrings(skills)
}

func extractYearsOfExperience(content string) int {
	re := regexp.MustCompile(`(\d+)\+?\s*(?:years?|yrs?)`)
	matches := re.FindAllStringSubmatch(strings.ToLower(content), -1)
	maxYears := 0
	for _, m := range matches {
		if len(m) > 1 {
			if y, err := strconv.Atoi(m[1]); err == nil && y > maxYears {
				maxYears = y
			}
		}
	}
	if maxYears == 0 {
		re2 := regexp.MustCompile(`(\d+)\+?\s*(?:years?|yrs?)\s+of\s+experience`)
		matches2 := re2.FindAllStringSubmatch(strings.ToLower(content), -1)
		for _, m := range matches2 {
			if len(m) > 1 {
				if y, err := strconv.Atoi(m[1]); err == nil && y > maxYears {
					maxYears = y
				}
			}
		}
	}
	return maxYears
}

func extractEducationLevel(content string) string {
	lower := strings.ToLower(content)
	if strings.Contains(lower, "phd") || strings.Contains(lower, "ph.d") || strings.Contains(lower, "doctorate") {
		return "phd"
	}
	if strings.Contains(lower, "master") || strings.Contains(lower, "msc") || strings.Contains(lower, "m.s.") || strings.Contains(lower, "mba") {
		return "master"
	}
	if strings.Contains(lower, "bachelor") || strings.Contains(lower, "b.s.") || strings.Contains(lower, "bsc") || strings.Contains(lower, "bs ") || strings.Contains(lower, "b.a.") || strings.Contains(lower, "degree") {
		return "bachelor"
	}
	return "unknown"
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

func norm(score float64) float64 {
	if score > 1.0 {
		return score / 100.0
	}
	return score
}

func truncateStr(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

package model

import (
	"time"

	"github.com/google/uuid"
)

type MatchResult struct {
	ID              string    `json:"id"`
	UserID          string    `json:"user_id"`
	CVID            string    `json:"cv_id"`
	JDID            string    `json:"jd_id"`
	CVTitle         string    `json:"cv_title,omitempty"`
	JDTitle         string    `json:"jd_title,omitempty"`
	CVFileName      string    `json:"cv_file_name,omitempty"`
	OverallScore    float64   `json:"overall_score"`
	SkillMatchScore float64   `json:"skill_match_score"`
	ExperienceScore float64   `json:"experience_score"`
	EducationScore  float64   `json:"education_score"`
	LLMAnalysis     string    `json:"analysis,omitempty"`
	MatchedSkills   []string  `json:"matched_skills,omitempty"`
	MissingSkills   []string  `json:"missing_skills,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
}

type DashboardStats struct {
	TotalCVs      int            `json:"total_cvs"`
	TotalJDs      int            `json:"total_jds"`
	TotalMatches  int            `json:"total_matches"`
	AverageScore  float64        `json:"average_score"`
	TopSkill      string         `json:"top_skill"`
	TopSkillCount int            `json:"top_skill_count"`
	MatchRate     float64        `json:"match_rate"`
	RecentMatches []*MatchResult `json:"recent_matches,omitempty"`
}

func NewMatchResult(userID, cvID, jdID string) *MatchResult {
	return &MatchResult{
		ID:     uuid.New().String(),
		UserID: userID,
		CVID:   cvID,
		JDID:   jdID,
	}
}

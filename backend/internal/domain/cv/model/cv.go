package model

import (
	"time"

	"github.com/google/uuid"
)

type ParsedExperience struct {
	Title       string `json:"title"`
	Company     string `json:"company"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
	Description string `json:"description"`
}

type ParsedEducation struct {
	Degree      string `json:"degree"`
	Field       string `json:"field"`
	Institution string `json:"institution"`
	StartYear   string `json:"start_year"`
	EndYear     string `json:"end_year"`
}

type CV struct {
	ID             string              `json:"id"`
	UserID         string              `json:"user_id"`
	Title          string              `json:"title"`
	Content        string              `json:"content"`
	ParsedSkills   []string            `json:"parsed_skills,omitempty"`
	ParsedExperience []ParsedExperience `json:"parsed_experience,omitempty"`
	ParsedEducation  []ParsedEducation  `json:"parsed_education,omitempty"`
	ParsedSummary  string              `json:"parsed_summary,omitempty"`
	Status         string              `json:"status"`
	FileName       string              `json:"file_name,omitempty"`
	FileSize       int64               `json:"file_size,omitempty"`
	CreatedAt      time.Time           `json:"created_at"`
	UpdatedAt      time.Time           `json:"updated_at"`
}

func NewCV(userID, title, content string) *CV {
	return &CV{
		ID:      uuid.New().String(),
		UserID:  userID,
		Title:   title,
		Content: content,
		Status:  "pending",
	}
}

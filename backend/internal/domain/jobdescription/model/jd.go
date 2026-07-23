package model

import (
	"time"

	"github.com/google/uuid"
)

type JobDescription struct {
	ID              string    `json:"id"`
	UserID          string    `json:"user_id"`
	Title           string    `json:"title"`
	Content         string    `json:"content"`
	RequiredSkills  []string  `json:"required_skills,omitempty"`
	PreferredSkills []string  `json:"preferred_skills,omitempty"`
	ExperienceLevel string    `json:"experience_level,omitempty"`
	EmploymentType  string    `json:"employment_type,omitempty"`
	Location        string    `json:"location,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func NewJobDescription(userID, title, content string) *JobDescription {
	return &JobDescription{
		ID:      uuid.New().String(),
		UserID:  userID,
		Title:   title,
		Content: content,
	}
}

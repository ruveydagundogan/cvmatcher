package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Role struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Permissions []string `json:"permissions"`
}

type UserRole struct {
	UserID    string `json:"user_id"`
	RoleID    string `json:"role_id"`
	RoleName  string `json:"role_name"`
}

func NewUser(email, passwordHash, firstName, lastName string) *User {
	now := time.Now().UTC()
	return &User{
		ID:           uuid.New().String(),
		Email:        email,
		PasswordHash: passwordHash,
		FirstName:    firstName,
		LastName:     lastName,
		Status:       "active",
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

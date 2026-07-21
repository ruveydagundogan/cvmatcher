package model

import (
	"time"

	"github.com/google/uuid"
)

type AuditLog struct {
	ID         string      `json:"id"`
	UserID     string      `json:"user_id"`
	Action     string      `json:"action"`
	Resource   string      `json:"resource"`
	ResourceID string      `json:"resource_id"`
	Details    interface{} `json:"details,omitempty"`
	IPAddress  string      `json:"ip_address"`
	UserAgent  string      `json:"user_agent"`
	CreatedAt  time.Time   `json:"created_at"`
}

func NewAuditLog(userID, action, resource, resourceID, ipAddress, userAgent string, details interface{}) *AuditLog {
	return &AuditLog{
		ID:         uuid.New().String(),
		UserID:     userID,
		Action:     action,
		Resource:   resource,
		ResourceID: resourceID,
		Details:    details,
		IPAddress:  ipAddress,
		UserAgent:  userAgent,
		CreatedAt:  time.Now().UTC(),
	}
}

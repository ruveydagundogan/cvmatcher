package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/ruveydagundogan/llm-decision-score/backend/internal/shared/response"
)

type TokenValidator interface {
	ValidateToken(token string) (userID, email, role string, err error)
}

const (
	contextKeyUserID    contextKey = "user_id"
	contextKeyUserEmail contextKey = "user_email"
	contextKeyUserRole  contextKey = "user_role"
)

func JWTAuth(validator TokenValidator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				response.Unauthorized(w, "authorization header required")
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
				response.Unauthorized(w, "invalid authorization format, expected: Bearer <token>")
				return
			}

			userID, email, role, err := validator.ValidateToken(parts[1])
			if err != nil {
				response.Unauthorized(w, "invalid or expired token")
				return
			}

			ctx := context.WithValue(r.Context(), contextKeyUserID, userID)
			ctx = context.WithValue(ctx, contextKeyUserEmail, email)
			ctx = context.WithValue(ctx, contextKeyUserRole, role)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequirePermission(permission string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role, ok := r.Context().Value(contextKeyUserRole).(string)
			if !ok {
				response.Unauthorized(w, "authentication required")
				return
			}

			if role == "admin" {
				next.ServeHTTP(w, r)
				return
			}

			if !hasPermission(role, permission) {
				response.Error(w, response.ForbiddenError("insufficient permissions"))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func GetUserID(ctx context.Context) string {
	if userID, ok := ctx.Value(contextKeyUserID).(string); ok {
		return userID
	}
	return ""
}

func GetUserEmail(ctx context.Context) string {
	if email, ok := ctx.Value(contextKeyUserEmail).(string); ok {
		return email
	}
	return ""
}

func GetUserRole(ctx context.Context) string {
	if role, ok := ctx.Value(contextKeyUserRole).(string); ok {
		return role
	}
	return ""
}

var rolePermissions = map[string][]string{
	"user":  {"score:write", "history:read", "history:write", "profile:read", "profile:write"},
	"admin": {"score:write", "history:read", "history:write", "profile:read", "profile:write", "user:read", "user:write"},
}

func hasPermission(role, permission string) bool {
	rolePerms, ok := rolePermissions[role]
	if !ok {
		return false
	}

	for _, p := range rolePerms {
		if p == permission || p == "*" {
			return true
		}
	}

	return false
}

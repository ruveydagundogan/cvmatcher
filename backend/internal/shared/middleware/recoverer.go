package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"

	"github.com/ruveydagundogan/llm-decision-score/backend/internal/shared/response"
)

type contextKey string

const (
	contextKeyRequestID contextKey = "request_id"
)

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Header.Get("X-Request-ID")
		if reqID == "" {
			reqID = uuid.New().String()
		}
		ctx := context.WithValue(r.Context(), contextKeyRequestID, reqID)
		w.Header().Set("X-Request-ID", reqID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func Recoverer(logger interface{ Error(string, ...interface{}) }) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					if logger != nil {
						logger.Error("panic recovered", "error", rec)
					}
					response.Error(w, response.InternalError("internal server error"))
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

func GetRequestID(ctx context.Context) string {
	if reqID, ok := ctx.Value(contextKeyRequestID).(string); ok {
		return reqID
	}
	return ""
}

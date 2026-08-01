package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"

	"github.com/ruveydagundogan/cvmatcher/backend/internal/shared/response"
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

type writeCatcher struct {
	http.ResponseWriter
	wroteHeader bool
}

func (w *writeCatcher) WriteHeader(code int) {
	if w.wroteHeader {
		return
	}
	w.wroteHeader = true
	w.ResponseWriter.WriteHeader(code)
}

func (w *writeCatcher) Write(b []byte) (int, error) {
	if !w.wroteHeader {
		w.wroteHeader = true
	}
	return w.ResponseWriter.Write(b)
}

func Recoverer(logger interface{ Error(string, ...interface{}) }) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			wrapped := &writeCatcher{ResponseWriter: w}
			defer func() {
				if rec := recover(); rec != nil {
					if logger != nil {
						logger.Error("panic recovered", "error", rec, "path", r.URL.Path)
					}
					if !wrapped.wroteHeader {
						response.Error(w, response.InternalError("internal server error"))
					}
				}
			}()
			next.ServeHTTP(wrapped, r)
		})
	}
}

func GetRequestID(ctx context.Context) string {
	if reqID, ok := ctx.Value(contextKeyRequestID).(string); ok {
		return reqID
	}
	return ""
}

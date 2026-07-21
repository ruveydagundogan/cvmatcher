package middleware

import (
	"log/slog"
	"net/http"
	"strings"
	"time"
)

func Logging(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			wrapped := &statusResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

			next.ServeHTTP(wrapped, r)

			if logger != nil {
				logger.Info("request",
					"method", r.Method,
					"path", r.URL.Path,
					"status", wrapped.statusCode,
					"duration", time.Since(start).String(),
					"remote_addr", maskIP(r.RemoteAddr),
				)
			}
		})
	}
}

type statusResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *statusResponseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func maskIP(ip string) string {
	if idx := strings.LastIndex(ip, ":"); idx != -1 {
		return ip[:idx] + ":***"
	}
	return ip
}

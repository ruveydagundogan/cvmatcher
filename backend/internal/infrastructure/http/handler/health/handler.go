package health

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	pool   *pgxpool.Pool
	logger *slog.Logger
}

func NewHandler(pool *pgxpool.Pool, logger *slog.Logger) *Handler {
	return &Handler{pool: pool, logger: logger}
}

func (h *Handler) Live(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{"status": "alive"}); err != nil {
		if h.logger != nil {
			h.logger.Error("json encode error", "error", err)
		}
	}
}

func (h *Handler) Ready(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	services := map[string]string{}

	if h.pool != nil {
		if err := h.pool.Ping(ctx); err != nil {
			services["postgres"] = "unhealthy"
		} else {
			services["postgres"] = "healthy"
		}
	} else {
		services["postgres"] = "unavailable"
	}

	status := "ready"
	for _, v := range services {
		if v != "healthy" && v != "unavailable" {
			status = "degraded"
			break
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"status":   status,
		"services": services,
	}); err != nil {
		if h.logger != nil {
			h.logger.Error("json encode error", "error", err)
		}
	}
}

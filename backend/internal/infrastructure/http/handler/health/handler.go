package health

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	pool *pgxpool.Pool
}

func NewHandler(pool *pgxpool.Pool) *Handler {
	return &Handler{pool: pool}
}

func (h *Handler) Live(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "alive"})
}

func (h *Handler) Ready(w http.ResponseWriter, r *http.Request) {
	services := map[string]string{}

	if h.pool != nil {
		if err := h.pool.Ping(context.Background()); err != nil {
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
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":   status,
		"services": services,
	})
}

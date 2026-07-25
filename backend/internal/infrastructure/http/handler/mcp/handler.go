package mcp

import (
	"encoding/json"
	"net/http"

	mcpengine "github.com/ruveydagundogan/cvmatcher/backend/internal/infrastructure/mcp"
	"github.com/ruveydagundogan/cvmatcher/backend/internal/shared/response"
)

type Handler struct {
	engine *mcpengine.Engine
}

func NewHandler(engine *mcpengine.Engine) *Handler {
	return &Handler{engine: engine}
}

func (h *Handler) Query(w http.ResponseWriter, r *http.Request) {
	var req mcpengine.Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	mcpResp, richResult, err := h.engine.Query(r.Context(), &req)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.Success(w, map[string]any{
		"response":    mcpResp,
		"rich_result": richResult,
	})
}

func (h *Handler) ListAdapters(w http.ResponseWriter, r *http.Request) {
	adapters := h.engine.ListAdapters()
	response.Success(w, adapters)
}

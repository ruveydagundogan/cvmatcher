package admin

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"

	adminmodel "github.com/ruveydagundogan/cvmatcher/backend/internal/domain/admin/model"
	adminrepo "github.com/ruveydagundogan/cvmatcher/backend/internal/domain/admin/repository"
	adminuc "github.com/ruveydagundogan/cvmatcher/backend/internal/application/admin/usecase"
	mcpengine "github.com/ruveydagundogan/cvmatcher/backend/internal/infrastructure/mcp"
	"github.com/ruveydagundogan/cvmatcher/backend/internal/shared/response"
)

type Handler struct {
	adminUC   *adminuc.AdminUseCase
	mcpEngine *mcpengine.Engine
}

func NewHandler(adminUC *adminuc.AdminUseCase, mcpEngine *mcpengine.Engine) *Handler {
	return &Handler{adminUC: adminUC, mcpEngine: mcpEngine}
}

func NewHandlerFromRepos(repo adminrepo.AdminRepository, mcpEngine *mcpengine.Engine, log *slog.Logger) *Handler {
	return &Handler{adminUC: adminuc.NewAdminUseCase(repo, log), mcpEngine: mcpEngine}
}

type createAdapterRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	FilePath    string `json:"file_path"`
	ModelName   string `json:"model_name"`
}

func (h *Handler) ListAdapters(w http.ResponseWriter, r *http.Request) {
	adapters, err := h.adminUC.ListAdapters(r.Context())
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Success(w, adapters)
}

func (h *Handler) CreateAdapter(w http.ResponseWriter, r *http.Request) {
	var req createAdapterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}
	adapter, err := h.adminUC.CreateAdapter(r.Context(), req.Name, req.Description, req.FilePath, req.ModelName)
	if err != nil {
		response.Error(w, err)
		return
	}
	h.mcpEngine.LoadAdapter(req.Name, req.Description)
	response.Created(w, adapter)
}

func (h *Handler) DeleteAdapter(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.adminUC.DeleteAdapter(r.Context(), id); err != nil {
		response.Error(w, err)
		return
	}
	h.mcpEngine.UnloadAdapter(id)
	w.WriteHeader(http.StatusNoContent)
}

type promptRequest struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

func (h *Handler) ListPrompts(w http.ResponseWriter, r *http.Request) {
	prompts, err := h.adminUC.ListPrompts(r.Context())
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Success(w, prompts)
}

func (h *Handler) CreatePrompt(w http.ResponseWriter, r *http.Request) {
	var req promptRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}
	prompt, err := h.adminUC.CreatePrompt(r.Context(), req.Name, req.Content)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Created(w, prompt)
}

func (h *Handler) UpdatePrompt(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req promptRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}
	prompt := &adminmodel.SystemPrompt{ID: id, Name: req.Name, Content: req.Content}
	if err := h.adminUC.UpdatePrompt(r.Context(), prompt); err != nil {
		response.Error(w, err)
		return
	}
	response.Success(w, prompt)
}

func (h *Handler) ActivatePrompt(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	prompt, err := h.adminUC.ActivatePrompt(r.Context(), id)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Success(w, prompt)
}

func (h *Handler) DeletePrompt(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.adminUC.DeletePrompt(r.Context(), id); err != nil {
		response.Error(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetSettings(w http.ResponseWriter, r *http.Request) {
	settings, err := h.adminUC.GetSettings(r.Context())
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Success(w, settings)
}

func (h *Handler) SaveSettings(w http.ResponseWriter, r *http.Request) {
	var settings adminmodel.LLMSettings
	if err := json.NewDecoder(r.Body).Decode(&settings); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}
	if err := h.adminUC.SaveSettings(r.Context(), &settings); err != nil {
		response.Error(w, err)
		return
	}
	response.Success(w, &settings)
}

func (h *Handler) ListLogs(w http.ResponseWriter, r *http.Request) {
	offset := 0
	limit := 50
	logs, total, err := h.adminUC.ListLogs(r.Context(), offset, limit)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Success(w, map[string]any{
		"items": logs,
		"total": total,
	})
}

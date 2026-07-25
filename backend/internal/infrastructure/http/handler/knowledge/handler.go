package knowledge

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	knowledgeuc "github.com/ruveydagundogan/cvmatcher/backend/internal/application/knowledge/usecase"
	"github.com/ruveydagundogan/cvmatcher/backend/internal/shared/middleware"
	"github.com/ruveydagundogan/cvmatcher/backend/internal/shared/response"
)

type Handler struct {
	knowledgeUC *knowledgeuc.KnowledgeUseCase
}

func NewHandler(knowledgeUC *knowledgeuc.KnowledgeUseCase) *Handler {
	return &Handler{knowledgeUC: knowledgeUC}
}

type createRequest struct {
	Title    string   `json:"title"`
	Content  string   `json:"content"`
	Category string   `json:"category"`
	Source   string   `json:"source"`
	Tags     []string `json:"tags"`
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	var req createRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}
	entry, err := h.knowledgeUC.Create(r.Context(), userID, req.Title, req.Content, req.Category, req.Source, req.Tags)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Created(w, entry)
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	entry, err := h.knowledgeUC.GetByID(r.Context(), id)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Success(w, entry)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	items, total, err := h.knowledgeUC.List(r.Context(), userID, offset, limit)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Success(w, map[string]any{
		"items": items,
		"total": total,
	})
}

func (h *Handler) Search(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	tags := r.URL.Query()["tag"]
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 {
		limit = 10
	}
	results, err := h.knowledgeUC.Search(r.Context(), query, tags, limit)
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Success(w, results)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.knowledgeUC.Delete(r.Context(), id); err != nil {
		response.Error(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ListCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.knowledgeUC.ListCategories(r.Context())
	if err != nil {
		response.Error(w, err)
		return
	}
	response.Success(w, categories)
}

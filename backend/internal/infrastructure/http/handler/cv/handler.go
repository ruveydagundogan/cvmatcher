package cv

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	cvmodel "github.com/ruveydagundogan/cvmatcher/backend/internal/domain/cv/model"
	"github.com/go-chi/chi/v5"
	"github.com/ruveydagundogan/cvmatcher/backend/internal/shared/middleware"
	"github.com/ruveydagundogan/cvmatcher/backend/internal/shared/response"
)

type CVUseCase interface {
	Create(ctx context.Context, userID, title, content string) (*cvmodel.CV, error)
	GetByID(ctx context.Context, id string) (*cvmodel.CV, error)
	ListByUser(ctx context.Context, userID string, offset, limit int) ([]*cvmodel.CV, int, error)
	Delete(ctx context.Context, id string) error
	ParseWithLLM(ctx context.Context, cvID string) (*cvmodel.CV, error)
}

type Handler struct {
	cvUC CVUseCase
}

func NewHandler(cvUC CVUseCase) *Handler {
	return &Handler{cvUC: cvUC}
}

type createCVRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type listResponse struct {
	Items      []*cvmodel.CV `json:"items"`
	Total      int           `json:"total"`
	Page       int           `json:"page"`
	Limit      int           `json:"limit"`
	TotalPages int           `json:"total_pages"`
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		response.Unauthorized(w, "authentication required")
		return
	}

	var req createCVRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}
	if req.Title == "" || req.Content == "" {
		response.BadRequest(w, "title and content are required")
		return
	}

	cv, err := h.cvUC.Create(r.Context(), userID, req.Title, req.Content)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.Created(w, cv)
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.BadRequest(w, "id is required")
		return
	}

	cv, err := h.cvUC.GetByID(r.Context(), id)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.Success(w, cv)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		response.Unauthorized(w, "authentication required")
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	offset := (page - 1) * limit
	items, total, err := h.cvUC.ListByUser(r.Context(), userID, offset, limit)
	if err != nil {
		response.Error(w, err)
		return
	}

	totalPages := total / limit
	if total%limit > 0 {
		totalPages++
	}

	response.Success(w, listResponse{
		Items:      items,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	})
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.BadRequest(w, "id is required")
		return
	}

	if err := h.cvUC.Delete(r.Context(), id); err != nil {
		response.Error(w, err)
		return
	}

	response.SuccessWithMessage(w, "cv deleted")
}

func (h *Handler) Parse(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.BadRequest(w, "id is required")
		return
	}

	cv, err := h.cvUC.ParseWithLLM(r.Context(), id)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.Success(w, cv)
}

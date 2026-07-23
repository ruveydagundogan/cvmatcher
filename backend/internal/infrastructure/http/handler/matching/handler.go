package matching

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	matchmodel "github.com/ruveydagundogan/cvmatcher/backend/internal/domain/matching/model"
	"github.com/go-chi/chi/v5"
	"github.com/ruveydagundogan/cvmatcher/backend/internal/shared/middleware"
	"github.com/ruveydagundogan/cvmatcher/backend/internal/shared/response"
)

type MatchUseCase interface {
	RunMatch(ctx context.Context, userID, cvID, jdID string) (*matchmodel.MatchResult, error)
	GetByID(ctx context.Context, id string) (*matchmodel.MatchResult, error)
	ListByUser(ctx context.Context, userID string, offset, limit int) ([]*matchmodel.MatchResult, int, error)
	GetDashboardStats(ctx context.Context, userID string) (*matchmodel.DashboardStats, error)
}

type Handler struct {
	matchUC MatchUseCase
}

func NewHandler(matchUC MatchUseCase) *Handler {
	return &Handler{matchUC: matchUC}
}

type matchRequest struct {
	CVID string `json:"cv_id"`
	JDID string `json:"jd_id"`
}

type listResponse struct {
	Items      []*matchmodel.MatchResult `json:"items"`
	Total      int                       `json:"total"`
	Page       int                       `json:"page"`
	Limit      int                       `json:"limit"`
	TotalPages int                       `json:"total_pages"`
}

func (h *Handler) RunMatch(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		response.Unauthorized(w, "authentication required")
		return
	}

	var req matchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}
	if req.CVID == "" || req.JDID == "" {
		response.BadRequest(w, "cv_id and jd_id are required")
		return
	}

	result, err := h.matchUC.RunMatch(r.Context(), userID, req.CVID, req.JDID)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.Created(w, result)
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.BadRequest(w, "id is required")
		return
	}

	result, err := h.matchUC.GetByID(r.Context(), id)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.Success(w, result)
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
	items, total, err := h.matchUC.ListByUser(r.Context(), userID, offset, limit)
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

func (h *Handler) GetDashboardStats(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		response.Unauthorized(w, "authentication required")
		return
	}

	stats, err := h.matchUC.GetDashboardStats(r.Context(), userID)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.Success(w, stats)
}

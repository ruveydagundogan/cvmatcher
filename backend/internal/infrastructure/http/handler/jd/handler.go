package jd

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	jdmodel "github.com/ruveydagundogan/cvmatcher/backend/internal/domain/jobdescription/model"
	"github.com/go-chi/chi/v5"
	"github.com/ruveydagundogan/cvmatcher/backend/internal/shared/middleware"
	"github.com/ruveydagundogan/cvmatcher/backend/internal/shared/response"
)

type JDUseCase interface {
	Create(ctx context.Context, userID, title, content string) (*jdmodel.JobDescription, error)
	GetByID(ctx context.Context, id string) (*jdmodel.JobDescription, error)
	ListByUser(ctx context.Context, userID string, offset, limit int) ([]*jdmodel.JobDescription, int, error)
	Update(ctx context.Context, jd *jdmodel.JobDescription) error
	Delete(ctx context.Context, id string) error
	AnalyzeWithLLM(ctx context.Context, jdID string) (*jdmodel.JobDescription, error)
}

type Handler struct {
	jdUC JDUseCase
}

func NewHandler(jdUC JDUseCase) *Handler {
	return &Handler{jdUC: jdUC}
}

type createJDRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type updateJDRequest struct {
	Title           string   `json:"title,omitempty"`
	Content         string   `json:"content,omitempty"`
	RequiredSkills  []string `json:"required_skills,omitempty"`
	PreferredSkills []string `json:"preferred_skills,omitempty"`
	ExperienceLevel string   `json:"experience_level,omitempty"`
	EmploymentType  string   `json:"employment_type,omitempty"`
	Location        string   `json:"location,omitempty"`
}

type listResponse struct {
	Items      []*jdmodel.JobDescription `json:"items"`
	Total      int                       `json:"total"`
	Page       int                       `json:"page"`
	Limit      int                       `json:"limit"`
	TotalPages int                       `json:"total_pages"`
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		response.Unauthorized(w, "authentication required")
		return
	}

	var req createJDRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}
	if req.Title == "" || req.Content == "" {
		response.BadRequest(w, "title and content are required")
		return
	}

	jd, err := h.jdUC.Create(r.Context(), userID, req.Title, req.Content)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.Created(w, jd)
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.BadRequest(w, "id is required")
		return
	}

	jd, err := h.jdUC.GetByID(r.Context(), id)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.Success(w, jd)
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
	items, total, err := h.jdUC.ListByUser(r.Context(), userID, offset, limit)
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

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.BadRequest(w, "id is required")
		return
	}

	jd, err := h.jdUC.GetByID(r.Context(), id)
	if err != nil {
		response.Error(w, err)
		return
	}

	var req updateJDRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	if req.Title != "" {
		jd.Title = req.Title
	}
	if req.Content != "" {
		jd.Content = req.Content
	}
	if req.RequiredSkills != nil {
		jd.RequiredSkills = req.RequiredSkills
	}
	if req.PreferredSkills != nil {
		jd.PreferredSkills = req.PreferredSkills
	}
	if req.ExperienceLevel != "" {
		jd.ExperienceLevel = req.ExperienceLevel
	}
	if req.EmploymentType != "" {
		jd.EmploymentType = req.EmploymentType
	}
	if req.Location != "" {
		jd.Location = req.Location
	}

	if err := h.jdUC.Update(r.Context(), jd); err != nil {
		response.Error(w, err)
		return
	}

	response.Success(w, jd)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.BadRequest(w, "id is required")
		return
	}

	if err := h.jdUC.Delete(r.Context(), id); err != nil {
		response.Error(w, err)
		return
	}

	response.SuccessWithMessage(w, "jd deleted")
}

func (h *Handler) Analyze(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		response.BadRequest(w, "id is required")
		return
	}

	jd, err := h.jdUC.AnalyzeWithLLM(r.Context(), id)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.Success(w, jd)
}

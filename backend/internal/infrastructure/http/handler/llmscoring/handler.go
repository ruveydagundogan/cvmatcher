package llmscoring

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/ruveydagundogan/llm-decision-score/backend/internal/application/llmscoring/dto"
	"github.com/ruveydagundogan/llm-decision-score/backend/internal/shared/middleware"
	"github.com/ruveydagundogan/llm-decision-score/backend/internal/shared/response"
)

const (
	maxPromptLen   = 10000
	maxResponseLen = 50000
)

type ScoreUseCase interface {
	Execute(ctx context.Context, userID string, req dto.ScoreRequestDTO) (*dto.ScoreResponseDTO, error)
}

type GetHistoryUseCase interface {
	Execute(ctx context.Context, userID string, page, limit int) (*dto.HistoryListResponseDTO, error)
}

type DeleteHistoryUseCase interface {
	Execute(ctx context.Context, userID string) error
}

type GetStatsUseCase interface {
	Execute(ctx context.Context, userID string) (*dto.StatsResponseDTO, error)
}

type Handler struct {
	scoreUC      ScoreUseCase
	historyUC    GetHistoryUseCase
	deleteHistUC DeleteHistoryUseCase
	statsUC      GetStatsUseCase
}

func NewHandler(
	scoreUC ScoreUseCase,
	historyUC GetHistoryUseCase,
	deleteHistUC DeleteHistoryUseCase,
	statsUC GetStatsUseCase,
) *Handler {
	return &Handler{
		scoreUC:      scoreUC,
		historyUC:    historyUC,
		deleteHistUC: deleteHistUC,
		statsUC:      statsUC,
	}
}

func (h *Handler) Score(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	var req dto.ScoreRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	if req.Prompt == "" || req.Response == "" {
		response.BadRequest(w, "prompt and response are required")
		return
	}

	if len(req.Prompt) > maxPromptLen {
		response.BadRequest(w, "prompt exceeds maximum length")
		return
	}
	if len(req.Response) > maxResponseLen {
		response.BadRequest(w, "response exceeds maximum length")
		return
	}

	for i := 0; i < len(req.Prompt); i++ {
		if req.Prompt[i] == 0 {
			response.BadRequest(w, "invalid character in prompt")
			return
		}
	}
	for i := 0; i < len(req.Response); i++ {
		if req.Response[i] == 0 {
			response.BadRequest(w, "invalid character in response")
			return
		}
	}

	req.Prompt = strings.TrimSpace(req.Prompt)
	req.Response = strings.TrimSpace(req.Response)

	if req.Model == "" {
		req.Model = "Gemma-2B"
	}

	result, err := h.scoreUC.Execute(r.Context(), userID, req)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.Success(w, result)
}

func (h *Handler) GetHistory(w http.ResponseWriter, r *http.Request) {
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

	result, err := h.historyUC.Execute(r.Context(), userID, page, limit)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.Success(w, result)
}

func (h *Handler) DeleteHistory(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		response.Unauthorized(w, "authentication required")
		return
	}

	if err := h.deleteHistUC.Execute(r.Context(), userID); err != nil {
		response.Error(w, err)
		return
	}

	response.SuccessWithMessage(w, "history deleted")
}

func (h *Handler) GetStats(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		response.Unauthorized(w, "authentication required")
		return
	}

	result, err := h.statsUC.Execute(r.Context(), userID)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.Success(w, result)
}

func (h *Handler) GetModels(w http.ResponseWriter, r *http.Request) {
	response.Success(w, map[string]interface{}{
		"models": []string{"Gemma-2B", "Gemma-7B"},
	})
}

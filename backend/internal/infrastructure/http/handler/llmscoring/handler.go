package llmscoring

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ruveydagundogan/llm-decision-score/backend/internal/application/llmscoring/dto"
	llmusecase "github.com/ruveydagundogan/llm-decision-score/backend/internal/application/llmscoring/usecase"
	"github.com/ruveydagundogan/llm-decision-score/backend/internal/shared/middleware"
	"github.com/ruveydagundogan/llm-decision-score/backend/internal/shared/response"
)

type Handler struct {
	scoreUC       *llmusecase.ScoreUseCase
	historyUC     *llmusecase.GetHistoryUseCase
	deleteHistUC  *llmusecase.DeleteHistoryUseCase
	statsUC       *llmusecase.GetStatsUseCase
}

func NewHandler(
	scoreUC *llmusecase.ScoreUseCase,
	historyUC *llmusecase.GetHistoryUseCase,
	deleteHistUC *llmusecase.DeleteHistoryUseCase,
	statsUC *llmusecase.GetStatsUseCase,
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
	if userID == "" {
		userID = "anonymous"
	}

	var req dto.ScoreRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	if req.Prompt == "" || req.Response == "" {
		response.BadRequest(w, "prompt and response are required")
		return
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

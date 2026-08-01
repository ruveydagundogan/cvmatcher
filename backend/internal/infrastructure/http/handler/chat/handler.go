package chat

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	chatuc "github.com/ruveydagundogan/cvmatcher/backend/internal/application/chat/usecase"
	chatmodel "github.com/ruveydagundogan/cvmatcher/backend/internal/domain/chat/model"
	"github.com/ruveydagundogan/cvmatcher/backend/internal/shared/middleware"
	"github.com/ruveydagundogan/cvmatcher/backend/internal/shared/response"
)

type ChatUseCase interface {
	CreateConversation(ctx context.Context, userID, title, cvID string) (*chatmodel.Conversation, error)
	ListConversations(ctx context.Context, userID string) ([]*chatmodel.Conversation, error)
	GetConversation(ctx context.Context, userID, convID string) (*chatuc.ConversationDetail, error)
	DeleteConversation(ctx context.Context, userID, convID string) error
	SendMessage(ctx context.Context, userID, convID, content string) (*chatuc.SendMessageResult, error)
}

type Handler struct {
	chatUC ChatUseCase
}

func NewHandler(chatUC ChatUseCase) *Handler {
	return &Handler{chatUC: chatUC}
}

type createConversationRequest struct {
	Title string `json:"title"`
	CVID  string `json:"cv_id"`
}

type sendMessageRequest struct {
	Content string `json:"content"`
}

func (h *Handler) CreateConversation(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		response.Unauthorized(w, "authentication required")
		return
	}

	var req createConversationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	conv, err := h.chatUC.CreateConversation(r.Context(), userID, req.Title, req.CVID)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.Created(w, conv)
}

func (h *Handler) ListConversations(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		response.Unauthorized(w, "authentication required")
		return
	}

	convs, err := h.chatUC.ListConversations(r.Context(), userID)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.Success(w, convs)
}

func (h *Handler) GetConversation(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		response.Unauthorized(w, "authentication required")
		return
	}

	id := chi.URLParam(r, "id")
	conv, err := h.chatUC.GetConversation(r.Context(), userID, id)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.Success(w, conv)
}

func (h *Handler) DeleteConversation(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		response.Unauthorized(w, "authentication required")
		return
	}

	id := chi.URLParam(r, "id")
	if err := h.chatUC.DeleteConversation(r.Context(), userID, id); err != nil {
		response.Error(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) SendMessage(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		response.Unauthorized(w, "authentication required")
		return
	}

	id := chi.URLParam(r, "id")
	var req sendMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	result, err := h.chatUC.SendMessage(r.Context(), userID, id, req.Content)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.Success(w, result)
}

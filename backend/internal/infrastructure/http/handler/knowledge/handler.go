package knowledge

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"

	knowledgeuc "github.com/ruveydagundogan/cvmatcher/backend/internal/application/knowledge/usecase"
	mcpengine "github.com/ruveydagundogan/cvmatcher/backend/internal/infrastructure/mcp"
	"github.com/ruveydagundogan/cvmatcher/backend/internal/shared/middleware"
	"github.com/ruveydagundogan/cvmatcher/backend/internal/shared/response"
)

type Handler struct {
	knowledgeUC *knowledgeuc.KnowledgeUseCase
	mcpEngine   *mcpengine.Engine
}

func NewHandler(knowledgeUC *knowledgeuc.KnowledgeUseCase, mcpEngine *mcpengine.Engine) *Handler {
	return &Handler{knowledgeUC: knowledgeUC, mcpEngine: mcpEngine}
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

type queryAIRequest struct {
	Question string   `json:"question"`
	EntryIDs []string `json:"entry_ids,omitempty"`
}

func (h *Handler) QueryAI(w http.ResponseWriter, r *http.Request) {
	var req queryAIRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}
	if req.Question == "" {
		response.BadRequest(w, "question is required")
		return
	}

	var contextStr string
	if len(req.EntryIDs) > 0 {
		var entries []string
		for _, id := range req.EntryIDs {
			entry, err := h.knowledgeUC.GetByID(r.Context(), id)
			if err == nil {
				entries = append(entries, fmt.Sprintf("Title: %s\nCategory: %s\nTags: %s\nContent: %s",
					entry.Title, entry.Category, strings.Join(entry.Tags, ", "), entry.Content))
			}
		}
		if len(entries) > 0 {
			contextStr = "Knowledge Base Context:\n" + strings.Join(entries, "\n---\n")
		}
	}

	systemPrompt := "You are a knowledge base assistant. Answer the user's question based on the provided knowledge context. If the context is empty or insufficient, say so clearly."
	fullContext := req.Question
	if contextStr != "" {
		fullContext = fmt.Sprintf("%s\n\nUser Question: %s", contextStr, req.Question)
	}

	mcpReq := &mcpengine.Request{
		SystemPrompt: systemPrompt,
		Messages:     []mcpengine.ChatMessage{{Role: "user", Content: fullContext}},
		MaxTokens:    1024,
		Temperature:  0.3,
		TopP:         0.9,
	}

	mcpResp, _, err := h.mcpEngine.Query(context.Background(), mcpReq)
	if err != nil {
		response.Error(w, fmt.Errorf("ai query failed: %w", err))
		return
	}

	answer := ""
	if len(mcpResp.Choices) > 0 {
		answer = mcpResp.Choices[0].Message.Content
	}

	response.Success(w, map[string]any{
		"answer":      answer,
		"context_ids": req.EntryIDs,
	})
}

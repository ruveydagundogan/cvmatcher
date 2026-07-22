package backendllm

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ruveydagundogan/llm-decision-score/backend/internal/infrastructure/llm"
	"github.com/ruveydagundogan/llm-decision-score/backend/internal/shared/response"
)

type ChatRequest struct {
	Prompt   string `json:"prompt"`
	MaxTokens int   `json:"max_tokens"`
}

type Handler struct {
	llmClient *llm.Client
}

func NewHandler(client *llm.Client) *Handler {
	return &Handler{llmClient: client}
}

func (h *Handler) Chat(w http.ResponseWriter, r *http.Request) {
	if h.llmClient == nil {
		response.Error(w, errors.New("server-side LLM is not enabled"))
		return
	}

	var req ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	if req.Prompt == "" {
		response.BadRequest(w, "prompt is required")
		return
	}

	if req.MaxTokens <= 0 {
		req.MaxTokens = 256
	}

	content, err := h.llmClient.ChatCompletion(r.Context(), []llm.ChatMessage{
		{Role: "user", Content: req.Prompt},
	}, req.MaxTokens)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.Success(w, map[string]string{"response": content})
}

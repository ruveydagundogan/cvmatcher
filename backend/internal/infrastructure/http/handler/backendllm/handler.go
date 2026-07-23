package backendllm

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ruveydagundogan/llm-decision-score/backend/internal/infrastructure/llm"
	"github.com/ruveydagundogan/llm-decision-score/backend/internal/infrastructure/metrics"
	"github.com/ruveydagundogan/llm-decision-score/backend/internal/shared/response"
)

type ChatRequest struct {
	Prompt    string `json:"prompt"`
	MaxTokens int    `json:"max_tokens"`
}

type Handler struct {
	llmClient *llm.Client
	metrics   *metrics.Metrics
}

func NewHandler(client *llm.Client, m *metrics.Metrics) *Handler {
	return &Handler{llmClient: client, metrics: m}
}

func (h *Handler) Chat(w http.ResponseWriter, r *http.Request) {
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

	var content string
	var err error

	if h.llmClient != nil {
		start := time.Now()
		content, err = h.llmClient.ChatCompletion(r.Context(), []llm.ChatMessage{
			{Role: "user", Content: req.Prompt},
		}, req.MaxTokens)
		h.metrics.LLM.Observe(time.Since(start), err == nil)
		if err != nil {
			response.Error(w, err)
			return
		}
	} else {
		content = mockResponse(req.Prompt)
	}

	response.Success(w, map[string]string{"response": content})
}

func mockResponse(prompt string) string {
	lower := strings.ToLower(prompt)
	words := strings.Fields(prompt)

	var responses []string

	if strings.Contains(lower, "merhaba") || strings.Contains(lower, "selam") || strings.Contains(lower, "hello") {
		responses = append(responses, "Merhaba! Size nasıl yardımcı olabilirim?")
	}
	if strings.Contains(lower, "nasılsın") || strings.Contains(lower, "how are you") {
		responses = append(responses, "Ben bir AI asistanıyım, her zaman iyiyim!")
	}
	if strings.Contains(lower, "yardım") || strings.Contains(lower, "help") {
		responses = append(responses, "Size nasıl yardımcı olabilirim? Lütfen sorunuzu detaylandırın.")
	}
	if strings.Contains(lower, "hava") || strings.Contains(lower, "weather") {
		responses = append(responses, "Hava durumu hakkında güncel bilgim yok, ancak bugünün harika bir gün olduğunu söyleyebilirim!")
	}
	if strings.Contains(lower, "teşekkür") || strings.Contains(lower, "thanks") {
		responses = append(responses, "Rica ederim! Başka bir şeye yardımcı olabilir miyim?")
	}

	if len(responses) > 0 {
		return strings.Join(responses, " ")
	}

	if len(words) <= 3 {
		return fmt.Sprintf("'%s' hakkında düşünmeme izin verin. Bu ilginç bir konu. Daha fazla bilgi verebilir misiniz?", prompt)
	}

	return fmt.Sprintf("'%s' konusunu değerlendiriyorum. Anladığım kadarıyla bu, üzerinde düşünülmesi gereken bir konu. Şu açıdan bakılabilir: öncelikle ana faktörleri belirlemek, ardından olası sonuçları değerlendirmek faydalı olacaktır. Sizin bu konudaki düşünceleriniz neler?", prompt)
}

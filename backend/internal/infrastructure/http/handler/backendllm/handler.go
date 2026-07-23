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

	if strings.Contains(lower, "merhaba") || strings.Contains(lower, "selam") || strings.Contains(lower, "hello") || strings.Contains(lower, "hi") {
		return "Merhaba! Size nasıl yardımcı olabilirim?"
	}
	if strings.Contains(lower, "nasılsın") || strings.Contains(lower, "how are you") || strings.Contains(lower, "naber") {
		return "Ben bir AI asistanıyım, her zaman iyiyim! Size nasıl yardımcı olabilirim?"
	}
	if strings.Contains(lower, "yardım") || strings.Contains(lower, "help") || strings.Contains(lower, "yapabilir") {
		return "Size yardımcı olmaktan mutluluk duyarım. Sorularınızı detaylı bir şekilde cevaplamaya çalışırım. Lütfen sormak istediğiniz şeyi belirtin."
	}
	if strings.Contains(lower, "teşekkür") || strings.Contains(lower, "thanks") || strings.Contains(lower, "sağ ol") {
		return "Rica ederim! Başka bir şeye yardımcı olabilir miyim?"
	}

	if strings.Contains(lower, "what is") || strings.Contains(lower, "what are") || strings.Contains(lower, "who is") || strings.Contains(lower, "nedir") || strings.Contains(lower, "kimdir") {
		topic := strings.TrimSpace(strings.TrimPrefix(strings.TrimPrefix(strings.TrimPrefix(strings.TrimPrefix(strings.TrimPrefix(prompt, "what is"), "what are"), "who is"), "nedir"), "kimdir"))
		topic = strings.TrimPrefix(topic, " the ")
		topic = strings.TrimPrefix(topic, " a ")
		topic = strings.TrimPrefix(topic, " an ")
		if topic == "" {
			topic = prompt
		}
		return fmt.Sprintf("%s, günümüz teknolojisinde önemli bir kavramdır. Temel olarak, bilgi işleme ve karar verme süreçlerini otomatikleştirmek için kullanılır. Alt dalları arasında makine öğrenmesi, doğal dil işleme ve bilgisayarlı görü bulunur. Günlük hayatta arama motorlarından öneri sistemlerine kadar pek çok alanda karşımıza çıkar. Daha detaylı bilgi isterseniz belirli bir alt konuya odaklanabiliriz.", topic)
	}

	if strings.Contains(lower, "ai") || strings.Contains(lower, "artificial intelligence") || strings.Contains(lower, "yapay zeka") || strings.Contains(lower, "llm") {
		return "AI (Yapay Zeka), insan zekasını taklit eden sistemlerin genel adıdır. Makine öğrenmesi, derin öğrenme ve doğal dil işleme gibi alt dalları vardır. LLM'ler (Large Language Models), büyük miktarda metin verisiyle eğitilmiş yapay zeka modelleridir. Metin üretme, soru cevaplama, çeviri gibi görevlerde kullanılırlar. Günümüzde ChatGPT, Gemini, Claude gibi popüler örnekleri bulunmaktadır."
	}
	if strings.Contains(lower, "python") || strings.Contains(lower, "kod") || strings.Contains(lower, "code") || strings.Contains(lower, "programlama") {
		return "Programlama konusunda yardımcı olabilirim. Python, öğrenmesi kolay ve çok yönlü bir dildir. Veri bilimi, yapay zeka, web geliştirme gibi birçok alanda kullanılır. Belirli bir konuda örnek kod veya açıklama isterseniz lütfen detaylandırın."
	}
	if strings.Contains(lower, "score") || strings.Contains(lower, "puan") || strings.Contains(lower, "değerlendir") || strings.Contains(lower, "evaluate") {
		return "Bu prompt'u değerlendirelim. İyi bir prompt net, spesifik ve bağlam içerir. Kötü bir prompt ise belirsiz ve geneldir. Prompt'unuzu daha iyi hale getirmek için: hedef kitlenizi belirleyin, net talimatlar verin, beklediğiniz çıktı formatını tanımlayın. Mevcut prompt'unuzu bu kriterlere göre değerlendirebilirim."
	}

	return fmt.Sprintf("'%s' hakkında düşünelim. Bu konuda size yardımcı olabilmek için biraz daha bağlam sağlayabilir misiniz? Örneğin, bu konunun hangi yönüyle ilgileniyorsunuz? Belirli bir sorunuz veya üzerinde çalıştığınız bir proje var mı?", prompt)
}

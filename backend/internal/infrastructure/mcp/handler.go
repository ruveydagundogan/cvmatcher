package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/ruveydagundogan/cvmatcher/backend/internal/infrastructure/llm"
)

type Engine struct {
	llmClient       *llm.Client
	log             *slog.Logger
	defaultModel    string
	loadedAdapters  []AdapterInfo
}

type AdapterInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Active      bool   `json:"active"`
}

func NewEngine(llmClient *llm.Client, log *slog.Logger) *Engine {
	return &Engine{
		llmClient:    llmClient,
		log:          log,
		defaultModel: "gemma:2b",
	}
}

func (e *Engine) Query(ctx context.Context, req *Request) (*Response, *RichResult, error) {
	model := req.Model
	if model == "" {
		model = e.defaultModel
	}

	systemMsg := ""
	if req.SystemPrompt != "" {
		systemMsg = req.SystemPrompt
	}

	maxTokens := req.MaxTokens
	if maxTokens <= 0 {
		maxTokens = 2048
	}

	temp := req.Temperature
	if temp == 0 {
		temp = 0.7
	}

	topP := req.TopP
	if topP == 0 {
		topP = 0.9
	}

	var llmMessages []llm.ChatMessage
	if systemMsg != "" {
		llmMessages = append(llmMessages, llm.ChatMessage{Role: "system", Content: systemMsg})
	}
	for _, msg := range req.Messages {
		llmMessages = append(llmMessages, llm.ChatMessage{Role: msg.Role, Content: msg.Content})
	}

	start := time.Now()
	response, err := e.llmClient.ChatCompletion(ctx, llmMessages, maxTokens)
	elapsed := time.Since(start)
	if err != nil {
		return nil, nil, fmt.Errorf("mcp query failed: %w", err)
	}

	e.log.Info("mcp query completed",
		"model", model,
		"tokens", maxTokens,
		"elapsed_ms", elapsed.Milliseconds(),
	)

	mcpResp := &Response{
		ID:    uuid.New().String(),
		Model: model,
		Choices: []Choice{
			{
				Index:        0,
				Message:      Message{Role: "assistant", Content: response},
				FinishReason: "stop",
			},
		},
		Usage: Usage{
			PromptTokens:     estimateTokens(llmMessages),
			CompletionTokens: estimateTokensStr(response),
			TotalTokens:      estimateTokens(llmMessages) + estimateTokensStr(response),
		},
	}

	richResult := &RichResult{
		Text: response,
		Type: determineResultType(response),
	}
	richResult.Sections = parseSections(response)

	return mcpResp, richResult, nil
}

func (e *Engine) ListAdapters() []AdapterInfo {
	return e.loadedAdapters
}

func (e *Engine) LLMClient() *llm.Client {
	return e.llmClient
}

func (e *Engine) LoadAdapter(name, description string) error {
	for i, a := range e.loadedAdapters {
		if a.Name == name {
			e.loadedAdapters[i].Active = true
			return nil
		}
	}
	e.loadedAdapters = append(e.loadedAdapters, AdapterInfo{
		Name:        name,
		Description: description,
		Active:      true,
	})
	e.log.Info("adapter loaded", "name", name)
	return nil
}

func (e *Engine) UnloadAdapter(name string) error {
	for i, a := range e.loadedAdapters {
		if a.Name == name {
			e.loadedAdapters[i].Active = false
			return nil
		}
	}
	return fmt.Errorf("adapter not found: %s", name)
}

func estimateTokens(messages []llm.ChatMessage) int {
	total := 0
	for _, m := range messages {
		total += len(strings.Fields(m.Content))
	}
	return total
}

func estimateTokensStr(s string) int {
	return len(strings.Fields(s))
}

func determineResultType(text string) string {
	lower := strings.ToLower(text)
	if strings.Contains(lower, "|") && strings.Contains(lower, "---") {
		return "table"
	}
	if strings.Contains(lower, "```") {
		return "code"
	}
	if strings.Contains(lower, "1.") || strings.Contains(lower, "- ") {
		return "list"
	}
	return "text"
}

func parseSections(text string) []Section {
	var sections []Section
	lines := strings.Split(text, "\n")
	var current Section

	for _, line := range lines {
		if strings.HasPrefix(line, "## ") || strings.HasPrefix(line, "### ") {
			if current.Title != "" || current.Content != "" {
				sections = append(sections, current)
			}
			current = Section{
				Title: strings.TrimSpace(strings.TrimPrefix(strings.TrimPrefix(line, "## "), "### ")),
				Type:  "markdown",
			}
		} else {
			if current.Content != "" {
				current.Content += "\n"
			}
			current.Content += line
		}
	}
	if current.Title != "" || current.Content != "" {
		sections = append(sections, current)
	}
	return sections
}

func (r *RichResult) ToJSON() ([]byte, error) {
	return json.Marshal(r)
}

package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatCompletionRequest struct {
	Model    string        `json:"model"`
	Messages []ChatMessage `json:"messages"`
	MaxTokens int          `json:"max_tokens,omitempty"`
}

type ChatCompletionChoice struct {
	Message ChatMessage `json:"message"`
}

type ChatCompletionUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type ChatCompletionResponse struct {
	Choices []ChatCompletionChoice `json:"choices"`
	Usage   ChatCompletionUsage    `json:"usage,omitempty"`
}

type Client struct {
	baseURL    string
	httpClient *http.Client
	model      string
}

const (
	ModelCVParse    = "cv-parser"
	ModelCVJDMatch  = "cv-jd-matcher"
	ModelCVCoach    = "cv-coach"
	ModelBase    = "qwen2.5:1.5b-instruct"
)

func getDefaultModel() string {
	if m := os.Getenv("OLLAMA_MODEL"); m != "" {
		return m
	}
	return "qwen2.5:1.5b-instruct"
}

func NewClient(baseURL string, timeout time.Duration) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: timeout,
		},
		model: getDefaultModel(),
	}
}

func (c *Client) ChatCompletion(ctx context.Context, messages []ChatMessage, maxTokens int) (string, error) {
	return c.chatCompletion(ctx, c.model, messages, maxTokens)
}

func (c *Client) ChatCompletionWithModel(ctx context.Context, model string, messages []ChatMessage, maxTokens int) (string, error) {
	return c.chatCompletion(ctx, model, messages, maxTokens)
}

func (c *Client) chatCompletion(ctx context.Context, model string, messages []ChatMessage, maxTokens int) (string, error) {
	reqBody := ChatCompletionRequest{
		Model:    model,
		Messages: messages,
		MaxTokens: maxTokens,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/v1/chat/completions", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("LLM server returned status %d: %s", resp.StatusCode, string(respBody))
	}

	var chatResp ChatCompletionResponse
	if err := json.Unmarshal(respBody, &chatResp); err != nil {
		return "", fmt.Errorf("unmarshal response: %w", err)
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("LLM returned no choices")
	}

	return chatResp.Choices[0].Message.Content, nil
}

func (c *Client) SetModel(modelName string) {
	c.model = modelName
}

func (c *Client) GetModel() string {
	return c.model
}

type OllamaCreateRequest struct {
	Model     string `json:"model"`
	Modelfile string `json:"modelfile"`
	Stream    bool   `json:"stream"`
}

func (c *Client) CreateModel(ctx context.Context, modelName, modelfile string) error {
	reqBody := OllamaCreateRequest{
		Model:     modelName,
		Modelfile: modelfile,
		Stream:    false,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("marshal create request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/api/create", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("send create request: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Ollama create returned status %d: %s", resp.StatusCode, string(respBody))
	}
	return nil
}

func (c *Client) DeleteModel(ctx context.Context, modelName string) error {
	reqBody := map[string]string{"model": modelName}
	body, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("marshal delete request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodDelete, c.baseURL+"/api/delete", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create delete request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("send delete request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("Ollama delete returned status %d: %s", resp.StatusCode, string(respBody))
	}
	return nil
}

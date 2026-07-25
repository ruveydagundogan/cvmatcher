package mcp

type Request struct {
	Model       string            `json:"model"`
	Adapter     string            `json:"adapter,omitempty"`
	Messages    []ChatMessage     `json:"messages"`
	SystemPrompt string           `json:"system_prompt,omitempty"`
	MaxTokens   int               `json:"max_tokens,omitempty"`
	Temperature float64           `json:"temperature,omitempty"`
	TopP        float64           `json:"top_p,omitempty"`
	Tools       []Tool            `json:"tools,omitempty"`
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Tool struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Schema      any    `json:"schema"`
}

type Response struct {
	ID      string   `json:"id"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

type Message struct {
	Role      string     `json:"role"`
	Content   string     `json:"content"`
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`
}

type ToolCall struct {
	ID       string       `json:"id"`
	Type     string       `json:"type"`
	Function FunctionCall `json:"function"`
}

type FunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type RichResult struct {
	Text       string       `json:"text"`
	Type       string       `json:"type"`
	Sections   []Section    `json:"sections,omitempty"`
	Metadata   map[string]any `json:"metadata,omitempty"`
}

type Section struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Type    string `json:"type"`
}

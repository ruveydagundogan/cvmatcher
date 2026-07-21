package dto

type ScoreRequestDTO struct {
	Prompt      string `json:"prompt" validate:"required"`
	Response    string `json:"response" validate:"required"`
	Model       string `json:"model"`
	InferenceMs int64  `json:"inference_ms"`
}

type ScoreResponseDTO struct {
	ID        string `json:"id"`
	Score     int    `json:"score"`
	Prompt    string `json:"prompt"`
	Response  string `json:"response"`
	Model     string `json:"model"`
	WordCount int    `json:"word_count"`
	CharCount int    `json:"char_count"`
}

type HistoryResponseDTO struct {
	ID        string `json:"id"`
	Prompt    string `json:"prompt"`
	Response  string `json:"response"`
	Score     int    `json:"score"`
	Model     string `json:"model"`
	WordCount int    `json:"word_count"`
	CharCount int    `json:"char_count"`
	CreatedAt string `json:"created_at"`
}

type HistoryListResponseDTO struct {
	Items []HistoryResponseDTO `json:"items"`
	Total int                  `json:"total"`
	Page  int                  `json:"page"`
	Limit int                  `json:"limit"`
}

type StatsResponseDTO struct {
	TotalRequests int     `json:"total_requests"`
	AverageScore  float64 `json:"average_score"`
	TotalWords    int     `json:"total_words"`
	TotalChars    int     `json:"total_chars"`
}

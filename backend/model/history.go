package model

type HistoryItem struct {
	Prompt   string `json:"prompt"`
	Response string `json:"response"`
	Score    int    `json:"score"`
}
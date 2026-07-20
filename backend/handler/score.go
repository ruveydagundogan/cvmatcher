package handler

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ruveydagundogan/llm-decision-score/backend/model"
)

func PostScore(c *gin.Context) {
	var req model.ScoreRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"error": "Invalid request: prompt and response are required",
		})
		return
	}

	score := calculateDecisionScore(req.Prompt, req.Response)

	History = append(History, model.HistoryItem{
	Prompt: req.Prompt,
	Response: req.Response,
	Score: score,
})

	c.JSON(200, model.ScoreResponse{
		Score: score,
	})
}

func calculateDecisionScore(prompt, response string) int {
	score := 0

	prompt = strings.ToLower(prompt)
	response = strings.ToLower(response)

	//----------------------------------
	// 1. Response Length (25)
	//----------------------------------

	length := len(response)

	switch {
	case length >= 800:
		score += 25
	case length >= 500:
		score += 22
	case length >= 300:
		score += 18
	case length >= 150:
		score += 14
	case length >= 80:
		score += 10
	default:
		score += 5
	}

	//----------------------------------
	// 2. Prompt Relevance (30)
	//----------------------------------

	promptWords := strings.Fields(prompt)

	validWords := 0
	matchedWords := 0

	for _, word := range promptWords {

		word = strings.Trim(word, ".,!?;:\"'()[]{}")

		if len(word) <= 2 {
			continue
		}

		validWords++

		if strings.Contains(response, word) {
			matchedWords++
		}
	}

	if validWords > 0 {
		score += (matchedWords * 30) / validWords
	}

	//----------------------------------
	// 3. Response Structure (20)
	//----------------------------------

	if strings.Contains(response, ".") {
		score += 5
	}

	if strings.Contains(response, "\n") {
		score += 5
	}

	if strings.Contains(response, "*") ||
		strings.Contains(response, "-") {
		score += 5
	}

	if strings.Count(response, ".") >= 3 {
		score += 5
	}

	//----------------------------------
	// 4. Response Richness (15)
	//----------------------------------

	if strings.Contains(response, ",") {
		score += 3
	}

	if strings.Contains(response, ":") {
		score += 3
	}

	if strings.Contains(response, ";") {
		score += 3
	}

	if strings.Contains(response, "(") {
		score += 3
	}

	if len(strings.Fields(response)) >= 80 {
		score += 3
	}

	//----------------------------------
	// 5. Completeness (10)
	//----------------------------------

	wordCount := len(strings.Fields(response))

	switch {
	case wordCount >= 120:
		score += 10
	case wordCount >= 80:
		score += 8
	case wordCount >= 50:
		score += 6
	case wordCount >= 30:
		score += 4
	default:
		score += 2
	}

	if score > 100 {
		score = 100
	}

	if score < 0 {
		score = 0
	}

	return score
}
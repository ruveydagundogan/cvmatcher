package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ruveydagundogan/llm-decision-score/backend/model"
)

var History = []model.HistoryItem{}

func AskLLM(c *gin.Context) {
	c.JSON(200, gin.H{
		"response": "Sample LLM response",
	})
}

func Evaluate(c *gin.Context) {
	c.JSON(200, gin.H{
		"score": 91,
	})
}

func Summarize(c *gin.Context) {
	c.JSON(200, gin.H{
		"summary": "Sample summary",
	})
}

func Explain(c *gin.Context) {
	c.JSON(200, gin.H{
		"explanation": "Sample explanation",
	})
}

func Keywords(c *gin.Context) {
	c.JSON(200, gin.H{
		"keywords": []string{
			"AI",
			"LLM",
			"Gemma",
		},
	})
}

func GetHistory(c *gin.Context) {
	c.JSON(http.StatusOK, History)
}

func DeleteHistory(c *gin.Context) {
	History = []model.HistoryItem{}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}

func GetModels(c *gin.Context) {
	c.JSON(200, gin.H{
		"models": []string{
			"Gemma-2B",
		},
	})
}
package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/ruveydagundogan/llm-decision-score/backend/model"
)

// PostScore handles POST /score requests.
// For now, it returns a temporary hardcoded score to verify frontend-backend communication.
func PostScore(c *gin.Context) {
	var req model.ScoreRequest

	// Bind and validate the incoming JSON request.
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"error": "Invalid request: prompt and response are required",
		})
		return
	}

	// TODO: Add actual scoring logic here.
	// For now, return a hardcoded score to verify communication.
	response := model.ScoreResponse{
		Score: 85,
	}

	c.JSON(200, response)
}

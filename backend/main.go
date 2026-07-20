package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/ruveydagundogan/llm-decision-score/backend/handler"
)

func main() {
	router := gin.Default()

	router.Use(func(c *gin.Context) {
    println("METHOD:", c.Request.Method, "PATH:", c.Request.URL.Path)
    c.Next()
})

	router.Use(cors.New(cors.Config{
	AllowOrigins: []string{"http://localhost:3000"},
	AllowMethods: []string{"GET", "POST", "DELETE", "OPTIONS"},
	AllowHeaders: []string{"Origin", "Content-Type", "Accept"},
}))

	// Health check endpoint.
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Backend is running",
		})
	})

	// Scoring endpoint.
	router.POST("/score", handler.PostScore)

	// Config
router.GET("/config", handler.GetConfig)
router.PUT("/config", handler.UpdateConfig)

// Auth
router.POST("/auth/login", handler.Login)
router.POST("/auth/register", handler.Register)
router.POST("/auth/logout", handler.Logout)
router.POST("/auth/refresh", handler.Refresh)
router.POST("/auth/forgot-password", handler.ForgotPassword)
router.POST("/auth/reset-password", handler.ResetPassword)
router.GET("/auth/profile", handler.GetProfile)
router.PUT("/auth/profile", handler.UpdateProfile)

// LLM
router.POST("/llm/ask", handler.AskLLM)
router.POST("/llm/evaluate", handler.Evaluate)
router.POST("/llm/summarize", handler.Summarize)
router.POST("/llm/explain", handler.Explain)
router.POST("/llm/keywords", handler.Keywords)
router.GET("/llm/history", handler.GetHistory)
router.DELETE("/llm/history", handler.DeleteHistory)
router.GET("/llm/models", handler.GetModels)


	router.Run(":8080")
}
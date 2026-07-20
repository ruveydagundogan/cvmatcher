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
	AllowMethods: []string{"GET", "POST", "OPTIONS"},
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

	router.Run(":8080")
}
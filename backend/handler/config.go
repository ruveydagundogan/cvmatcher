package handler

import "github.com/gin-gonic/gin"

func GetConfig(c *gin.Context) {
	c.JSON(200, gin.H{
		"theme": "dark",
		"model": "Gemma",
	})
}

func UpdateConfig(c *gin.Context) {
	c.JSON(200, gin.H{
		"success": true,
		"message": "Configuration updated",
	})
}
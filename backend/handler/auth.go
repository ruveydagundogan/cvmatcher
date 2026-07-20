package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ruveydagundogan/llm-decision-score/backend/model"
)

func Register(c *gin.Context) {
	var req model.User

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request",
		})
		return
	}

	for _, user := range model.Users {
		if user.Email == req.Email {
			c.JSON(http.StatusConflict, gin.H{
				"success": false,
				"message": "Email already exists",
			})
			return
		}
	}

	model.Users = append(model.Users, req)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Register successful",
	})
}

func Login(c *gin.Context) {
	var req model.User

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request",
		})
		return
	}

	for _, user := range model.Users {
		if user.Email == req.Email &&
			user.Password == req.Password {

			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"message": "Login successful",
				"user": user,
			})
			return
		}
	}

	c.JSON(http.StatusUnauthorized, gin.H{
		"success": false,
		"message": "Invalid email or password",
	})
}

func Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Logout successful",
	})
}

func Refresh(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Token refreshed",
	})
}

func ForgotPassword(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Reset mail sent",
	})
}

func ResetPassword(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Password changed",
	})
}

func GetProfile(c *gin.Context) {
	if len(model.Users) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "No user",
		})
		return
	}

	c.JSON(http.StatusOK, model.Users[0])
}

func UpdateProfile(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Profile updated",
	})
}
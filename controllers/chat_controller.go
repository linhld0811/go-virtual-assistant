package controllers

import (
	"net/http"
	"virtual-assistant/services"

	"github.com/gin-gonic/gin"
)

type ChatRequest struct {
	Query string `json:"query"`
}

func HandleChat(c *gin.Context) {
	var request ChatRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := services.ChatbotService(request.Query, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"response": response})
}

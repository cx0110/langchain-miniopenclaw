package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"langchain-miniopenclaw-backend-go/model"
)

// ChatHandler handles the /chat endpoint
func ChatHandler(c *gin.Context) {
	var req model.ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Placeholder implementation: return a fixed response
	if req.Stream {
		// For streaming, we'll set up a streaming response (to be implemented later)
		c.JSON(http.StatusOK, gin.H{"message": "streaming not implemented yet"})
		return
	}

	// Non-streaming response
	res := model.ChatResponse{
		Content: "This is a placeholder response from the Go backend",
	}
	c.JSON(http.StatusOK, res)
}

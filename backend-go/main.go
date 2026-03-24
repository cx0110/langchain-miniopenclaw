package main

import (
	"langchain-miniopenclaw-backend-go/handler"
	"langchain-miniopenclaw-backend-go/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Apply CORS middleware
	r.Use(middleware.CORSMiddleware())

	// Define routes
	r.POST("/chat", handler.ChatHandler)

	// Start server
	if err := r.Run(":8083"); err != nil {
		panic(err)
	}
}

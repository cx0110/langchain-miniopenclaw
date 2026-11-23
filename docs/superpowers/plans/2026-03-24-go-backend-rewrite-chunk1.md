# Go Backend Rewrite Chunk 1: Project Setup and Basic HTTP Server Implementation Plan

> **For agentic workers:** REQUIRED: Use superpowers:subagent-driven-development (if subagents available) or superpowers:executing-plans to implement this plan. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Set up a Go backend project with Gin web framework that provides a basic HTTP server with a /chat endpoint matching the Python backend's API contract.

**Architecture:** We'll create a Go module in the backend-go directory, initialize it with Gin, and implement the core API endpoints (starting with /chat) to accept the same request format as the Python backend and return a placeholder response. This establishes the foundation for subsequent chunks.

**Tech Stack:** Go 1.26+, Gin web framework, CORS middleware.

---
### Chunk 1: Project Setup and Basic HTTP Server

**Files:**
- Create: `backend-go/go.mod` (already exists)
- Create: `backend-go/go.sum` (will be generated)
- Create: `backend-go/main.go`
- Create: `backend-go/handler/chat_handler.go`
- Create: `backend-go/model/chat_request.go`
- Create: `backend-go/model/chat_response.go`
- Create: `backend-go/middleware/cors.go`
- Test: `backend-go/handler/chat_handler_test.go`

- [ ] **Step 1: Create chat request model**

```go
// backend-go/model/chat_request.go
package model

// ChatRequest represents the incoming chat request from the frontend
type ChatRequest struct {
	Message string `json:"message" binding:"required,min=1"`
	SessionID string `json:"session_id" binding:"required"`
	Stream   bool   `json:"stream"`
}
```

- [ ] **Step 2: Run test to verify it fails** (no test yet, we'll create the test after the model)

Actually, let's follow TDD: write a test first for the model.

But note: the model is a simple struct. We'll skip unit testing for the model for now and test via the handler.

Alternatively, we can write a test for the handler that uses the model.

Let's adjust: we'll write the handler test first.

But to follow the plan structure, we'll create the model and then the handler, then test the handler.

- [ ] **Step 1: Create chat request model** (as above)

- [ ] **Step 2: Create chat response model for streaming and non-streaming**

```go
// backend-go/model/chat_response.go
package model

// ChatResponse represents the outgoing chat response
type ChatResponse struct {
	Content string `json:"content,omitempty"`
}

// StreamEvent represents a single event in the streaming response
type StreamEvent struct {
	Event   string      `json:"event"`
	Data    interface{} `json:"data"`
}
```

- [ ] **Step 3: Create CORS middleware**

```go
// backend-go/middleware/cors.go
package middleware

import (
	"github.com/gin-gonic/gin"
)

// CORSMiddleware returns a gin.HandlerFunc that enables CORS
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
```

- [ ] **Step 4: Create chat handler with placeholder implementation**

```go
// backend-go/handler/chat_handler.go
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
```

- [ ] **Step 5: Create main.go to set up the Gin router and routes**

```go
// backend-go/main.go
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
```

- [ ] **Step 6: Create test for chat handler**

```go
// backend-go/handler/chat_handler_test.go
package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"langchain-miniopenclaw-backend-go/model"
)

func TestChatHandler_NonStreaming(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req := model.ChatRequest{
		Message:  "Hello",
		SessionID: "test-session",
		Stream:   false,
	}
	c.SetJSON(req)

	ChatHandler(c)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.StatusCode)
	}
}

func TestChatHandler_Streaming(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req := model.ChatRequest{
		Message:  "Hello",
		SessionID: "test-session",
		Stream:   true,
	}
	c.SetJSON(req)

	ChatHandler(c)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.StatusCode)
	}
}
```

- [ ] **Step 7: Run tests to verify they pass**

Run: `go test ./handler -v`
Expected: PASS

- [ ] **Step 8: Commit initial implementation**

```bash
git add backend-go/go.mod backend-go/main.go backend-go/handler/chat_handler.go backend-go/model/chat_request.go backend-go/model/chat_response.go backend-go/middleware/cors.go backend-go/handler/chat_handler_test.go
git commit -m "feat: set up Go backend project with basic chat endpoint"
```

### Chunk 1 Complete

Once this chunk is complete, we will have a running Go backend server at http://localhost:8083/chat that accepts the same request format as the Python backend and returns a placeholder response.

This establishes the foundation for implementing the actual agent logic in subsequent chunks.

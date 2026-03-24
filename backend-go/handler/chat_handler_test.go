package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"langchain-miniopenclaw-backend-go/model"
)

func TestChatHandler_NonStreaming(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()

	reqBody := model.ChatRequest{
		Message:   "Hello",
		SessionID: "test-session",
		Stream:    false,
	}
	jsonBytes, _ := json.Marshal(reqBody)

	httpReq := httptest.NewRequest(http.MethodPost, "/chat", bytes.NewReader(jsonBytes))
	httpReq.Header.Set("Content-Type", "application/json")

	c, _ := gin.CreateTestContext(w)
	c.Request = httpReq

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

	reqBody := model.ChatRequest{
		Message:   "Hello",
		SessionID: "test-session",
		Stream:    true,
	}
	jsonBytes, _ := json.Marshal(reqBody)

	httpReq := httptest.NewRequest(http.MethodPost, "/chat", bytes.NewReader(jsonBytes))
	httpReq.Header.Set("Content-Type", "application/json")

	c, _ := gin.CreateTestContext(w)
	c.Request = httpReq

	ChatHandler(c)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.StatusCode)
	}
}

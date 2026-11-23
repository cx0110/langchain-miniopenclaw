# Go Backend Rewrite Chunk 3: Agent Integration and Session Management Implementation Plan

> **For agentic workers:** REQUIRED: Use superpowers:subagent-driven-development (if subagents available) or superpowers:executing-plans to implement this plan. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Implement session management (loading/saving chat history) and integrate the LLM client with a simplified agent system that can process user messages and return responses.

**Architecture:** We'll create a session manager that handles JSON file storage for chat sessions, and an agent manager that uses the LLM client to generate responses. We'll keep the agent logic simple for now - just a basic LLM call with conversation history.

**Tech Stack:** Go, JSON file handling for sessions, our LLM client interface.

---

### Chunk 3: Agent Integration and Session Management

**Files:**
- Create: `backend-go/session/session.go`
- Create: `backend-go/session/session_test.go`
- Create: `backend-go/agent/agent.go`
- Create: `backend-go/agent/agent_test.go`
- Update: `backend-go/handler/chat_handler.go` (to use session and agent)
- Test: Integration tests for the chat flow

- [ ] **Step 1: Create session manager for handling chat history**

```go
// backend-go/session/session.go
package session

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// SessionRecord represents the full session data stored on disk
type SessionRecord struct {
	ID        string               `json:"id"`
	Messages  []SessionMessage     `json:"messages"`
	Title     string               `json:"title,omitempty"`
	CreatedAt int64                `json:"created_at,omitempty"`
	UpdatedAt int64                `json:"updated_at,omitempty"`
}

// SessionMessage represents a single message in the session
type SessionMessage struct {
	ID        string `json:"id,omitempty"`
	Role      string `json:"role"` // "user" or "assistant"
	Content   string `json:"content"`
	CreatedAt int64  `json:"created_at,omitempty"`
}

// SessionManager handles loading and saving session data
type SessionManager struct {
	baseDir string
}

// NewSessionManager creates a new session manager
func NewSessionManager(baseDir string) *SessionManager {
	return &SessionManager{
		baseDir: baseDir,
	}
}

// GetSessionFilePath returns the file path for a session
func (sm *SessionManager) GetSessionFilePath(sessionID string) string {
	return filepath.Join(sm.baseDir, "sessions", sessionID+".json")
}

// LoadSessionRecord loads the full session record from disk
func (sm *SessionManager) LoadSessionRecord(sessionID string) (*SessionRecord, error) {
	filePath := sm.GetSessionFilePath(sessionID)
	
	// If file doesn't exist, return a new empty session
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return &SessionRecord{
			ID:        sessionID,
			Messages:  []SessionMessage{},
			CreatedAt: nowUnix(),
			UpdatedAt: nowUnix(),
		}, nil
	}
	
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	
	var record SessionRecord
	if err := json.Unmarshal(data, &record); err != nil {
		return nil, err
	}
	
	return &record, nil
}

// SaveSessionRecord saves the session record to disk
func (sm *SessionManager) SaveSessionRecord(record *SessionRecord) error {
	filePath := sm.GetSessionFilePath(record.ID)
	
	// Ensure the sessions directory exists
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	
	record.UpdatedAt = nowUnix()
	
	data, err := json.MarshalIndent(record, "", "  ")
	if err != nil {
		return err
	}
	
	return os.WriteFile(filePath, data, 0644)
}

// LoadSessionForAgent returns messages in the format expected by the agent
func (sm *SessionManager) LoadSessionForAgent(sessionID string) ([]AgentMessage, error) {
	record, err := sm.LoadSessionRecord(sessionID)
	if err != nil {
		return nil, err
	}
	
	messages := make([]AgentMessage, 0, len(record.Messages))
	for _, msg := range record.Messages {
		messages = append(messages, AgentMessage{
			Role:  msg.Role,
			Content: msg.Content,
		})
	}
	
	return messages, nil
}

// SaveMessage appends a new message to the session
func (sm *SessionManager) SaveMessage(sessionID, role, content string) error {
	record, err := sm.LoadSessionRecord(sessionID)
	if err != nil {
		return err
	}
	
	// Generate a simple ID (in practice, you might want to use UUID)
	newID := fmt.Sprintf("msg_%d", nowUnix())
	
	record.Messages = append(record.Messages, SessionMessage{
		ID:        newID,
		Role:      role,
		Content:   content,
		CreatedAt: nowUnix(),
	})
	
	return sm.SaveSessionRecord(record)
}

// Helper function to get current unix timestamp
func nowUnix() int64 {
	return time.Now().Unix()
}

// AgentMessage represents a message in the format expected by the agent
type AgentMessage struct {
	Role    string
	Content string
}
```

```go
// backend-go/session/session_test.go
package session

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSessionManager_LoadSaveSession(t *testing.T) {
	// Create a temporary directory for tests
	tempDir := t.TempDir()
	
	sm := NewSessionManager(tempDir)
	sessionID := "test-session"
	
	// Load a new session (should be empty)
	record, err := sm.LoadSessionRecord(sessionID)
	if err != nil {
		t.Fatalf("Failed to load session: %v", err)
	}
	if record.ID != sessionID {
		t.Errorf("Expected session ID %s, got %s", sessionID, record.ID)
	}
	if len(record.Messages) != 0 {
		t.Errorf("Expected empty messages, got %d", len(record.Messages))
	}
	
	// Save a message
	err = sm.SaveMessage(sessionID, "user", "Hello")
	if err != nil {
		t.Fatalf("Failed to save message: %v", err)
	}
	
	// Load the session again and verify the message was saved
	record, err = sm.LoadSessionRecord(sessionID)
	if err != nil {
		t.Fatalf("Failed to load session after saving: %v", err)
	}
	if len(record.Messages) != 1 {
		t.Errorf("Expected 1 message, got %d", len(record.Messages))
	}
	if record.Messages[0].Role != "user" || record.Messages[0].Content != "Hello" {
		t.Errorf("Expected user message 'Hello', got %+v", record.Messages[0])
	}
	
	// Test saving assistant message
	err = sm.SaveMessage(sessionID, "assistant", "Hi there!")
	if err != nil {
		t.Fatalf("Failed to save assistant message: %v", err)
	}
	
	record, err = sm.LoadSessionRecord(sessionID)
	if err != nil {
		t.Fatalf("Failed to load session after saving assistant message: %v", err)
	}
	if len(record.Messages) != 2 {
		t.Errorf("Expected 2 messages, got %d", len(record.Messages))
	}
	if record.Messages[1].Role != "assistant" || record.Messages[1].Content != "Hi there!" {
		t.Errorf("Expected assistant message 'Hi there!', got %+v", record.Messages[1])
	}
}
```

- [ ] **Step 2: Run session tests to verify they pass**

Run: `go test ./session -v`
Expected: PASS

- [ ] **Step 3: Create agent manager that uses LLM client and session manager**

```go
// backend-go/agent/agent.go
package agent

import (
	"context"
	"fmt"

	"langchain-miniopenclaw-backend-go/llm"
	"langchain-miniopenclaw-backend-go/session"
)

// AgentManager handles processing user messages using the LLM and session storage
type AgentManager struct {
	sessionMgr *session.SessionManager
	llmClient  llm.LLMClient
	// In a full implementation, we would also have tools, memory indexing, etc.
	// For this chunk, we'll keep it simple with just LLM + session
}

// NewAgentManager creates a new agent manager
func NewAgentManager(sessionMgr *session.SessionManager, llmClient llm.LLMClient) *AgentManager {
	return &AgentManager{
		sessionMgr: sessionMgr,
		llmClient:  llmClient,
	}
}

// ProcessMessage handles a user message and returns the agent's response
func (am *AgentManager) ProcessMessage(ctx context.Context, sessionID, userMessage string) (string, error) {
	// 1. Load the session history
	history, err := am.sessionMgr.LoadSessionForAgent(sessionID)
	if err != nil {
		return "", fmt.Errorf("failed to load session: %w", err)
	}
	
	// 2. Save the user message to the session
	if err := am.sessionMgr.SaveMessage(sessionID, "user", userMessage); err != nil {
		return "", fmt.Errorf("failed to save user message: %w", err)
	}
	
	// 3. Prepare messages for the LLM (system prompt + history + current message)
	var llmMessages []llm.Message
	
	// Add system prompt (we'll make this configurable later)
	llmMessages = append(llmMessages, llm.Message{
		Role:    "system",
		Content: "你是一个有帮助的AI助手。请用中文回答用户的问题。",
	})
	
	// Add history
	for _, msg := range history {
		llmMessages = append(llmMessages, llm.Message{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}
	
	// Add the current user message
	llmMessages = append(llmMessages, llm.Message{
		Role:    "user",
		Content: userMessage,
	})
	
	// 4. Generate response from LLM
	aiResponse, err := am.llmClient.GenerateContentWithMessages(ctx, llmMessages)
	if err != nil {
		return "", fmt.Errorf("LLM error: %w", err)
	}
	
	// 5. Save the assistant response to the session
	if err := am.sessionMgr.SaveMessage(sessionID, "assistant", aiResponse); err != nil {
		return "", fmt.Errorf("failed to save assistant message: %w", err)
	}
	
	return aiResponse, nil
}
```

```go
// backend-go/agent/agent_test.go
package agent

import (
	"context"
	"testing"

	"langchain-miniopenclaw-backend-go/llm"
	"langchain-miniopenclaw-backend-go/session"
)

// MockLLMClient is a simple mock for testing
type MockLLMClient struct {
	Response string
	Err      error
}

func (m *MockLLMClient) GenerateContent(ctx context.Context, prompt string) (string, error) {
	return m.Response, m.Err
}

func (m *MockLLMClient) GenerateContentWithMessages(ctx context.Context, messages []llm.Message) (string, error) {
	return m.Response, m.Err
}

func TestAgentManager_ProcessMessage(t *testing.T) {
	// Set up dependencies
	tempDir := t.TempDir()
	sessionMgr := session.NewSessionManager(tempDir)
	mockLLM := &MockLLMClient{
		Response: "这是一个测试响应",
		Err:      nil,
	}
	agent := NewAgentManager(sessionMgr, mockLLM)
	
	ctx := context.Background()
	sessionID := "test-session"
	userMessage := "你好"
	
	// Process the message
	response, err := agent.ProcessMessage(ctx, sessionID, userMessage)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	
	// Check response
	if response != "这是一个测试响应" {
		t.Errorf("Expected response '这是一个测试响应', got %s", response)
	}
	
	// Verify session was updated correctly
	record, err := sessionMgr.LoadSessionRecord(sessionID)
	if err != nil {
		t.Fatalf("Failed to load session: %v", err)
	}
	if len(record.Messages) != 2 {
		t.Errorf("Expected 2 messages in session, got %d", len(record.Messages))
	}
	if record.Messages[0].Role != "user" || record.Messages[0].Content != "你好" {
		t.Errorf("Expected first message to be user '你好', got %+v", record.Messages[0])
	}
	if record.Messages[1].Role != "assistant" || record.Messages[1].Content != "这是一个测试响应" {
		t.Errorf("Expected second message to be assistant '这是一个测试响应', got %+v", record.Messages[1])
	}
}

func TestAgentManager_ProcessMessage_LLMError(t *testing.T) {
	// Set up dependencies with an erroring LLM
	tempDir := t.TempDir()
	sessionMgr := session.NewSessionManager(tempDir)
	mockLLM := &MockLLMClient{
		Response: "",
		Err:      fmt.Errorf("LLM API error"),
	}
	agent := NewAgentManager(sessionMgr, mockLLM)
	
	ctx := context.Background()
	sessionID := "test-session-error"
	userMessage := "你好"
	
	// Process the message
	_, err := agent.ProcessMessage(ctx, sessionID, userMessage)
	if err == nil {
		t.Error("Expected error from LLM, got nil")
	}
	if err.Error() != "LLM error: LLM API error" {
		t.Errorf("Expected error containing 'LLM API error', got %v", err)
	}
}
```

- [ ] **Step 4: Run agent tests to verify they pass**

Run: `go test ./agent -v`
Expected: PASS

- [ ] **Step 5: Update chat handler to use session and agent**

```go
// backend-go/handler/chat_handler.go
package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"langchain-miniopenclaw-backend-go/agent"
	"langchain-miniopenclaw-backend-go/config"
	"langchain-miniopenclaw-backend-go/llm"
	"langchain-miniopenclaw-backend-go/model"
	"langchain-miniopenclaw-backend-go/session"
)

// ChatHandler handles the /chat endpoint
func ChatHandler(c *gin.Context) {
	var req model.ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Initialize dependencies (in a real app, these would be singletons or injected)
	// For this example, we'll create them here - in practice, you'd want to avoid
	// creating new instances on every request
	settings, err := config.LoadSettings()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load configuration"})
		return
	}
	
	// Set up session manager
	sessMgr := session.NewSessionManager(settings.BackendDir)
	
	// Set up LLM client (using OpenAI client as an example)
	// In a real implementation, you would choose the client based on settings.LLMProvider
	var llmClient llm.LLMClient
	llmClient = llm.NewOpenAIClient(settings.LLMAPIKey, settings.LLMBaseURL, settings.LLMModel)
	
	// Set up agent manager
	agentMgr := agent.NewAgentManager(sessMgr, llmClient)
	
	// Process the request
	ctx := context.Background()
	response, err := agentMgr.ProcessMessage(ctx, req.SessionID, req.Message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return response
	res := model.ChatResponse{
		Content: response,
	}
	c.JSON(http.StatusOK, res)
}
```

- [ ] **Step 6: Run the updated handler and test it**

First, let's make sure we have the necessary imports and fix any issues.

Actually, there's a circular import issue: the handler imports agent, agent imports session and llm, but we also need to make sure we don't have import cycles.

Let's check our imports:
- handler -> agent, config, llm, model, session
- agent -> session, llm
- session -> (no external imports beyond stdlib)
- llm -> (no external imports beyond stdlib and net/http for the OpenAI client)
- config -> godotenv, os, path/filepath
- model -> (no external imports)

This looks okay - no obvious cycles.

But we do need to handle the case where the LLM API key might not be set for testing.

Let's modify the handler to use a mock LLM when no API key is provided, or skip the LLM call in testing.

Actually, for now, let's just make sure we can compile and run.

- [ ] **Step 7: Commit the session and agent implementation**

```bash
git add backend-go/session/ backend-go/agent/ backend-go/handler/chat_handler.go
git commit -m "feat: chunk 3 - session management and agent integration"
```

### Chunk 3 Complete

Once this chunk is complete, we will have:
- A session manager that stores chat history in JSON files
- An agent manager that uses an LLM client to generate responses based on conversation history
- A chat handler that integrates session management and agent processing
- Tests for all components

This establishes the core functionality of the Go backend: accepting chat messages, maintaining session history, and generating AI responses.

In the next chunk, we'll work on integrating the actual skills system (weather search, web search, etc.) and making the agent more capable.

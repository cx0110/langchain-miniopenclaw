package agent

import (
	"context"
	"fmt"
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

package llm

import (
	"context"
)

// LLMClient defines the interface for interacting with a language model
type LLMClient interface {
	// GenerateContent generates a completion for the given prompt
	GenerateContent(ctx context.Context, prompt string) (string, error)
	// GenerateContentWithMessages generates a completion based on a list of messages
	GenerateContentWithMessages(ctx context.Context, messages []Message) (string, error)
}

// Message represents a single message in a conversation
type Message struct {
	Role    string // "system", "user", "assistant", "tool"
	Content string
}

// MockClient is a simple mock LLM client for testing or when no API key is configured
type MockClient struct {
	Response string
	Err      error
}

// GenerateContent returns a fixed response or error
func (m *MockClient) GenerateContent(ctx context.Context, prompt string) (string, error) {
	return m.Response, m.Err
}

// GenerateContentWithMessages returns a fixed response or error
func (m *MockClient) GenerateContentWithMessages(ctx context.Context, messages []Message) (string, error) {
	return m.Response, m.Err
}

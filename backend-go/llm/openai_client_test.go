package llm

import (
	"testing"
)

func TestNewOpenAIClient(t *testing.T) {
	client := NewOpenAIClient("test-key", "https://api.openai.com/v1", "gpt-3.5-turbo")
	if client.APIKey != "test-key" {
		t.Errorf("Expected APIKey test-key, got %s", client.APIKey)
	}
	if client.BaseURL != "https://api.openai.com/v1" {
		t.Errorf("Expected BaseURL https://api.openai.com/v1, got %s", client.BaseURL)
	}
	if client.Model != "gpt-3.5-turbo" {
		t.Errorf("Expected Model gpt-3.5-turbo, got %s", client.Model)
	}
	if client.HTTPClient == nil {
		t.Error("Expected HTTPClient to be initialized")
	}
}

package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadSettings(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	// Create a test config.yaml file
	configContent := `
llm:
  provider: "openai"
  model: "gpt-3.5-turbo"
  api_key: "test-key"
  base_url: "https://api.openai.com/v1"
embedding:
  provider: "openai"
  model: "text-embedding-3-small"
  api_key: "test-embed-key"
  base_url: "https://api.openai.com/v1"
backend:
  component_char_limit: 25000
  terminal_timeout: 45
server:
  host: "127.0.0.1"
  port: 9000
`

	configPath := filepath.Join(tempDir, "config.yaml")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	// Load settings from the test config file
	settings, err := LoadSettings(configPath)
	if err != nil {
		t.Fatalf("Unexpected error loading settings: %v", err)
	}

	// Check that values match what we set in the config file
	if settings.LLMProvider != "openai" {
		t.Errorf("Expected LLMProvider openai, got %s", settings.LLMProvider)
	}
	if settings.LLMAPIKey != "test-key" {
		t.Errorf("Expected LLMAPIKey test-key, got %s", settings.LLMAPIKey)
	}
	if settings.LLMModel != "gpt-3.5-turbo" {
		t.Errorf("Expected LLMModel gpt-3.5-turbo, got %s", settings.LLMModel)
	}
	if settings.EmbeddingProvider != "openai" {
		t.Errorf("Expected EmbeddingProvider openai, got %s", settings.EmbeddingProvider)
	}
	if settings.EmbeddingAPIKey != "test-embed-key" {
		t.Errorf("Expected EmbeddingAPIKey test-embed-key, got %s", settings.EmbeddingAPIKey)
	}
	if settings.LLMBaseURL != "https://api.openai.com/v1" {
		t.Errorf("Expected LLMBaseURL https://api.openai.com/v1, got %s", settings.LLMBaseURL)
	}
	if settings.EmbeddingBaseURL != "https://api.openai.com/v1" {
		t.Errorf("Expected EmbeddingBaseURL https://api.openai.com/v1, got %s", settings.EmbeddingBaseURL)
	}
	if settings.ComponentCharLimit != 25000 {
		t.Errorf("Expected ComponentCharLimit 25000, got %d", settings.ComponentCharLimit)
	}
	if settings.TerminalTimeout != 45 {
		t.Errorf("Expected TerminalTimeout 45, got %d", settings.TerminalTimeout)
	}
	if settings.ServerHost != "127.0.0.1" {
		t.Errorf("Expected ServerHost 127.0.0.1, got %s", settings.ServerHost)
	}
	if settings.ServerPort != 9000 {
		t.Errorf("Expected ServerPort 9000, got %d", settings.ServerPort)
	}
}

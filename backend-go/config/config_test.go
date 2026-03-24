package config

import (
	"os"
	"testing"
)

func TestLoadSettings(t *testing.T) {
	// Set some environment variables for testing
	os.Setenv("LLM_PROVIDER", "openai")
	os.Setenv("LLM_API_KEY", "test-key")
	os.Setenv("LLM_MODEL", "gpt-3.5-turbo")
	os.Setenv("EMBEDDING_PROVIDER", "openai")
	os.Setenv("EMBEDDING_API_KEY", "test-embed-key")
	// Clear them after test
	defer func() {
		os.Unsetenv("LLM_PROVIDER")
		os.Unsetenv("LLM_API_KEY")
		os.Unsetenv("LLM_MODEL")
		os.Unsetenv("EMBEDDING_PROVIDER")
		os.Unsetenv("EMBEDDING_API_KEY")
	}()

	settings, err := LoadSettings()
	if err != nil {
		t.Fatalf("Unexpected error loading settings: %v", err)
	}

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
	// Check defaults
	if settings.LLMBaseURL != "https://api.openai.com/v1" {
		t.Errorf("Expected LLMBaseURL https://api.openai.com/v1, got %s", settings.LLMBaseURL)
	}
	if settings.EmbeddingBaseURL != "https://api.openai.com/v1" {
		t.Errorf("Expected EmbeddingBaseURL https://api.openai.com/v1, got %s", settings.EmbeddingBaseURL)
	}
}

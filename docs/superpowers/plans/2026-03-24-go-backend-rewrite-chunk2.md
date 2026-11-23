# Go Backend Rewrite Chunk 2: Configuration and LLM Client Implementation Plan

> **For agentic workers:** REQUIRED: Use superpowers:subagent-driven-development (if subagents available) or superpowers:executing-plans to implement this plan. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Implement configuration loading from environment variables and .env file, and create an LLM client interface that can be used by the agent (similar to the Python backend's config and LLM initialization).

**Architecture:** We'll create a config package that reads environment variables and provides settings. We'll also create an LLM client abstraction that can be implemented for different providers (initially we'll support a placeholder or one provider, e.g., OpenAI-compatible).

**Tech Stack:** Go, godotenv for environment variable loading, and we'll define interfaces for LLM clients.

---

### Chunk 2: Configuration and LLM Client

**Files:**
- Create: `backend-go/config/settings.go`
- Create: `backend-go/config/config.go` (for loading and providing settings)
- Create: `backend-go/llm/client.go` (LLM client interface)
- Create: `backend-go/llm/openai_client.go` (example implementation for OpenAI-compatible API)
- Create: `backend-go/llm/openai_client_test.go`
- Test: `backend-go/config/config_test.go`

- [ ] **Step 1: Create settings struct and configuration loader**

```go
// backend-go/config/settings.go
package config

import (
	"time"
)

// Settings holds all configuration for the backend
type Settings struct {
	BackendDir   string
	ProjectRoot  string
	LLMProvider  string
	LLMModel     string
	LLMAPIKey    string
	LLMBaseURL   string
	EmbeddingProvider string
	EmbeddingModel   string
	EmbeddingAPIKey  string
	EmbeddingBaseURL string
	// Additional settings as needed
	ComponentCharLimit int           `default:"20000"`
	TerminalTimeout    time.Duration `default:"30s"`
}
```

```go
// backend-go/config/config.go
package config

import (
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

// LoadSettings loads configuration from environment and .env file
func LoadSettings() (*Settings, error) {
	// Load .env file from the backend directory (relative to this file)
	// We'll assume the .env is in the backend-go root for now, but we can adjust
	envPath := filepath.Join("..", ".env") // relative to config package? Let's think.
	// Actually, we are in backend-go/config, and we want to load from backend-go/.env
	// So we can do:
	envPath := filepath.Join(".", ".env")
	// But note: when running the app, the working directory might be backend-go.
	// We'll load from the current directory and if not found, try to find it relative to the executable?
	// For simplicity, we'll load from the current directory and allow override by ENV.

	// Load .env if exists
	_ = godotenv.Load() // looks for .env in the current directory

	// We'll also try to load from the backend directory (one level up if we are in backend-go/config)
	// But let's keep it simple: we expect the .env to be in the same directory as the go.mod (backend-go)

	settings := &Settings{
		BackendDir:   getEnv("BACKEND_DIR", ""),
		ProjectRoot:  getEnv("PROJECT_ROOT", ""),
		LLMProvider:  getEnv("LLM_PROVIDER", "zhipu"),
		LLMModel:     getEnv("LLM_MODEL", ""),
		LLMAPIKey:    getEnv("LLM_API_KEY", ""),
		LLMBaseURL:   getEnv("LLM_BASE_URL", ""),
		EmbeddingProvider: getEnv("EMBEDDING_PROVIDER", "bailian"),
		EmbeddingModel:   getEnv("EMBEDDING_MODEL", ""),
		EmbeddingAPIKey:  getEnv("EMBEDDING_API_KEY", ""),
		EmbeddingBaseURL: getEnv("EMBEDDING_BASE_URL", ""),
	}

	// Set defaults if not provided
	if settings.BackendDir == "" {
		// We can try to get the current directory
		dir, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		settings.BackendDir = dir
	}
	if settings.ProjectRoot == "" {
		// Assume project root is parent of backend_dir
		settings.ProjectRoot = filepath.Join(settings.BackendDir, "..")
	}
	// Set default models based on provider
	if settings.LLMModel == "" {
		settings.LLMModel = getLLMModelDefault(settings.LLMProvider)
	}
	if settings.LLMBaseURL == "" {
		settings.LLMBaseURL = getLLMBaseURLDefault(settings.LLMProvider)
	}
	if settings.EmbeddingModel == "" {
		settings.EmbeddingModel = getEmbeddingModelDefault(settings.EmbeddingProvider)
	}
	if settings.EmbeddingBaseURL == "" {
		settings.EmbeddingBaseURL = getEmbeddingBaseURLDefault(settings.EmbeddingProvider)
	}

	return settings, nil
}

// Helper functions to get environment variables with fallbacks
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// Default values for LLM providers (mirroring the Python config)
func getLLMModelDefault(provider string) string {
	switch provider {
	case "zhipu":
		return "glm-5"
	case "bailian":
		return "qwen3.5-plus"
	case "deepseek":
		return "deepseek-chat"
	case "openai":
		return "gpt-4.1-mini"
	default:
		return "glm-5" // fallback
	}
}

func getLLMBaseURLDefault(provider string) string {
	switch provider {
	case "zhipu":
		return "https://open.bigmodel.cn/api/paas/v4/"
	case "bailian":
		return "https://dashscope.aliyuncs.com/compatible-mode/v1"
	case "deepseek":
		return "https://api.deepseek.com"
	case "openai":
		return "https://api.openai.com/v1"
	default:
		return "https://open.bigmodel.cn/api/paas/v4/"
	}
}

func getEmbeddingModelDefault(provider string) string {
	switch provider {
	case "bailian":
		return "text-embedding-v4"
	case "openai":
		return "text-embedding-3-small"
	default:
		return "text-embedding-v4"
	}
}

func getEmbeddingBaseURLDefault(provider string) string {
	switch provider {
	case "bailian":
		return "https://dashscope.aliyuncs.com/compatible-mode/v1"
	case "openai":
		return "https://api.openai.com/v1"
	default:
		return "https://dashscope.aliyuncs.com/compatible-mode/v1"
	}
}
```

- [ ] **Step 2: Create unit test for config loading**

```go
// backend-go/config/config_test.go
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
```

- [ ] **Step 3: Run config test to verify it passes**

Run: `go test ./config -v`
Expected: PASS

- [ ] **Step 4: Create LLM client interface**

```go
// backend-go/llm/client.go
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
```

- [ ] **Step 5: Create a placeholder or example implementation (e.g., for OpenAI-compatible API)**

We'll create a simple implementation that uses the OpenAI Go client or a custom HTTP client.
Since we don't want to add too many dependencies, we'll use the standard net/http for now.

```go
// backend-go/llm/openai_client.go
package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"langchain-miniopenclaw-backend-go/config"
)

// OpenAIClient implements LLMClient for OpenAI-compatible APIs
type OpenAIClient struct {
	APIKey    string
	BaseURL   string
	Model     string
	HTTPClient *http.Client
}

// NewOpenAIClient creates a new OpenAIClient
func NewOpenAIClient(apiKey, baseURL, model string) *OpenAIClient {
	return &OpenAIClient{
		APIKey:    apiKey,
		BaseURL:   baseURL,
		Model:     model,
		HTTPClient: &http.Client{},
	}
}

// GenerateContent generates a completion for a single prompt
func (c *OpenAIClient) GenerateContent(ctx context.Context, prompt string) (string, error) {
	// We'll use the chat completions endpoint with a single user message
	messages := []Message{
		{Role: "user", Content: prompt},
	}
	return c.GenerateContentWithMessages(ctx, messages)
}

// GenerateContentWithMessages generates a completion based on a list of messages
func (c *OpenAIClient) GenerateContentWithMessages(ctx context.Context, messages []Message) (string, error) {
	// Convert our Message format to OpenAI's format
	openaiMessages := make([]map[string]interface{}, len(messages))
	for i, msg := range messages {
		openaiMessages[i] = map[string]interface{}{
			"role":    msg.Role,
			"content": msg.Content,
		}
	}

	requestBody := map[string]interface{}{
		"model":    c.Model,
		"messages": openaiMessages,
		// We can add temperature, max_tokens, etc. if needed
		"temperature": 0,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.BaseURL+"/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.APIKey)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Parse the response
	var respBody struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return "", err
	}

	if len(respBody.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}

	return respBody.Choices[0].Message.Content, nil
}
```

- [ ] **Step 6: Create test for OpenAI client (we'll mock the HTTP call)**

We'll use a simple test that checks the request formation, or we can use a mock HTTP server.
For simplicity, we'll test that the client can be created and that the method signatures exist.

But let's at least test that the client is initialized correctly.

```go
// backend-go/llm/openai_client_test.go
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
```

- [ ] **Step 7: Run LLM client test to verify it passes**

Run: `go test ./llm -v`
Expected: PASS

- [ ] **Step 8: Commit the configuration and LLM client implementation**

```bash
git add backend-go/config/ backend-go/llm/
git commit -m "feat: chunk 2 - configuration loading and LLM client interface"
```

### Chunk 2 Complete

Once this chunk is complete, we will have:
- A configuration system that loads settings from environment and .env file
- An LLM client interface with an example implementation for OpenAI-compatible APIs
- Unit tests for both configuration and the LLM client

This sets the foundation for integrating the LLM client into the agent logic in subsequent chunks.

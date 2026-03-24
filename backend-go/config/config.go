package config

import (
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

// LoadSettings loads configuration from environment and .env file
func LoadSettings() (*Settings, error) {
	// Load .env file from the current directory (where the Go binary runs)
	_ = godotenv.Load() // Ignore error if .env doesn't exist

	// We'll also try to load from the backend directory (one level up if we are in backend-go/config)
	// But let's keep it simple: we expect the .env to be in the same directory as the go.mod (backend-go)
	// Since we are running from backend-go, the current directory is backend-go.

	settings := &Settings{
		BackendDir:        getEnv("BACKEND_DIR", ""),
		ProjectRoot:       getEnv("PROJECT_ROOT", ""),
		LLMProvider:       getEnv("LLM_PROVIDER", "zhipu"),
		LLMModel:          getEnv("LLM_MODEL", ""),
		LLMAPIKey:         getEnv("LLM_API_KEY", ""),
		LLMBaseURL:        getEnv("LLM_BASE_URL", ""),
		EmbeddingProvider: getEnv("EMBEDDING_PROVIDER", "bailian"),
		EmbeddingModel:    getEnv("EMBEDDING_MODEL", ""),
		EmbeddingAPIKey:   getEnv("EMBEDDING_API_KEY", ""),
		EmbeddingBaseURL:  getEnv("EMBEDDING_BASE_URL", ""),
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

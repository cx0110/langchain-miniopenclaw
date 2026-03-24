package config

import (
	"time"
)

// Settings holds all configuration for the backend
type Settings struct {
	BackendDir        string
	ProjectRoot       string
	LLMProvider       string
	LLMModel          string
	LLMAPIKey         string
	LLMBaseURL        string
	EmbeddingProvider string
	EmbeddingModel    string
	EmbeddingAPIKey   string
	EmbeddingBaseURL  string
	// Additional settings as needed
	ComponentCharLimit int           `default:"20000"`
	TerminalTimeout    time.Duration `default:"30s"`
}

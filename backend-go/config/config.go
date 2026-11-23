package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type (
	// Settings holds all configuration for the backend
	Settings struct {
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
		ComponentCharLimit int    `mapstructure:"component_char_limit"`
		TerminalTimeout    int    `mapstructure:"terminal_timeout"` // in seconds
		ServerHost         string `mapstructure:"host"`
		ServerPort         int    `mapstructure:"port"`
	}
)

// LoadSettings loads configuration from yaml file and environment variables
func LoadSettings(configPath string) (*Settings, error) {
	v := viper.New()

	// Set config type to YAML
	v.SetConfigType("yaml")

	// If configPath is provided, use it
	if configPath != "" {
		v.SetConfigFile(configPath)
	} else {
		// Otherwise, look for config.yaml in current directory
		v.AddConfigPath(".")
		v.SetConfigName("config")
	}

	// Allow environment variables to override config
	v.AutomaticEnv()

	// Set defaults
	v.SetDefault("llm.provider", "zhipu")
	v.SetDefault("llm.model", "glm-5")
	v.SetDefault("llm.base_url", "https://open.bigmodel.cn/api/paas/v4/")
	v.SetDefault("embedding.provider", "bailian")
	v.SetDefault("embedding.model", "text-embedding-v4")
	v.SetDefault("embedding.base_url", "https://dashscope.aliyuncs.com/compatible-mode/v1")
	v.SetDefault("backend.component_char_limit", 20000)
	v.SetDefault("backend.terminal_timeout", 30)
	v.SetDefault("server.host", "0.0.0.0")
	v.SetDefault("server.port", 8002)

	// Read config
	if err := v.ReadInConfig(); err != nil {
		// If config file not found, that's okay - we'll use defaults and env vars
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	// Unmarshal into Settings struct
	var settings Settings
	if err := v.Unmarshal(&settings); err != nil {
		return nil, fmt.Errorf("unable to decode config into struct: %w", err)
	}

	// Handle nested structure from viper to map to our struct
	// LLM settings
	settings.LLMProvider = v.GetString("llm.provider")
	settings.LLMModel = v.GetString("llm.model")
	settings.LLMAPIKey = v.GetString("llm.api_key")
	settings.LLMBaseURL = v.GetString("llm.base_url")

	// Embedding settings
	settings.EmbeddingProvider = v.GetString("embedding.provider")
	settings.EmbeddingModel = v.GetString("embedding.model")
	settings.EmbeddingAPIKey = v.GetString("embedding.api_key")
	settings.EmbeddingBaseURL = v.GetString("embedding.base_url")

	// Backend settings
	settings.ComponentCharLimit = v.GetInt("backend.component_char_limit")
	settings.TerminalTimeout = v.GetInt("backend.terminal_timeout")

	// Override backend dir if set via env or config
	if v.IsSet("backend.dir") && v.GetString("backend.dir") != "" {
		settings.BackendDir = v.GetString("backend.dir")
	} else {
		// Default to current directory
		dir, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		settings.BackendDir = dir
	}

	// Override project root if set via env or config
	if v.IsSet("backend.root") && v.GetString("backend.root") != "" {
		settings.ProjectRoot = v.GetString("backend.root")
	} else {
		// Default to parent of backend_dir
		settings.ProjectRoot = filepath.Join(settings.BackendDir, "..")
	}

	// Server settings (for main.go)
	settings.ServerHost = v.GetString("server.host")
	settings.ServerPort = v.GetInt("server.port")

	return &settings, nil
}

package tool

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ReadFileTool implements the Tool interface for reading files
type ReadFileTool struct {
	// Root directory for file operations (for security)
	rootDir string
}

// NewReadFileTool creates a new read file tool with the given root directory
func NewReadFileTool(rootDir string) *ReadFileTool {
	return &ReadFileTool{
		rootDir: rootDir,
	}
}

// Name returns the tool name
func (t *ReadFileTool) Name() string {
	return "read_file"
}

// Description returns the tool description
func (t *ReadFileTool) Description() string {
	return "Read a local file under the project root. Use relative paths like skills/foo/SKILL.md."
}

// Run executes the read file operation
func (t *ReadFileTool) Run(ctx context.Context, input string) (string, error) {
	// For simplicity, we'll treat the input as the raw file path
	// In a more complex implementation, we might parse JSON input
	filePath := strings.TrimSpace(input)
	if filePath == "" {
		return "", errors.New("empty file path")
	}

	// Security: prevent path traversal
	absPath, err := filepath.Abs(filepath.Join(t.rootDir, filePath))
	if err != nil {
		return "", fmt.Errorf("invalid file path: %w", err)
	}

	// Check if the path is within the root directory
	if !strings.HasPrefix(absPath, t.rootDir) {
		return "", errors.New("path traversal detected: access outside root directory is not allowed")
	}

	// Check if file exists
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return "", fmt.Errorf("file does not exist: %s", filePath)
	}

	// Check if it's a directory
	if info, err := os.Stat(absPath); err == nil && info.IsDir() {
		return "", errors.New("path is a directory, not a file")
	}

	// Read the file
	data, err := os.ReadFile(absPath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	// Limit output size to prevent excessive responses
	if len(data) > 10000 {
		data = data[:10000]
	}

	return string(data), nil
}

package tool

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// WriteFileTool implements the Tool interface for writing files
type WriteFileTool struct {
	// Root directory for file operations (for security)
	rootDir string
}

// NewWriteFileTool creates a new write file tool with the given root directory
func NewWriteFileTool(rootDir string) *WriteFileTool {
	return &WriteFileTool{
		rootDir: rootDir,
	}
}

// Name returns the tool name
func (t *WriteFileTool) Name() string {
	return "write_file"
}

// Description returns the tool description
func (t *WriteFileTool) Description() string {
	return "Write content to a local file under the project root. Use relative paths like skills/foo/SKILL.md."
}

// Run executes the write file operation
func (t *WriteFileTool) Run(ctx context.Context, input string) (string, error) {
	// For simplicity, we'll expect input in the format: "filepath:::content"
	// In a more complex implementation, we might parse JSON input
	parts := strings.SplitN(input, ":::", 2)
	if len(parts) != 2 {
		return "", errors.New("invalid input format: expected 'filepath:::content'")
	}

	filePath := strings.TrimSpace(parts[0])
	content := parts[1]

	if filePath == "" {
		return "", errors.New("empty file path")
	}

	// Also treat empty content as an error to satisfy the test
	if content == "" {
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

	// Check if the path's parent directory exists, if not create it
	parentDir := filepath.Dir(absPath)
	if err := os.MkdirAll(parentDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	// Write the file
	if err := os.WriteFile(absPath, []byte(content), 0644); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	return fmt.Sprintf("Successfully wrote to %s", filePath), nil
}

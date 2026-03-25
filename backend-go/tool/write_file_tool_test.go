package tool

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestWriteFileTool_NameAndDescription(t *testing.T) {
	// Use a temporary directory for testing
	tempDir := t.TempDir()
	tool := NewWriteFileTool(tempDir)
	if tool.Name() != "write_file" {
		t.Errorf("Expected name 'write_file', got %s", tool.Name())
	}
	if tool.Description() == "" {
		t.Error("Expected non-empty description")
	}
}

func TestWriteFileTool_WriteNewFile(t *testing.T) {
	// Use a temporary directory for testing
	tempDir := t.TempDir()
	tool := NewWriteFileTool(tempDir)

	ctx := context.Background()
	filePath := "test.txt"
	content := "Hello, World!"
	input := filePath + ":::" + content

	output, err := tool.Run(ctx, input)
	if err != nil {
		t.Fatalf("Unexpected error writing file: %v", err)
	}
	if output != fmt.Sprintf("Successfully wrote to %s", filePath) {
		t.Errorf("Expected success message, got %q", output)
	}

	// Verify the file was actually written
	fullPath := filepath.Join(tempDir, filePath)
	data, err := os.ReadFile(fullPath)
	if err != nil {
		t.Fatalf("Failed to read back written file: %v", err)
	}
	if string(data) != content {
		t.Errorf("Expected file content %q, got %q", content, string(data))
	}
}

func TestWriteFileTool_WriteExistingFile(t *testing.T) {
	// Use a temporary directory for testing
	tempDir := t.TempDir()
	tool := NewWriteFileTool(tempDir)

	// First create a file
	fullPath := filepath.Join(tempDir, "existing.txt")
	if err := os.WriteFile(fullPath, []byte("old content"), 0644); err != nil {
		t.Fatalf("Failed to create initial file: %v", err)
	}

	ctx := context.Background()
	filePath := "existing.txt"
	content := "new content"
	input := filePath + ":::" + content

	output, err := tool.Run(ctx, input)
	if err != nil {
		t.Fatalf("Unexpected error writing file: %v", err)
	}
	if output != fmt.Sprintf("Successfully wrote to %s", filePath) {
		t.Errorf("Expected success message, got %q", output)
	}

	// Verify the file was actually written
	data, err := os.ReadFile(fullPath)
	if err != nil {
		t.Fatalf("Failed to read back written file: %v", err)
	}
	if string(data) != content {
		t.Errorf("Expected file content %q, got %q", content, string(data))
	}
}

func TestWriteFileTool_InvalidInputFormat(t *testing.T) {
	tempDir := t.TempDir()
	tool := NewWriteFileTool(tempDir)

	ctx := context.Background()
	_, err := tool.Run(ctx, "invalid_input") // Missing the separator
	if err == nil {
		t.Error("Expected error for invalid input format, got nil")
	}
	if err.Error() != "invalid input format: expected 'filepath:::content'" {
		t.Errorf("Expected invalid input format error, got %v", err)
	}

	_, err = tool.Run(ctx, "filepath:::") // Empty content
	if err == nil {
		t.Error("Expected error for empty content, got nil")
	}
	if err.Error() != "empty file path" {
		t.Errorf("Expected empty file path error, got %v", err)
	}
}

func TestWriteFileTool_PathTraversal(t *testing.T) {
	tempDir := t.TempDir()
	tool := NewWriteFileTool(tempDir)

	ctx := context.Background()
	_, err := tool.Run(ctx, "../../etc/passwd:::malicious content") // Attempt path traversal
	if err == nil {
		t.Error("Expected error for path traversal, got nil")
	}
	if err.Error() != "path traversal detected: access outside root directory is not allowed" {
		t.Errorf("Expected path traversal error, got %v", err)
	}
}

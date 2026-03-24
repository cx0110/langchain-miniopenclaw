package tool

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestReadFileTool_NameAndDescription(t *testing.T) {
	// Use a temporary directory for testing
	tempDir := t.TempDir()
	tool := NewReadFileTool(tempDir)
	if tool.Name() != "read_file" {
		t.Errorf("Expected name 'read_file', got %s", tool.Name())
	}
	if tool.Description() == "" {
		t.Error("Expected non-empty description")
	}
}

func TestReadFileTool_ReadExistingFile(t *testing.T) {
	// Create a temporary directory and a test file
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.txt")
	content := "Hello, World!"
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	tool := NewReadFileTool(tempDir)

	ctx := context.Background()
	output, err := tool.Run(ctx, "test.txt")
	if err != nil {
		t.Fatalf("Unexpected error reading file: %v", err)
	}
	if output != content {
		t.Errorf("Expected file content %q, got %q", content, output)
	}
}

func TestReadFileTool_ReadNonexistentFile(t *testing.T) {
	tempDir := t.TempDir()
	tool := NewReadFileTool(tempDir)

	ctx := context.Background()
	_, err := tool.Run(ctx, "nonexistent.txt")
	if err == nil {
		t.Error("Expected error for nonexistent file, got nil")
	}
	if err.Error() != "file does not exist: nonexistent.txt" {
		t.Errorf("Expected 'file does not exist' error, got %v", err)
	}
}

func TestReadFileTool_ReadDirectory(t *testing.T) {
	tempDir := t.TempDir()
	tool := NewReadFileTool(tempDir)

	ctx := context.Background()
	_, err := tool.Run(ctx, ".") // Try to read the directory itself
	if err == nil {
		t.Error("Expected error for reading directory, got nil")
	}
	if err.Error() != "path is a directory, not a file" {
		t.Errorf("Expected 'path is a directory' error, got %v", err)
	}
}

func TestReadFileTool_PathTraversal(t *testing.T) {
	tempDir := t.TempDir()
	tool := NewReadFileTool(tempDir)

	ctx := context.Background()
	_, err := tool.Run(ctx, "../../etc/passwd") // Attempt path traversal
	if err == nil {
		t.Error("Expected error for path traversal, got nil")
	}
	if err.Error() != "path traversal detected: access outside root directory is not allowed" {
		t.Errorf("Expected path traversal error, got %v", err)
	}
}

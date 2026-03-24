package tool

import (
	"context"
	"testing"
	"time"
)

func TestTerminalTool_NameAndDescription(t *testing.T) {
	tool := NewTerminalTool()
	if tool.Name() != "terminal" {
		t.Errorf("Expected name 'terminal', got %s", tool.Name())
	}
	if tool.Description() == "" {
		t.Error("Expected non-empty description")
	}
}

func TestTerminalTool_RunSimpleCommand(t *testing.T) {
	tool := NewTerminalTool()

	// Test a simple command that should work on most systems
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	output, err := tool.Run(ctx, "echo hello")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if output != "hello\n" {
		t.Errorf("Expected 'hello\\n', got %q", output)
	}
}

func TestTerminalTool_RunBlockedCommand(t *testing.T) {
	tool := NewTerminalTool()

	// Test a blocked command
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := tool.Run(ctx, "rm -rf /")
	if err == nil {
		t.Error("Expected error for blocked command, got nil")
	}
	if err.Error() != "blocked: command matches the terminal blacklist" {
		t.Errorf("Expected blocked command error, got %v", err)
	}
}

func TestTerminalTool_RunEmptyCommand(t *testing.T) {
	tool := NewTerminalTool()

	// Test an empty command
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := tool.Run(ctx, "")
	if err == nil {
		t.Error("Expected error for empty command, got nil")
	}
	if err.Error() != "empty command" {
		t.Errorf("Expected 'empty command' error, got %v", err)
	}
}

package tool

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

// TerminalTool implements the Tool interface for executing shell commands
type TerminalTool struct {
	// Configuration could be added here
}

// NewTerminalTool creates a new terminal tool
func NewTerminalTool() *TerminalTool {
	return &TerminalTool{}
}

// Name returns the tool name
func (t *TerminalTool) Name() string {
	return "terminal"
}

// Description returns the tool description
func (t *TerminalTool) Description() string {
	return "Execute shell commands. Use this for quick inspection, building, or local commands."
}

// Run executes the terminal command
func (t *TerminalTool) Run(ctx context.Context, input string) (string, error) {
	// For simplicity, we'll treat the input as the raw command
	// In a more complex implementation, we might parse JSON input
	command := strings.TrimSpace(input)
	if command == "" {
		return "", errors.New("empty command")
	}

	// Basic security: block dangerous commands
	blockedPatterns := []string{
		"rm -rf /",
		"shutdown",
		"reboot",
		"mkfs",
		"format ",
		":(){:|:&};:",
	}
	lowered := strings.ToLower(command)
	for _, pattern := range blockedPatterns {
		if strings.Contains(lowered, pattern) {
			return "", errors.New("blocked: command matches the terminal blacklist")
		}
	}

	// Determine the shell to use
	bashPath, err := exec.LookPath("bash")
	if err == nil && bashPath != "" {
		// Use bash
		cmd := exec.CommandContext(ctx, "bash", "-lc", command)
		output, err := cmd.CombinedOutput()
		if err != nil {
			// If it's a context deadline exceeded error, return timeout message
			if ctx.Err() == context.DeadlineExceeded {
				return "", errors.New("command timed out")
			}
			// Combine stdout and stderr for error message
			if exitError, ok := err.(*exec.ExitError); ok {
				if len(exitError.Stderr) > 0 {
					return "", fmt.Errorf("command failed with exit code %d: %s", exitError.ExitCode(), string(exitError.Stderr))
				}
			}
			return "", fmt.Errorf("command failed: %w", err)
		}

		// Limit output size to prevent excessive responses
		if len(output) > 5000 {
			output = output[:5000]
		}

		return string(output), nil
	}

	// Try powershell as fallback
	powershellPath, err := exec.LookPath("powershell")
	if err == nil && powershellPath != "" {
		// Use powershell
		cmd := exec.CommandContext(ctx, "powershell", "-NoProfile", "-Command", command)
		output, err := cmd.CombinedOutput()
		if err != nil {
			// If it's a context deadline exceeded error, return timeout message
			if ctx.Err() == context.DeadlineExceeded {
				return "", errors.New("command timed out")
			}
			// Combine stdout and stderr for error message
			if exitError, ok := err.(*exec.ExitError); ok {
				if len(exitError.Stderr) > 0 {
					return "", fmt.Errorf("command failed with exit code %d: %s", exitError.ExitCode(), string(exitError.Stderr))
				}
			}
			return "", fmt.Errorf("command failed: %w", err)
		}

		// Limit output size to prevent excessive responses
		if len(output) > 5000 {
			output = output[:5000]
		}

		return string(output), nil
	}

	return "", errors.New("no suitable shell found (bash or powershell required)")
}

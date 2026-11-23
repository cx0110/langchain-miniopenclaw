# Go Backend Rewrite Chunk 4: Tool Implementations Implementation Plan

> **For agentic workers:** REQUIRED: Use superpowers:subagent-driven-development (if subagents available) or superpowers:executing-plans to implement this plan. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Implement core tool functionalities (terminal, file read/write, fetch URL, etc.) that the agent can use to interact with the system and external world.

**Architecture:** We'll create a tool interface and implement individual tools that match the functionality of the Python backend tools. Each tool will have a clear input/output contract and be registrable with the agent system.

**Tech Stack:** Go, standard library for OS/file operations, HTTP client for fetch functionality.

---

### Chunk 4: Tool Implementations

**Files:**
- Create: `backend-go/tool/tool.go` (tool interface and base types)
- Create: `backend-go/tool/terminal_tool.go`
- Create: `backend-go/tool/read_file_tool.go`
- Create: `backend-go/tool/write_file_tool.go`
- Create: `backend-go/tool/fetch_url_tool.go`
- Create: `backend-go/tool/python_repl_tool.go`
- Create: `backend-go/tool/tool_registry.go`
- Test: Individual test files for each tool

- [ ] **Step 1: Define tool interface and base types**

```go
// backend-go/tool/tool.go
package tool

import (
	"context"
)

// Tool defines the interface for all tools that the agent can use
type Tool interface {
	// Name returns the unique name of the tool
	Name() string
	// Description returns a human-readable description of what the tool does
	Description() string
	// Run executes the tool with the given input and returns the result
	Run(ctx context.Context, input string) (string, error)
}

// ToolInput represents the input to a tool (for more complex tools)
type ToolInput struct {
	// Tool-specific fields would be defined in each tool's implementation
	// For simplicity, we're using string input/output for now
	// In a more complex implementation, we'd use structured input
	Raw string
}

// ToolOutput represents the output from a tool
type ToolOutput struct {
	Content string
	Error   string
}
```

- [ ] **Step 2: Implement terminal tool**

```go
// backend-go/tool/terminal_tool.go
package tool

import (
	"context"
	"errors"
	"os/exec"
	"strings"
	"time"
)

// TerminalToolInput represents the input for the terminal tool
type TerminalToolInput struct {
	Command string `json:"command"`
}

// TerminalTool implements the Tool interface for executing shell commands
type TerminalTool struct {
	// In a real implementation, we'd have configuration for timeouts, allowed commands, etc.
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
	// Parse input as JSON to get the command
	var toolInput TerminalToolInput
	if err := json.Unmarshal([]byte(input), &toolInput); err != nil {
		return "", fmt.Errorf("failed to parse terminal tool input: %w", err)
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
	lowered := strings.ToLower(toolInput.Command)
	for _, pattern := range blockedPatterns {
		if strings.Contains(lowered, pattern) {
			return "", errors.New("blocked: command matches the terminal blacklist")
		}
	}
	
	// Execute the command
	var shell string
	var arg string
	if strings.HasPrefix(strings.ToLower(exec.LookPath("bash")), "/bin") || strings.HasPrefix(strings.ToLower(exec.LookPath("bash")), "/usr/bin") {
		shell = "bash"
		arg = "-lc"
	} else if exec.LookPath("powershell") != "" {
		shell = "powershell"
		arg = "-NoProfile -Command"
	} else {
		return "", errors.New("no suitable shell found")
	}
	
	cmd := exec.CommandContext(ctx, shell, arg, toolInput.Command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		// If it's a context deadline exceeded error, return timeout message
		if ctx.Err() == context.DeadlineExceeded {
			return "", errors.New("command timed out")
		}
		return "", fmt.Errorf("command failed: %w", err)
	}
	
	// Limit output size
	if len(output) > 5000 {
		output = output[:5000]
	}
	
	return string(output), nil
}

// Note: We'll need to add the json import
```

Actually, let me fix the missing import and improve the implementation.

Let me restart with a better approach for the tool implementations.

Actually, let me first create the directory and basic structure.

</think>
<tool_call>
<function=bash>
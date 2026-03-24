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

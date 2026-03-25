package tool

import "sync"

// ToolRegistry holds a collection of tools and provides lookup by name
type ToolRegistry struct {
	tools map[string]Tool
	mu    sync.RWMutex
}

// NewToolRegistry creates a new tool registry
func NewToolRegistry() *ToolRegistry {
	return &ToolRegistry{
		tools: make(map[string]Tool),
	}
}

// Register adds a tool to the registry
func (tr *ToolRegistry) Register(tool Tool) {
	tr.mu.Lock()
	defer tr.mu.Unlock()
	tr.tools[tool.Name()] = tool
}

// Get returns a tool by name, or nil if not found
func (tr *ToolRegistry) Get(name string) Tool {
	tr.mu.RLock()
	defer tr.mu.RUnlock()
	return tr.tools[name]
}

// List returns all registered tool names
func (tr *ToolRegistry) List() []string {
	tr.mu.RLock()
	defer tr.mu.RUnlock()
	names := make([]string, 0, len(tr.tools))
	for name := range tr.tools {
		names = append(names, name)
	}
	return names
}

package tool

import "testing"

func TestToolRegistry(t *testing.T) {
	registry := NewToolRegistry()

	// Initially empty
	if len(registry.List()) != 0 {
		t.Errorf("Expected empty registry, got %v", registry.List())
	}

	// Register a terminal tool
	terminalTool := NewTerminalTool()
	registry.Register(terminalTool)

	// Check that it's registered
	if retrieved := registry.Get("terminal"); retrieved == nil {
		t.Error("Expected to find terminal tool, got nil")
	} else if retrieved.Name() != "terminal" {
		t.Errorf("Expected terminal tool, got %s", retrieved.Name())
	}

	// List should contain terminal
	names := registry.List()
	if len(names) != 1 {
		t.Errorf("Expected 1 tool in registry, got %d", len(names))
	}
	if names[0] != "terminal" {
		t.Errorf("Expected 'terminal' in list, got %s", names[0])
	}

	// Register another tool
	readFileTool := NewReadFileTool("/tmp")
	registry.Register(readFileTool)

	// Check both are present
	if retrieved := registry.Get("terminal"); retrieved == nil {
		t.Error("Expected to find terminal tool after adding another, got nil")
	}
	if retrieved := registry.Get("read_file"); retrieved == nil {
		t.Error("Expected to find read_file tool, got nil")
	}

	// List should contain both
	names = registry.List()
	if len(names) != 2 {
		t.Errorf("Expected 2 tools in registry, got %d", len(names))
	}
	// Check that both are present (order doesn't matter)
	foundTerminal := false
	foundReadFile := false
	for _, name := range names {
		if name == "terminal" {
			foundTerminal = true
		}
		if name == "read_file" {
			foundReadFile = true
		}
	}
	if !foundTerminal {
		t.Error("Expected to find terminal tool in registry list")
	}
	if !foundReadFile {
		t.Error("Expected to find read_file tool in registry list")
	}
}

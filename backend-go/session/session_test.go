package session

import "testing"

func TestSessionManager_LoadSaveSession(t *testing.T) {
	// Create a temporary directory for tests
	tempDir := t.TempDir()

	sm := NewSessionManager(tempDir)
	sessionID := "test-session"

	// Load a new session (should be empty)
	record, err := sm.LoadSessionRecord(sessionID)
	if err != nil {
		t.Fatalf("Failed to load session: %v", err)
	}
	if record.ID != sessionID {
		t.Errorf("Expected session ID %s, got %s", sessionID, record.ID)
	}
	if len(record.Messages) != 0 {
		t.Errorf("Expected empty messages, got %d", len(record.Messages))
	}

	// Save a message
	err = sm.SaveMessage(sessionID, "user", "Hello")
	if err != nil {
		t.Fatalf("Failed to save message: %v", err)
	}

	// Load the session again and verify the message was saved
	record, err = sm.LoadSessionRecord(sessionID)
	if err != nil {
		t.Fatalf("Failed to load session after saving: %v", err)
	}
	if len(record.Messages) != 1 {
		t.Errorf("Expected 1 message, got %d", len(record.Messages))
	}
	if record.Messages[0].Role != "user" || record.Messages[0].Content != "Hello" {
		t.Errorf("Expected user message 'Hello', got %+v", record.Messages[0])
	}

	// Test saving assistant message
	err = sm.SaveMessage(sessionID, "assistant", "Hi there!")
	if err != nil {
		t.Fatalf("Failed to save assistant message: %v", err)
	}

	record, err = sm.LoadSessionRecord(sessionID)
	if err != nil {
		t.Fatalf("Failed to load session after saving assistant message: %v", err)
	}
	if len(record.Messages) != 2 {
		t.Errorf("Expected 2 messages, got %d", len(record.Messages))
	}
	if record.Messages[1].Role != "assistant" || record.Messages[1].Content != "Hi there!" {
		t.Errorf("Expected assistant message 'Hi there!', got %+v", record.Messages[1])
	}
}

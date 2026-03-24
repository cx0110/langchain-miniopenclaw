package session

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// SessionRecord represents the full session data stored on disk
type SessionRecord struct {
	ID        string           `json:"id"`
	Messages  []SessionMessage `json:"messages"`
	Title     string           `json:"title,omitempty"`
	CreatedAt int64            `json:"created_at,omitempty"`
	UpdatedAt int64            `json:"updated_at,omitempty"`
}

// SessionMessage represents a single message in the session
type SessionMessage struct {
	ID        string `json:"id,omitempty"`
	Role      string `json:"role"` // "user" or "assistant"
	Content   string `json:"content"`
	CreatedAt int64  `json:"created_at,omitempty"`
}

// SessionManager handles loading and saving session data
type SessionManager struct {
	baseDir string
}

// NewSessionManager creates a new session manager
func NewSessionManager(baseDir string) *SessionManager {
	return &SessionManager{
		baseDir: baseDir,
	}
}

// GetSessionFilePath returns the file path for a session
func (sm *SessionManager) GetSessionFilePath(sessionID string) string {
	return filepath.Join(sm.baseDir, "sessions", sessionID+".json")
}

// LoadSessionRecord loads the full session record from disk
func (sm *SessionManager) LoadSessionRecord(sessionID string) (*SessionRecord, error) {
	filePath := sm.GetSessionFilePath(sessionID)

	// If file doesn't exist, return a new empty session
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return &SessionRecord{
			ID:        sessionID,
			Messages:  []SessionMessage{},
			CreatedAt: nowUnix(),
			UpdatedAt: nowUnix(),
		}, nil
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var record SessionRecord
	if err := json.Unmarshal(data, &record); err != nil {
		return nil, err
	}

	return &record, nil
}

// SaveSessionRecord saves the session record to disk
func (sm *SessionManager) SaveSessionRecord(record *SessionRecord) error {
	filePath := sm.GetSessionFilePath(record.ID)

	// Ensure the sessions directory exists
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	record.UpdatedAt = nowUnix()

	data, err := json.MarshalIndent(record, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, data, 0644)
}

// LoadSessionForAgent returns messages in the format expected by the agent
func (sm *SessionManager) LoadSessionForAgent(sessionID string) ([]AgentMessage, error) {
	record, err := sm.LoadSessionRecord(sessionID)
	if err != nil {
		return nil, err
	}

	messages := make([]AgentMessage, 0, len(record.Messages))
	for _, msg := range record.Messages {
		messages = append(messages, AgentMessage{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}

	return messages, nil
}

// SaveMessage appends a new message to the session
func (sm *SessionManager) SaveMessage(sessionID, role, content string) error {
	record, err := sm.LoadSessionRecord(sessionID)
	if err != nil {
		return err
	}

	// Generate a simple ID (in practice, you might want to use UUID)
	newID := fmt.Sprintf("msg_%d", nowUnix())

	record.Messages = append(record.Messages, SessionMessage{
		ID:        newID,
		Role:      role,
		Content:   content,
		CreatedAt: nowUnix(),
	})

	return sm.SaveSessionRecord(record)
}

// Helper function to get current unix timestamp
func nowUnix() int64 {
	return time.Now().Unix()
}

// AgentMessage represents a message in the format expected by the agent
type AgentMessage struct {
	Role    string
	Content string
}

package agent

import (
	"context"
	"fmt"

	"langchain-miniopenclaw-backend-go/llm"
	"langchain-miniopenclaw-backend-go/session"
)

// AgentManager handles processing user messages using the LLM and session storage
type AgentManager struct {
	sessionMgr *session.SessionManager
	llmClient  llm.LLMClient
}

// NewAgentManager creates a new agent manager
func NewAgentManager(sessionMgr *session.SessionManager, llmClient llm.LLMClient) *AgentManager {
	return &AgentManager{
		sessionMgr: sessionMgr,
		llmClient:  llmClient,
	}
}

// ProcessMessage handles a user message and returns the agent's response
func (am *AgentManager) ProcessMessage(ctx context.Context, sessionID, userMessage string) (string, error) {
	// 1. Load the session history
	history, err := am.sessionMgr.LoadSessionForAgent(sessionID)
	if err != nil {
		return "", fmt.Errorf("failed to load session: %w", err)
	}

	// 2. Save the user message to the session
	if err := am.sessionMgr.SaveMessage(sessionID, "user", userMessage); err != nil {
		return "", fmt.Errorf("failed to save user message: %w", err)
	}

	// 3. Prepare messages for the LLM (system prompt + history + current message)
	var llmMessages []llm.Message

	// Add system prompt (we'll make this configurable later)
	llmMessages = append(llmMessages, llm.Message{
		Role:    "system",
		Content: "你是一个有帮助的AI助手。请用中文回答用户的问题。",
	})

	// Add history
	for _, msg := range history {
		llmMessages = append(llmMessages, llm.Message{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}

	// Add the current user message
	llmMessages = append(llmMessages, llm.Message{
		Role:    "user",
		Content: userMessage,
	})

	// 4. Generate response from LLM
	aiResponse, err := am.llmClient.GenerateContentWithMessages(ctx, llmMessages)
	if err != nil {
		return "", fmt.Errorf("LLM error: %w", err)
	}

	// 5. Save the assistant response to the session
	if err := am.sessionMgr.SaveMessage(sessionID, "assistant", aiResponse); err != nil {
		return "", fmt.Errorf("failed to save assistant message: %w", err)
	}

	return aiResponse, nil
}

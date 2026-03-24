package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"langchain-miniopenclaw-backend-go/agent"
	"langchain-miniopenclaw-backend-go/config"
	"langchain-miniopenclaw-backend-go/llm"
	"langchain-miniopenclaw-backend-go/model"
	"langchain-miniopenclaw-backend-go/session"
)

// ChatHandler handles the /chat endpoint
func ChatHandler(c *gin.Context) {
	var req model.ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Initialize dependencies (in a real app, these would be singletons or injected)
	// For this example, we'll create them here - in practice, you'd want to avoid
	// creating new instances on every request
	settings, err := config.LoadSettings()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load configuration"})
		return
	}

	// Set up session manager
	sessMgr := session.NewSessionManager(settings.BackendDir)

	// Set up LLM client (using OpenAI client as an example)
	// In a real implementation, you would choose the client based on settings.LLMProvider
	var llmClient llm.LLMClient
	if settings.LLMAPIKey != "" {
		llmClient = llm.NewOpenAIClient(settings.LLMAPIKey, settings.LLMBaseURL, settings.LLMModel)
	} else {
		// If no API key is set, create a mock client for testing/development
		llmClient = &llm.MockClient{
			Response: "这是一个测试响应（未配置LLM API密钥）",
		}
	}

	// Set up agent manager
	agentMgr := agent.NewAgentManager(sessMgr, llmClient)

	// Process the request
	ctx := context.Background()
	response, err := agentMgr.ProcessMessage(ctx, req.SessionID, req.Message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return response
	res := model.ChatResponse{
		Content: response,
	}
	c.JSON(http.StatusOK, res)
}

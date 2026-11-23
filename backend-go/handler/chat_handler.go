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

// Handler encapsulates HTTP handler dependencies
type Handler struct {
	cfg *config.Settings
}

// NewHandler creates a new Handler with the given configuration
func NewHandler(cfg *config.Settings) *Handler {
	return &Handler{
		cfg: cfg,
	}
}

// ChatHandler handles the /chat endpoint
func (h *Handler) ChatHandler(c *gin.Context) {
	var req model.ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set up session manager
	sessMgr := session.NewSessionManager(h.cfg.BackendDir)

	// Set up LLM client
	var llmClient llm.LLMClient
	if h.cfg.LLMAPIKey != "" {
		llmClient = llm.NewOpenAIClient(h.cfg.LLMAPIKey, h.cfg.LLMBaseURL, h.cfg.LLMModel)
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

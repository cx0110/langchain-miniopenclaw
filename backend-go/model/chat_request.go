package model

// ChatRequest represents the incoming chat request from the frontend
type ChatRequest struct {
	Message   string `json:"message" binding:"required,min=1"`
	SessionID string `json:"session_id" binding:"required"`
	Stream    bool   `json:"stream"`
}

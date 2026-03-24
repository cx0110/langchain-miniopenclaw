package model

// ChatResponse represents the outgoing chat response
type ChatResponse struct {
	Content string `json:"content,omitempty"`
}

// StreamEvent represents a single event in the streaming response
type StreamEvent struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}

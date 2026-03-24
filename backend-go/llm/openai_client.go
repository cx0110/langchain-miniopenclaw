package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// OpenAIClient implements LLMClient for OpenAI-compatible APIs
type OpenAIClient struct {
	APIKey     string
	BaseURL    string
	Model      string
	HTTPClient *http.Client
}

// NewOpenAIClient creates a new OpenAIClient
func NewOpenAIClient(apiKey, baseURL, model string) *OpenAIClient {
	return &OpenAIClient{
		APIKey:     apiKey,
		BaseURL:    baseURL,
		Model:      model,
		HTTPClient: &http.Client{},
	}
}

// GenerateContent generates a completion for a single prompt
func (c *OpenAIClient) GenerateContent(ctx context.Context, prompt string) (string, error) {
	// We'll use the chat completions endpoint with a single user message
	messages := []Message{
		{Role: "user", Content: prompt},
	}
	return c.GenerateContentWithMessages(ctx, messages)
}

// GenerateContentWithMessages generates a completion based on a list of messages
func (c *OpenAIClient) GenerateContentWithMessages(ctx context.Context, messages []Message) (string, error) {
	// Convert our Message format to OpenAI's format
	openaiMessages := make([]map[string]interface{}, len(messages))
	for i, msg := range messages {
		openaiMessages[i] = map[string]interface{}{
			"role":    msg.Role,
			"content": msg.Content,
		}
	}

	requestBody := map[string]interface{}{
		"model":    c.Model,
		"messages": openaiMessages,
		// We can add temperature, max_tokens, etc. if needed
		"temperature": 0,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.BaseURL+"/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.APIKey)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Parse the response
	var respBody struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return "", err
	}

	if len(respBody.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}

	return respBody.Choices[0].Message.Content, nil
}

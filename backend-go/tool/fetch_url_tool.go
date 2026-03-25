package tool

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// FetchURLTool implements the Tool interface for fetching URLs
type FetchURLTool struct {
	// Timeout for HTTP requests
	timeout time.Duration
}

// NewFetchURLTool creates a new fetch URL tool with the given timeout
func NewFetchURLTool(timeout time.Duration) *FetchURLTool {
	return &FetchURLTool{
		timeout: timeout,
	}
}

// Name returns the tool name
func (t *FetchURLTool) Name() string {
	return "fetch_url"
}

// Description returns the tool description
func (t *FetchURLTool) Description() string {
	return "Fetch content from a URL. Use this to get information from the web."
}

// Run executes the fetch URL operation
func (t *FetchURLTool) Run(ctx context.Context, input string) (string, error) {
	// For simplicity, we'll treat the input as the raw URL
	// In a more complex implementation, we might parse JSON input
	urlStr := strings.TrimSpace(input)
	if urlStr == "" {
		return "", errors.New("empty URL")
	}

	// Validate URL format
	_, err := url.ParseRequestURI(urlStr)
	if err != nil {
		return "", fmt.Errorf("invalid URL: %w", err)
	}

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: t.timeout,
	}

	// Create request with context
	req, err := http.NewRequestWithContext(ctx, "GET", urlStr, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Set a reasonable user agent
	req.Header.Set("User-Agent", "langchain-miniopenclaw-go/1.0")

	// Execute the request
	resp, err := client.Do(req)
	if err != nil {
		// Check if it's a context cancellation or timeout
		if ctx.Err() == context.DeadlineExceeded {
			return "", errors.New("request timed out")
		}
		if ctx.Err() == context.Canceled {
			return "", errors.New("request canceled")
		}
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check if the response status is OK
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	// Limit output size to prevent excessive responses
	if len(body) > 10000 {
		body = body[:10000]
	}

	return string(body), nil
}

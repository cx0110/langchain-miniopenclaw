package tool

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestFetchURLTool_NameAndDescription(t *testing.T) {
	tool := NewFetchURLTool(30 * time.Second)
	if tool.Name() != "fetch_url" {
		t.Errorf("Expected name 'fetch_url', got %s", tool.Name())
	}
	if tool.Description() == "" {
		t.Error("Expected non-empty description")
	}
}

func TestFetchURLTool_FetchValidURL(t *testing.T) {
	// Create a test server to mock an external URL
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"test": "success"}`))
	}))
	defer ts.Close()

	tool := NewFetchURLTool(5 * time.Second)

	ctx := context.Background()
	output, err := tool.Run(ctx, ts.URL)
	if err != nil {
		t.Fatalf("Unexpected error fetching URL: %v", err)
	}
	if output != `{"test": "success"}` {
		t.Errorf("Expected response body, got %q", output)
	}
}

func TestFetchURLTool_InvalidURL(t *testing.T) {
	tool := NewFetchURLTool(5 * time.Second)

	ctx := context.Background()
	_, err := tool.Run(ctx, "not a valid url")
	if err == nil {
		t.Error("Expected error for invalid URL, got nil")
	}
	// Accept any error that indicates invalid URL
	if err == nil || !strings.Contains(err.Error(), "invalid URL") {
		t.Errorf("Expected invalid URL error, got %v", err)
	}
}

func TestFetchURLTool_EmptyURL(t *testing.T) {
	tool := NewFetchURLTool(5 * time.Second)

	ctx := context.Background()
	_, err := tool.Run(ctx, "")
	if err == nil {
		t.Error("Expected error for empty URL, got nil")
	}
	if err.Error() != "empty URL" {
		t.Errorf("Expected 'empty URL' error, got %v", err)
	}
}

func TestFetchURLTool_Timout(t *testing.T) {
	// Create a server that delays response
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate slow response
		<-time.After(2 * time.Second)
		w.Write([]byte("slow response"))
	}))
	defer ts.Close()

	// Use a very short timeout
	tool := NewFetchURLTool(100 * time.Millisecond)

	ctx := context.Background()
	_, err := tool.Run(ctx, ts.URL)
	if err == nil {
		t.Error("Expected timeout error, got nil")
	}
	// Check for timeout error
	if !strings.Contains(err.Error(), "context deadline exceeded") && !strings.Contains(err.Error(), "timeout") {
		t.Errorf("Expected timeout error, got %v", err)
	}
}

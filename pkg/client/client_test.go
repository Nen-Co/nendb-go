package client

import (
	"context"
	"net/http"
	"testing"
	"time"
)

func TestClientConfig(t *testing.T) {
	// Test default config
	config := DefaultConfig()
	if config.BaseURL != "http://localhost:8080" {
		t.Errorf("Expected default BaseURL to be 'http://localhost:8080', got '%s'", config.BaseURL)
	}
	if config.Timeout != 30*time.Second {
		t.Errorf("Expected default Timeout to be 30s, got %v", config.Timeout)
	}
	if config.MaxRetries != 3 {
		t.Errorf("Expected default MaxRetries to be 3, got %d", config.MaxRetries)
	}

	// Test custom config
	customConfig := &ClientConfig{
		BaseURL:    "http://example.com:9090",
		Timeout:    60 * time.Second,
		MaxRetries: 5,
	}
	if customConfig.BaseURL != "http://example.com:9090" {
		t.Errorf("Expected custom BaseURL to be 'http://example.com:9090', got '%s'", customConfig.BaseURL)
	}
}

func TestNewClient(t *testing.T) {
	// Test with nil config (should use defaults and try to connect to localhost:8080)
	// Since there might be a server running, we'll test with a non-existent server
	config := &ClientConfig{
		BaseURL:        "http://localhost:9999", // Non-existent server
		Timeout:        1 * time.Second,         // Short timeout for test
		MaxRetries:     1,                       // Few retries for test
		SkipValidation: false,                   // Should try to validate
	}
	
	client, err := NewClient(config)
	if err == nil {
		t.Error("Expected error when creating client without server, got nil")
	}
	// The error should be a connection error or timeout error
	if err != nil {
		t.Logf("Got expected error: %v", err)
	}

	// Test with custom config and skip validation
	config2 := &ClientConfig{
		BaseURL:        "http://localhost:9999", // Non-existent server
		Timeout:        5 * time.Second,
		MaxRetries:     1,
		SkipValidation: true,
	}
	
	client, err = NewClient(config2)
	if err != nil {
		t.Errorf("Expected no error when skip validation is true, got %v", err)
	}
	if client == nil {
		t.Error("Expected client to be created, got nil")
	}

	// Test with nil config but skip validation
	config3 := &ClientConfig{
		BaseURL:        "http://localhost:9999",
		SkipValidation: true,
	}
	
	client2, err := NewClient(config3)
	if err != nil {
		t.Errorf("Expected no error when skip validation is true, got %v", err)
	}
	if client2 == nil {
		t.Error("Expected client to be created, got nil")
	}

	// Test with nil config but skip validation - should work
	config4 := &ClientConfig{
		SkipValidation: true,
	}
	
	client3, err := NewClient(config4)
	if err != nil {
		t.Errorf("Expected no error when skip validation is true, got %v", err)
	}
	if client3 == nil {
		t.Error("Expected client to be created, got nil")
	}
}

func TestClientBaseURL(t *testing.T) {
	// Test with trailing slash
	config := &ClientConfig{
		BaseURL:        "http://localhost:8080/",
		SkipValidation: true,
	}
	
	client, err := NewClient(config)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	
	// BaseURL should be trimmed of trailing slash
	if client.baseURL != "http://localhost:8080" {
		t.Errorf("Expected baseURL to be 'http://localhost:8080', got '%s'", client.baseURL)
	}
}

func TestClientHTTPClient(t *testing.T) {
	// Test with custom HTTP client
	customHTTPClient := &http.Client{
		Timeout: 10 * time.Second,
	}
	
	config := &ClientConfig{
		BaseURL:        "http://localhost:8080",
		HTTPClient:     customHTTPClient,
		SkipValidation: true,
	}
	
	client, err := NewClient(config)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	
	// Should use the custom HTTP client
	if client.httpClient != customHTTPClient {
		t.Error("Expected client to use custom HTTP client")
	}
}

func TestContextTimeout(t *testing.T) {
	// Test context timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	// Wait for context to timeout
	time.Sleep(2 * time.Millisecond)

	select {
	case <-ctx.Done():
		// Expected
	default:
		t.Error("Expected context to be done after timeout")
	}
}

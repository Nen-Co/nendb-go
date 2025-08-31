package tests

import (
	"context"
	"testing"
	"time"

	"github.com/nen-co/nendb-go-driver/pkg/client"
	"github.com/nen-co/nendb-go-driver/pkg/types"
)

func TestClientConfig(t *testing.T) {
	// Test default config
	config := client.DefaultConfig()
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
	customConfig := &client.ClientConfig{
		BaseURL:    "http://example.com:9090",
		Timeout:    60 * time.Second,
		MaxRetries: 5,
	}
	if customConfig.BaseURL != "http://example.com:9090" {
		t.Errorf("Expected custom BaseURL to be 'http://example.com:9090', got '%s'", customConfig.BaseURL)
	}
}

func TestNewClient(t *testing.T) {
	// Test with nil config (should use defaults)
	client, err := client.NewClient(nil)
	if err == nil {
		t.Error("Expected error when creating client without server, got nil")
	}

	// Test with custom config and skip validation
	config := &client.ClientConfig{
		BaseURL:        "http://localhost:9999", // Non-existent server
		Timeout:        5 * time.Second,
		MaxRetries:     1,
		SkipValidation: true,
	}
	
	client, err = client.NewClient(config)
	if err != nil {
		t.Errorf("Expected no error when skip validation is true, got %v", err)
	}
	if client == nil {
		t.Error("Expected client to be created, got nil")
	}
}

func TestGraphNodeValidation(t *testing.T) {
	// Test valid node
	node, err := types.NewGraphNode(1, []string{"Person"}, map[string]interface{}{"name": "Alice"})
	if err != nil {
		t.Errorf("Expected no error creating valid node, got %v", err)
	}
	if node.ID != 1 {
		t.Errorf("Expected node ID to be 1, got %d", node.ID)
	}

	// Test invalid node ID
	_, err = types.NewGraphNode(-1, []string{"Person"}, map[string]interface{}{"name": "Alice"})
	if err == nil {
		t.Error("Expected error for negative node ID, got nil")
	}

	// Test validation
	if err := node.Validate(); err != nil {
		t.Errorf("Expected no validation error, got %v", err)
	}
}

func TestGraphEdgeValidation(t *testing.T) {
	// Test valid edge
	edge, err := types.NewGraphEdge(1, 1, 2, "KNOWS", map[string]interface{}{"since": "2022"})
	if err != nil {
		t.Errorf("Expected no error creating valid edge, got %v", err)
	}
	if edge.ID != 1 {
		t.Errorf("Expected edge ID to be 1, got %d", edge.ID)
	}

	// Test invalid edge ID
	_, err = types.NewGraphEdge(-1, 1, 2, "KNOWS", map[string]interface{}{"since": "2022"})
	if err == nil {
		t.Error("Expected error for negative edge ID, got nil")
	}

	// Test invalid source node
	_, err = types.NewGraphEdge(1, -1, 2, "KNOWS", map[string]interface{}{"since": "2022"})
	if err == nil {
		t.Error("Expected error for negative source node ID, got nil")
	}

	// Test invalid target node
	_, err = types.NewGraphEdge(1, 1, -1, "KNOWS", map[string]interface{}{"since": "2022"})
	if err == nil {
		t.Error("Expected error for negative target node ID, got nil")
	}

	// Test empty edge type
	_, err = types.NewGraphEdge(1, 1, 2, "", map[string]interface{}{"since": "2022"})
	if err == nil {
		t.Error("Expected error for empty edge type, got nil")
	}

	// Test validation
	if err := edge.Validate(); err != nil {
		t.Errorf("Expected no validation error, got %v", err)
	}
}

func TestAlgorithmResultValidation(t *testing.T) {
	// Test valid algorithm result
	result, err := types.NewAlgorithmResult("BFS", types.StatusCompleted, "Algorithm completed successfully", nil)
	if err != nil {
		t.Errorf("Expected no error creating valid algorithm result, got %v", err)
	}
	if result.Algorithm != "BFS" {
		t.Errorf("Expected algorithm to be 'BFS', got '%s'", result.Algorithm)
	}

	// Test empty algorithm name
	_, err = types.NewAlgorithmResult("", types.StatusCompleted, "Algorithm completed successfully", nil)
	if err == nil {
		t.Error("Expected error for empty algorithm name, got nil")
	}

	// Test empty message
	_, err = types.NewAlgorithmResult("BFS", types.StatusCompleted, "", nil)
	if err == nil {
		t.Error("Expected error for empty message, got nil")
	}

	// Test validation
	if err := result.Validate(); err != nil {
		t.Errorf("Expected no validation error, got %v", err)
	}
}

func TestBFSResult(t *testing.T) {
	baseResult, err := types.NewAlgorithmResult("BFS", types.StatusCompleted, "BFS completed", nil)
	if err != nil {
		t.Fatalf("Failed to create base result: %v", err)
	}

	bfsResult := types.NewBFSResult(baseResult, []int{1, 2, 3}, []int{1, 2, 3}, 2)
	if bfsResult.VisitedNodes == nil {
		t.Error("Expected visited nodes to be initialized, got nil")
	}
	if bfsResult.Path == nil {
		t.Error("Expected path to be initialized, got nil")
	}
	if bfsResult.Depth != 2 {
		t.Errorf("Expected depth to be 2, got %d", bfsResult.Depth)
	}
}

func TestDijkstraResult(t *testing.T) {
	baseResult, err := types.NewAlgorithmResult("Dijkstra", types.StatusCompleted, "Dijkstra completed", nil)
	if err != nil {
		t.Fatalf("Failed to create base result: %v", err)
	}

	dijkstraResult := types.NewDijkstraResult(baseResult, []int{1, 2, 3}, 15.5, []map[string]interface{}{})
	if dijkstraResult.ShortestPath == nil {
		t.Error("Expected shortest path to be initialized, got nil")
	}
	if dijkstraResult.PathDetails == nil {
		t.Error("Expected path details to be initialized, got nil")
	}
	if dijkstraResult.TotalCost != 15.5 {
		t.Errorf("Expected total cost to be 15.5, got %f", dijkstraResult.TotalCost)
	}
}

func TestPageRankResult(t *testing.T) {
	baseResult, err := types.NewAlgorithmResult("PageRank", types.StatusCompleted, "PageRank completed", nil)
	if err != nil {
		t.Fatalf("Failed to create base result: %v", err)
	}

	pagerankResult := types.NewPageRankResult(baseResult, map[int]float64{1: 0.5, 2: 0.3}, 50, true)
	if pagerankResult.NodeScores == nil {
		t.Error("Expected node scores to be initialized, got nil")
	}
	if pagerankResult.Iterations != 50 {
		t.Errorf("Expected iterations to be 50, got %d", pagerankResult.Iterations)
	}
	if !pagerankResult.Convergence {
		t.Error("Expected convergence to be true, got false")
	}
}

func TestPropertyValueValidation(t *testing.T) {
	// Test valid property values
	validValues := []interface{}{
		"string",
		42,
		int8(8),
		int16(16),
		int32(32),
		int64(64),
		uint(42),
		uint8(8),
		uint16(16),
		uint32(32),
		uint64(64),
		float32(3.14),
		float64(3.14),
		true,
		false,
		nil,
	}

	for _, value := range validValues {
		if !types.IsValidPropertyValue(value) {
			t.Errorf("Expected %v to be a valid property value", value)
		}
	}

	// Test invalid property values
	invalidValues := []interface{}{
		[]int{1, 2, 3},
		map[string]int{"a": 1},
		struct{}{},
		func() {},
	}

	for _, value := range invalidValues {
		if types.IsValidPropertyValue(value) {
			t.Errorf("Expected %v to be an invalid property value", value)
		}
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

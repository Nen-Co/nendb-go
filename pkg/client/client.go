package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/nen-co/nendb-go-driver/pkg/errors"
	"github.com/nen-co/nendb-go-driver/pkg/types"
)

// ClientConfig holds configuration for the NenDB client
type ClientConfig struct {
	BaseURL        string
	Timeout        time.Duration
	MaxRetries     int
	RetryDelay     time.Duration
	SkipValidation bool
	HTTPClient     *http.Client
}

// DefaultConfig returns a default client configuration
func DefaultConfig() *ClientConfig {
	return &ClientConfig{
		BaseURL:    "http://localhost:8080",
		Timeout:    30 * time.Second,
		MaxRetries: 3,
		RetryDelay: 1 * time.Second,
	}
}

// NenDBClient is the main client for interacting with NenDB
type NenDBClient struct {
	config     *ClientConfig
	httpClient *http.Client
	baseURL    string
}

// NewClient creates a new NenDB client
func NewClient(config *ClientConfig) (*NenDBClient, error) {
	if config == nil {
		config = DefaultConfig()
	}

	// Clean up base URL
	baseURL := strings.TrimRight(config.BaseURL, "/")

	// Create HTTP client if not provided
	httpClient := config.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: config.Timeout,
		}
	}

	client := &NenDBClient{
		config:     config,
		httpClient: httpClient,
		baseURL:    baseURL,
	}

	// Validate connection if not skipped
	if !config.SkipValidation {
		if err := client.Health(); err != nil {
			return nil, errors.NewConnectionError(
				fmt.Sprintf("Failed to connect to NenDB server at %s", baseURL),
				map[string]interface{}{"error": err.Error()},
			)
		}
	}

	return client, nil
}

// makeRequest performs an HTTP request with retry logic
func (c *NenDBClient) makeRequest(ctx context.Context, method, endpoint string, data interface{}, params map[string]string) ([]byte, error) {
	// Build URL
	requestURL := c.baseURL + endpoint
	if len(params) > 0 {
		u, err := url.Parse(requestURL)
		if err != nil {
			return nil, errors.NewValidationError("Invalid URL", map[string]interface{}{"url": requestURL, "error": err.Error()})
		}
		q := u.Query()
		for key, value := range params {
			q.Set(key, value)
		}
		u.RawQuery = q.Encode()
		requestURL = u.String()
	}

	// Prepare request body
	var body io.Reader
	if data != nil {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return nil, errors.NewValidationError("Failed to marshal request data", map[string]interface{}{"error": err.Error()})
		}
		body = bytes.NewBuffer(jsonData)
	}

	// Create request
	req, err := http.NewRequestWithContext(ctx, method, requestURL, body)
	if err != nil {
		return nil, errors.NewValidationError("Failed to create request", map[string]interface{}{"error": err.Error()})
	}

	// Set headers
	if data != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("User-Agent", "nendb-go-driver/0.1.0")

	// Perform request with retries
	var lastErr error
	for attempt := 0; attempt <= c.config.MaxRetries; attempt++ {
		if attempt > 0 {
			time.Sleep(c.config.RetryDelay * time.Duration(attempt))
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = err
			continue
		}

		defer resp.Body.Close()

		// Read response body
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			lastErr = err
			continue
		}

		// Check response status
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return respBody, nil
		}

		// Handle error responses
		if resp.StatusCode >= 400 {
			var errorResp map[string]interface{}
			if json.Unmarshal(respBody, &errorResp) == nil {
				message := fmt.Sprintf("HTTP %d: %s", resp.StatusCode, resp.Status)
				if msg, ok := errorResp["message"].(string); ok {
					message = msg
				}
				return nil, errors.NewResponseError(message, errorResp)
			}
			return nil, errors.NewResponseError(fmt.Sprintf("HTTP %d: %s", resp.StatusCode, resp.Status), nil)
		}

		// For 3xx status codes, continue with retry
		lastErr = fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// All retries exhausted
	if lastErr != nil {
		return nil, errors.NewTimeoutError("Request failed after all retries", map[string]interface{}{"error": lastErr.Error()})
	}

	return nil, errors.NewTimeoutError("Request failed after all retries", nil)
}

// Health checks the health of the NenDB server
func (c *NenDBClient) Health() error {
	ctx, cancel := context.WithTimeout(context.Background(), c.config.Timeout)
	defer cancel()

	_, err := c.makeRequest(ctx, "GET", "/health", nil, nil)
	return err
}

// GetNode retrieves a node by ID
func (c *NenDBClient) GetNode(ctx context.Context, nodeID int) (*types.GraphNode, error) {
	endpoint := fmt.Sprintf("/nodes/%d", nodeID)
	
	respBody, err := c.makeRequest(ctx, "GET", endpoint, nil, nil)
	if err != nil {
		return nil, err
	}

	var node types.GraphNode
	if err := json.Unmarshal(respBody, &node); err != nil {
		return nil, errors.NewResponseError("Failed to parse node response", map[string]interface{}{"error": err.Error()})
	}

	return &node, nil
}

// CreateNode creates a new node
func (c *NenDBClient) CreateNode(ctx context.Context, labels []string, properties map[string]interface{}) (*types.GraphNode, error) {
	data := map[string]interface{}{
		"labels":     labels,
		"properties": properties,
	}

	respBody, err := c.makeRequest(ctx, "POST", "/nodes", data, nil)
	if err != nil {
		return nil, err
	}

	var node types.GraphNode
	if err := json.Unmarshal(respBody, &node); err != nil {
		return nil, errors.NewResponseError("Failed to parse node response", map[string]interface{}{"error": err.Error()})
	}

	return &node, nil
}

// UpdateNode updates an existing node
func (c *NenDBClient) UpdateNode(ctx context.Context, nodeID int, labels []string, properties map[string]interface{}) (*types.GraphNode, error) {
	endpoint := fmt.Sprintf("/nodes/%d", nodeID)
	data := map[string]interface{}{
		"labels":     labels,
		"properties": properties,
	}

	respBody, err := c.makeRequest(ctx, "PUT", endpoint, data, nil)
	if err != nil {
		return nil, err
	}

	var node types.GraphNode
	if err := json.Unmarshal(respBody, &node); err != nil {
		return nil, errors.NewResponseError("Failed to parse node response", map[string]interface{}{"error": err.Error()})
	}

	return &node, nil
}

// DeleteNode deletes a node by ID
func (c *NenDBClient) DeleteNode(ctx context.Context, nodeID int) error {
	endpoint := fmt.Sprintf("/nodes/%d", nodeID)
	_, err := c.makeRequest(ctx, "DELETE", endpoint, nil, nil)
	return err
}

// GetEdge retrieves an edge by ID
func (c *NenDBClient) GetEdge(ctx context.Context, edgeID int) (*types.GraphEdge, error) {
	endpoint := fmt.Sprintf("/edges/%d", edgeID)
	
	respBody, err := c.makeRequest(ctx, "GET", endpoint, nil, nil)
	if err != nil {
		return nil, err
	}

	var edge types.GraphEdge
	if err := json.Unmarshal(respBody, &edge); err != nil {
		return nil, errors.NewResponseError("Failed to parse edge response", map[string]interface{}{"error": err.Error()})
	}

	return &edge, nil
}

// CreateEdge creates a new edge
func (c *NenDBClient) CreateEdge(ctx context.Context, source, target int, edgeType string, properties map[string]interface{}) (*types.GraphEdge, error) {
	data := map[string]interface{}{
		"source":     source,
		"target":     target,
		"type":       edgeType,
		"properties": properties,
	}

	respBody, err := c.makeRequest(ctx, "POST", "/edges", data, nil)
	if err != nil {
		return nil, err
	}

	var edge types.GraphEdge
	if err := json.Unmarshal(respBody, &edge); err != nil {
		return nil, errors.NewResponseError("Failed to parse edge response", map[string]interface{}{"error": err.Error()})
	}

	return &edge, nil
}

// UpdateEdge updates an existing edge
func (c *NenDBClient) UpdateEdge(ctx context.Context, edgeID int, edgeType string, properties map[string]interface{}) (*types.GraphEdge, error) {
	endpoint := fmt.Sprintf("/edges/%d", edgeID)
	data := map[string]interface{}{
		"type":       edgeType,
		"properties": properties,
	}

	respBody, err := c.makeRequest(ctx, "PUT", endpoint, data, nil)
	if err != nil {
		return nil, err
	}

	var edge types.GraphEdge
	if err := json.Unmarshal(respBody, &edge); err != nil {
		return nil, errors.NewResponseError("Failed to parse edge response", map[string]interface{}{"error": err.Error()})
	}

	return &edge, nil
}

// DeleteEdge deletes an edge by ID
func (c *NenDBClient) DeleteEdge(ctx context.Context, edgeID int) error {
	endpoint := fmt.Sprintf("/edges/%d", edgeID)
	_, err := c.makeRequest(ctx, "DELETE", endpoint, nil, nil)
	return err
}

// RunBFS runs the BFS algorithm
func (c *NenDBClient) RunBFS(ctx context.Context, startNode, targetNode int, maxDepth int) (*types.BFSResult, error) {
	data := map[string]interface{}{
		"start_node": startNode,
		"target_node": targetNode,
		"max_depth":  maxDepth,
	}

	respBody, err := c.makeRequest(ctx, "POST", "/algorithms/bfs", data, nil)
	if err != nil {
		return nil, err
	}

	var result types.BFSResult
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, errors.NewResponseError("Failed to parse BFS result", map[string]interface{}{"error": err.Error()})
	}

	return &result, nil
}

// RunDijkstra runs the Dijkstra shortest path algorithm
func (c *NenDBClient) RunDijkstra(ctx context.Context, startNode, targetNode int) (*types.DijkstraResult, error) {
	data := map[string]interface{}{
		"start_node": startNode,
		"target_node": targetNode,
	}

	respBody, err := c.makeRequest(ctx, "POST", "/algorithms/dijkstra", data, nil)
	if err != nil {
		return nil, err
	}

	var result types.DijkstraResult
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, errors.NewResponseError("Failed to parse Dijkstra result", map[string]interface{}{"error": err.Error()})
	}

	return &result, nil
}

// RunPageRank runs the PageRank algorithm
func (c *NenDBClient) RunPageRank(ctx context.Context, maxIterations int, tolerance float64) (*types.PageRankResult, error) {
	data := map[string]interface{}{
		"max_iterations": maxIterations,
		"tolerance":      tolerance,
	}

	respBody, err := c.makeRequest(ctx, "POST", "/algorithms/pagerank", data, nil)
	if err != nil {
		return nil, err
	}

	var result types.PageRankResult
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, errors.NewResponseError("Failed to parse PageRank result", map[string]interface{}{"error": err.Error()})
	}

	return &result, nil
}

// Query executes a custom Cypher-like query
func (c *NenDBClient) Query(ctx context.Context, query string, params map[string]interface{}) (interface{}, error) {
	data := map[string]interface{}{
		"query":  query,
		"params": params,
	}

	respBody, err := c.makeRequest(ctx, "POST", "/query", data, nil)
	if err != nil {
		return nil, err
	}

	var result interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, errors.NewResponseError("Failed to parse query result", map[string]interface{}{"error": err.Error()})
	}

	return result, nil
}

// GetStatistics retrieves database statistics
func (c *NenDBClient) GetStatistics(ctx context.Context) (map[string]interface{}, error) {
	respBody, err := c.makeRequest(ctx, "GET", "/statistics", nil, nil)
	if err != nil {
		return nil, err
	}

	var stats map[string]interface{}
	if err := json.Unmarshal(respBody, &stats); err != nil {
		return nil, errors.NewResponseError("Failed to parse statistics", map[string]interface{}{"error": err.Error()})
	}

	return stats, nil
}

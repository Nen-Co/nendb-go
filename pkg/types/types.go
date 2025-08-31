package types

import (
	"fmt"
	"reflect"
)

// AlgorithmStatus represents the status of algorithm execution
type AlgorithmStatus string

const (
	StatusQueued    AlgorithmStatus = "queued"
	StatusRunning   AlgorithmStatus = "running"
	StatusCompleted AlgorithmStatus = "completed"
	StatusFailed    AlgorithmStatus = "failed"
	StatusCancelled AlgorithmStatus = "cancelled"
)

// GraphNode represents a node in the graph
type GraphNode struct {
	ID         int                    `json:"id"`
	Labels     []string               `json:"labels"`
	Properties map[string]interface{} `json:"properties"`
}

// NewGraphNode creates a new GraphNode with validation
func NewGraphNode(id int, labels []string, properties map[string]interface{}) (*GraphNode, error) {
	if id < 0 {
		return nil, fmt.Errorf("node ID must be a non-negative integer")
	}
	if labels == nil {
		labels = []string{}
	}
	if properties == nil {
		properties = make(map[string]interface{})
	}
	
	return &GraphNode{
		ID:         id,
		Labels:     labels,
		Properties: properties,
	}, nil
}

// Validate validates the GraphNode
func (n *GraphNode) Validate() error {
	if n.ID < 0 {
		return fmt.Errorf("node ID must be a non-negative integer")
	}
	if n.Labels == nil {
		return fmt.Errorf("labels cannot be nil")
	}
	if n.Properties == nil {
		return fmt.Errorf("properties cannot be nil")
	}
	return nil
}

// GraphEdge represents an edge in the graph
type GraphEdge struct {
	ID         int                    `json:"id"`
	Source     int                    `json:"source"`
	Target     int                    `json:"target"`
	Type       string                 `json:"type"`
	Properties map[string]interface{} `json:"properties"`
}

// NewGraphEdge creates a new GraphEdge with validation
func NewGraphEdge(id, source, target int, edgeType string, properties map[string]interface{}) (*GraphEdge, error) {
	if id < 0 {
		return nil, fmt.Errorf("edge ID must be a non-negative integer")
	}
	if source < 0 {
		return nil, fmt.Errorf("source node ID must be a non-negative integer")
	}
	if target < 0 {
		return nil, fmt.Errorf("target node ID must be a non-negative integer")
	}
	if edgeType == "" {
		return nil, fmt.Errorf("edge type cannot be empty")
	}
	if properties == nil {
		properties = make(map[string]interface{})
	}
	
	return &GraphEdge{
		ID:         id,
		Source:     source,
		Target:     target,
		Type:       edgeType,
		Properties: properties,
	}, nil
}

// Validate validates the GraphEdge
func (e *GraphEdge) Validate() error {
	if e.ID < 0 {
		return fmt.Errorf("edge ID must be a non-negative integer")
	}
	if e.Source < 0 {
		return fmt.Errorf("source node ID must be a non-negative integer")
	}
	if e.Target < 0 {
		return fmt.Errorf("target node ID must be a non-negative integer")
	}
	if e.Type == "" {
		return fmt.Errorf("edge type cannot be empty")
	}
	if e.Properties == nil {
		return fmt.Errorf("properties cannot be nil")
	}
	return nil
}

// AlgorithmResult represents the base result for algorithm execution
type AlgorithmResult struct {
	Algorithm string                 `json:"algorithm"`
	Status    AlgorithmStatus        `json:"status"`
	Message   string                 `json:"message"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// NewAlgorithmResult creates a new AlgorithmResult with validation
func NewAlgorithmResult(algorithm string, status AlgorithmStatus, message string, metadata map[string]interface{}) (*AlgorithmResult, error) {
	if algorithm == "" {
		return nil, fmt.Errorf("algorithm name cannot be empty")
	}
	if message == "" {
		return nil, fmt.Errorf("message cannot be empty")
	}
	if metadata == nil {
		metadata = make(map[string]interface{})
	}
	
	return &AlgorithmResult{
		Algorithm: algorithm,
		Status:    status,
		Message:   message,
		Metadata:  metadata,
	}, nil
}

// Validate validates the AlgorithmResult
func (r *AlgorithmResult) Validate() error {
	if r.Algorithm == "" {
		return fmt.Errorf("algorithm name cannot be empty")
	}
	if r.Message == "" {
		return fmt.Errorf("message cannot be empty")
	}
	if r.Metadata == nil {
		return fmt.Errorf("metadata cannot be nil")
	}
	return nil
}

// BFSResult represents the result of BFS algorithm execution
type BFSResult struct {
	*AlgorithmResult
	VisitedNodes []int `json:"visited_nodes"`
	Path         []int `json:"path"`
	Depth        int   `json:"depth"`
}

// NewBFSResult creates a new BFSResult
func NewBFSResult(base *AlgorithmResult, visitedNodes []int, path []int, depth int) *BFSResult {
	if visitedNodes == nil {
		visitedNodes = []int{}
	}
	if path == nil {
		path = []int{}
	}
	
	return &BFSResult{
		AlgorithmResult: base,
		VisitedNodes:    visitedNodes,
		Path:            path,
		Depth:           depth,
	}
}

// DijkstraResult represents the result of Dijkstra algorithm execution
type DijkstraResult struct {
	*AlgorithmResult
	ShortestPath []int                    `json:"shortest_path"`
	TotalCost    float64                  `json:"total_cost"`
	PathDetails  []map[string]interface{} `json:"path_details"`
}

// NewDijkstraResult creates a new DijkstraResult
func NewDijkstraResult(base *AlgorithmResult, shortestPath []int, totalCost float64, pathDetails []map[string]interface{}) *DijkstraResult {
	if shortestPath == nil {
		shortestPath = []int{}
	}
	if pathDetails == nil {
		pathDetails = []map[string]interface{}{}
	}
	
	return &DijkstraResult{
		AlgorithmResult: base,
		ShortestPath:    shortestPath,
		TotalCost:       totalCost,
		PathDetails:     pathDetails,
	}
}

// PageRankResult represents the result of PageRank algorithm execution
type PageRankResult struct {
	*AlgorithmResult
	NodeScores  map[int]float64 `json:"node_scores"`
	Iterations  int             `json:"iterations"`
	Convergence bool            `json:"convergence"`
}

// NewPageRankResult creates a new PageRankResult
func NewPageRankResult(base *AlgorithmResult, nodeScores map[int]float64, iterations int, convergence bool) *PageRankResult {
	if nodeScores == nil {
		nodeScores = make(map[int]float64)
	}
	
	return &PageRankResult{
		AlgorithmResult: base,
		NodeScores:      nodeScores,
		Iterations:      iterations,
		Convergence:     convergence,
	}
}

// Type aliases for convenience
type NodeID = int
type EdgeID = int
type PropertyValue = interface{}
type PropertyMap = map[string]interface{}

// IsValidPropertyValue checks if a value is a valid property value
func IsValidPropertyValue(value interface{}) bool {
	if value == nil {
		return true
	}
	
	switch reflect.TypeOf(value).Kind() {
	case reflect.String, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		 reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		 reflect.Float32, reflect.Float64, reflect.Bool:
		return true
	default:
		return false
	}
}

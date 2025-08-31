# NenDB Go Driver

A high-performance Go client for the NenDB graph database, built with the same design principles as the Python driver but optimized for Go applications.

## Features

- **High Performance**: Built with Go's efficient HTTP client and JSON handling
- **Full API Coverage**: Complete support for all NenDB operations including nodes, edges, and algorithms
- **Robust Error Handling**: Comprehensive error types with detailed context
- **Retry Logic**: Built-in retry mechanism with configurable backoff
- **Context Support**: Full context.Context support for timeouts and cancellation
- **Type Safety**: Strong typing with validation for all graph entities
- **CLI Tool**: Command-line interface for testing and administration

## Installation

```bash
go get github.com/nen-co/nendb-go-driver
```

## Quick Start

### Basic Usage

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/nen-co/nendb-go-driver/pkg/client"
)

func main() {
    // Create client configuration
    config := &client.ClientConfig{
        BaseURL:    "http://localhost:8080",
        Timeout:    30 * time.Second,
        MaxRetries: 3,
    }

    // Create client
    client, err := client.NewClient(config)
    if err != nil {
        log.Fatalf("Failed to create client: %v", err)
    }

    // Create context with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
    defer cancel()

    // Check server health
    if err := client.Health(); err != nil {
        log.Fatalf("Health check failed: %v", err)
    }

    // Create a node
    node, err := client.CreateNode(ctx, []string{"Person"}, map[string]interface{}{
        "name": "Alice",
        "age":  30,
    })
    if err != nil {
        log.Fatalf("Failed to create node: %v", err)
    }

    fmt.Printf("Created node with ID: %d\n", node.ID)
}
```

### Working with Nodes

```go
// Create a node with labels and properties
node, err := client.CreateNode(ctx, []string{"Person", "Employee"}, map[string]interface{}{
    "name":     "Bob",
    "age":      25,
    "position": "Developer",
    "salary":   75000,
})

// Get a node by ID
retrievedNode, err := client.GetNode(ctx, node.ID)

// Update a node
updatedNode, err := client.UpdateNode(ctx, node.ID, []string{"Person", "Manager"}, map[string]interface{}{
    "name":     "Bob",
    "age":      26,
    "position": "Senior Developer",
    "salary":   85000,
})

// Delete a node
err = client.DeleteNode(ctx, node.ID)
```

### Working with Edges

```go
// Create an edge between two nodes
edge, err := client.CreateEdge(ctx, sourceNodeID, targetNodeID, "KNOWS", map[string]interface{}{
    "since":    "2022-01-15",
    "strength": "strong",
})

// Get an edge by ID
retrievedEdge, err := client.GetEdge(ctx, edge.ID)

// Update an edge
updatedEdge, err := client.UpdateEdge(ctx, edge.ID, "KNOWS", map[string]interface{}{
    "since":    "2022-01-15",
    "strength": "very strong",
    "notes":    "Close friends",
})

// Delete an edge
err = client.DeleteEdge(ctx, edge.ID)
```

### Running Algorithms

```go
// Run BFS algorithm
bfsResult, err := client.RunBFS(ctx, startNodeID, targetNodeID, 5)
if err != nil {
    log.Printf("BFS failed: %v", err)
} else {
    fmt.Printf("BFS visited %d nodes, path length: %d\n", 
        len(bfsResult.VisitedNodes), len(bfsResult.Path))
}

// Run Dijkstra shortest path
dijkstraResult, err := client.RunDijkstra(ctx, startNodeID, targetNodeID)
if err != nil {
    log.Printf("Dijkstra failed: %v", err)
} else {
    fmt.Printf("Shortest path cost: %f\n", dijkstraResult.TotalCost)
}

// Run PageRank
pagerankResult, err := client.RunPageRank(ctx, 100, 0.001)
if err != nil {
    log.Printf("PageRank failed: %v", err)
} else {
    fmt.Printf("PageRank completed in %d iterations\n", pagerankResult.Iterations)
}
```

### Custom Queries

```go
// Execute a custom Cypher-like query
result, err := client.Query(ctx, "MATCH (n:Person) WHERE n.age > 25 RETURN n LIMIT 10", nil)
if err != nil {
    log.Printf("Query failed: %v", err)
} else {
    fmt.Printf("Query returned: %+v\n", result)
}

// Query with parameters
params := map[string]interface{}{
    "minAge": 25,
    "limit":  10,
}
result, err = client.Query(ctx, "MATCH (n:Person) WHERE n.age > $minAge RETURN n LIMIT $limit", params)
```

## CLI Usage

The driver includes a command-line interface for testing and administration:

```bash
# Check server health
nendb -command health

# Get a node by ID
nendb -command node 1

# Get an edge by ID
nendb -command edge 1

# Run BFS algorithm
nendb -command algorithm bfs 1 5 3

# Run Dijkstra algorithm
nendb -command algorithm dijkstra 1 5

# Run PageRank algorithm
nendb -command algorithm pagerank 100 0.001

# Execute custom query
nendb -command query "MATCH (n) RETURN n LIMIT 5"

# Get database statistics
nendb -command stats

# Use custom server URL
nendb -url http://localhost:9090 -command health

# Skip health check on startup
nendb -skip-health -command health
```

## Configuration

### ClientConfig Options

- **BaseURL**: NenDB server base URL (default: "http://localhost:8080")
- **Timeout**: Request timeout (default: 30s)
- **MaxRetries**: Maximum number of retries (default: 3)
- **RetryDelay**: Delay between retries (default: 1s)
- **SkipValidation**: Skip health check on startup (default: false)
- **HTTPClient**: Custom HTTP client (optional)

### Environment Variables

You can also configure the client using environment variables:

```bash
export NENDB_URL=http://localhost:8080
export NENDB_TIMEOUT=60s
export NENDB_MAX_RETRIES=5
```

## Error Handling

The driver provides comprehensive error types:

```go
import "github.com/nen-co/nendb-go-driver/pkg/errors"

// Check error types
if connErr, ok := err.(*errors.NenDBConnectionError); ok {
    log.Printf("Connection error: %v", connErr)
} else if timeoutErr, ok := err.(*errors.NenDBTimeoutError); ok {
    log.Printf("Timeout error: %v", timeoutErr)
} else if validationErr, ok := err.(*errors.NenDBValidationError); ok {
    log.Printf("Validation error: %v", validationErr)
} else if algoErr, ok := err.(*errors.NenDBAlgorithmError); ok {
    log.Printf("Algorithm error: %v", algoErr)
} else if respErr, ok := err.(*errors.NenDBResponseError); ok {
    log.Printf("Response error: %v", respErr)
}
```

## Testing

Run the test suite:

```bash
go test ./...
```

Run tests with verbose output:

```bash
go test -v ./...
```

Run tests with coverage:

```bash
go test -cover ./...
```

## Examples

See the `examples/` directory for complete working examples:

- `basic_usage.go` - Basic client operations
- Additional examples coming soon

## Performance Considerations

- **Connection Pooling**: The client uses Go's built-in HTTP connection pooling
- **JSON Marshaling**: Efficient JSON handling with Go's standard library
- **Memory Management**: Minimal allocations in hot paths
- **Context Usage**: Proper context handling for timeouts and cancellation

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass
6. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

- **Documentation**: [https://docs.nen.co](https://docs.nen.co)
- **Issues**: [GitHub Issues](https://github.com/nen-co/nendb-go-driver/issues)
- **Discussions**: [GitHub Discussions](https://github.com/nen-co/nendb-go-driver/discussions)

## Version History

- **v0.1.0** - Initial release with core functionality
  - Full CRUD operations for nodes and edges
  - Algorithm support (BFS, Dijkstra, PageRank)
  - Custom query execution
  - CLI tool
  - Comprehensive error handling
  - Full test coverage

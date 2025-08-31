# NenDB Go Driver

[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org)
[![Go Module](https://img.shields.io/badge/Go%20Module-v0.1.0-green.svg)](https://pkg.go.dev/github.com/nen-co/nendb-go-driver)
[![License](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/nen-co/nendb-go-driver)](https://goreportcard.com/report/github.com/nen-co/nendb-go-driver)

A high-performance Go client for the NenDB graph database, built with the same design principles as the Python driver but optimized for Go applications.

> **🚀 Officially published as a Go module** - Available via `go get github.com/nen-co/nendb-go-driver`

## Features

- **High Performance**: Built with Go's efficient HTTP client and JSON handling
- **Full API Coverage**: Complete support for all NenDB operations including nodes, edges, and algorithms
- **Robust Error Handling**: Comprehensive error types with detailed context
- **Retry Logic**: Built-in retry mechanism with configurable backoff
- **Context Support**: Full context.Context support for timeouts and cancellation
- **Type Safety**: Strong typing with validation for all graph entities
- **CLI Tool**: Command-line interface for testing and administration

## Installation

The NenDB Go Driver is now officially published as a Go module and available for installation:

```bash
# Get the latest version
go get github.com/nen-co/nendb-go-driver

# Get a specific version
go get github.com/nen-co/nendb-go-driver@v0.1.0
```

### Module Information

- **Module Path**: `github.com/nen-co/nendb-go-driver`
- **Latest Version**: `v0.1.0`
- **Go Version**: 1.21+
- **Repository**: [https://github.com/Nen-Co/nendb-go.git](https://github.com/Nen-Co/nendb-go.git)

## Quick Start

### Prerequisites

Before using the Go driver, you need to have the NenDB server running. The NenDB server is built in Zig and provides the HTTP API that the Go driver connects to.

#### Running NenDB Server

1. **Clone the NenDB repository**:
   ```bash
   git clone https://github.com/Nen-Co/nen-db.git
   cd nen-db
   ```

2. **Build the server**:
   ```bash
   zig build
   ```

3. **Run the server**:
   ```bash
   # Run with default configuration (port 8080)
   ./zig-out/bin/nendb
   
   # Or run the server directly
   zig build run
   ```

4. **Verify server is running**:
   ```bash
   curl http://localhost:8080/health
   # Should return: {"status": "healthy", "service": "nendb", "version": "0.0.1"}
   ```

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
        BaseURL:    "http://localhost:8080", // NenDB server address
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

## NenDB Server Integration

The Go driver connects to the NenDB server, which is built in Zig and provides a high-performance HTTP API for graph database operations.

### Server Architecture

- **Language**: Built in Zig for maximum performance
- **HTTP API**: RESTful endpoints for all operations
- **Port**: Default 8080 (configurable)
- **Memory**: Statically allocated with efficient memory management
- **Networking**: Uses custom nen-net library for high-performance I/O

### Available API Endpoints

The NenDB server provides these endpoints that the Go driver uses:

#### Health & Status
- `GET /health` - Server health check
- `GET /statistics` - Database statistics

#### Graph Operations
- `GET /nodes/{id}` - Retrieve node by ID
- `POST /nodes` - Create new node
- `PUT /nodes/{id}` - Update existing node
- `DELETE /nodes/{id}` - Delete node

- `GET /edges/{id}` - Retrieve edge by ID
- `POST /edges` - Create new edge
- `PUT /edges/{id}` - Update existing edge
- `DELETE /edges/{id}` - Delete edge

#### Algorithms
- `POST /algorithms/bfs` - Breadth-First Search
- `POST /algorithms/dijkstra` - Shortest Path (Dijkstra)
- `POST /algorithms/pagerank` - PageRank algorithm

#### Query
- `POST /query` - Execute custom Cypher-like queries

### Server Configuration

The NenDB server can be configured with various options:

```bash
# Run with custom port
zig build run -- --port 9090

# Run with custom host
zig build run -- --host 127.0.0.1

# Run with custom buffer size
zig build run -- --buffer-size 16384
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

### Development Workflow

When developing with the Go driver and NenDB server:

1. **Start the NenDB server** in one terminal:
   ```bash
   cd nen-db
   zig build run
   ```

2. **Run your Go application** in another terminal:
   ```bash
   cd your-go-project
   go run main.go
   ```

3. **Test the integration**:
   ```bash
   # Test server health
   curl http://localhost:8080/health
   
   # Test Go driver
   go test ./...
   ```

### Testing with NenDB Server

The Go driver includes tests that can run against a live NenDB server:

```bash
# Set environment variable to skip health check in tests
export NENDB_SKIP_HEALTH=true

# Run tests
go test ./...

# Or run specific package tests
go test ./pkg/client
go test ./pkg/types
go test ./pkg/errors
```

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

## Go Module Ecosystem

This driver is fully integrated into the Go module ecosystem:

### Module Discovery
- **Go.dev**: [github.com/nen-co/nendb-go-driver](https://pkg.go.dev/github.com/nen-co/nendb-go-driver)
- **Go Modules**: Available via `go get` command
- **Proxy Support**: Compatible with Go module proxies

### Version Management
```bash
# Check available versions
go list -m -versions github.com/nen-co/nendb-go-driver

# Update to latest version
go get -u github.com/nen-co/nendb-go-driver

# Pin to specific version
go get github.com/nen-co/nendb-go-driver@v0.1.0
```

### Go Workspace Support
```bash
# Add to go.work file
go work use ./path/to/nendb-go-driver

# Or use directly in projects
go mod edit -require=github.com/nen-co/nendb-go-driver@v0.1.0
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

- **Documentation**: [https://docs.nen.co](https://docs.nen.co)
- **Issues**: [GitHub Issues](https://github.com/Nen-Co/nendb-go/issues)
- **Discussions**: [GitHub Discussions](https://github.com/Nen-Co/nendb-go/discussions)
- **Go Module**: [pkg.go.dev](https://pkg.go.dev/github.com/nen-co/nendb-go-driver)

## Official Go Module

This driver follows the official Go module publishing workflow and is available through:

- **Go Module Registry**: `github.com/nen-co/nendb-go-driver`
- **Version Control**: Git tags for each release
- **Go Tools**: Full support for `go get`, `go mod`, and `go work`
- **Proxy Compatibility**: Works with all Go module proxies

## Module Publishing

The NenDB Go Driver is officially published as a Go module and follows the standard Go module publishing workflow:

### Publishing Status
- ✅ **Module Published**: Available at `github.com/nen-co/nendb-go-driver`
- ✅ **Version Tagged**: `v0.1.0` released and tagged
- ✅ **Go Module Index**: Registered with Go's module system
- ✅ **Repository**: [GitHub Repository](https://github.com/Nen-Co/nendb-go.git)

### Installation for Users
```bash
# Install the latest version
go get github.com/nen-co/nendb-go-driver

# Install specific version
go get github.com/nen-co/nendb-go-driver@v0.1.0

# Import in your Go code
import "github.com/nen-co/nendb-go-driver/pkg/client"
```

## Version History

- **v0.1.0** - Initial release with core functionality
  - Full CRUD operations for nodes and edges
  - Algorithm support (BFS, Dijkstra, PageRank)
  - Custom query execution
  - CLI tool
  - Comprehensive error handling
  - Full test coverage
  - **Officially published as Go module**

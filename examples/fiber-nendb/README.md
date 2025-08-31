# 🍳 Fiber + NenDB Recipe

This is a cookbook recipe for setting up a high-performance GraphQL-like API using Fiber web framework and NenDB graph database. 🚀

## 🌟 Features

- **High Performance**: Built with Fiber for fast HTTP handling and NenDB for efficient graph operations
- **Full CRUD Operations**: Complete node and edge management
- **Graph Algorithms**: BFS, Dijkstra, and PageRank implementations
- **Custom Queries**: Execute custom Cypher-like queries
- **RESTful API**: Clean, intuitive endpoints for all operations
- **Error Handling**: Comprehensive error handling with proper HTTP status codes
- **Middleware**: CORS, logging, and recovery middleware included

## 📋 Prerequisites

Go is an obvious prerequisite. Make sure it is installed and configured properly.

After that you need two Go packages: Fiber and NenDB Go driver. You can install them with the following commands:

```bash
go get -u github.com/gofiber/fiber/v2
go get github.com/nen-co/nendb-go
```

## 🚀 Run NenDB

The easiest way to run NenDB is to use the prebuilt static binary. Once you have it on your machine, you can run NenDB with the following command:

### Linux/macOS (Quick Install)
```bash
curl -fsSL https://github.com/Nen-Co/nen-db/releases/latest/download/nen-linux-x86_64.tar.gz | tar -xz
./nen-linux-x86_64
```

### Windows PowerShell
```powershell
Invoke-WebRequest -Uri "https://github.com/Nen-Co/nen-db/releases/latest/download/nen-windows-x86_64.zip" -OutFile "nen-windows.zip"
Expand-Archive -Path "nen-windows.zip" -DestinationPath "."
```

### Docker (Optional)
```bash
docker run --rm -p 9000:9000 --name nendb \
  -v $(pwd)/data:/var/lib/nendb \
  nenco/nendb:latest ./nendb serve
```

> **📖 For complete installation details, see the [Official NenDB Documentation](https://nen-co.github.io/docs/nendb)**

## 🏃‍♂️ Run the Recipe

After you have installed all the prerequisites, you can run the recipe with the following command:

```bash
cd examples/fiber-nendb
go run ./main.go
```

This will do the following:

1. **Connect Fiber backend to NenDB database**
2. **Start HTTP server on port 3000**
3. **Provide comprehensive REST API endpoints**
4. **Enable real-time graph operations**

## 🧪 Test the Recipe

Once the Fiber app is running, you can test the recipe by sending requests to the following endpoints:

### 🏠 Home & Documentation
- `GET http://localhost:3000/` - API overview and endpoint documentation

### 📊 Graph Operations
- `GET http://localhost:3000/graph` - Get graph structure and statistics
- `GET http://localhost:3000/stats` - Get detailed graph statistics

### 🔗 Node Operations
- `GET http://localhost:3000/nodes` - Get all nodes (placeholder)
- `GET http://localhost:3000/nodes/:id` - Get specific node by ID
- `POST http://localhost:3000/nodes` - Create new node
- `PUT http://localhost:3000/nodes/:id` - Update existing node
- `DELETE http://localhost:3000/nodes/:id` - Delete node

### 🔗 Edge Operations
- `GET http://localhost:3000/edges` - Get all edges (placeholder)
- `GET http://localhost:3000/edges/:id` - Get specific edge by ID
- `POST http://localhost:3000/edges` - Create new edge
- `PUT http://localhost:3000/edges/:id` - Update existing edge
- `DELETE http://localhost:3000/edges/:id` - Delete edge

### 🧮 Algorithm Operations
- `POST http://localhost:3000/algorithms/bfs` - Run Breadth-First Search
- `POST http://localhost:3000/algorithms/dijkstra` - Run Dijkstra Shortest Path
- `POST http://localhost:3000/algorithms/pagerank` - Run PageRank Algorithm

### 🔍 Query Operations
- `POST http://localhost:3000/query` - Execute custom Cypher-like queries

## 📝 Example Requests

### Create a Person Node
```bash
curl -X POST http://localhost:3000/nodes \
  -H "Content-Type: application/json" \
  -d '{
    "labels": ["Person", "Employee"],
    "properties": {
      "name": "Alice Johnson",
      "age": 30,
      "position": "Software Engineer",
      "department": "Engineering"
    }
  }'
```

### Create a Relationship
```bash
curl -X POST http://localhost:3000/edges \
  -H "Content-Type: application/json" \
  -d '{
    "source_id": 1,
    "target_id": 2,
    "type": "WORKS_WITH",
    "properties": {
      "since": "2023-01-15",
      "project": "Graph Database API"
    }
  }'
```

### Run BFS Algorithm
```bash
curl -X POST http://localhost:3000/algorithms/bfs \
  -H "Content-Type: application/json" \
  -d '{
    "start_node": 1,
    "max_depth": 3
  }'
```

### Execute Custom Query
```bash
curl -X POST http://localhost:3000/query \
  -H "Content-Type: application/json" \
  -d '{
    "query": "MATCH (n:Person) WHERE n.age > 25 RETURN n LIMIT 10",
    "params": {}
  }'
```

## 🏗️ Architecture

```
┌─────────────────┐    HTTP/JSON    ┌─────────────────┐
│   Fiber App     │ ──────────────► │  NenDB Server   │
│                 │                 │   (Zig + HTTP)  │
│  REST API      │ ◄────────────── │                 │
│  Endpoints     │                 │                 │
└─────────────────┘                 └─────────────────┘
```

The Fiber application provides a RESTful HTTP API that communicates with the NenDB server, which is built in Zig for maximum performance. This gives you the best of both worlds: Fiber's ease of use with NenDB's performance.

## 🔧 Configuration

### NenDB Server Configuration
- **Default Port**: 8080 (static binary) or 9000 (Docker)
- **Base URL**: `http://localhost:8080` (configurable in the code)
- **Timeout**: 30 seconds for most operations, 60 seconds for algorithms
- **Retries**: 3 attempts with exponential backoff

### Fiber App Configuration
- **Port**: 3000
- **Middleware**: CORS, logging, recovery
- **Error Handling**: Custom error handler with proper HTTP status codes
- **Context**: Timeout management for all database operations

## 🚀 Performance Features

- **Static Memory**: NenDB uses zero dynamic allocations
- **HTTP/2 Ready**: Fiber supports modern HTTP protocols
- **Connection Pooling**: Efficient database connection management
- **Context Timeouts**: Prevents hanging requests
- **Error Recovery**: Graceful error handling and logging

## 🧪 Testing

### Manual Testing
```bash
# Test server health
curl http://localhost:3000/

# Test graph statistics
curl http://localhost:3000/stats

# Test node creation
curl -X POST http://localhost:3000/nodes \
  -H "Content-Type: application/json" \
  -d '{"labels": ["Test"], "properties": {"name": "Test Node"}}'
```

### Automated Testing
```bash
# Run tests
go test ./...

# Run with coverage
go test -cover ./...
```

## 🔍 Troubleshooting

### Common Issues

1. **Connection Refused**: Make sure NenDB server is running
   ```bash
   curl http://localhost:8080/health
   ```

2. **Port Already in Use**: Change the port in the code or stop conflicting services
   ```go
   app.Listen(":3001") // Change port
   ```

3. **Database Errors**: Check NenDB server logs and health endpoint

4. **CORS Issues**: The app includes CORS middleware, but you can customize it if needed

### Debug Mode
```go
// Enable debug logging
app.Use(logger.New(logger.Config{
    Format: "[${time}] ${status} - ${latency} ${method} ${path}\n",
}))
```

## 📚 Additional Resources

For extra information use the documentation on the following links:

- **Fiber**: [https://docs.gofiber.io/](https://docs.gofiber.io/)
- **NenDB**: [https://nen-co.github.io/docs/nendb](https://nen-co.github.io/docs/nendb)
- **NenDB Go Driver**: [https://github.com/Nen-Co/nendb-go](https://github.com/Nen-Co/nendb-go)
- **Go Modules**: [https://go.dev/doc/modules](https://go.dev/doc/modules)

## 🌟 Contributing

If you have found an amazing recipe for Fiber + NenDB — share it with others! We are ready to accept your PR and add your recipe to the cookbook.

## ⭐ Star Us

🌟 If you like this recipe, don't forget to give us a star on Github 🌟

- **NenDB**: [https://github.com/Nen-Co/nen-db](https://github.com/Nen-Co/nen-db)
- **NenDB Go Driver**: [https://github.com/Nen-Co/nendb-go](https://github.com/Nen-Co/nendb-go)
- **Fiber**: [https://github.com/gofiber/fiber](https://github.com/gofiber/fiber)

---

**Happy coding! 🚀🍳**

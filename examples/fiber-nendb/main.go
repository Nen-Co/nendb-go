package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/nen-co/nendb-go/pkg/client"
)

// Recipe: Fiber + NenDB Integration
// This recipe demonstrates how to build a high-performance GraphQL-like API
// using Fiber web framework and NenDB graph database.

func main() {
	// Initialize NenDB client
	nendbClient, err := client.NewClient(&client.ClientConfig{
		BaseURL:    "http://localhost:8080", // NenDB server address
		Timeout:    30 * time.Second,
		MaxRetries: 3,
	})
	if err != nil {
		log.Fatalf("Failed to create NenDB client: %v", err)
	}

	// Check NenDB server health
	if err := nendbClient.Health(); err != nil {
		log.Fatalf("NenDB server health check failed: %v", err)
	}
	log.Println("✅ Connected to NenDB server successfully!")

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		AppName: "NenDB Fiber Recipe",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error":   err.Error(),
				"success": false,
			})
		},
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New())

	// Routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "🍳 NenDB + Fiber Recipe",
			"version": "1.0.0",
			"endpoints": fiber.Map{
				"GET  /graph":           "Get entire graph structure",
				"GET  /nodes":           "Get all nodes",
				"GET  /nodes/:id":       "Get node by ID",
				"POST /nodes":           "Create new node",
				"PUT  /nodes/:id":       "Update node",
				"DELETE /nodes/:id":     "Delete node",
				"GET  /edges":           "Get all edges",
				"GET  /edges/:id":       "Get edge by ID",
				"POST /edges":           "Create new edge",
				"PUT  /edges/:id":       "Update edge",
				"DELETE /edges/:id":     "Delete edge",
				"POST /algorithms/bfs":  "Run BFS algorithm",
				"POST /algorithms/dijkstra": "Run Dijkstra algorithm",
				"POST /algorithms/pagerank": "Run PageRank algorithm",
				"POST /query":           "Execute custom query",
				"GET  /stats":           "Get graph statistics",
			},
		})
	})

	// Graph operations
	app.Get("/graph", func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Get graph statistics
		stats, err := nendbClient.GetStatistics(ctx)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to get graph statistics",
			})
		}

		return c.JSON(fiber.Map{
			"success": true,
			"data": fiber.Map{
				"statistics": stats,
				"message":    "Graph structure retrieved successfully",
			},
		})
	})

	// Node operations
	app.Get("/nodes", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success": true,
			"message": "Get all nodes endpoint - implement pagination for large graphs",
		})
	})

	app.Get("/nodes/:id", func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		nodeID := c.Params("id")
		if nodeID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Node ID is required",
			})
		}

		// Convert string ID to int (you might want to add validation)
		var id int
		if _, err := fmt.Sscanf(nodeID, "%d", &id); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid node ID format",
			})
		}

		node, err := nendbClient.GetNode(ctx, id)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Node not found",
			})
		}

		return c.JSON(fiber.Map{
			"success": true,
			"data":    node,
		})
	})

	app.Post("/nodes", func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		var request struct {
			Labels    []string               `json:"labels"`
			Properties map[string]interface{} `json:"properties"`
		}

		if err := c.BodyParser(&request); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		if len(request.Labels) == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "At least one label is required",
			})
		}

		node, err := nendbClient.CreateNode(ctx, request.Labels, request.Properties)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to create node",
			})
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"success": true,
			"data":    node,
			"message": "Node created successfully",
		})
	})

	app.Put("/nodes/:id", func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		nodeID := c.Params("id")
		if nodeID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Node ID is required",
			})
		}

		// Convert string ID to int (you might want to add validation)
		var id int
		if _, err := fmt.Sscanf(nodeID, "%d", &id); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid node ID format",
			})
		}

		var request struct {
			Labels    []string               `json:"labels"`
			Properties map[string]interface{} `json:"properties"`
		}

		if err := c.BodyParser(&request); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		node, err := nendbClient.UpdateNode(ctx, id, request.Labels, request.Properties)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to update node",
			})
		}

		return c.JSON(fiber.Map{
			"success": true,
			"data":    node,
			"message": "Node updated successfully",
		})
	})

	app.Delete("/nodes/:id", func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		nodeID := c.Params("id")
		if nodeID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Node ID is required",
			})
		}

		var id int
		if _, err := fmt.Sscanf(nodeID, "%d", &id); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid node ID format",
			})
		}

		if err := nendbClient.DeleteNode(ctx, id); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to delete node",
			})
		}

		return c.JSON(fiber.Map{
			"success": true,
			"message": "Node deleted successfully",
		})
	})

	// Edge operations
	app.Get("/edges", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success": true,
			"message": "Get all edges endpoint - implement pagination for large graphs",
		})
	})

	app.Get("/edges/:id", func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		edgeID := c.Params("id")
		if edgeID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Edge ID is required",
			})
		}

		var id int
		if _, err := fmt.Sscanf(edgeID, "%d", &id); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid edge ID format",
			})
		}

		edge, err := nendbClient.GetEdge(ctx, id)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Edge not found",
			})
		}

		return c.JSON(fiber.Map{
			"success": true,
			"data":    edge,
		})
	})

	app.Post("/edges", func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		var request struct {
			SourceID   int                    `json:"source_id"`
			TargetID   int                    `json:"target_id"`
			Type       string                 `json:"type"`
			Properties map[string]interface{} `json:"properties"`
		}

		if err := c.BodyParser(&request); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		if request.Type == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Edge type is required",
			})
		}

		edge, err := nendbClient.CreateEdge(ctx, request.SourceID, request.TargetID, request.Type, request.Properties)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to create edge",
			})
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"success": true,
			"data":    edge,
			"message": "Edge created successfully",
		})
	})

	// Algorithm endpoints
	app.Post("/algorithms/bfs", func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		var request struct {
			StartNode int `json:"start_node"`
			MaxDepth  int `json:"max_depth"`
		}

		if err := c.BodyParser(&request); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		result, err := nendbClient.RunBFS(ctx, request.StartNode, 0, request.MaxDepth)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "BFS algorithm failed",
			})
		}

		return c.JSON(fiber.Map{
			"success": true,
			"data":    result,
			"message": "BFS algorithm completed successfully",
		})
	})

	app.Post("/algorithms/dijkstra", func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		var request struct {
			StartNode int `json:"start_node"`
			EndNode   int `json:"end_node"`
		}

		if err := c.BodyParser(&request); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		result, err := nendbClient.RunDijkstra(ctx, request.StartNode, request.EndNode)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Dijkstra algorithm failed",
			})
		}

		return c.JSON(fiber.Map{
			"success": true,
			"data":    result,
			"message": "Dijkstra algorithm completed successfully",
		})
	})

	app.Post("/algorithms/pagerank", func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		var request struct {
			Iterations     int     `json:"iterations"`
			DampingFactor  float64 `json:"damping_factor"`
		}

		if err := c.BodyParser(&request); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		// Set defaults if not provided
		if request.Iterations == 0 {
			request.Iterations = 100
		}
		if request.DampingFactor == 0 {
			request.DampingFactor = 0.85
		}

		result, err := nendbClient.RunPageRank(ctx, request.Iterations, request.DampingFactor)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "PageRank algorithm failed",
			})
		}

		return c.JSON(fiber.Map{
			"success": true,
			"data":    result,
			"message": "PageRank algorithm completed successfully",
		})
	})

	// Custom query endpoint
	app.Post("/query", func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		var request struct {
			Query  string                 `json:"query"`
			Params map[string]interface{} `json:"params"`
		}

		if err := c.BodyParser(&request); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		if request.Query == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Query is required",
			})
		}

		result, err := nendbClient.Query(ctx, request.Query, request.Params)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Query execution failed",
			})
		}

		return c.JSON(fiber.Map{
			"success": true,
			"data":    result,
			"message": "Query executed successfully",
		})
	})

	// Statistics endpoint
	app.Get("/stats", func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		stats, err := nendbClient.GetStatistics(ctx)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to get graph statistics",
			})
		}

		return c.JSON(fiber.Map{
			"success": true,
			"data":    stats,
		})
	})

	// Start the server
	log.Println("🚀 Starting Fiber + NenDB recipe server on :3000")
	log.Println("📖 API documentation available at http://localhost:3000")
	log.Fatal(app.Listen(":3000"))
}

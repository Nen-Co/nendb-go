package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/nen-co/nendb-go-driver/pkg/client"
)

const version = "0.1.0"

func main() {
	// Parse command line flags
	var (
		baseURL    = flag.String("url", "http://localhost:8080", "NenDB server base URL")
		timeout    = flag.Duration("timeout", 30*time.Second, "Request timeout")
		maxRetries = flag.Int("retries", 3, "Maximum number of retries")
		skipHealth = flag.Bool("skip-health", false, "Skip health check on startup")
		command    = flag.String("command", "", "Command to execute (health, node, edge, algorithm, query, stats)")
		help       = flag.Bool("help", false, "Show help")
		showVer    = flag.Bool("version", false, "Show version")
	)
	flag.Parse()

	// Show version
	if *showVer {
		fmt.Printf("nendb-go-driver version %s\n", version)
		os.Exit(0)
	}

	// Show help
	if *help || flag.NFlag() == 0 {
		showHelp()
		os.Exit(0)
	}

	// Create client configuration
	config := &client.ClientConfig{
		BaseURL:        *baseURL,
		Timeout:        *timeout,
		MaxRetries:     *maxRetries,
		SkipValidation: *skipHealth,
	}

	// Create client
	client, err := client.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Execute command
	if err := executeCommand(client, *command, flag.Args()); err != nil {
		log.Fatalf("Command failed: %v", err)
	}
}

func showHelp() {
	fmt.Printf(`NenDB Go Driver v%s

Usage: nendb [flags] -command <command> [args...]

Flags:
  -url string        NenDB server base URL (default "http://localhost:8080")
  -timeout duration  Request timeout (default 30s)
  -retries int       Maximum number of retries (default 3)
  -skip-health       Skip health check on startup
  -help              Show this help message
  -version           Show version

Commands:
  health             Check server health
  node <id>          Get node by ID
  edge <id>          Get edge by ID
  algorithm <type>   Run algorithm (bfs, dijkstra, pagerank)
  query <query>      Execute custom query
  stats              Get database statistics

Examples:
  nendb -command health
  nendb -command node 1
  nendb -command algorithm bfs -url http://localhost:9090
  nendb -command query "MATCH (n) RETURN n LIMIT 5"
`, version)
}

func executeCommand(client *client.NenDBClient, command string, args []string) error {
	ctx := context.Background()

	switch command {
	case "health":
		return executeHealth(client, ctx)
	case "node":
		if len(args) < 1 {
			return fmt.Errorf("node command requires an ID")
		}
		return executeGetNode(client, ctx, args[0])
	case "edge":
		if len(args) < 1 {
			return fmt.Errorf("edge command requires an ID")
		}
		return executeGetEdge(client, ctx, args[0])
	case "algorithm":
		if len(args) < 1 {
			return fmt.Errorf("algorithm command requires a type")
		}
		return executeAlgorithm(client, ctx, args[0], args[1:])
	case "query":
		if len(args) < 1 {
			return fmt.Errorf("query command requires a query string")
		}
		return executeQuery(client, ctx, args[0])
	case "stats":
		return executeStats(client, ctx)
	default:
		return fmt.Errorf("unknown command: %s", command)
	}
}

func executeHealth(client *client.NenDBClient, ctx context.Context) error {
	fmt.Println("Checking NenDB server health...")
	if err := client.Health(); err != nil {
		return fmt.Errorf("health check failed: %v", err)
	}
	fmt.Println("✓ Server is healthy")
	return nil
}

func executeGetNode(client *client.NenDBClient, ctx context.Context, nodeIDStr string) error {
	var nodeID int
	if _, err := fmt.Sscanf(nodeIDStr, "%d", &nodeID); err != nil {
		return fmt.Errorf("invalid node ID: %s", nodeIDStr)
	}

	fmt.Printf("Getting node %d...\n", nodeID)
	node, err := client.GetNode(ctx, nodeID)
	if err != nil {
		return fmt.Errorf("failed to get node: %v", err)
	}

	output, err := json.MarshalIndent(node, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal node: %v", err)
	}
	fmt.Println(string(output))
	return nil
}

func executeGetEdge(client *client.NenDBClient, ctx context.Context, edgeIDStr string) error {
	var edgeID int
	if _, err := fmt.Sscanf(edgeIDStr, "%d", &edgeID); err != nil {
		return fmt.Errorf("invalid edge ID: %s", edgeIDStr)
	}

	fmt.Printf("Getting edge %d...\n", edgeID)
	edge, err := client.GetEdge(ctx, edgeID)
	if err != nil {
		return fmt.Errorf("failed to get edge: %v", err)
	}

	output, err := json.MarshalIndent(edge, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal edge: %v", err)
	}
	fmt.Println(string(output))
	return nil
}

func executeAlgorithm(client *client.NenDBClient, ctx context.Context, algoType string, args []string) error {
	fmt.Printf("Running %s algorithm...\n", algoType)

	switch algoType {
	case "bfs":
		if len(args) < 2 {
			return fmt.Errorf("bfs algorithm requires start and target node IDs")
		}
		var startNode, targetNode int
		if _, err := fmt.Sscanf(args[0], "%d", &startNode); err != nil {
			return fmt.Errorf("invalid start node ID: %s", args[0])
		}
		if _, err := fmt.Sscanf(args[1], "%d", &targetNode); err != nil {
			return fmt.Errorf("invalid target node ID: %s", args[1])
		}

		maxDepth := 10
		if len(args) > 2 {
			if _, err := fmt.Sscanf(args[2], "%d", &maxDepth); err != nil {
				return fmt.Errorf("invalid max depth: %s", args[2])
			}
		}

		result, err := client.RunBFS(ctx, startNode, targetNode, maxDepth)
		if err != nil {
			return fmt.Errorf("bfs algorithm failed: %v", err)
		}

		output, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal result: %v", err)
		}
		fmt.Println(string(output))

	case "dijkstra":
		if len(args) < 2 {
			return fmt.Errorf("dijkstra algorithm requires start and target node IDs")
		}
		var startNode, targetNode int
		if _, err := fmt.Sscanf(args[0], "%d", &startNode); err != nil {
			return fmt.Errorf("invalid start node ID: %s", args[0])
		}
		if _, err := fmt.Sscanf(args[1], "%d", &targetNode); err != nil {
			return fmt.Errorf("invalid target node ID: %s", args[1])
		}

		result, err := client.RunDijkstra(ctx, startNode, targetNode)
		if err != nil {
			return fmt.Errorf("dijkstra algorithm failed: %v", err)
		}

		output, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal result: %v", err)
		}
		fmt.Println(string(output))

	case "pagerank":
		maxIterations := 100
		tolerance := 0.001

		if len(args) > 0 {
			if _, err := fmt.Sscanf(args[0], "%d", &maxIterations); err != nil {
				return fmt.Errorf("invalid max iterations: %s", args[0])
			}
		}
		if len(args) > 1 {
			if _, err := fmt.Sscanf(args[1], "%f", &tolerance); err != nil {
				return fmt.Errorf("invalid tolerance: %s", args[1])
			}
		}

		result, err := client.RunPageRank(ctx, maxIterations, tolerance)
		if err != nil {
			return fmt.Errorf("pagerank algorithm failed: %v", err)
		}

		output, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal result: %v", err)
		}
		fmt.Println(string(output))

	default:
		return fmt.Errorf("unknown algorithm type: %s", algoType)
	}

	return nil
}

func executeQuery(client *client.NenDBClient, ctx context.Context, query string) error {
	fmt.Printf("Executing query: %s\n", query)
	result, err := client.Query(ctx, query, nil)
	if err != nil {
		return fmt.Errorf("query failed: %v", err)
	}

	output, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal result: %v", err)
	}
	fmt.Println(string(output))
	return nil
}

func executeStats(client *client.NenDBClient, ctx context.Context) error {
	fmt.Println("Getting database statistics...")
	stats, err := client.GetStatistics(ctx)
	if err != nil {
		return fmt.Errorf("failed to get statistics: %v", err)
	}

	output, err := json.MarshalIndent(stats, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal statistics: %v", err)
	}
	fmt.Println(string(output))
	return nil
}

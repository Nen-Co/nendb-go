package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/nen-co/nendb-go/pkg/client"
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

	fmt.Println("=== NenDB Go Driver Basic Usage Example ===")
	fmt.Println()

	// Check server health
	fmt.Println("1. Checking server health...")
	if err := client.Health(); err != nil {
		log.Fatalf("Health check failed: %v", err)
	}
	fmt.Println("✓ Server is healthy")
	fmt.Println()

	// Create nodes
	fmt.Println("2. Creating nodes...")
	
	// Create a person node
	personLabels := []string{"Person"}
	personProps := map[string]interface{}{
		"name":  "Alice",
		"age":   30,
		"email": "alice@example.com",
	}
	personNode, err := client.CreateNode(ctx, personLabels, personProps)
	if err != nil {
		log.Fatalf("Failed to create person node: %v", err)
	}
	fmt.Printf("✓ Created person node with ID: %d\n", personNode.ID)

	// Create another person node
	person2Props := map[string]interface{}{
		"name":  "Bob",
		"age":   25,
		"email": "bob@example.com",
	}
	person2Node, err := client.CreateNode(ctx, personLabels, person2Props)
	if err != nil {
		log.Fatalf("Failed to create second person node: %v", err)
	}
	fmt.Printf("✓ Created person node with ID: %d\n", person2Node.ID)

	// Create a company node
	companyLabels := []string{"Company"}
	companyProps := map[string]interface{}{
		"name":        "TechCorp",
		"industry":    "Technology",
		"founded":     2020,
		"employees":   100,
		"is_public":   false,
	}
	companyNode, err := client.CreateNode(ctx, companyLabels, companyProps)
	if err != nil {
		log.Fatalf("Failed to create company node: %v", err)
	}
	fmt.Printf("✓ Created company node with ID: %d\n", companyNode.ID)

	// Create edges
	fmt.Println("\n3. Creating edges...")
	
	// Create KNOWS relationship between Alice and Bob
	knowsProps := map[string]interface{}{
		"since":     "2022-01-15",
		"strength":  "strong",
	}
	knowsEdge, err := client.CreateEdge(ctx, personNode.ID, person2Node.ID, "KNOWS", knowsProps)
	if err != nil {
		log.Fatalf("Failed to create KNOWS edge: %v", err)
	}
	fmt.Printf("✓ Created KNOWS edge with ID: %d\n", knowsEdge.ID)

	// Create WORKS_AT relationship between Alice and company
	worksAtProps := map[string]interface{}{
		"position":  "Software Engineer",
		"start_date": "2022-03-01",
		"salary":    75000,
	}
	worksAtEdge, err := client.CreateEdge(ctx, personNode.ID, companyNode.ID, "WORKS_AT", worksAtProps)
	if err != nil {
		log.Fatalf("Failed to create WORKS_AT edge: %v", err)
	}
	fmt.Printf("✓ Created WORKS_AT edge with ID: %d\n", worksAtEdge.ID)

	// Create FOUNDED relationship between Bob and company
	foundedProps := map[string]interface{}{
		"role":      "Founder",
		"equity":    0.8,
	}
	foundedEdge, err := client.CreateEdge(ctx, person2Node.ID, companyNode.ID, "FOUNDED", foundedProps)
	if err != nil {
		log.Fatalf("Failed to create FOUNDED edge: %v", err)
	}
	fmt.Printf("✓ Created FOUNDED edge with ID: %d\n", foundedEdge.ID)

	// Retrieve and display nodes
	fmt.Println("\n4. Retrieving nodes...")
	
	retrievedPerson, err := client.GetNode(ctx, personNode.ID)
	if err != nil {
		log.Fatalf("Failed to retrieve person node: %v", err)
	}
	fmt.Printf("✓ Retrieved person node: %+v\n", retrievedPerson)

	retrievedCompany, err := client.GetNode(ctx, companyNode.ID)
	if err != nil {
		log.Fatalf("Failed to retrieve company node: %v", err)
	}
	fmt.Printf("✓ Retrieved company node: %+v\n", retrievedCompany)

	// Run algorithms
	fmt.Println("\n5. Running algorithms...")
	
	// Run BFS to find path between Alice and Bob
	fmt.Println("Running BFS algorithm...")
	bfsResult, err := client.RunBFS(ctx, personNode.ID, person2Node.ID, 5)
	if err != nil {
		log.Printf("BFS algorithm failed: %v", err)
	} else {
		fmt.Printf("✓ BFS result: %+v\n", bfsResult)
	}

	// Run PageRank to find important nodes
	fmt.Println("Running PageRank algorithm...")
	pagerankResult, err := client.RunPageRank(ctx, 100, 0.001)
	if err != nil {
		log.Printf("PageRank algorithm failed: %v", err)
	} else {
		fmt.Printf("✓ PageRank result: %+v\n", pagerankResult)
	}

	// Execute custom query
	fmt.Println("\n6. Executing custom query...")
	query := "MATCH (n) RETURN n LIMIT 5"
	result, err := client.Query(ctx, query, nil)
	if err != nil {
		log.Printf("Custom query failed: %v", err)
	} else {
		fmt.Printf("✓ Query result: %+v\n", result)
	}

	// Get database statistics
	fmt.Println("\n7. Getting database statistics...")
	stats, err := client.GetStatistics(ctx)
	if err != nil {
		log.Printf("Failed to get statistics: %v", err)
	} else {
		fmt.Printf("✓ Database statistics: %+v\n", stats)
	}

	fmt.Println("\n=== Example completed successfully! ===")
}

# 🍳 Fiber + NenDB Recipe Makefile
# Easy commands for building, testing, and running the recipe

.PHONY: help build test run clean deps install

# Default target
help:
	@echo "🍳 Fiber + NenDB Recipe - Available Commands:"
	@echo ""
	@echo "📦 Setup:"
	@echo "  deps      - Download dependencies"
	@echo "  install   - Install dependencies"
	@echo ""
	@echo "🔨 Build:"
	@echo "  build     - Build the recipe"
	@echo "  clean     - Clean build artifacts"
	@echo ""
	@echo "🧪 Test:"
	@echo "  test      - Run tests"
	@echo "  test-v    - Run tests with verbose output"
	@echo ""
	@echo "🚀 Run:"
	@echo "  run       - Run the recipe"
	@echo "  dev       - Run with live reload (requires air)"
	@echo ""
	@echo "📚 Documentation:"
	@echo "  readme    - Show recipe README"
	@echo "  help      - Show this help message"

# Download dependencies
deps:
	@echo "📦 Downloading dependencies..."
	go mod download
	@echo "✅ Dependencies downloaded"

# Install dependencies
install:
	@echo "📦 Installing dependencies..."
	go mod tidy
	go mod download
	@echo "✅ Dependencies installed"

# Build the recipe
build:
	@echo "🔨 Building Fiber + NenDB recipe..."
	go build -o fiber-nendb-recipe .
	@echo "✅ Recipe built successfully"

# Clean build artifacts
clean:
	@echo "🧹 Cleaning build artifacts..."
	rm -f fiber-nendb-recipe
	@echo "✅ Cleaned successfully"

# Run tests
test:
	@echo "🧪 Running tests..."
	go test ./...
	@echo "✅ Tests completed"

# Run tests with verbose output
test-v:
	@echo "🧪 Running tests with verbose output..."
	go test -v ./...
	@echo "✅ Tests completed"

# Run the recipe
run:
	@echo "🚀 Starting Fiber + NenDB recipe..."
	@echo "📖 Make sure NenDB server is running on localhost:8080"
	@echo "🌐 Recipe will be available at http://localhost:3000"
	@echo ""
	go run main.go

# Run with live reload (requires air)
dev:
	@echo "🚀 Starting Fiber + NenDB recipe with live reload..."
	@echo "📖 Make sure NenDB server is running on localhost:8080"
	@echo "🌐 Recipe will be available at http://localhost:3000"
	@echo ""
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "❌ Air not found. Install with: go install github.com/cosmtrek/air@latest"; \
		echo "💡 Or run without live reload: make run"; \
	fi

# Show recipe README
readme:
	@echo "📚 Fiber + NenDB Recipe README:"
	@echo "=================================="
	@cat README.md

# Quick start - setup everything and run
quickstart: deps build test run

# Development workflow
dev-setup: deps install test
	@echo "🎯 Development environment ready!"
	@echo "💡 Run 'make run' to start the recipe"
	@echo "💡 Run 'make dev' for live reload (requires air)"

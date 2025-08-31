.PHONY: all build test clean run-example run-cli help

# Default target
all: build test

# Build the CLI and examples
build:
	@echo "Building NenDB Go Driver..."
	go build -o bin/nendb cmd/nendb/main.go
	go build -o bin/basic_usage examples/basic_usage.go
	@echo "Build complete!"

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...
	@echo "Tests complete!"

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run tests with race detection
test-race:
	@echo "Running tests with race detection..."
	go test -race ./...
	@echo "Race detection tests complete!"

# Run benchmarks
bench:
	@echo "Running benchmarks..."
	go test -bench=. ./...
	@echo "Benchmarks complete!"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -f coverage.out coverage.html
	go clean
	@echo "Clean complete!"

# Run the basic usage example
run-example: build
	@echo "Running basic usage example..."
	@echo "Note: Make sure NenDB server is running on localhost:8080"
	./bin/basic_usage

# Run the CLI tool
run-cli: build
	@echo "Running CLI tool..."
	@echo "Usage: ./bin/nendb -help"
	./bin/nendb -help

# Install the CLI tool
install: build
	@echo "Installing CLI tool..."
	cp bin/nendb /usr/local/bin/
	@echo "CLI tool installed to /usr/local/bin/nendb"

# Uninstall the CLI tool
uninstall:
	@echo "Uninstalling CLI tool..."
	rm -f /usr/local/bin/nendb
	@echo "CLI tool uninstalled"

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...
	@echo "Code formatting complete!"

# Run linter
lint:
	@echo "Running linter..."
	golangci-lint run
	@echo "Linting complete!"

# Check for security vulnerabilities
security:
	@echo "Checking for security vulnerabilities..."
	gosec ./...
	@echo "Security check complete!"

# Generate documentation
docs:
	@echo "Generating documentation..."
	godoc -http=:6060 &
	@echo "Documentation server started at http://localhost:6060"
	@echo "Press Ctrl+C to stop"

# Dependencies
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy
	@echo "Dependencies installed!"

# Update dependencies
deps-update:
	@echo "Updating dependencies..."
	go get -u ./...
	go mod tidy
	@echo "Dependencies updated!"

# Show help
help:
	@echo "NenDB Go Driver Makefile"
	@echo ""
	@echo "Available targets:"
	@echo "  all              - Build and test (default)"
	@echo "  build            - Build the CLI and examples"
	@echo "  test             - Run tests"
	@echo "  test-coverage    - Run tests with coverage report"
	@echo "  test-race        - Run tests with race detection"
	@echo "  bench            - Run benchmarks"
	@echo "  clean            - Clean build artifacts"
	@echo "  run-example      - Run the basic usage example"
	@echo "  run-cli          - Run the CLI tool"
	@echo "  install          - Install CLI tool to /usr/local/bin"
	@echo "  uninstall        - Uninstall CLI tool"
	@echo "  fmt              - Format code"
	@echo "  lint             - Run linter"
	@echo "  security         - Check for security vulnerabilities"
	@echo "  docs             - Start documentation server"
	@echo "  deps             - Install dependencies"
	@echo "  deps-update      - Update dependencies"
	@echo "  help             - Show this help message"
	@echo ""
	@echo "Examples:"
	@echo "  make build       - Build the project"
	@echo "  make test        - Run tests"
	@echo "  make run-example - Run example (requires NenDB server)"
	@echo "  make install     - Install CLI tool"

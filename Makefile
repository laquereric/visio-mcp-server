.PHONY: all build test clean install run debug help

# Variables
BINARY_NAME=visio-mcp-server
GO_FILES=$(shell find . -name '*.go' -type f)
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")

# Default target
all: build

# Build the project
build: build-go build-ts

# Build Go binary
build-go:
	@echo "Building Go binary..."
	go build -ldflags="-X main.version=$(VERSION)" -o $(BINARY_NAME) ./cmd/visio-mcp-server

# Build TypeScript
build-ts:
	@echo "Building TypeScript..."
	npm run build

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -cover -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -f $(BINARY_NAME)
	rm -rf dist/
	rm -f coverage.out coverage.html

# Install dependencies
install:
	@echo "Installing dependencies..."
	go mod download
	npm install

# Run the server
run: build-go
	@echo "Running server..."
	./$(BINARY_NAME)

# Debug with MCP Inspector
debug: build
	@echo "Starting MCP Inspector..."
	npm run debug

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Lint code
lint:
	@echo "Linting code..."
	go vet ./...

# Build snapshot release
snapshot:
	@echo "Building snapshot release..."
	goreleaser build --snapshot --clean

# Show help
help:
	@echo "Available targets:"
	@echo "  all           - Build everything (default)"
	@echo "  build         - Build Go binary and TypeScript"
	@echo "  build-go      - Build Go binary only"
	@echo "  build-ts      - Build TypeScript only"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage"
	@echo "  clean         - Remove build artifacts"
	@echo "  install       - Install dependencies"
	@echo "  run           - Build and run the server"
	@echo "  debug         - Start MCP Inspector"
	@echo "  fmt           - Format Go code"
	@echo "  lint          - Lint Go code"
	@echo "  snapshot      - Build snapshot release"
	@echo "  help          - Show this help message"

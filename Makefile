# Go Fiber Application
.PHONY: help build run dev test clean install docker-build docker-run

# Variables
BINARY_NAME=go_server
GO_VERSION=1.24.7
PORT=3000

# Default target - show help
default: help

# Install dependencies
install:
	@echo "📦 Installing dependencies..."
	go mod download
	go mod verify

# Build the application
build:
	@echo "🔨 Building application..."
	go build -o $(BINARY_NAME) .

# Run the application
run: build
	@echo "🚀 Starting application..."
	./$(BINARY_NAME)

# Development mode with hot reload (requires Air)
dev:
	@echo "🔥 Starting development server with hot reload..."
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "⚠️  Air is not installed. Installing now..."; \
		go install github.com/cosmtrek/air@latest; \
		air; \
	fi

# Run tests
test:
	@echo "🧪 Running tests..."
	go test -v ./...

# Clean build artifacts
clean:
	@echo "🧹 Cleaning build artifacts..."
	go clean
	rm -f $(BINARY_NAME)
	rm -rf tmp/

# Docker commands
docker-build:
	@echo "🐳 Building Docker image..."
	docker build -t $(BINARY_NAME):latest .

docker-run:
	@echo "🐳 Running Docker container..."
	docker run -p $(PORT):$(PORT) --rm $(BINARY_NAME):latest

# Docker Compose commands
compose-up:
	@echo "🐳 Starting services with Docker Compose..."
	docker-compose up

compose-down:
	@echo "🐳 Stopping services..."
	docker-compose down

compose-build:
	@echo "🐳 Building services with Docker Compose..."
	docker-compose build

# Format code
fmt:
	@echo "✨ Formatting code..."
	go fmt ./...

# Lint code (requires golangci-lint)
lint:
	@echo "🔍 Linting code..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "⚠️  golangci-lint is not installed. Install from https://golangci-lint.run/usage/install/"; \
	fi

# Show help
help:
	@echo "go_server - Available commands:"
	@echo ""
	@echo "📋 Basic Commands:"
	@echo "  make install        Install dependencies"
	@echo "  make build          Build the application"
	@echo "  make run            Build and run the application"
	@echo "  make dev            Run with hot reload (installs Air if needed)"
	@echo "  make test           Run tests"
	@echo "  make clean          Clean build artifacts"
	@echo ""
	@echo "🐳 Docker Commands:"
	@echo "  make docker-build   Build Docker image"
	@echo "  make docker-run     Run in Docker container"
	@echo "  make compose-up     Start with Docker Compose"
	@echo "  make compose-down   Stop Docker Compose services"
	@echo "  make compose-build  Build Docker Compose services"
	@echo ""
	@echo "🛠️  Development Commands:"
	@echo "  make fmt            Format code"
	@echo "  make lint           Lint code (requires golangci-lint)"
	@echo ""
	@echo "💡 Quick Start:"
	@echo "  1. make install     # Install dependencies"
	@echo "  2. make dev         # Start development server with hot reload"
	@echo ""
	@echo "🚀 Creating a new app with create-fiber-app:"
	@echo "  npx create-fiber-app my-app"
	@echo "  cd my-app"
	@echo "  make dev"
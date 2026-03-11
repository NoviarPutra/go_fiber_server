include .env
export

# Go Fiber Application
.PHONY: help build run dev test clean install docker-build docker-run

# Variables
BINARY_NAME=go_server
GO_VERSION=1.24.7
PORT=3000
MIGRATE_PATH=migrations
DB_URL=$(DATABASE_URL)

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
	@AIR_PATH=$$(go env GOPATH)/bin/air; \
	if [ -f "$$AIR_PATH" ]; then \
		$$AIR_PATH; \
	else \
		echo "⚠️  Air not found at $$AIR_PATH. Installing..."; \
		go install github.com/air-verse/air@latest; \
		$$AIR_PATH; \
	fi

# Run tests
test:
	@echo "🧪 Running tests..."
	go test -v ./...

# --- Database Migration Targets ---
.PHONY: migrate-create migrate-up migrate-down migrate-force db-status

## migrate-create: Create new migration (usage: make migrate-create name=create_users)
migrate-create:
	@migrate create -ext sql -dir $(MIGRATE_PATH) -seq $(name)

## migrate-up: Run all up migrations
migrate-up:
	@migrate -path $(MIGRATE_PATH) -database "$(DB_URL)" up

## migrate-down: Rollback 1 migration
migrate-down:
	@migrate -path $(MIGRATE_PATH) -database "$(DB_URL)" down 1

## migrate-force: Force migration version (usage: make migrate-force version=1)
migrate-force:
	@migrate -path $(MIGRATE_PATH) -database "$(DB_URL)" force $(version)

## db-status: Check current migration version
db-status:
	@migrate -path $(MIGRATE_PATH) -database "$(DB_URL)" version

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
include .env
export

.PHONY: help install build run dev test clean fmt lint \
        migrate-create migrate-up migrate-down migrate-force db-status \
        docker-build docker-run compose-up compose-down compose-build

# ─── Variables ────────────────────────────────────────────────────────────────
BINARY_NAME  = go_server
APP_PORT     = $(PORT)
MIGRATE_PATH = migrations
DB_URL       = $(DATABASE_URL)

# ─── Default ──────────────────────────────────────────────────────────────────
default: help

# ─── Dependencies ─────────────────────────────────────────────────────────────
install:
	@echo "📦 Installing dependencies..."
	@go mod download
	@go mod verify
	@echo "✅ Dependencies installed"

# ─── Build ────────────────────────────────────────────────────────────────────
vet:
	@echo "🔍 Running go vet..."
	@go vet ./...

build: vet
	@echo "🔨 Building application..."
	@go build -ldflags="-s -w" -o $(BINARY_NAME) .
	@echo "✅ Build complete: ./$(BINARY_NAME)"

run: build
	@echo "🚀 Starting application on port $(APP_PORT)..."
	@./$(BINARY_NAME)

# ─── Development ──────────────────────────────────────────────────────────────
dev:
	@echo "🔥 Starting development server with hot reload..."
	@AIR_PATH=$$(go env GOPATH)/bin/air; \
	if [ ! -f "$$AIR_PATH" ]; then \
		echo "⚠️  Air not found. Installing..."; \
		go install github.com/air-verse/air@latest; \
	fi; \
	if [ ! -f ".air.toml" ]; then \
		echo "⚠️  .air.toml tidak ditemukan. Membuat default config..."; \
		$$AIR_PATH init; \
	fi; \
	$$AIR_PATH

# ─── Test ─────────────────────────────────────────────────────────────────────
test:
	@echo "🧪 Running tests..."
	@go test -v -race -count=1 ./...

test-cover:
	@echo "🧪 Running tests with coverage..."
	@go test -v -race -count=1 -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Coverage report: coverage.html"

# ─── Database Migration ───────────────────────────────────────────────────────
check-migrate:
	@if ! command -v migrate > /dev/null; then \
		echo "❌ golang-migrate tidak ditemukan."; \
		echo "   Install: https://github.com/golang-migrate/migrate/tree/master/cmd/migrate"; \
		exit 1; \
	fi

migrate-create: check-migrate
	@if [ -z "$(name)" ]; then \
		echo "❌ Usage: make migrate-create name=create_users"; \
		exit 1; \
	fi
	@migrate create -ext sql -dir $(MIGRATE_PATH) -seq $(name)
	@echo "✅ Migration created: $(MIGRATE_PATH)"

migrate-up: check-migrate
	@echo "⬆️  Running all pending migrations..."
	@migrate -path $(MIGRATE_PATH) -database "$(DB_URL)" up
	@echo "✅ Migrations applied"

migrate-down: check-migrate
	@echo "⬇️  Rolling back 1 migration..."
	@migrate -path $(MIGRATE_PATH) -database "$(DB_URL)" down 1
	@echo "✅ Rollback complete"

migrate-force: check-migrate
	@if [ -z "$(version)" ]; then \
		echo "❌ Usage: make migrate-force version=1"; \
		exit 1; \
	fi
	@migrate -path $(MIGRATE_PATH) -database "$(DB_URL)" force $(version)
	@echo "✅ Migration forced to version $(version)"

db-status: check-migrate
	@echo "📊 Current migration version:"
	@migrate -path $(MIGRATE_PATH) -database "$(DB_URL)" version

# ─── Utilities ────────────────────────────────────────────────────────────────
clean:
	@echo "🧹 Cleaning build artifacts..."
	@go clean
	@rm -f $(BINARY_NAME)
	@rm -f coverage.out coverage.html
	@rm -rf tmp/
	@echo "✅ Clean complete"

fmt:
	@echo "✨ Formatting code..."
	@go fmt ./...
	@echo "✅ Format complete"

lint:
	@echo "🔍 Linting code..."
	@LINT_PATH=$$(go env GOPATH)/bin/golangci-lint; \
	if [ ! -f "$$LINT_PATH" ]; then \
		echo "⚠️  golangci-lint tidak ditemukan. Installing..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
	fi; \
	$$LINT_PATH run
	@echo "✅ Lint complete"

# ─── Docker ───────────────────────────────────────────────────────────────────
docker-build:
	@echo "🐳 Building Docker image..."
	@docker build -t $(BINARY_NAME):latest .
	@echo "✅ Image built: $(BINARY_NAME):latest"

docker-run:
	@echo "🐳 Running Docker container on port $(APP_PORT)..."
	@docker run -p $(APP_PORT):$(APP_PORT) --env-file .env --rm $(BINARY_NAME):latest

compose-up:
	@echo "🐳 Starting services..."
	@docker-compose up

compose-down:
	@echo "🐳 Stopping services..."
	@docker-compose down

compose-build:
	@echo "🐳 Building services..."
	@docker-compose build

# ─── Help ─────────────────────────────────────────────────────────────────────
help:
	@echo ""
	@echo "  $(BINARY_NAME) — Available Commands"
	@echo ""
	@echo "  📋 Basic"
	@echo "    make install                        Install dependencies"
	@echo "    make build                          Vet + build (optimized)"
	@echo "    make run                            Build and run"
	@echo "    make dev                            Hot reload (auto-install Air)"
	@echo "    make test                           Run tests (with -race)"
	@echo "    make test-cover                     Run tests + coverage report"
	@echo "    make clean                          Clean all artifacts"
	@echo ""
	@echo "  🗄️  Database"
	@echo "    make migrate-create name=<name>     Create new migration"
	@echo "    make migrate-up                     Run all pending migrations"
	@echo "    make migrate-down                   Rollback 1 migration"
	@echo "    make migrate-force version=<ver>    Force migration version"
	@echo "    make db-status                      Check current version"
	@echo ""
	@echo "  🐳 Docker"
	@echo "    make docker-build                   Build Docker image"
	@echo "    make docker-run                     Run in Docker"
	@echo "    make compose-up                     Start Docker Compose"
	@echo "    make compose-down                   Stop Docker Compose"
	@echo "    make compose-build                  Build Docker Compose"
	@echo ""
	@echo "  🛠️  Code Quality"
	@echo "    make fmt                            Format code"
	@echo "    make lint                           Lint (requires golangci-lint)"
	@echo "    make vet                            Run go vet"
	@echo ""
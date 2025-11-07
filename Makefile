.PHONY: help build run test clean docker-build docker-run docker-stop install deps lint fmt

# Variables
BINARY_NAME=gourl
MAIN_PATH=./cmd/server
DOCKER_IMAGE=gourl:latest

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

install: ## Install dependencies
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy

build: ## Build the application
	@echo "Building $(BINARY_NAME)..."
	@go build -o $(BINARY_NAME) $(MAIN_PATH)
	@echo "Build complete: $(BINARY_NAME)"

run: ## Run the application
	@echo "Running $(BINARY_NAME)..."
	@go run $(MAIN_PATH)

test: ## Run tests
	@echo "Running tests..."
	@go test -v ./...

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -f $(BINARY_NAME)
	@rm -f gourl.db
	@go clean
	@echo "Clean complete"

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	@go mod download

lint: ## Run linter
	@echo "Running linter..."
	@golangci-lint run ./... || echo "Install golangci-lint: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"

fmt: ## Format code
	@echo "Formatting code..."
	@go fmt ./...

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	@docker build -t $(DOCKER_IMAGE) .
	@echo "Docker image built: $(DOCKER_IMAGE)"

docker-run: ## Run Docker container
	@echo "Running Docker container..."
	@docker-compose up -d
	@echo "Container running. Check logs with: docker-compose logs -f"

docker-stop: ## Stop Docker container
	@echo "Stopping Docker container..."
	@docker-compose down
	@echo "Container stopped"

docker-logs: ## View Docker logs
	@docker-compose logs -f

dev: ## Run in development mode
	@echo "Running in development mode..."
	@ENV=development go run $(MAIN_PATH)

prod: build ## Run in production mode
	@echo "Running in production mode..."
	@ENV=production ./$(BINARY_NAME)

db-reset: ## Reset database (WARNING: deletes all data)
	@echo "WARNING: This will delete gourl.db"
	@read -p "Are you sure? [y/N] " -n 1 -r; \
	echo; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		rm -f gourl.db; \
		echo "Database reset complete"; \
	fi


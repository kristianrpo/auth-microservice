.PHONY: help build run test clean docker-build docker-up docker-down lint fmt tidy

# Variables
APP_NAME=auth-service
DOCKER_IMAGE=ghcr.io/kristianrpo/auth-microservice
DOCKER_TAG=latest
GO_FILES=$(shell find . -name '*.go' -not -path "./vendor/*")

# Help command
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Build the application
build: ## Build the application binary
	@echo "Building $(APP_NAME)..."
	@go build -o bin/$(APP_NAME) ./cmd/server

# Run the application
run: ## Run the application
	@echo "Running $(APP_NAME)..."
	@go run ./cmd/server/main.go

# Run tests
test: ## Run all tests
	@echo "Running tests..."
	@go test -v -race -coverprofile=coverage.out -covermode=atomic ./...

# Run tests with coverage report
test-coverage: test ## Run tests and generate coverage report
	@echo "Generating coverage report..."
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Clean build artifacts
clean: ## Remove build artifacts
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@go clean

# Docker commands
docker-build: ## Build Docker image
	@echo "Building Docker image..."
	@docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

docker-push: ## Push Docker image
	@echo "Pushing Docker image..."
	@docker push $(DOCKER_IMAGE):$(DOCKER_TAG)

docker-up: ## Start all services with docker-compose
	@echo "Starting services..."
	@docker-compose up -d

docker-down: ## Stop all services
	@echo "Stopping services..."
	@docker-compose down

docker-logs: ## Show docker-compose logs
	@docker-compose logs -f

docker-restart: docker-down docker-up ## Restart all services

# Linting and formatting
lint: ## Run linter
	@echo "Running linter..."
	@golangci-lint run ./...

fmt: ## Format code
	@echo "Formatting code..."
	@gofmt -s -w $(GO_FILES)
	@goimports -w $(GO_FILES)

# Go mod commands
tidy: ## Tidy go modules
	@echo "Tidying go modules..."
	@go mod tidy

download: ## Download dependencies
	@echo "Downloading dependencies..."
	@go mod download

verify: ## Verify dependencies
	@echo "Verifying dependencies..."
	@go mod verify

# Database commands
db-migrate-up: ## Run database migrations up
	@echo "Running migrations up..."
	@migrate -path migrations -database "postgresql://authuser:authpassword@localhost:5432/authdb?sslmode=disable" up

db-migrate-down: ## Rollback last migration
	@echo "Rolling back migration..."
	@migrate -path migrations -database "postgresql://authuser:authpassword@localhost:5432/authdb?sslmode=disable" down 1

# Kubernetes commands
k8s-apply-dev: ## Apply Kubernetes manifests for dev
	@echo "Applying Kubernetes manifests for dev..."
	@kubectl apply -k k8s/overlays/dev

k8s-apply-staging: ## Apply Kubernetes manifests for staging
	@echo "Applying Kubernetes manifests for staging..."
	@kubectl apply -k k8s/overlays/staging

k8s-apply-prod: ## Apply Kubernetes manifests for production
	@echo "Applying Kubernetes manifests for production..."
	@kubectl apply -k k8s/overlays/production

k8s-delete-dev: ## Delete Kubernetes resources for dev
	@echo "Deleting Kubernetes resources for dev..."
	@kubectl delete -k k8s/overlays/dev

# Swagger documentation
swagger: ## Generate Swagger documentation
	@echo "Generating Swagger documentation..."
	@swag init -g cmd/server/main.go -o docs --parseDependency --parseInternal
	@echo "Swagger documentation generated in docs/"

swagger-fmt: ## Format Swagger comments
	@echo "Formatting Swagger comments..."
	@swag fmt

# Install development tools
install-tools: ## Install development tools
	@echo "Installing development tools..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install golang.org/x/tools/cmd/goimports@latest
	@go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@go install github.com/swaggo/swag/cmd/swag@latest

# CI/CD commands
ci: lint test ## Run CI checks

# Show project info
info: ## Show project information
	@echo "Project: $(APP_NAME)"
	@echo "Docker Image: $(DOCKER_IMAGE):$(DOCKER_TAG)"
	@echo "Go Version: $(shell go version)"
	@echo "Go Files: $(words $(GO_FILES))"


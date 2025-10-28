# Golang RabbitMQ v2 - Makefile
# 
# Available commands:
#   make help        - Show this help message
#   make dev-setup   - Setup development environment
#   make deps        - Install dependencies
#   make build       - Build the application
#   make run         - Run the application locally
#   make test        - Run tests
#   make clean       - Clean build artifacts
#   make docker-*    - Docker related commands

.PHONY: help dev-setup deps build run test clean docker-up docker-down docker-build docker-logs docker-ps

# Variables
APP_NAME := rabbitmq-app
GO_VERSION := 1.23
DOCKER_COMPOSE := docker-compose
BUILD_DIR := bin
MAIN_PATH := cmd/api/main.go

# Default target
help: ## Show this help message
	@echo "Golang RabbitMQ v2 - Available Commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
	@echo ""

# Development Setup
dev-setup: ## Setup development environment
	@echo "Setting up development environment..."
	@cp .env.example .env
	@echo "✅ Created .env file from .env.example"
	@echo "📝 Please edit .env file with your configuration"
	@go mod tidy
	@echo "✅ Dependencies installed"

deps: ## Install/update dependencies
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy
	@echo "✅ Dependencies updated"

# Build Commands
build: ## Build the application
	@echo "Building application..."
	@mkdir -p $(BUILD_DIR)
	@CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PATH)
	@echo "✅ Application built: $(BUILD_DIR)/$(APP_NAME)"

build-local: ## Build for local platform
	@echo "Building application for local platform..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PATH)
	@echo "✅ Application built: $(BUILD_DIR)/$(APP_NAME)"

# Run Commands
run: ## Run the application locally
	@echo "Running application..."
	@go run $(MAIN_PATH)

run-watch: ## Run with live reload (requires air)
	@echo "Running with live reload..."
	@air

# Testing
test: ## Run tests
	@echo "Running tests..."
	@go test ./... -v

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	@go test ./... -v -coverprofile=coverage.out
	@go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Coverage report: coverage.html"

benchmark: ## Run benchmarks
	@echo "Running benchmarks..."
	@go test ./... -bench=. -benchmem

# Code Quality
fmt: ## Format code
	@echo "Formatting code..."
	@go fmt ./...
	@echo "✅ Code formatted"

lint: ## Lint code (requires golangci-lint)
	@echo "Linting code..."
	@golangci-lint run
	@echo "✅ Code linted"

vet: ## Run go vet
	@echo "Running go vet..."
	@go vet ./...
	@echo "✅ Code vetted"

# Cleanup
clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@echo "✅ Cleaned"

clean-logs: ## Clean log files
	@echo "Cleaning log files..."
	@rm -rf logs/*.log logs/*.gz
	@echo "✅ Log files cleaned"

# Docker Commands
docker-up: ## Start all services with Docker Compose
	@echo "Starting Docker services..."
	@$(DOCKER_COMPOSE) up -d
	@echo "✅ Services started"
	@make docker-ps

docker-down: ## Stop all services
	@echo "Stopping Docker services..."
	@$(DOCKER_COMPOSE) down
	@echo "✅ Services stopped"

docker-build: ## Build Docker images
	@echo "Building Docker images..."
	@$(DOCKER_COMPOSE) build --no-cache
	@echo "✅ Images built"

docker-rebuild: ## Rebuild and restart services
	@echo "Rebuilding and restarting services..."
	@$(DOCKER_COMPOSE) down
	@$(DOCKER_COMPOSE) build --no-cache
	@$(DOCKER_COMPOSE) up -d
	@echo "✅ Services rebuilt and restarted"

docker-ps: ## Show container status
	@echo "Container Status:"
	@$(DOCKER_COMPOSE) ps

docker-logs: ## Show application logs
	@echo "Application logs:"
	@$(DOCKER_COMPOSE) logs -f app

docker-logs-all: ## Show all service logs
	@echo "All service logs:"
	@$(DOCKER_COMPOSE) logs -f

docker-shell: ## Access application container shell
	@echo "Accessing application container..."
	@$(DOCKER_COMPOSE) exec app sh

docker-clean: ## Remove all containers, networks, and volumes
	@echo "Cleaning Docker resources..."
	@$(DOCKER_COMPOSE) down -v
	@docker system prune -f
	@echo "✅ Docker resources cleaned"

# Database Commands
db-migrate: ## Run database migrations (inside container)
	@echo "Running database migrations..."
	@$(DOCKER_COMPOSE) exec app ./$(APP_NAME) migrate
	@echo "✅ Migrations completed"

db-reset: ## Reset database (recreate containers)
	@echo "Resetting database..."
	@$(DOCKER_COMPOSE) stop postgres
	@$(DOCKER_COMPOSE) rm -f postgres
	@docker volume rm golang-rabbitmq-v2_postgres_data || true
	@$(DOCKER_COMPOSE) up -d postgres
	@echo "✅ Database reset"

# RabbitMQ Commands
rabbitmq-ui: ## Open RabbitMQ Management UI
	@echo "Opening RabbitMQ Management UI..."
	@echo "URL: http://localhost:15672"
	@echo "Login: guest/guest"

rabbitmq-reset: ## Reset RabbitMQ (recreate container)
	@echo "Resetting RabbitMQ..."
	@$(DOCKER_COMPOSE) stop rabbitmq
	@$(DOCKER_COMPOSE) rm -f rabbitmq
	@docker volume rm golang-rabbitmq-v2_rabbitmq_data || true
	@$(DOCKER_COMPOSE) up -d rabbitmq
	@echo "✅ RabbitMQ reset"

# Dashboard Commands
dashboard: ## Open monitoring dashboard
	@echo "Opening RabbitMQ Monitoring Dashboard..."
	@cd web && ./open-dashboard.bat || ./open-dashboard.sh

dashboard-direct: ## Open dashboard directly (no HTTP server)
	@echo "Opening dashboard directly..."
	@start web/dashboard.html || open web/dashboard.html || xdg-open web/dashboard.html

# Health Checks
health: ## Check application health
	@echo "Checking application health..."
	@curl -s http://localhost:8080/api/v1/monitoring/health | jq . || echo "❌ Health check failed"

metrics: ## Show RabbitMQ metrics
	@echo "RabbitMQ Metrics:"
	@curl -s http://localhost:8080/api/v1/monitoring/metrics | jq . || echo "❌ Metrics unavailable"

# API Testing
test-api: ## Run basic API tests
	@echo "Testing API endpoints..."
	@echo "1. Health Check:"
	@curl -s http://localhost:8080/api/v1/monitoring/health | jq .status || echo "❌ Health check failed"
	@echo ""
	@echo "2. Create Test User:"
	@curl -s -X POST http://localhost:8080/api/v1/users \
		-H "Content-Type: application/json" \
		-d '{"name":"Test User","email":"test@example.com","phone":"1234567890"}' | jq . || echo "❌ User creation failed"
	@echo ""
	@echo "3. Message Stats:"
	@curl -s http://localhost:8080/api/v1/messages/stats | jq . || echo "❌ Stats unavailable"

# Development Tools
install-tools: ## Install development tools
	@echo "Installing development tools..."
	@go install github.com/cosmtrek/air@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "✅ Development tools installed"

# Project Info
info: ## Show project information
	@echo "Project Information:"
	@echo "  Name: $(APP_NAME)"
	@echo "  Go Version: $(GO_VERSION)"
	@echo "  Main Path: $(MAIN_PATH)"
	@echo "  Build Dir: $(BUILD_DIR)"
	@echo ""
	@echo "Services:"
	@echo "  Application: http://localhost:8080"
	@echo "  RabbitMQ UI: http://localhost:15672 (guest/guest)"
	@echo "  PostgreSQL: localhost:5432"
	@echo "  Dashboard: web/dashboard.html (run 'make dashboard')"
	@echo ""
	@echo "API Endpoints:"
	@echo "  Health: GET /api/v1/monitoring/health"
	@echo "  Metrics: GET /api/v1/monitoring/metrics"
	@echo "  Activity: GET /api/v1/monitoring/activity"
	@echo "  Events: GET /api/v1/monitoring/events"
	@echo "  Users: POST|GET /api/v1/users"
	@echo "  Messages: POST|GET /api/v1/messages"

# Quick start for new developers
quick-start: dev-setup docker-up ## Quick start for new developers
	@echo ""
	@echo "🚀 Quick Start Complete!"
	@echo ""
	@echo "Next steps:"
	@echo "  1. Check services: make docker-ps"
	@echo "  2. Test health: make health"
	@echo "  3. View logs: make docker-logs"
	@echo "  4. Open RabbitMQ UI: make rabbitmq-ui"
	@echo ""

# Full demo start (start services, wait, test, open dashboard)
start: ## Start everything (services + dashboard + test)
	@echo "🚀 Starting Full Quick Start Demo..."
	@echo ""
	@echo "Step 1/4: Starting Docker services..."
	@make docker-up
	@echo ""
	@echo "Step 2/4: Waiting for services to be ready..."
	@echo "⏳ Please wait 30 seconds..."
	@sleep 30 || timeout 30
	@echo ""
	@echo "Step 3/4: Testing API..."
	@make test-api
	@echo ""
	@echo "Step 4/4: Opening dashboard..."
	@echo "🎨 Dashboard will open in your browser..."
	@echo "📊 You can also run: make dashboard"
	@echo ""
	@echo "✅ Quick Start Complete!"
	@echo ""
	@echo "📍 Services Running:"
	@echo "  • Application: http://localhost:8080"
	@echo "  • Dashboard: Run 'make dashboard' to open"
	@echo "  • RabbitMQ UI: http://localhost:15672 (guest/guest)"
	@echo ""
	@echo "🔧 Useful Commands:"
	@echo "  • make dashboard    - Open monitoring dashboard"
	@echo "  • make health       - Check application health"
	@echo "  • make metrics      - Show RabbitMQ metrics"
	@echo "  • make docker-logs  - View application logs"
	@echo "  • make docker-down  - Stop all services"
	@echo ""
	@echo "🎉 Ready to monitor! Run 'make dashboard' now!"


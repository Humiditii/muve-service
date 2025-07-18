.PHONY: run build test clean docker-build docker-run docker-stop setup deps

# Go parameters
BINARY_NAME=movie-api
MAIN_PATH=./cmd/main.go

# Build the application
build:
	go build -o $(BINARY_NAME) $(MAIN_PATH)

# Run the application
run:
	go run $(MAIN_PATH)

# Install dependencies
deps:
	go mod tidy
	go mod download

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	go clean
	rm -f $(BINARY_NAME)

# Setup development environment
setup:
	@echo "Setting up development environment..."
	@if ! command -v go &> /dev/null; then \
		echo "Go is not installed. Please install Go 1.21+"; \
		exit 1; \
	fi
	@echo "Installing dependencies..."
	go mod tidy
	@echo "Creating upload directory..."
	mkdir -p uploads/movies
	@echo "✅ Setup complete!"

# Docker commands
docker-build:
	docker build -t $(BINARY_NAME) .

docker-run:
	docker-compose up -d

docker-stop:
	docker-compose down

docker-logs:
	docker-compose logs -f

# Development database commands
db-up:
	docker run --name moviedb-dev -e POSTGRES_DB=moviedb -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=password -p 5432:5432 -d postgres:15

db-down:
	docker stop moviedb-dev || true
	docker rm moviedb-dev || true

# API testing
test-api:
	@if ! command -v jq &> /dev/null; then \
		echo "jq is not installed. Please install jq to run API tests"; \
		exit 1; \
	fi
	./test_api.sh

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	golangci-lint run

# Generate API documentation
docs:
	@echo "API documentation is available in README.md"

# Production build
build-prod:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o $(BINARY_NAME) $(MAIN_PATH)

# Help
help:
	@echo "Available commands:"
	@echo "  setup        - Setup development environment"
	@echo "  deps         - Install dependencies"
	@echo "  run          - Run the application"
	@echo "  build        - Build the application"
	@echo "  build-prod   - Build for production"
	@echo "  test         - Run tests"
	@echo "  test-api     - Run API tests"
	@echo "  clean        - Clean build artifacts"
	@echo "  fmt          - Format code"
	@echo "  lint         - Lint code"
	@echo "  db-up        - Start development database"
	@echo "  db-down      - Stop development database"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-run   - Run with Docker Compose"
	@echo "  docker-stop  - Stop Docker containers"
	@echo "  docker-logs  - View Docker logs"
	@echo "  help         - Show this help"

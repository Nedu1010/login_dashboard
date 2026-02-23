.PHONY: run build test migrate-up migrate-down docker-up docker-down lint clean

# Run the application
run:
	@echo "Starting server..."
	go run cmd/server/main.go

# Build the application
build:
	@echo "Building binary..."
	go build -o bin/server cmd/server/main.go

# Run tests with coverage
test:
	@echo "Running tests..."
	go test -v -cover ./...

# Run tests with coverage report
test-coverage:
	@echo "Generating coverage report..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Install golang-migrate tool
install-migrate:
	@echo "Installing golang-migrate..."
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Run migrations up
migrate-up:
	@echo "Running migrations..."
	migrate -path migrations -database "${DATABASE_URL}" up

# Rollback migrations
migrate-down:
	@echo "Rolling back migrations..."
	migrate -path migrations -database "${DATABASE_URL}" down

# Start PostgreSQL in Docker
docker-up:
	@echo "Starting PostgreSQL..."
	docker run --name auth-postgres \
		-e POSTGRES_USER=postgres \
		-e POSTGRES_PASSWORD=postgres \
		-e POSTGRES_DB=auth_db \
		-p 5432:5432 \
		-d postgres:15-alpine

# Stop and remove PostgreSQL container
docker-down:
	@echo "Stopping PostgreSQL..."
	docker stop auth-postgres || true
	docker rm auth-postgres || true

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# Run linter
lint:
	@echo "Running linter..."
	golangci-lint run

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf bin/
	rm -f coverage.out coverage.html

# Show help
help:
	@echo "Available commands:"
	@echo "  make run            - Run the application"
	@echo "  make build          - Build the binary"
	@echo "  make test           - Run tests"
	@echo "  make test-coverage  - Generate coverage report"
	@echo "  make migrate-up     - Run database migrations"
	@echo "  make migrate-down   - Rollback migrations"
	@echo "  make docker-up      - Start PostgreSQL in Docker"
	@echo "  make docker-down    - Stop PostgreSQL container"
	@echo "  make deps           - Install dependencies"
	@echo "  make lint           - Run linter"
	@echo "  make clean          - Clean build artifacts"

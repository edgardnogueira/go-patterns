.PHONY: help build test up down lint proto clean

# Default target
help:
	@echo "Usage:"
	@echo "  make build     - Build all services"
	@echo "  make test      - Run tests for all services"
	@echo "  make up        - Start all services with docker-compose"
	@echo "  make down      - Stop all services"
	@echo "  make lint      - Run linters on all services"
	@echo "  make proto     - Generate protobuf code (requires protoc)"
	@echo "  make clean     - Clean build artifacts"

# Build all services
build:
	@echo "Building API Gateway..."
	@cd api-gateway && go build -o bin/api-gateway ./main.go
	@echo "Building Order Service API..."
	@cd order-service/cmd/api && go build -o ../../../bin/order-service-api ./main.go
	@echo "Building Order Service Worker..."
	@cd order-service/cmd/worker && go build -o ../../../bin/order-service-worker ./main.go
	@echo "Building Inventory Service API..."
	@cd inventory-service/cmd/api && go build -o ../../../bin/inventory-service-api ./main.go
	@echo "Building Inventory Service Worker..."
	@cd inventory-service/cmd/worker && go build -o ../../../bin/inventory-service-worker ./main.go
	@echo "Building Notification Service..."
	@cd notification-service && go build -o bin/notification-service ./main.go

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./api-gateway/...
	@go test -v ./order-service/...
	@go test -v ./inventory-service/...
	@go test -v ./notification-service/...
	@go test -v ./pkg/...

# Start services with docker-compose
up:
	@echo "Starting services..."
	@docker-compose up -d

# Stop services
down:
	@echo "Stopping services..."
	@docker-compose down

# Run linters
lint:
	@echo "Running linters..."
	@golangci-lint run ./...

# Generate protobuf code
proto:
	@echo "Generating protobuf code..."
	@protoc --go_out=. --go-grpc_out=. ./pkg/proto/*.proto

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@find . -name "*.pb.go" -type f -delete

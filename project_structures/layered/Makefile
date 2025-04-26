.PHONY: build run-api run-worker test lint clean help

# Build variables
BUILD_DIR := ./bin
API_NAME := blog-api
WORKER_NAME := blog-worker

help:
	@echo "Usage:"
	@echo "  make build       - Build both API and worker binaries"
	@echo "  make run-api     - Run the API server"
	@echo "  make run-worker  - Run the worker process"
	@echo "  make test        - Run tests"
	@echo "  make lint        - Run linters"
	@echo "  make clean       - Remove build artifacts"

build:
	@echo "Building API and worker..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(API_NAME) ./cmd/api
	go build -o $(BUILD_DIR)/$(WORKER_NAME) ./cmd/worker
	@echo "Build complete! Binaries available in $(BUILD_DIR) directory"

run-api:
	@echo "Starting API server..."
	go run ./cmd/api/main.go

run-worker:
	@echo "Starting worker process..."
	go run ./cmd/worker/main.go

test:
	@echo "Running tests..."
	go test ./... -v

lint:
	@echo "Running linters..."
	golangci-lint run ./...

clean:
	@echo "Cleaning up..."
	rm -rf $(BUILD_DIR)

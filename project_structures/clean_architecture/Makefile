.PHONY: all build clean test run-api run-worker

# Variables
API_BINARY=api
WORKER_BINARY=worker
BUILD_DIR=./build

all: clean build

build: build-api build-worker

build-api:
	@echo "Building API binary..."
	@go build -o $(BUILD_DIR)/$(API_BINARY) ./cmd/api

build-worker:
	@echo "Building Worker binary..."
	@go build -o $(BUILD_DIR)/$(WORKER_BINARY) ./cmd/worker

clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)
	@mkdir -p $(BUILD_DIR)

test:
	@echo "Running tests..."
	@go test -v ./...

run-api: build-api
	@echo "Running API server..."
	@$(BUILD_DIR)/$(API_BINARY)

run-worker: build-worker
	@echo "Running Worker..."
	@$(BUILD_DIR)/$(WORKER_BINARY)

fmt:
	@echo "Formatting code..."
	@go fmt ./...

vet:
	@echo "Vetting code..."
	@go vet ./...

lint:
	@echo "Linting code..."
	@if command -v golangci-lint &> /dev/null; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not installed, skipping lint"; \
	fi

help:
	@echo "Available targets:"
	@echo "  all        - Clean and build all binaries"
	@echo "  build      - Build all binaries"
	@echo "  build-api  - Build API binary"
	@echo "  build-worker - Build Worker binary"
	@echo "  clean      - Clean build directory"
	@echo "  test       - Run tests"
	@echo "  run-api    - Run API server"
	@echo "  run-worker - Run Worker"
	@echo "  fmt        - Format code using go fmt"
	@echo "  vet        - Vet code using go vet"
	@echo "  lint       - Lint code using golangci-lint if available"

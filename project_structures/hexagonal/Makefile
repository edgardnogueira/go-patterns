.PHONY: build run-api run-worker test clean

# Default Go build flags
GO_BUILD_FLAGS := -v

# Directories
BIN_DIR := ./bin
API_BIN := $(BIN_DIR)/api
WORKER_BIN := $(BIN_DIR)/worker

# Go source directories
API_SRC := ./cmd/api
WORKER_SRC := ./cmd/worker

# Create bin directory if it doesn't exist
$(BIN_DIR):
	mkdir -p $(BIN_DIR)

# Build all binaries
build: $(BIN_DIR) build-api build-worker

# Build API server
build-api: $(BIN_DIR)
	go build $(GO_BUILD_FLAGS) -o $(API_BIN) $(API_SRC)

# Build worker
build-worker: $(BIN_DIR)
	go build $(GO_BUILD_FLAGS) -o $(WORKER_BIN) $(WORKER_SRC)

# Run API server
run-api:
	go run $(API_SRC)/main.go

# Run worker
run-worker:
	go run $(WORKER_SRC)/main.go

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf $(BIN_DIR)
	go clean

# Show help
help:
	@echo "Available targets:"
	@echo "  build       - Build all binaries"
	@echo "  build-api   - Build API server"
	@echo "  build-worker - Build worker"
	@echo "  run-api     - Run API server"
	@echo "  run-worker  - Run worker"
	@echo "  test        - Run tests"
	@echo "  clean       - Remove build artifacts"
	@echo "  help        - Show this help message"

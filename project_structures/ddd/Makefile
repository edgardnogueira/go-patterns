# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=myapp
API_BINARY=api
WORKER_BINARY=worker

# Main targets
.PHONY: all build test clean run fmt vet tidy download api worker

all: test build

# Initialize the module
init:
	$(GOMOD) init github.com/edgardnogueira/go-patterns/project_structures/ddd
	$(GOMOD) tidy

# Build all binaries
build: 
	mkdir -p bin
	$(GOBUILD) -o bin/$(API_BINARY) ./cmd/api
	$(GOBUILD) -o bin/$(WORKER_BINARY) ./cmd/worker

# Build API server
api:
	$(GOBUILD) -o bin/$(API_BINARY) ./cmd/api

# Build Worker
worker:
	$(GOBUILD) -o bin/$(WORKER_BINARY) ./cmd/worker

# Run tests
test:
	$(GOTEST) -v ./...

# Run tests with coverage
test-cover:
	$(GOTEST) -v -cover ./...

# Run API server
run-api:
	$(GOCMD) run ./cmd/api

# Run worker
run-worker:
	$(GOCMD) run ./cmd/worker

# Clean build artifacts
clean:
	$(GOCLEAN)
	rm -rf bin

# Format the code
fmt:
	$(GOCMD) fmt ./...

# Vet the code
vet:
	$(GOCMD) vet ./...

# Tidy the dependencies
tidy:
	$(GOMOD) tidy

# Ensure all dependencies are downloaded
download:
	$(GOMOD) download

FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/api-gateway ./project_structures/microservices/api-gateway

# Create a minimal production image
FROM alpine:3.19

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/bin/api-gateway /app/api-gateway

# Copy configuration files if needed
COPY --from=builder /app/project_structures/microservices/api-gateway/docs /app/docs

# Set environment variables
ENV PORT=8080
ENV ORDER_SERVICE_URL=http://order-service:8081
ENV INVENTORY_SERVICE_URL=http://inventory-service:8082

# Expose the port
EXPOSE 8080

# Run the application
ENTRYPOINT ["/app/api-gateway"]

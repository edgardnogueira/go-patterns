FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/order-service-worker ./project_structures/microservices/order-service/cmd/worker

# Create a minimal production image
FROM alpine:3.19

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/bin/order-service-worker /app/order-service-worker

# Set environment variables
ENV POSTGRES_DSN=postgres://postgres:postgres@postgres:5432/orders?sslmode=disable
ENV NATS_URL=nats://nats:4222

# Run the application
ENTRYPOINT ["/app/order-service-worker"]

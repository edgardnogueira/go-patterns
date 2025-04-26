FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/notification-service ./project_structures/microservices/notification-service

# Create a minimal production image
FROM alpine:3.19

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/bin/notification-service /app/notification-service

# Set environment variables
ENV PORT=8083
ENV NATS_URL=nats://nats:4222
ENV SMTP_HOST=mailhog
ENV SMTP_PORT=1025

# Expose the port
EXPOSE 8083

# Run the application
ENTRYPOINT ["/app/notification-service"]

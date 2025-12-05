# Build stage
FROM golang:1.24.4-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the main application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/main.go

# Build the cron application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o cron ./cmd/cron/main.go

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy the binaries from builder
COPY --from=builder /app/main .
COPY --from=builder /app/cron .

# Copy migration files
COPY --from=builder /app/internal/config/di/migrations ./migrations

# Expose port
EXPOSE 8083

# Default command (can be overridden in docker-compose)
CMD ["./main"]

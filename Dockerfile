# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o fintech-service ./cmd/main.go

# Runtime stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/fintech-service .
COPY --from=builder /app/configs ./configs
COPY --from=builder /app/migrations ./migrations

# Expose port and set environment variables
EXPOSE 8080
ENV GIN_MODE=release

# Run the application
CMD ["./fintech-service"]
# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Copy go mod files and shared directory
COPY services/livekit-service/go.mod ./
COPY services/livekit-service/go.sum ./
COPY services/livekit-service/shared/ ./shared/

# Download dependencies
RUN go mod download

# Generate go.sum if it doesn't exist (for fresh builds)
RUN test -f go.sum || go mod tidy

# Copy source code
COPY services/livekit-service/ .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/main.go

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates curl

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Expose port
EXPOSE 8089

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:8089/api/v1/health || exit 1

# Run the binary
CMD ["./main"]
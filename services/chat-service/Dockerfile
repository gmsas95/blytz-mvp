# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY services/chat-service/go.mod ./
COPY services/chat-service/go.sum* ./

# Copy shared package
COPY shared ./shared

# Download dependencies
RUN go mod download || go mod tidy

# Copy service source code
COPY services/chat-service/ ./services/chat-service/

# Build the application from the correct directory
RUN cd services/chat-service && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o chat-service ./cmd/main.go

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates curl

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/services/chat-service/chat-service .

# Expose port
EXPOSE 8088

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3   CMD curl -f http://localhost:8088/health || exit 1

# Run the application
CMD ["./chat-service"]

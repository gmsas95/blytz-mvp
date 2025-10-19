# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app/services/gateway

# Install dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY services/gateway/go.mod ./
COPY services/gateway/go.sum* ./

# Copy shared package
COPY shared ./shared

# Download dependencies
RUN go mod download || go mod tidy

# Copy service source code
COPY services/gateway/ .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o gateway ./cmd/main.go

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates curl

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/services/gateway/gateway .

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3   CMD curl -f http://localhost:8080/health || exit 1

# Run the application
CMD ["./gateway"]

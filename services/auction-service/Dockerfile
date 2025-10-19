# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Create proper directory structure
RUN mkdir -p services shared

# Copy go mod files to service directory
COPY services/auction-service/go.mod ./services/auction-service/
COPY services/auction-service/go.sum* ./services/auction-service/

# Copy shared package to root level
COPY shared ./shared

# Copy service source code
COPY services/auction-service/ ./services/auction-service/

# Create symlink for shared package in the correct location
RUN ln -sf /app/shared /shared

# Download dependencies and build from service directory
RUN cd services/auction-service && go mod download || go mod tidy
RUN cd services/auction-service && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o auction-service ./cmd/main.go

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates curl

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/services/auction-service/auction-service .

# Expose port
EXPOSE 8083

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:8083/health || exit 1

# Run the application
CMD ["./auction-service"]
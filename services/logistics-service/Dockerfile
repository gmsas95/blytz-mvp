# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Create proper directory structure
RUN mkdir -p services shared

# Copy go mod files to service directory
COPY services/logistics-service/go.mod ./services/logistics-service/
COPY services/logistics-service/go.sum* ./services/logistics-service/

# Copy shared package to root level
COPY shared ./shared

# Copy service source code
COPY services/logistics-service/ ./services/logistics-service/

# Create symlink for shared package in the correct location
RUN ln -sf /app/shared /shared

# Download dependencies and build from service directory
RUN cd services/logistics-service && go mod download || go mod tidy
RUN cd services/logistics-service && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o logistics-service ./cmd/main.go

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates curl

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/services/logistics-service/logistics-service .

# Expose port
EXPOSE 8087

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3   CMD curl -f http://localhost:8087/health || exit 1

# Run the application
CMD ["./logistics-service"]

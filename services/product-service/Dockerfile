# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app/services/product-service

# Install dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY services/product-service/go.mod ./
COPY services/product-service/go.sum* ./

# Copy shared package
COPY shared ./shared

# Download dependencies
RUN go mod download || go mod tidy

# Copy service source code
COPY services/product-service/ .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o product-service ./cmd/main.go

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates curl

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/services/product-service/product-service .

# Expose port
EXPOSE 8082

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3   CMD curl -f http://localhost:8082/health || exit 1

# Run the application
CMD ["./product-service"]

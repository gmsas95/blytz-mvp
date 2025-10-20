# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy all go modules and vendor directory
COPY go.mod go.sum go.work go.work.sum ./
COPY vendor ./vendor

# Copy all source code
COPY . .

# Build the application from the correct directory using vendored modules
RUN cd services/payment-service && CGO_ENABLED=0 GOOS=linux go build -mod=vendor -a -installsuffix cgo -o payment-service ./cmd/main.go

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates curl

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/services/payment-service/payment-service .

# Expose port
EXPOSE 8086

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3   CMD curl -f http://localhost:8086/health || exit 1

# Run the application
CMD ["./payment-service"]

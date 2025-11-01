# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy the entire repository
COPY . .

# Copy shared directory from parent
COPY ../shared ./shared

# Build the application
RUN cd services/payment-service && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o payment-service ./cmd/main.go

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

# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Create proper directory structure
RUN mkdir -p services/auth-service

# Copy go mod files to service directory
COPY services/auth-service/go.mod services/auth-service/go.sum* ./services/auth-service/

# Copy shared package to root level
COPY shared ./shared

# Copy service source code
COPY services/auth-service/ ./services/auth-service/

# Download dependencies and build from service directory
RUN cd services/auth-service && go mod tidy
RUN cd services/auth-service && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o auth-service ./cmd/main.go

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN addgroup -g 1001 -S appuser && \
    adduser -S -u 1001 -G appuser appuser

# Set working directory
WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/services/auth-service/auth-service .

# Change ownership to non-root user
RUN chown appuser:appuser auth-service

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the application
CMD ["./auth-service"]
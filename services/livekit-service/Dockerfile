# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy the entire repository
COPY . .

# Build the application
RUN cd services/livekit-service && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/main.go

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates curl

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/services/livekit-service/main .

# Expose port
EXPOSE 8089

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:8089/api/v1/health || exit 1

# Run the binary
CMD ["./main"]
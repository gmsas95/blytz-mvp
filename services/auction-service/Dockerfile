# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy module files and download dependencies
COPY go.work go.work.sum ./
COPY services/auction-service/go.mod services/auction-service/go.sum ./services/auction-service/
COPY shared/go.mod shared/go.sum ./shared/
RUN cd services/auction-service && go mod download

# Copy the rest of the source code
COPY . .

# Build the application
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
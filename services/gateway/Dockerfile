# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy module files and download dependencies
COPY go.work go.work.sum ./
COPY services/gateway/go.mod services/gateway/go.sum ./services/gateway/
COPY shared/go.mod shared/go.sum ./shared/
RUN cd services/gateway && go mod download

# Copy the rest of the source code
COPY . .

# Build the application
RUN cd services/gateway && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o gateway ./cmd/main.go

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

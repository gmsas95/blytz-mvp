# 🚀 Blytz MVP - Local Development Makefile
# Mirrors the GitHub Actions pipeline for local testing

.PHONY: help test test-all build build-all deploy-local load-test clean

# Default values
GO_VERSION := 1.25
SERVICES := auth-service product-service auction-service
SERVICE_PATHS := backend/auth-service backend/product-service backend/auction-service

# Colors for output
RED := \033[31m
GREEN := \033[32m
YELLOW := \033[33m
BLUE := \033[34m
RESET := \033[0m

help: ## 📖 Show this help message
	@echo "$(BLUE)🚀 Blytz MVP - Local Development Commands$(RESET)"
	@echo ""
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "$(GREEN)%-20s$(RESET) %s\n", $$1, $$2}' $(MAKEFILE_LIST)

test: ## 🧪 Test a specific service (make test SERVICE=auth-service)
	@if [ -z "$(SERVICE)" ]; then echo "$(RED)❌ Please specify SERVICE=$(RESET)"; exit 1; fi
	@echo "$(BLUE)🧪 Testing $(SERVICE)...$(RESET)"
	cd backend/$(SERVICE) && go mod download && go mod tidy
	cd backend/$(SERVICE) && go build -v ./...
	cd backend/$(SERVICE) && go test ./... -v -race -coverprofile=coverage.out
	cd backend/$(SERVICE) && go tool cover -func=coverage.out

test-all: ## 🧪 Test all services (parallel execution)
	@echo "$(BLUE)🧪 Testing all services...$(RESET)"
	@for service in $(SERVICES); do \
		echo "$(YELLOW)Testing $$service...$(RESET)"; \
		$(MAKE) test SERVICE=$$service & \
	done; \
	wait; \
	echo "$(GREEN)✅ All tests completed!$(RESET)"

build: ## 🏗️ Build Docker image for specific service (make build SERVICE=auth-service)
	@if [ -z "$(SERVICE)" ]; then echo "$(RED)❌ Please specify SERVICE=$(RESET)"; exit 1; fi
	@echo "$(BLUE)🏗️ Building Docker image for $(SERVICE)...$(RESET)"
	docker build -f backend/$(SERVICE)/Dockerfile -t blytz-$(SERVICE):local ./backend

build-all: ## 🏗️ Build all service Docker images
	@echo "$(BLUE)🏗️ Building all Docker images...$(RESET)"
	@for service in $(SERVICES); do \
		echo "$(YELLOW)Building $$service...$(RESET)"; \
		docker build -f backend/$$service/Dockerfile -t blytz-$$service:local ./backend; \
	done; \
	echo "$(GREEN)✅ All images built!$(RESET)"

deploy-local: ## 🚀 Deploy locally with docker-compose
	@echo "$(BLUE)🚀 Starting local deployment...$(RESET)"
	docker compose up -d --build
	@echo "$(YELLOW)⏳ Waiting for services to be healthy...$(RESET)"
	@sleep 10
	@echo "$(GREEN)✅ Running health checks...$(RESET)"
	@for port in 8081 8082 8083; do \
		if curl -s http://localhost:$$port/health | grep -q "ok"; then \
			echo "$(GREEN)✅ Port $$port: Healthy$(RESET)"; \
		else \
			echo "$(RED)❌ Port $$port: Unhealthy$(RESET)"; \
		fi; \
	done

health-check: ## 🔍 Quick health check of all services
	@echo "$(BLUE)🔍 Checking service health...$(RESET)"
	@for port in 8081 8082 8083; do \
		echo -n "Port $$port: "; \
		if curl -s http://localhost:$$port/health | grep -q "ok"; then \
			echo "$(GREEN)✅ Healthy$(RESET)"; \
		else \
			echo "$(RED)❌ Unhealthy$(RESET)"; \
		fi; \
	done

# Simple load test without k6 dependency
load-test: ## ⚡ Simple load test with curl
	@echo "$(BLUE)⚡ Running simple load test...$(RESET)"
	@for i in {1..10}; do \
		start_time=$$(date +%s%N); \
		response=$$(curl -s -w "%{http_code}" http://localhost:8083/api/v1/auctions); \
		end_time=$$(date +%s%N); \
		duration=$$(( (end_time - start_time) / 1000000 )); \
		if [ "$$response" = "200" ]; then \
			echo "Request $$i: ✅ $$duration ms"; \
		else \
			echo "Request $$i: ❌ $$response"; \
		fi; \
	done

clean: ## 🧹 Clean up Docker containers and images
	@echo "$(BLUE)🧹 Cleaning up...$(RESET)"
	docker compose down
	echo "$(GREEN)✅ Cleanup completed!$(RESET)"

logs: ## 📋 Show logs for all services
	@echo "$(BLUE)📋 Showing service logs...$(RESET)"
	docker compose logs -f

# 🔄 CI/CD Pipeline Equivalents
ci-test: test-all ## 🔄 Run CI test stage locally
ci-build: build-all ## 🔄 Run CI build stage locally
ci-deploy: deploy-local ## 🔄 Run CI deploy stage locally
ci-load-test: load-test ## 🔄 Run CI load test stage locally

# 🎯 Full Pipeline
ci-pipeline: test-all build-all deploy-local health-check load-test ## 🔄 Run complete CI pipeline locally

# 📊 Development Utilities
dev-setup: ## 🔧 Setup development environment
	@echo "$(BLUE)🔧 Setting up development environment...$(RESET)"
	@echo "Installing Go $(GO_VERSION)..."
	@echo "Installing Docker..."
	@echo "Installing k6..."
	@echo "$(GREEN)✅ Development environment ready!$(RESET)"

status: ## 📊 Show system status
	@echo "$(BLUE)📊 System Status:$(RESET)"
	docker compose ps
	@echo ""
	@echo "$(BLUE)Docker Images:$(RESET)"
	docker images | grep blytz || echo "No blytz images found"
	@echo ""
	@echo "$(BLUE)Resource Usage:$(RESET)"
	docker stats --no-stream --format "table {{.Name}}	{{.CPUPerc}}	{{.MemUsage}}" 2>/dev/null || echo "No running containers""file_path":"/home/sas/blytzmvp-clean/Makefile
#!/bin/bash
# 🚀 Local CI/CD Pipeline Script - Mirrors GitHub Actions workflow

set -e  # Exit on any error

# Colors for output
RED='\033[31m'
GREEN='\033[32m'
YELLOW='\033[33m'
BLUE='\033[34m'
RESET='\033[0m'

# Services to test
SERVICES=("auth-service" "product-service" "auction-service")
SERVICE_PATHS=("backend/auth-service" "backend/product-service" "backend/auction-service")
SERVICE_PORTS=("8081" "8082" "8083")

# Helper functions
log_info() {
    echo -e "${BLUE}[INFO]${RESET} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${RESET} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${RESET} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${RESET} $1"
}

# Check if required tools are installed
check_dependencies() {
    log_info "Checking dependencies..."

    local missing_deps=()

    if ! command -v go &> /dev/null; then
        missing_deps+=("go")
    fi

    if ! command -v docker &> /dev/null; then
        missing_deps+=("docker")
    fi

    if ! command -v curl &> /dev/null; then
        missing_deps+=("curl")
    fi

    if ! command -v k6 &> /dev/null; then
        missing_deps+=("k6")
    fi

    if [ ${#missing_deps[@]} -ne 0 ]; then
        log_error "Missing dependencies: ${missing_deps[*]}"
        log_info "Please install missing dependencies and try again"
        exit 1
    fi

    log_success "All dependencies found"
}

# Test individual service
test_service() {
    local service_name=$1
    local service_path=$2
    local service_port=$3

    log_info "Testing $service_name..."

    cd "$service_path" || exit 1

    # Install dependencies
    log_info "Installing dependencies for $service_name..."
    go mod download || { log_error "Failed to download dependencies for $service_name"; return 1; }
    go mod tidy || { log_error "Failed to tidy dependencies for $service_name"; return 1; }

    # Build service
    log_info "Building $service_name..."
    go build -v ./... || { log_error "Failed to build $service_name"; return 1; }

    # Run tests
    log_info "Running tests for $service_name..."
    go test ./... -v -race -coverprofile=coverage.out || { log_error "Tests failed for $service_name"; return 1; }

    # Show coverage
    log_info "Coverage report for $service_name:"
    go tool cover -func=coverage.out | tail -5

    cd - > /dev/null || exit 1
    log_success "$service_name tests passed"
}

# Test all services (parallel)
test_all_services() {
    log_info "Testing all services in parallel..."

    local pids=()
    local failed_services=()

    for i in "${!SERVICES[@]}"; do
        service_name=${SERVICES[$i]}
        service_path=${SERVICE_PATHS[$i]}
        service_port=${SERVICE_PORTS[$i]}

        test_service "$service_name" "$service_path" "$service_port" &
        pids+=($!)
    done

    # Wait for all background jobs
    for pid in "${pids[@]}"; do
        if ! wait "$pid"; then
            failed_services+=("${SERVICES[$i]}")
        fi
    done

    if [ ${#failed_services[@]} -ne 0 ]; then
        log_error "Some services failed tests: ${failed_services[*]}"
        return 1
    fi

    log_success "All services passed tests"
}

# Build Docker images
build_images() {
    log_info "Building Docker images..."

    for i in "${!SERVICES[@]}"; do
        service_name=${SERVICES[$i]}
        dockerfile="backend/$service_name/Dockerfile"

        log_info "Building Docker image for $service_name..."
        docker build -f "$dockerfile" -t "blytz-$service_name:local" ./backend || {
            log_error "Failed to build Docker image for $service_name"
            return 1
        }
    done

    log_success "All Docker images built successfully"
}

# Deploy locally
deploy_local() {
    log_info "Deploying locally with docker-compose..."

    docker compose up -d --build || {
        log_error "Failed to start services"
        return 1
    }

    log_info "Waiting for services to be healthy..."
    sleep 30

    # Health checks
    log_info "Running health checks..."
    for i in "${!SERVICES[@]}"; do
        service_name=${SERVICES[$i]}
        service_port=${SERVICE_PORTS[$i]}

        if curl -s "http://localhost:$service_port/health" | grep -q "ok"; then
            log_success "$service_name is healthy (port $service_port)"
        else
            log_error "$service_name health check failed (port $service_port)"
            return 1
        fi
    done
}

# Run load tests
load_test() {
    log_info "Running load tests..."

    # Create load test script
    cat > /tmp/blytz-load-test.js << 'EOF'
import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
  stages: [
    { duration: '10s', target: 10 },   // Ramp up to 10 users
    { duration: '20s', target: 25 },   // Stay at 25 users
    { duration: '10s', target: 0 },    // Ramp down to 0 users
  ],
  thresholds: {
    http_req_duration: ['p(95)<300'],  // 95% of requests under 300ms
    http_req_failed: ['rate<0.1'],     // Error rate under 10%
  },
};

export default function() {
  // Test auction listing endpoint
  let response = http.get('http://localhost:8083/api/v1/auctions');

  check(response, {
    'status is 200': (r) => r.status === 200,
    'response time < 300ms': (r) => r.timings.duration < 300,
    'has success field': (r) => JSON.parse(r.body).success === true,
  });

  sleep(1);
}
EOF

    if command -v k6 &> /dev/null; then
        k6 run /tmp/blytz-load-test.js
        log_success "Load test completed"
    else
        log_warning "k6 not found, skipping load test"
        log_info "Install k6 with: sudo apt-get install k6"
    fi
}

# Show system status
show_status() {
    log_info "System Status:"
    docker compose ps
    echo ""
    log_info "Health Summary:"
    for i in "${!SERVICES[@]}"; do
        service_name=${SERVICES[$i]}
        service_port=${SERVICE_PORTS[$i]}

        if curl -s "http://localhost:$service_port/health" | grep -q "ok"; then
            log_success "$service_name (port $service_port): Healthy"
        else
            log_error "$service_name (port $service_port): Unhealthy"
        fi
    done
}

# Cleanup
cleanup() {
    log_info "Cleaning up..."
    docker compose down
    docker system prune -f
    rm -f /tmp/blytz-load-test.js
    log_success "Cleanup completed"
}

# Full CI pipeline
run_pipeline() {
    log_info "🚀 Starting local CI/CD pipeline..."

    check_dependencies
    test_all_services
    build_images
    deploy_local
    load_test
    show_status

    log_success "🎉 Local CI/CD pipeline completed successfully!"
}

# Parse command line arguments
case "${1:-pipeline}" in
    "test")
        test_service "${SERVICES[0]}" "${SERVICE_PATHS[0]}" "${SERVICE_PORTS[0]}"
        ;;
    "test-all")
        test_all_services
        ;;
    "build")
        build_images
        ;;
    "deploy")
        deploy_local
        ;;
    "load-test")
        load_test
        ;;
    "status")
        show_status
        ;;
    "clean")
        cleanup
        ;;
    "pipeline"|"")
        run_pipeline
        ;;
    "help"|"-h"|"--help")
        echo "🚀 Blytz MVP Local CI/CD Pipeline"
        echo ""
        echo "Usage: $0 [command]"
        echo ""
        echo "Commands:"
        echo "  pipeline      Run complete CI/CD pipeline (default)"
        echo "  test          Test a specific service"
        echo "  test-all      Test all services"
        echo "  build         Build all Docker images"
        echo "  deploy        Deploy locally with docker-compose"
        echo "  load-test     Run load tests"
        echo "  status        Show system status"
        echo "  clean         Clean up resources"
        echo "  help          Show this help message"
        echo ""
        echo "Examples:"
        echo "  $0                    # Run complete pipeline"
        echo "  $0 test-all           # Test all services"
        echo "  $0 build              # Build Docker images"
        echo "  $0 deploy             # Deploy locally"
        ;;
    *)
        log_error "Unknown command: $1"
        echo "Use '$0 help' for usage information"
        exit 1
        ;;
esac

exit 0
EOF

chmod +x /home/sas/blytzmvp-clean/scripts/local-ci.sh

# Create a simpler wrapper script
Write-Content -FilePath /home/sas/blytzmvp-clean/scripts/test.sh -Content "#!/bin/bash
# Quick test script for individual services

SERVICE=${1:-all}

if [ "$SERVICE" == "all" ]; then
    ./scripts/local-ci.sh test-all
else
    ./scripts/local-ci.sh test SERVICE=$SERVICE
fi
"
chmod +x /home/sas/blytzmvp-clean/scripts/test.sh

log_success "Local CI/CD scripts created successfully!"
log_info "Usage: ./scripts/local-ci.sh [command]"
log_info "Try: make help" "for Makefile commands"
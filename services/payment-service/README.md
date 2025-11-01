# Blytz Payment Service with Fiuu Integration

A production-ready payment processing service built with Go, featuring comprehensive Fiuu payment gateway integration, advanced monitoring, security compliance, and operational excellence.

## 🚀 Features

### Core Functionality
- **Complete Fiuu Integration**: 110+ payment channels across 6 Southeast Asian countries
- **Multiple Payment Methods**: FPX, e-wallets (GrabPay, Touch 'n Go, ShopeePay), credit cards, BNPL, QR codes
- **Seamless Integration**: JavaScript-based frontend integration with no redirects
- **Real-time Processing**: Sub-200ms payment processing with Redis caching
- **Comprehensive Refunds**: Full refund lifecycle management

### Advanced Features
- **Circuit Breaker Pattern**: Automatic failover and recovery for external dependencies
- **Retry Mechanisms**: Intelligent retry with exponential backoff and jitter
- **Resilient Architecture**: Graceful degradation under load
- **Rate Limiting**: DDoS protection and abuse prevention
- **Webhook Processing**: Asynchronous payment status updates with retry logic

### Security & Compliance
- **PCI DSS Compliance**: Full compliance validation and enforcement
- **Data Encryption**: AES-256 encryption for sensitive data at rest
- **PII Masking**: Automatic masking of personally identifiable information
- **Audit Logging**: Comprehensive audit trails for security events
- **Access Control**: Role-based access control with least privilege principle

### Monitoring & Observability
- **Prometheus Metrics**: 50+ custom metrics for comprehensive monitoring
- **Structured Logging**: JSON-based logging with trace IDs and correlation
- **Performance Monitoring**: Response time percentiles and error rate tracking
- **Business Metrics**: Revenue, success rates, and payment method analytics
- **Health Checks**: Multi-level health checks for dependencies

### Database & Data Management
- **Enhanced Schema**: Comprehensive database schema with Fiuu-specific fields
- **Data Retention**: Automated data retention policies
- **Migration System**: Version-controlled database migrations
- **Seed Data**: Complete seed data for testing and development

## 📊 Architecture Overview

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Frontend      │    │   API Gateway   │    │  Payment Service│
│   React Native  │────│   Nginx/Go      │────│   Go + Gin      │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                                        │
                                                        ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   PostgreSQL    │    │      Redis      │    │     Fiuu API    │
│   Primary DB    │    │    Cache/Store  │    │  Payment Gateway│
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

### Service Components

1. **Payment API**: RESTful API for payment processing
2. **Fiuu Client**: Resilient client with circuit breaker and retries
3. **Database Layer**: PostgreSQL with comprehensive schema
4. **Cache Layer**: Redis for performance and session management
5. **Monitoring Layer**: Prometheus metrics and structured logging
6. **Security Layer**: Compliance validation and audit logging

## 🛠️ Technology Stack

- **Backend**: Go 1.23+, Gin Web Framework
- **Database**: PostgreSQL 13+ with GORM ORM
- **Cache**: Redis 6.0+
- **Monitoring**: Prometheus, Grafana, Loki
- **Containerization**: Docker, Kubernetes
- **Security**: PCI DSS compliance, AES-256 encryption
- **Testing**: Comprehensive unit and integration tests
- **CI/CD**: GitHub Actions with automated deployment

## 📦 Installation

### Prerequisites

- Go 1.23+
- PostgreSQL 13+
- Redis 6.0+
- Docker (optional)
- Kubernetes (for production)

### Quick Start

```bash
# Clone the repository
git clone https://github.com/gmsas95/blytz-mvp.git
cd blytzmvp-clean/services/payment-service

# Install dependencies
go mod download

# Set up environment variables
cp .env.example .env
# Edit .env with your configuration

# Run database migrations
go run cmd/migrate/main.go up

# Seed test data
go run cmd/seed/main.go

# Start the service
go run cmd/main.go
```

### Docker Setup

```bash
# Build Docker image
docker build -t payment-service .

# Run with Docker Compose
docker-compose up -d

# Check service health
curl http://localhost:8086/health
```

## ⚙️ Configuration

### Environment Variables

```bash
# Application
PORT=8086
LOG_LEVEL=info
ENVIRONMENT=production

# Database
DB_HOST=localhost
DB_PORT=5432
DB_NAME=payment_service
DB_USER=payment_service
DB_PASSWORD=your_password
DB_SSL_MODE=require

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=your_redis_password

# Fiuu Payment Gateway
FIUU_MERCHANT_ID=your_merchant_id
FIUU_VERIFY_KEY=your_verify_key
FIUU_SANDBOX=false
FIUU_TIMEOUT=30s

# Security
ENCRYPTION_KEY=your_32_character_encryption_key
JWT_SECRET=your_jwt_secret
```

### Payment Method Configuration

The service supports 110+ payment methods across Southeast Asia:

#### Malaysia (MYR)
- **Bank Transfer**: FPX, FPX Business Banking
- **E-wallets**: GrabPay, Touch 'n Go, ShopeePay, Boost
- **QR Codes**: MAE QR, DuitNow QR
- **BNPL**: Atome, Rely
- **Credit Cards**: Visa, Mastercard

#### Singapore (SGD)
- **QR Payment**: PayNow
- **Bank Transfer**: eNETS
- **Credit Cards**: Visa, Mastercard

#### Philippines (PHP)
- **E-wallets**: GCash, Maya
- **Bank Transfer**: BPI, UnionBank

## 🧪 Testing

### Run All Tests

```bash
# Run comprehensive test suite
./tests/run_tests.sh

# Run specific test categories
go test ./pkg/fiuu/ -v
go test ./internal/services/ -v
go test ./internal/api/handlers/ -v
```

### Test Coverage

The service maintains >90% test coverage across all modules:

- **Unit Tests**: Individual component testing
- **Integration Tests**: API endpoint testing
- **Performance Tests**: Load testing with k6
- **Security Tests**: Compliance validation testing

### Load Testing

```bash
# Install k6
curl https://github.com/grafana/k6/releases/download/v0.47.0/k6-v0.47.0-linux-amd64.tar.gz -L | tar xvz

# Run load tests
k6 run tests/load/payment-api-test.js
```

## 📈 Monitoring

### Key Metrics

#### Application Metrics
- `payments_total`: Total payment attempts by method and status
- `payment_duration_seconds`: Payment processing time distribution
- `payment_success_rate`: Success rate by payment method
- `fiuu_circuit_breaker_state`: Circuit breaker state (0=CLOSED, 1=OPEN)

#### Infrastructure Metrics
- `http_requests_total`: HTTP request count by endpoint and status
- `http_request_duration_seconds`: Request processing time
- `db_connections_active`: Active database connections
- `memory_usage_bytes`: Memory consumption

#### Business Metrics
- `revenue_total`: Total revenue processed
- `active_users`: Number of active users
- `top_payment_methods`: Most used payment methods

### Grafana Dashboard

Access the comprehensive monitoring dashboard at:
- **Production**: `https://grafana.yourdomain.com/d/payment-service`
- **Staging**: `https://grafana-staging.yourdomain.com/d/payment-service`

### Alerting

Key alerts configured:

- **Service Down**: Payment service unavailable
- **High Error Rate**: Error rate >5%
- **Slow Response**: P95 response time >500ms
- **Fiuu Circuit Breaker Open**: Fiuu gateway circuit breaker open
- **Database Connection Issues**: Connection pool exhaustion

## 🔒 Security

### PCI DSS Compliance

The service implements comprehensive PCI DSS compliance:

- **Data Encryption**: AES-256 encryption for sensitive data
- **Access Control**: Role-based access with audit logging
- **Network Security**: TLS 1.2+ for all communications
- **Vulnerability Management**: Regular security scans and patches

### Security Features

- **PII Masking**: Automatic masking of sensitive data in logs
- **Audit Logging**: Complete audit trail for all payment operations
- **Rate Limiting**: Protection against abuse and DDoS attacks
- **Input Validation**: Comprehensive validation of all inputs
- **Secure Headers**: Security headers for all HTTP responses

### Security Monitoring

```bash
# Monitor security events
kubectl logs deployment/payment-service -n payment-service | grep security_event

# Check authentication failures
curl https://payments.yourdomain.com/metrics | grep failed_auth_total
```

## 🚀 Deployment

### Production Deployment

```bash
# Deploy to Kubernetes
kubectl apply -f k8s/namespace.yaml
kubectl apply -f k8s/configmaps.yaml
kubectl apply -f k8s/secrets.yaml
kubectl apply -f k8s/deployment.yaml
kubectl apply -f k8s/service.yaml
kubectl apply -f k8s/ingress.yaml

# Monitor deployment
kubectl rollout status deployment/payment-service -n payment-service
```

### Environment Configuration

The service supports multiple deployment environments:

- **Development**: Local development with Docker Compose
- **Staging**: Production-like environment for testing
- **Production**: High-availability deployment with 3 replicas

### Blue-Green Deployment

```bash
# Deploy to green environment
kubectl apply -f k8s/green-deployment.yaml

# Switch traffic to green
kubectl patch service payment-service -p '{"spec":{"selector":{"version":"green"}}}'

# Verify green deployment
kubectl rollout status deployment/payment-service-green -n payment-service

# Clean up blue environment
kubectl delete deployment payment-service-blue -n payment-service
```

## 📚 API Documentation

### Payment Endpoints

#### Create Payment
```http
POST /api/v1/payments
Content-Type: application/json
Authorization: Bearer <token>

{
  "order_id": "ORD123456",
  "amount": 100.50,
  "currency": "MYR",
  "payment_method": "FPX",
  "bill_name": "John Doe",
  "bill_email": "john@example.com",
  "bill_mobile": "01234567890",
  "bill_description": "Test Payment",
  "return_url": "https://example.com/return",
  "notify_url": "https://example.com/webhook"
}
```

#### Get Payment Status
```http
GET /api/v1/payments/{payment_id}
Authorization: Bearer <token>
```

#### Process Refund
```http
POST /api/v1/payments/{payment_id}/refund
Content-Type: application/json
Authorization: Bearer <token>

{
  "amount": 50.00,
  "refund_reason": "Customer requested refund"
}
```

#### Seamless Configuration
```http
GET /api/v1/payments/seamless/config?order_id=ORD123456&amount=100.50&payment_method=FPX
Authorization: Bearer <token>
```

### Webhook Handling

#### Payment Status Webhook
```http
POST /api/v1/webhooks/fiuu
Content-Type: application/json
X-Fiuu-Signature: <signature>

{
  "transaction_id": "TXN123456789",
  "order_id": "ORD123456",
  "amount": 100.50,
  "currency": "MYR",
  "payment_status": "1",
  "error_code": "",
  "error_description": "",
  "signature": "<calculated_signature>"
}
```

## 🛠️ Troubleshooting

### Common Issues

#### Service Unavailable
```bash
# Check pod status
kubectl get pods -n payment-service

# Check logs
kubectl logs deployment/payment-service -n payment-service --tail=100

# Check health endpoint
curl https://payments.yourdomain.com/health
```

#### High Error Rate
```bash
# Check error logs
kubectl logs deployment/payment-service -n payment-service | grep ERROR | tail -50

# Check metrics
curl https://payments.yourdomain.com/metrics | grep error_total

# Check circuit breaker state
curl https://payments.yourdomain.com/metrics | grep circuit_breaker
```

#### Database Issues
```bash
# Check database connections
kubectl exec -it <pod-name> -n payment-service -- psql -h $DB_HOST -U $DB_USER -d $DB_NAME -c "SELECT count(*) FROM pg_stat_activity;"

# Check slow queries
kubectl exec -it <pod-name> -n payment-service -- psql -h $DB_HOST -U $DB_USER -d $DB_NAME -c "SELECT query, mean_time FROM pg_stat_statements ORDER BY mean_time DESC LIMIT 10;"
```

### Performance Tuning

#### Database Optimization
```sql
-- Add missing indexes
CREATE INDEX CONCURRENTLY idx_payments_user_id_created_at ON payments(user_id, created_at);

-- Update statistics
ANALYZE payments;

-- Check query performance
EXPLAIN ANALYZE SELECT * FROM payments WHERE user_id = 'user123' ORDER BY created_at DESC LIMIT 10;
```

#### Application Optimization
```bash
# Enable GC tuning
export GOGC=100
export GOMEMLIMIT=500MiB

# Monitor memory usage
curl https://payments.yourdomain.com/metrics | grep memory_usage
```

## 📋 Operational Procedures

### Daily Operations

- **Health Checks**: Monitor service health and key metrics
- **Log Review**: Check for errors and unusual patterns
- **Performance Monitoring**: Review response times and error rates
- **Backup Verification**: Ensure backups are completing successfully

### Weekly Operations

- **Performance Review**: Analyze performance trends
- **Security Review**: Check security events and alerts
- **Capacity Planning**: Review resource utilization
- **Update Management**: Apply security patches and updates

### Monthly Operations

- **Full System Review**: Comprehensive system health check
- **Disaster Recovery Test**: Test backup and recovery procedures
- **Security Audit**: Review security controls and compliance
- **Documentation Update**: Update runbooks and documentation

## 📞 Support

### Emergency Contacts

- **On-call Engineer**: +1-XXX-XXX-XXXX
- **Engineering Manager**: +1-XXX-XXX-XXXX
- **Security Team**: security@yourdomain.com
- **Fiuu Support**: support@fiuu.com

### Documentation

- **Deployment Guide**: [DEPLOYMENT_GUIDE.md](docs/DEPLOYMENT_GUIDE.md)
- **Operational Runbook**: [OPERATIONAL_RUNBOOK.md](docs/OPERATIONAL_RUNBOOK.md)
- **API Documentation**: [API_REFERENCE.md](docs/API_REFERENCE.md)
- **Security Guidelines**: [SECURITY_GUIDELINES.md](docs/SECURITY_GUIDELINES.md)

### Monitoring Dashboards

- **Service Metrics**: Grafana Dashboard
- **Error Tracking**: Sentry Integration
- **Performance**: APM Integration
- **Security**: Security Information and Event Management (SIEM)

## 🤝 Contributing

### Development Setup

```bash
# Clone repository
git clone https://github.com/gmsas95/blytz-mvp.git
cd blytzmvp-clean/services/payment-service

# Install development dependencies
go mod download
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run linter
golangci-lint run

# Run tests
go test ./... -v -race -cover

# Run integration tests
go test ./tests/integration/... -v
```

### Code Standards

- Follow Go best practices and idiomatic code
- Maintain >90% test coverage
- Use structured logging with appropriate context
- Implement comprehensive error handling
- Write clear, concise documentation

### Pull Request Process

1. Create feature branch from `main`
2. Implement changes with tests
3. Ensure all tests pass and coverage is maintained
4. Update documentation as needed
5. Submit pull request with detailed description
6. Address review feedback promptly

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- **Fiuu**: Payment gateway integration and support
- **Go Team**: Excellent programming language and tools
- **PostgreSQL**: Reliable database for transaction processing
- **Redis**: High-performance caching solution
- **Prometheus**: Comprehensive monitoring solution
- **Docker/Kubernetes**: Container orchestration platform

---

**Version**: 1.0.0
**Last Updated**: 2024-01-15
**Maintainers**: Blytz Engineering Team
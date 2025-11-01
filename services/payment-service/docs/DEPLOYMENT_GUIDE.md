# Payment Service Deployment Guide

## Overview

This guide covers the deployment of the Blytz Payment Service with Fiuu integration for production environments. The service is designed for high availability, security, and compliance with PCI DSS requirements.

## Prerequisites

### Infrastructure Requirements

- **Kubernetes Cluster**: v1.24+ or Docker Swarm
- **Database**: PostgreSQL 13+ with TLS enabled
- **Redis**: v6.0+ for caching and session management
- **Load Balancer**: TLS termination with valid SSL certificates
- **Monitoring**: Prometheus + Grafana stack
- **Logging**: ELK stack or Loki + Grafana
- **Secrets Management**: HashiCorp Vault or Kubernetes Secrets

### Security Requirements

- **PCI DSS Compliance**: Network segmentation, access controls
- **Encryption**: TLS 1.2+ for all external communications
- **Data Encryption**: AES-256 for sensitive data at rest
- **Audit Logging**: Immutable logs for security events
- **Rate Limiting**: DDoS protection and abuse prevention

## Configuration

### Environment Variables

```bash
# Application Configuration
PORT=8086
LOG_LEVEL=info
LOG_FORMAT=json
ENVIRONMENT=production

# Database Configuration
DB_HOST=postgres.example.com
DB_PORT=5432
DB_NAME=payment_service
DB_USER=payment_service
DB_SSL_MODE=require
DB_MAX_CONNECTIONS=25
DB_MAX_IDLE_CONNECTIONS=5
DB_CONNECTION_LIFETIME=1h

# Redis Configuration
REDIS_HOST=redis.example.com
REDIS_PORT=6379
REDIS_PASSWORD=redis_password
REDIS_DB=0
REDIS_MAX_CONNECTIONS=10

# Fiuu Configuration
FIUU_MERCHANT_ID=your_merchant_id
FIUU_VERIFY_KEY=your_verify_key
FIUU_SANDBOX=false
FIUU_TIMEOUT=30s
FIUU_RETRY_ATTEMPTS=3
FIUU_RETRY_DELAY=1s

# Security Configuration
ENCRYPTION_KEY=your_32_character_encryption_key
JWT_SECRET=your_jwt_secret
ALLOWED_ORIGINS=https://yourdomain.com

# Monitoring Configuration
METRICS_ENABLED=true
METRICS_PORT=9090
HEALTH_CHECK_ENABLED=true
HEALTH_CHECK_PORT=8087

# Webhook Configuration
WEBHOOK_SECRET=webhook_signature_secret
WEBHOOK_TIMEOUT=30s
WEBHOOK_RETRY_ATTEMPTS=5
WEBHOOK_RETRY_DELAY=60s
```

### Kubernetes Deployment

**Namespace Configuration:**
```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: payment-service
  labels:
    name: payment-service
    environment: production
```

**Secret Configuration:**
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: payment-service-secrets
  namespace: payment-service
type: Opaque
data:
  db-password: <base64-encoded-password>
  redis-password: <base64-encoded-password>
  fiuu-verify-key: <base64-encoded-key>
  encryption-key: <base64-encoded-key>
  jwt-secret: <base64-encoded-secret>
  webhook-secret: <base64-encoded-secret>
```

**ConfigMap Configuration:**
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: payment-service-config
  namespace: payment-service
data:
  PORT: "8086"
  LOG_LEVEL: "info"
  LOG_FORMAT: "json"
  ENVIRONMENT: "production"
  DB_HOST: "postgres.payment-service.svc.cluster.local"
  DB_PORT: "5432"
  DB_NAME: "payment_service"
  DB_USER: "payment_service"
  DB_SSL_MODE: "require"
  REDIS_HOST: "redis.payment-service.svc.cluster.local"
  REDIS_PORT: "6379"
  FIUU_MERCHANT_ID: "your_merchant_id"
  FIUU_SANDBOX: "false"
  FIUU_TIMEOUT: "30s"
  METRICS_ENABLED: "true"
  METRICS_PORT: "9090"
```

**Deployment Configuration:**
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: payment-service
  namespace: payment-service
  labels:
    app: payment-service
    version: v1.0.0
spec:
  replicas: 3
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 0
  selector:
    matchLabels:
      app: payment-service
  template:
    metadata:
      labels:
        app: payment-service
        version: v1.0.0
      annotations:
        prometheus.io/scrape: "true"
        prometheus.io/port: "9090"
        prometheus.io/path: "/metrics"
    spec:
      securityContext:
        runAsNonRoot: true
        runAsUser: 1001
        fsGroup: 1001
      containers:
      - name: payment-service
        image: your-registry/payment-service:v1.0.0
        imagePullPolicy: Always
        ports:
        - containerPort: 8086
          name: http
          protocol: TCP
        - containerPort: 9090
          name: metrics
          protocol: TCP
        - containerPort: 8087
          name: health
          protocol: TCP
        envFrom:
        - configMapRef:
            name: payment-service-config
        env:
        - name: DB_PASSWORD
          valueFrom:
            secretKeyRef:
              name: payment-service-secrets
              key: db-password
        - name: REDIS_PASSWORD
          valueFrom:
            secretKeyRef:
              name: payment-service-secrets
              key: redis-password
        - name: FIUU_VERIFY_KEY
          valueFrom:
            secretKeyRef:
              name: payment-service-secrets
              key: fiuu-verify-key
        - name: ENCRYPTION_KEY
          valueFrom:
            secretKeyRef:
              name: payment-service-secrets
              key: encryption-key
        - name: JWT_SECRET
          valueFrom:
            secretKeyRef:
              name: payment-service-secrets
              key: jwt-secret
        - name: WEBHOOK_SECRET
          valueFrom:
            secretKeyRef:
              name: payment-service-secrets
              key: webhook-secret
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /health
            port: 8087
          initialDelaySeconds: 30
          periodSeconds: 10
          timeoutSeconds: 5
          failureThreshold: 3
        readinessProbe:
          httpGet:
            path: /ready
            port: 8087
          initialDelaySeconds: 5
          periodSeconds: 5
          timeoutSeconds: 3
          failureThreshold: 3
        securityContext:
          allowPrivilegeEscalation: false
          readOnlyRootFilesystem: true
          capabilities:
            drop:
            - ALL
        volumeMounts:
        - name: tmp
          mountPath: /tmp
      volumes:
      - name: tmp
        emptyDir: {}
      imagePullSecrets:
      - name: registry-secret
      restartPolicy: Always
      terminationGracePeriodSeconds: 30
```

**Service Configuration:**
```yaml
apiVersion: v1
kind: Service
metadata:
  name: payment-service
  namespace: payment-service
  labels:
    app: payment-service
  annotations:
    prometheus.io/scrape: "true"
    prometheus.io/port: "9090"
spec:
  type: ClusterIP
  ports:
  - port: 80
    targetPort: 8086
    protocol: TCP
    name: http
  - port: 9090
    targetPort: 9090
    protocol: TCP
    name: metrics
  - port: 8087
    targetPort: 8087
    protocol: TCP
    name: health
  selector:
    app: payment-service
```

**Ingress Configuration:**
```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: payment-service-ingress
  namespace: payment-service
  annotations:
    kubernetes.io/ingress.class: nginx
    cert-manager.io/cluster-issuer: letsencrypt-prod
    nginx.ingress.kubernetes.io/ssl-redirect: "true"
    nginx.ingress.kubernetes.io/proxy-body-size: "1m"
    nginx.ingress.kubernetes.io/rate-limit: "100"
    nginx.ingress.kubernetes.io/rate-limit-window: "1m"
spec:
  tls:
  - hosts:
    - payments.yourdomain.com
    secretName: payment-service-tls
  rules:
  - host: payments.yourdomain.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: payment-service
            port:
              number: 80
```

## Database Setup

### Initial Migration

```bash
# Run database migrations
go run cmd/migrate/main.go up

# Or use migration tool
migrate -path migrations -database "postgres://user:password@host:port/dbname?sslmode=require" up
```

### Seed Data

```bash
# Run seed data script
psql -h host -U user -d dbname -f migrations/seed_data.sql
```

### Database Security

```sql
-- Create database users
CREATE ROLE payment_service_app WITH LOGIN PASSWORD 'secure_password';
CREATE ROLE payment_service_readonly WITH LOGIN PASSWORD 'readonly_password';

-- Grant permissions
GRANT CONNECT ON DATABASE payment_service TO payment_service_app;
GRANT USAGE ON SCHEMA public TO payment_service_app;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO payment_service_app;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO payment_service_app;

-- Read-only user for analytics
GRANT CONNECT ON DATABASE payment_service TO payment_service_readonly;
GRANT USAGE ON SCHEMA public TO payment_service_readonly;
GRANT SELECT ON ALL TABLES IN SCHEMA public TO payment_service_readonly;
```

## Monitoring and Logging

### Prometheus Configuration

```yaml
# prometheus.yml
global:
  scrape_interval: 15s
  evaluation_interval: 15s

scrape_configs:
- job_name: 'payment-service'
  static_configs:
  - targets: ['payment-service:9090']
  metrics_path: /metrics
  scrape_interval: 15s
```

### Grafana Dashboard

Create a Grafana dashboard with the following panels:

1. **Payment Metrics**
   - Payments per minute
   - Success rate by payment method
   - Average payment amount
   - Payment processing duration

2. **System Metrics**
   - HTTP requests per second
   - Error rate
   - Response time percentiles
   - Memory and CPU usage

3. **Fiuu Integration**
   - Fiuu API request rate
   - Circuit breaker state
   - Retry attempts
   - Error breakdown

4. **Database Metrics**
   - Connection pool usage
   - Query duration
   - Database size growth

### Alerting Rules

```yaml
# alerts.yml
groups:
- name: payment-service-alerts
  rules:
  - alert: PaymentServiceDown
    expr: up{job="payment-service"} == 0
    for: 1m
    labels:
      severity: critical
    annotations:
      summary: "Payment service is down"
      description: "Payment service has been down for more than 1 minute"

  - alert: HighErrorRate
    expr: rate(http_requests_total{status_code=~"5.."}[5m]) / rate(http_requests_total[5m]) > 0.05
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: "High error rate detected"
      description: "Error rate is above 5% for more than 5 minutes"

  - alert: FiuuCircuitBreakerOpen
    expr: fiuu_circuit_breaker_state{merchant_id="your_merchant_id"} == 1
    for: 1m
    labels:
      severity: critical
    annotations:
      summary: "Fiuu circuit breaker is open"
      description: "Fiuu circuit breaker has been open for more than 1 minute"
```

## Deployment Process

### Pre-deployment Checklist

- [ ] Backup current database
- [ ] Verify all environment variables are set
- [ ] Test database connection
- [ ] Verify Fiuu credentials
- [ ] Check SSL certificates
- [ ] Validate monitoring configuration
- [ ] Run security scans
- [ ] Perform load testing

### Deployment Steps

1. **Prepare Release**
   ```bash
   # Build and push Docker image
   docker build -t your-registry/payment-service:v1.0.0 .
   docker push your-registry/payment-service:v1.0.0
   ```

2. **Run Database Migrations**
   ```bash
   # Backup database
   pg_dump payment_service > backup_$(date +%Y%m%d_%H%M%S).sql

   # Run migrations
   kubectl run migration --image=your-registry/payment-service:v1.0.0 --restart=Never -- migrate up
   ```

3. **Update Deployment**
   ```bash
   # Update image tag
   kubectl set image deployment/payment-service payment-service=your-registry/payment-service:v1.0.0 -n payment-service

   # Wait for rollout
   kubectl rollout status deployment/payment-service -n payment-service
   ```

4. **Verify Deployment**
   ```bash
   # Check pod status
   kubectl get pods -n payment-service

   # Check service health
   kubectl port-forward svc/payment-service 8086:80 -n payment-service
   curl http://localhost:8086/health

   # Check metrics
   curl http://localhost:9090/metrics
   ```

### Rollback Process

```bash
# If deployment fails, rollback to previous version
kubectl rollout undo deployment/payment-service -n payment-service

# Verify rollback
kubectl rollout status deployment/payment-service -n payment-service
```

## Security Considerations

### Network Security

- Implement network policies to restrict traffic
- Use TLS for all external communications
- Enable mTLS for service-to-service communication
- Implement DDoS protection

### Access Control

- Use role-based access control (RBAC)
- Implement least privilege principle
- Regular access reviews
- Multi-factor authentication for admin access

### Data Protection

- Encrypt sensitive data at rest
- Use secure key management
- Implement data retention policies
- Regular security audits

### Compliance

- PCI DSS compliance validation
- Regular penetration testing
- Security incident response plan
- Employee security training

## Troubleshooting

### Common Issues

1. **Database Connection Issues**
   - Check connection strings
   - Verify network connectivity
   - Check SSL certificates
   - Monitor connection pool

2. **Fiuu API Issues**
   - Verify credentials
   - Check API rate limits
   - Monitor circuit breaker state
   - Review error logs

3. **High Memory Usage**
   - Check for memory leaks
   - Monitor garbage collection
   - Review resource limits
   - Optimize database queries

4. **Performance Issues**
   - Monitor response times
   - Check database query performance
   - Review caching strategies
   - Analyze bottleneck patterns

### Debug Commands

```bash
# Check pod logs
kubectl logs -f deployment/payment-service -n payment-service

# Debug pod
kubectl exec -it deployment/payment-service -n payment-service -- /bin/sh

# Check resource usage
kubectl top pods -n payment-service

# Port forward for local debugging
kubectl port-forward svc/payment-service 8086:80 -n payment-service

# Check database connections
kubectl exec -it deployment/payment-service -n payment-service -- psql -h $DB_HOST -U $DB_USER -d $DB_NAME -c "SELECT count(*) FROM pg_stat_activity;"
```

## Performance Optimization

### Database Optimization

- Add appropriate indexes
- Optimize slow queries
- Implement connection pooling
- Regular maintenance tasks

### Application Optimization

- Implement caching strategies
- Use connection pooling
- Optimize garbage collection
- Enable compression

### Infrastructure Optimization

- Use SSD storage
- Optimize network configuration
- Implement horizontal scaling
- Use content delivery networks

## Maintenance

### Regular Tasks

- Daily: Check error logs and metrics
- Weekly: Review performance trends
- Monthly: Security updates and patches
- Quarterly: Performance optimization review

### Backup Strategy

- Database backups: Every 6 hours
- Configuration backups: Daily
- Log archival: Monthly
- Disaster recovery testing: Quarterly

### Updates and Patches

- Follow security patch management process
- Test updates in staging environment
- Schedule maintenance windows
- Communicate changes to stakeholders

## Support

For production support:

- **Emergency Contact**: ops-team@yourdomain.com
- **Documentation**: Internal knowledge base
- **Monitoring**: Grafana dashboards
- **Alerting**: PagerDuty integration
- **Escalation**: On-call rotation process
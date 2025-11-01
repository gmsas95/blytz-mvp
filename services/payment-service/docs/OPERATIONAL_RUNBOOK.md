# Payment Service Operational Runbook

## Overview

This runbook provides step-by-step procedures for operating, monitoring, and troubleshooting the Payment Service in production environments.

## Table of Contents

1. [Service Health Monitoring](#service-health-monitoring)
2. [Incident Response](#incident-response)
3. [Common Issues and Solutions](#common-issues-and-solutions)
4. [Performance Tuning](#performance-tuning)
5. [Maintenance Procedures](#maintenance-procedures)
6. [Security Incidents](#security-incidents)
7. [Disaster Recovery](#disaster-recovery)

## Service Health Monitoring

### Key Metrics to Monitor

#### Application Metrics
- **Payment Success Rate**: Should be >95%
- **API Response Time**: Should be <200ms (p95)
- **Error Rate**: Should be <1%
- **Circuit Breaker State**: Should be CLOSED (0)
- **Active Connections**: Database and Redis connections

#### Business Metrics
- **Revenue per Hour**: Track transaction volume
- **Top Payment Methods**: Monitor method usage
- **Failed Transactions**: Track failure reasons
- **Refund Rate**: Should be <5%

#### Infrastructure Metrics
- **CPU Usage**: Should be <70%
- **Memory Usage**: Should be <80%
- **Disk Usage**: Should be <85%
- **Network I/O**: Monitor throughput

### Monitoring Dashboard Access

- **Grafana Dashboard**: `https://grafana.yourdomain.com/d/payment-service`
- **Prometheus**: `https://prometheus.yourdomain.com`
- **Kibana Logs**: `https://kibana.yourdomain.com`

### Health Check Endpoints

```bash
# Application health
curl https://payments.yourdomain.com/health

# Readiness check
curl https://payments.yourdomain.com/ready

# Metrics endpoint
curl https://payments.yourdomain.com/metrics
```

## Incident Response

### Severity Levels

| Severity | Description | Response Time | Escalation |
|----------|-------------|----------------|------------|
| P0 - Critical | Service completely down, major revenue impact | 15 minutes | Immediate |
| P1 - High | Significant functionality loss, revenue impact | 30 minutes | 1 hour |
| P2 - Medium | Partial functionality loss, limited impact | 2 hours | 4 hours |
| P3 - Low | Minor issues, no revenue impact | 8 hours | 24 hours |

### Incident Response Process

#### 1. Incident Detection
- Automated alerts from monitoring systems
- Customer reports via support channels
- Manual observation by operations team

#### 2. Initial Assessment (First 15 minutes)
```bash
# Check service status
kubectl get pods -n payment-service
kubectl describe pod <pod-name> -n payment-service

# Check recent logs
kubectl logs --tail=100 -f deployment/payment-service -n payment-service

# Check system resources
kubectl top pods -n payment-service
kubectl top nodes
```

#### 3. Communication
- Create incident Slack channel
- Notify stakeholders based on severity
- Post status updates every 15 minutes
- Document all actions taken

#### 4. Investigation and Resolution
- Identify root cause
- Implement temporary fix if needed
- Plan permanent solution
- Test resolution thoroughly

#### 5. Post-Incident Review
- Conduct post-mortem within 24 hours
- Document timeline and root cause
- Identify preventive measures
- Update runbooks and monitoring

## Common Issues and Solutions

### 1. Service Unavailable

#### Symptoms
- All health checks failing
- 503 Service Unavailable errors
- No new payments being processed

#### Diagnosis Commands
```bash
# Check pod status
kubectl get pods -n payment-service -l app=payment-service

# Check recent events
kubectl get events -n payment-service --sort-by='.lastTimestamp' | tail -20

# Check resource constraints
kubectl describe pod <pod-name> -n payment-service

# Check logs for errors
kubectl logs deployment/payment-service -n payment-service --tail=200
```

#### Common Causes and Solutions

**Pod Crashing**
```bash
# Check pod logs for crash reason
kubectl logs <pod-name> -n payment-service --previous

# Common issues:
# - Out of memory: Increase memory limits
# - Configuration error: Check environment variables
# - Database connection failure: Verify DB connectivity
```

**Network Issues**
```bash
# Test database connectivity
kubectl exec -it <pod-name> -n payment-service -- nc -zv $DB_HOST $DB_PORT

# Test Redis connectivity
kubectl exec -it <pod-name> -n payment-service -- redis-cli -h $REDIS_HOST -p $REDIS_PORT ping

# Check service endpoints
kubectl get endpoints -n payment-service
```

**Resource Exhaustion**
```bash
# Check node resources
kubectl top nodes
kubectl describe nodes

# Check resource quotas
kubectl describe quota -n payment-service
```

### 2. High Error Rate

#### Symptoms
- Error rate >5%
- Failed payment transactions
- Customer complaints

#### Diagnosis Commands
```bash
# Check error logs
kubectl logs deployment/payment-service -n payment-service | grep ERROR | tail -50

# Check specific error patterns
kubectl logs deployment/payment-service -n payment-service | grep "payment.*failed" | tail -20

# Monitor real-time errors
kubectl logs -f deployment/payment-service -n payment-service | grep ERROR
```

#### Common Error Types and Solutions

**Fiuu API Errors**
```bash
# Check Fiuu connectivity
kubectl exec -it <pod-name> -n payment-service -- curl -X POST https://pay.fiuu.com/RMS/API/payment/PaymentRequest

# Verify credentials
kubectl exec -it <pod-name> -n payment-service -- env | grep FIUU

# Check circuit breaker status
curl https://payments.yourdomain.com/metrics | grep fiuu_circuit_breaker_state
```

**Database Errors**
```bash
# Check database connections
kubectl exec -it <pod-name> -n payment-service -- psql -h $DB_HOST -U $DB_USER -d $DB_NAME -c "SELECT count(*) FROM pg_stat_activity WHERE state = 'active';"

# Check database locks
kubectl exec -it <pod-name> -n payment-service -- psql -h $DB_HOST -U $DB_USER -d $DB_NAME -c "SELECT * FROM pg_locks WHERE NOT granted;"

# Check slow queries
kubectl exec -it <pod-name> -n payment-service -- psql -h $DB_HOST -U $DB_USER -d $DB_NAME -c "SELECT query, mean_time, calls FROM pg_stat_statements ORDER BY mean_time DESC LIMIT 10;"
```

**Rate Limiting Issues**
```bash
# Check rate limit metrics
curl https://payments.yourdomain.com/metrics | grep rate_limit

# Review recent requests
kubectl logs deployment/payment-service -n payment-service | grep "rate limit" | tail -20
```

### 3. Performance Degradation

#### Symptoms
- Response time >500ms
- Increasing queue lengths
- Customer complaints about slowness

#### Diagnosis Commands
```bash
# Check response time metrics
curl https://payments.yourdomain.com/metrics | grep http_request_duration_seconds

# Check database query performance
kubectl exec -it <pod-name> -n payment-service -- psql -h $DB_HOST -U $DB_USER -d $DB_NAME -c "SELECT * FROM pg_stat_activity WHERE query_start < now() - interval '1 minute';"

# Check system resources
kubectl top pods -n payment-service
kubectl top nodes
```

#### Common Causes and Solutions

**Database Performance Issues**
```bash
# Check indexes
kubectl exec -it <pod-name> -n payment-service -- psql -h $DB_HOST -U $DB_USER -d $DB_NAME -c "\d payments"

# Analyze slow queries
kubectl exec -it <pod-name> -n payment-service -- psql -h $DB_HOST -U $DB_USER -d $DB_NAME -c "SELECT query, mean_time, calls FROM pg_stat_statements ORDER BY mean_time DESC LIMIT 10;"

# Update table statistics
kubectl exec -it <pod-name> -n payment-service -- psql -h $DB_HOST -U $DB_USER -d $DB_NAME -c "ANALYZE payments;"
```

**Memory Leaks**
```bash
# Check memory usage over time
kubectl top pod <pod-name> -n payment-service --watch

# Check garbage collection
kubectl exec -it <pod-name> -n payment-service -- ps aux | grep payment-service
```

**External API Latency**
```bash
# Test Fiuu API response time
time kubectl exec -it <pod-name> -n payment-service -- curl -X POST https://pay.fiuu.com/RMS/API/payment/PaymentRequest

# Check network latency
kubectl exec -it <pod-name> -n payment-service -- ping -c 5 pay.fiuu.com
```

## Performance Tuning

### Database Optimization

#### Query Optimization
```sql
-- Identify slow queries
SELECT query, mean_time, calls, total_time
FROM pg_stat_statements
WHERE mean_time > 1000 -- queries taking more than 1 second
ORDER BY mean_time DESC;

-- Check missing indexes
SELECT schemaname, tablename, attname, n_distinct, correlation
FROM pg_stats
WHERE schemaname = 'public'
ORDER BY tablename, attname;

-- Analyze table statistics
ANALYZE payments;
ANALYZE refunds;
ANALYZE payment_attempts;
```

#### Connection Pool Tuning
```bash
# Monitor connection pool usage
curl https://payments.yourdomain.com/metrics | grep db_connections

# Adjust pool size in environment variables
export DB_MAX_CONNECTIONS=25
export DB_MAX_IDLE_CONNECTIONS=5
```

### Application Optimization

#### Memory Management
```bash
# Check heap profile
kubectl exec -it <pod-name> -n payment-service -- curl http://localhost:6060/debug/pprof/heap

# Enable GC tuning
export GOGC=100
export GOMEMLIMIT=500MiB
```

#### Concurrency Control
```bash
# Check goroutine leaks
kubectl exec -it <pod-name> -n payment-service -- curl http://localhost:6060/debug/pprof/goroutine?debug=1

# Monitor concurrency metrics
curl https://payments.yourdomain.com/metrics | grep goroutines
```

### Infrastructure Scaling

#### Horizontal Scaling
```bash
# Scale deployment
kubectl scale deployment payment-service --replicas=5 -n payment-service

# Monitor scaling impact
kubectl rollout status deployment/payment-service -n payment-service
```

#### Resource Optimization
```yaml
# Update resource limits
resources:
  requests:
    memory: "512Mi"
    cpu: "500m"
  limits:
    memory: "1Gi"
    cpu: "1000m"
```

## Maintenance Procedures

### Rolling Updates

#### Update Procedure
```bash
# 1. Backup current deployment
kubectl get deployment payment-service -n payment-service -o yaml > payment-service-backup.yaml

# 2. Update image
kubectl set image deployment/payment-service payment-service=your-registry/payment-service:v1.1.0 -n payment-service

# 3. Monitor rollout
kubectl rollout status deployment/payment-service -n payment-service
kubectl rollout history deployment/payment-service -n payment-service

# 4. Verify health
kubectl get pods -n payment-service -l app=payment-service
curl https://payments.yourdomain.com/health
```

#### Rollback Procedure
```bash
# Quick rollback to previous version
kubectl rollout undo deployment/payment-service -n payment-service

# Rollback to specific version
kubectl rollout undo deployment/payment-service --to-revision=2 -n payment-service

# Verify rollback
kubectl rollout status deployment/payment-service -n payment-service
```

### Database Maintenance

#### Regular Maintenance Tasks
```bash
# Weekly maintenance script
#!/bin/bash

# 1. Database backup
pg_dump -h $DB_HOST -U $DB_USER -d $DB_NAME > backup_$(date +%Y%m%d_%H%M%S).sql

# 2. Update statistics
psql -h $DB_HOST -U $DB_USER -d $DB_NAME -c "ANALYZE;"

# 3. Rebuild indexes if needed
psql -h $DB_HOST -U $DB_USER -d $DB_NAME -c "REINDEX DATABASE payment_service;"

# 4. Clean up old data (retention policy)
psql -h $DB_HOST -U $DB_USER -d $DB_NAME -c "DELETE FROM payment_attempts WHERE created_at < NOW() - INTERVAL '90 days';"
```

#### Schema Updates
```bash
# Run migrations
kubectl run migration --image=your-registry/payment-service:v1.1.0 --restart=Never -- migrate up

# Verify migration
kubectl logs migration -n payment-service

# Clean up migration job
kubectl delete job migration -n payment-service
```

### Log Management

#### Log Rotation
```bash
# Configure log rotation in Kubernetes
apiVersion: v1
kind: ConfigMap
metadata:
  name: logging-config
data:
  fluent.conf: |
    <source>
      @type tail
      path /var/log/containers/*payment-service*.log
      pos_file /var/log/fluentd-containers.log.pos
      tag kubernetes.*
      format json
      time_format %Y-%m-%dT%H:%M:%S.%NZ
    </source>
```

#### Log Archival
```bash
# Archive old logs
kubectl exec -it <pod-name> -n payment-service -- find /var/log -name "*.log" -mtime +30 -exec gzip {} \;

# Transfer to S3
aws s3 sync /var/log/archive s3://your-bucket/payment-service-logs/
```

## Security Incidents

### Incident Types

#### 1. Unauthorized Access
- Immediate action: Revoke compromised credentials
- Investigation: Audit logs for data access
- Prevention: Review access controls

#### 2. Data Breach
- Immediate action: Isolate affected systems
- Investigation: Scope of data exposure
- Prevention: Enhance security controls

#### 3. DDoS Attack
- Immediate action: Enable rate limiting
- Investigation: Attack source and patterns
- Prevention: Implement additional protections

### Security Response Checklist

#### Immediate Actions (First Hour)
- [ ] Identify and isolate affected systems
- [ ] Revoke compromised credentials
- [ ] Enable enhanced monitoring
- [ ] Notify security team
- [ ] Document timeline

#### Investigation Actions (First 24 Hours)
- [ ] Analyze audit logs
- [ ] Assess data exposure scope
- [ ] Identify attack vector
- [ ] Collect forensic evidence
- [ ] Plan remediation

#### Recovery Actions (Following Days)
- [ ] Patch vulnerabilities
- [ ] Restore from clean backups
- [ ] Implement additional controls
- [ ] Update incident response plan
- [ ] Conduct security training

### Security Monitoring

#### Real-time Monitoring
```bash
# Monitor failed authentication attempts
kubectl logs deployment/payment-service -n payment-service | grep "auth.*failed" | tail -10

# Monitor unusual API access patterns
kubectl logs deployment/payment-service -n payment-service | grep "payment.*failed" | tail -10

# Check for security events
curl https://payments.yourdomain.com/metrics | grep security_events
```

#### Alert Configuration
```yaml
# Security alerts in Prometheus
groups:
- name: security-alerts
  rules:
  - alert: HighFailedAuthRate
    expr: rate(failed_auth_total[5m]) > 10
    for: 2m
    labels:
      severity: warning
    annotations:
      summary: "High failed authentication rate detected"

  - alert: SuspiciousAPIAccess
    expr: rate(suspicious_api_calls[5m]) > 5
    for: 1m
    labels:
      severity: critical
    annotations:
      summary: "Suspicious API access pattern detected"
```

## Disaster Recovery

### Backup Strategy

#### Data Backups
```bash
# Database backup script
#!/bin/bash

BACKUP_DIR="/backups/payment-service"
DATE=$(date +%Y%m%d_%H%M%S)

# Full database backup
pg_dump -h $DB_HOST -U $DB_USER -d $DB_NAME | gzip > $BACKUP_DIR/db_backup_$DATE.sql.gz

# Configuration backup
kubectl get configmap payment-service-config -n payment-service -o yaml > $BACKUP_DIR/config_$DATE.yaml
kubectl get secret payment-service-secrets -n payment-service -o yaml > $BACKUP_DIR/secrets_$DATE.yaml

# Upload to S3
aws s3 cp $BACKUP_DIR/db_backup_$DATE.sql.gz s3://your-backup-bucket/payment-service/
aws s3 cp $BACKUP_DIR/config_$DATE.yaml s3://your-backup-bucket/payment-service/
aws s3 cp $BACKUP_DIR/secrets_$DATE.yaml s3://your-backup-bucket/payment-service/
```

#### Backup Verification
```bash
# Verify backup integrity
#!/bin/bash

BACKUP_FILE=$1
TEMP_DB="test_payment_service_$(date +%s)"

# Create test database
createdb $TEMP_DB

# Restore backup
gunzip -c $BACKUP_FILE | psql $TEMP_DB

# Verify data integrity
psql $TEMP_DB -c "SELECT count(*) FROM payments;"
psql $TEMP_DB -c "SELECT count(*) FROM refunds;"

# Cleanup
dropdb $TEMP_DB
```

### Recovery Procedures

#### Full Service Recovery
```bash
# 1. Restore database
gunzip -c /backups/db_backup_20240115_120000.sql.gz | psql -h $DB_HOST -U $DB_USER -d payment_service

# 2. Restore configuration
kubectl apply -f /backups/config_20240115_120000.yaml
kubectl apply -f /backups/secrets_20240115_120000.yaml

# 3. Restart services
kubectl rollout restart deployment/payment-service -n payment-service

# 4. Verify recovery
kubectl rollout status deployment/payment-service -n payment-service
curl https://payments.yourdomain.com/health
```

#### Partial Recovery
```bash
# Restore specific tables
pg_dump -h $DB_HOST -U $DB_USER -d backup_db -t payments | psql -h $DB_HOST -U $DB_USER -d payment_service

# Verify data consistency
psql -h $DB_HOST -U $DB_USER -d payment_service -c "SELECT COUNT(*) FROM payments WHERE created_at > '2024-01-15';"
```

### Business Continuity

#### RTO/RPO Targets
- **Recovery Time Objective (RTO)**: 4 hours
- **Recovery Point Objective (RPO)**: 15 minutes

#### Failover Testing
```bash
# Monthly failover test procedure
1. Schedule maintenance window
2. Notify stakeholders
3. Perform controlled failover
4. Verify service functionality
5. Document results
6. Address any issues found
```

#### Emergency Contacts
- **On-call Engineer**: +1-XXX-XXX-XXXX
- **Engineering Manager**: +1-XXX-XXX-XXXX
- **Security Team**: security@yourdomain.com
- **External Fiuu Support**: support@fiuu.com

## Escalation Procedures

### Escalation Matrix

| Issue Type | L1 Support | L2 Support | L3 Support | Escalation Time |
|------------|------------|------------|------------|-----------------|
| Service Outage | On-call Engineer | Senior Engineer | Architecture Team | 30 min |
| Security Incident | On-call Engineer | Security Team | CISO | 15 min |
| Performance Issue | On-call Engineer | Performance Team | Architecture Team | 1 hour |
| Data Issue | On-call Engineer | DBA Team | Data Engineering | 2 hours |

### Escalation Checklist

#### Before Escalation
- [ ] Document all troubleshooting steps taken
- [ ] Gather relevant logs and metrics
- [ ] Assess business impact
- [ ] Prepare incident summary

#### Escalation Process
1. Contact L2 support via Slack/Phone
2. Share incident context and timeline
3. Provide access to relevant systems
4. Monitor escalation progress
5. Document resolution

### Post-Incident Procedures

#### Post-Mortem Template
```markdown
# Post-Mortem: [Incident Title]

## Summary
[Brief description of the incident]

## Timeline
- [Time]: Incident detected
- [Time]: Investigation started
- [Time]: Root cause identified
- [Time]: Mitigation implemented
- [Time]: Service restored

## Impact
- Users affected: [Number]
- Duration: [Time]
- Revenue impact: [Amount]

## Root Cause
[Detailed explanation of what happened]

## Resolution
[Steps taken to fix the issue]

## Prevention
[Measures to prevent recurrence]

## Lessons Learned
[Key takeaways and improvements needed]
```

#### Follow-up Actions
- [ ] Complete post-mortem within 24 hours
- [ ] Update runbooks and documentation
- [ ] Implement preventive measures
- [ ] Schedule training if needed
- [ ] Review monitoring and alerting
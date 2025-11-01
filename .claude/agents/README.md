# Blytz Multi-Agent System

This directory contains specialized agents designed to help you develop, maintain, and operate the Blytz Live Auction MVP platform efficiently.

## Available Agents

### 🔐 Security Agent
**Specializes in**: Authentication, JWT tokens, security auditing, and vulnerability management
**Use for**: Authentication issues, security audits, JWT troubleshooting, CORS problems

### 🏗️ Infrastructure Agent
**Specializes in**: Docker, Kubernetes, database management, and deployment automation
**Use for**: Container issues, database problems, deployment troubleshooting, infrastructure setup

### ⚡ Performance Agent
**Specializes in**: Redis optimization, Go performance profiling, database tuning, and load testing
**Use for**: Performance issues, slow queries, Redis optimization, load testing

### 🧪 Testing Agent
**Specializes in**: Test automation, integration testing, load testing, and CI/CD testing
**Use for**: Test creation, debugging test failures, performance testing, CI/CD issues

### 🔧 Frontend Agent
**Specializes in**: React/Next.js development, responsive design, API integration, and UI optimization
**Use for**: Frontend development, API integration issues, responsive design problems

### 🔄 Deployment Agent
**Specializes in**: Environment management, zero-downtime deployments, rollback procedures, and infrastructure provisioning
**Use for**: Deployment issues, environment setup, rollback procedures, configuration management

### 📊 Analytics Agent
**Specializes in**: Metrics analysis, business intelligence, data visualization, and performance analytics
**Use for**: Performance analysis, business metrics, dashboard creation, data insights

### 🛠️ Debugging Agent
**Specializes in**: Issue diagnosis, log analysis, root cause analysis, and troubleshooting
**Use for**: Complex debugging, log analysis, system diagnostics, issue resolution

## Usage

To use an agent, simply reference it by name in your Claude Code conversation:

```
@security-agent Help me fix JWT token validation
@infrastructure-agent Set up Redis clustering
@performance-agent Optimize auction service response times
@testing-agent Create integration tests for payment flow
@frontend-agent Fix responsive design issues
@deployment-agent Set up production environment
@analytics-agent Create performance dashboard
@debugging-agent Diagnose why bids are failing
```

Each agent has access to relevant tools, documentation, and best practices specific to their domain of expertise.
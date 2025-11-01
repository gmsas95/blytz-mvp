# Frontend Environment Configuration Guide

## 🚀 Quick Setup

### 1. Copy Environment Template
```bash
# For each frontend directory, copy the example file:
cp frontend/.env.local.example frontend/.env.local
cp frontend-demo/.env.local.example frontend-demo/.env.local
cp frontend-seller/.env.local.example frontend-seller/.env.local
```

### 2. Development Environment Variables
Update each `.env.local` file with your local development settings:

```bash
# API Configuration
NEXT_PUBLIC_API_URL=http://localhost:8080
NEXT_PUBLIC_WS_URL=ws://localhost:8080

# LiveKit Configuration
NEXT_PUBLIC_LIVEKIT_URL=ws://localhost:7880

# Application Mode
NODE_ENV=development
MODE=mock
```

### 3. Production Environment Variables
In your deployment platform (Dokploy/Vercel/etc.):

```bash
# Main Frontend
NEXT_PUBLIC_API_URL=https://api.blytz.app
NEXT_PUBLIC_WS_URL=wss://api.blytz.app
NEXT_PUBLIC_LIVEKIT_URL=wss://livekit.blytz.app
NODE_ENV=production
MODE=remote

# Demo Frontend
NEXT_PUBLIC_API_URL=https://api.blytz.app
NEXT_PUBLIC_WS_URL=wss://api.blytz.app
NEXT_PUBLIC_LIVEKIT_URL=wss://livekit.blytz.app
NEXT_PUBLIC_DEMO_MODE=viewer
NODE_ENV=production

# Seller Frontend
NEXT_PUBLIC_API_URL=https://api.blytz.app
NEXT_PUBLIC_WS_URL=wss://api.blytz.app
NEXT_PUBLIC_LIVEKIT_URL=wss://livekit.blytz.app
NEXT_PUBLIC_DEMO_MODE=broadcaster
NODE_ENV=production
```

## 🔧 Environment Variables Explained

| Variable | Required | Description | Example |
|----------|-----------|-------------|----------|
| `NEXT_PUBLIC_API_URL` | ✅ | API base URL | `https://api.blytz.app` |
| `NEXT_PUBLIC_WS_URL` | ✅ | WebSocket URL | `wss://api.blytz.app` |
| `NEXT_PUBLIC_LIVEKIT_URL` | ✅ | LiveKit server URL | `wss://livekit.blytz.app` |
| `NODE_ENV` | ✅ | Environment mode | `development` or `production` |
| `MODE` | ✅ | API mode | `mock` or `remote` |
| `NEXT_PUBLIC_DEMO_MODE` | ❌ | Demo app mode | `viewer` or `broadcaster` |

## 🚨 Security Notes

### ✅ Safe Practices
- All public variables use `NEXT_PUBLIC_` prefix
- No hardcoded API keys or secrets
- Environment-specific configurations
- Input validation implemented

### ⚠️ Important
- Never commit `.env.local` files to version control
- Use different values for development vs production
- Regularly rotate API keys and secrets
- Validate all environment variables on startup

## 🛠️ Troubleshooting

### Common Issues

1. **"NEXT_PUBLIC_API_URL environment variable is required"**
   - Solution: Set the environment variable in your deployment platform
   - Check `.env.local` file exists for development

2. **LiveKit connection failed**
   - Solution: Verify `NEXT_PUBLIC_LIVEKIT_URL` is correct
   - Ensure WebSocket protocol (`wss://` for production)

3. **API requests failing**
   - Solution: Check `MODE` environment variable
   - Use `remote` for production, `mock` for development

### Validation Commands
```bash
# Check if environment variables are set
echo $NEXT_PUBLIC_API_URL

# Test API connectivity
curl $NEXT_PUBLIC_API_URL/health

# Test LiveKit connectivity
wscat -c $NEXT_PUBLIC_LIVEKIT_URL
```

## 📋 Deployment Checklist

### Pre-deployment
- [ ] All environment variables set
- [ ] API endpoints accessible
- [ ] LiveKit server reachable
- [ ] No hardcoded URLs in code
- [ ] Error handling implemented

### Post-deployment
- [ ] Test authentication flow
- [ ] Verify API connectivity
- [ ] Test LiveKit streaming
- [ ] Check browser console for errors
- [ ] Monitor network requests

## 🔗 Related Files

- `frontend/.env.local.example` - Main frontend template
- `frontend-demo/.env.local.example` - Demo app template  
- `frontend-seller/.env.local.example` - Seller app template
- `frontend/src/lib/api-adapter.ts` - API integration logic
- `docker-compose.yml` - Production environment variables
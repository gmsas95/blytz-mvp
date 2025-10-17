# 🚀 Blytz Auction MVP - Deployment Guide

This guide will help you deploy the Blytz Auction MVP to your VPS at blytz.app for your exhibition.

## ✅ What's Ready

### Backend Services (Complete)
- ✅ **Authentication Service** (Port 8084) - JWT-based auth with Better Auth
- ✅ **Auction Service** (Port 8083) - Real-time bidding with PostgreSQL + Redis
- ✅ **Firebase Functions** - Payment processing, notifications, auction management
- ✅ **Database Persistence** - PostgreSQL with auction and bid tables
- ✅ **API Gateway** (Port 8080) - Central routing with Nginx

### Frontend (Ready)
- ✅ **Web Interface** (`/frontend/index.html`) - Complete auction testing interface
- ✅ **Real-time Updates** - Auto-refresh every 10 seconds
- ✅ **Mobile Responsive** - Works on phones/tablets

## 🎯 MVP Features for Exhibition

### Core Auction Functionality
1. **User Registration/Login** - Visitors can create accounts
2. **Create Auctions** - Add items for bidding
3. **Browse Auctions** - View active auctions with time remaining
4. **Place Bids** - Real-time bidding system
5. **Auction Status** - Live updates on current price and bids

### Demo Workflow
1. Visitor registers → Login
2. Creates auction → Sets price/duration
3. Other visitors browse → Place bids
4. Real-time updates → Auction ends
5. Winner notification → Payment processing

## 🚀 Quick Start

### 1. Start All Services
```bash
# Start the complete stack
docker-compose up -d

# Verify services are running
docker-compose ps
curl http://localhost:8080/health
```

### 2. Initialize Database
```bash
# Run database setup with sample data
cd services/auction-service
./scripts/init-db.sh

# This creates:
# - Auctions table with indexes
# - Bids table with foreign keys
# - Sample auction data for demo
```

### 3. Test the API
```bash
# Test authentication
curl -X POST http://localhost:8084/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email": "demo@blytz.app", "password": "password123", "display_name": "Demo User"}'

# Test auction creation (after login)
curl -X POST http://localhost:8083/api/v1/auctions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "title": "Vintage Watch",
    "description": "Beautiful vintage Rolex",
    "starting_price": 100,
    "reserve_price": 500,
    "min_bid_increment": 25,
    "start_time": "2024-01-01T10:00:00Z",
    "end_time": "2024-01-01T12:00:00Z",
    "type": "live",
    "product_id": "demo_watch_001"
  }'
```

### 4. Access Web Interface
Open `http://localhost:8080/frontend/index.html` in your browser.

## 📋 Exhibition Setup Checklist

### Pre-Exhibition (Day Before)
- [ ] VPS is running and accessible at blytz.app
- [ ] Docker and Docker Compose installed
- [ ] Domain DNS pointing to VPS IP
- [ ] SSL certificates configured (Let's Encrypt)
- [ ] All services started and healthy
- [ ] Database initialized with sample data
- [ ] Web interface tested on mobile devices

### During Exhibition
- [ ] Monitor service health
- [ ] Have backup plan ready
- [ ] Test with multiple users simultaneously
- [ ] Verify real-time bid updates
- [ ] Check auction completion flow

### Services Health Check
```bash
# Check all services
curl http://localhost:8080/health      # Gateway
curl http://localhost:8083/health      # Auction Service
curl http://localhost:8084/health      # Auth Service
curl http://localhost:5432/health      # PostgreSQL (if exposed)
curl http://localhost:6379/ping        # Redis (if exposed)
```

## 🔧 Production Deployment

### 1. Server Setup
```bash
# Update system
sudo apt update && sudo apt upgrade -y

# Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# Install Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
```

### 2. Domain Configuration
Point your domain blytz.app to your VPS IP address in your DNS provider.

### 3. SSL Setup (Let's Encrypt)
```bash
# Install Certbot
sudo apt install certbot python3-certbot-nginx -y

# Get SSL certificate
sudo certbot --nginx -d blytz.app -d www.blytz.app
```

### 4. Environment Variables
Create `.env` file with production settings:
```env
ENVIRONMENT=production
JWT_SECRET=your-super-secret-jwt-key
BETTER_AUTH_SECRET=your-better-auth-secret
DATABASE_URL=postgres://user:password@localhost:5432/auction_db?sslmode=require
REDIS_URL=redis://localhost:6379
```

### 5. Deploy
```bash
# Clone repository
git clone https://github.com/gmsas95/blytz-mvp.git
cd blytz-mvp

# Start services
docker-compose -f docker-compose.prod.yml up -d

# Initialize database
./services/auction-service/scripts/init-db.sh
```

## 📱 Mobile Testing

Test the web interface on various devices:
- iPhone Safari
- Android Chrome
- iPad Safari
- Desktop browsers (Chrome, Firefox, Safari)

## 🚨 Troubleshooting

### Common Issues

1. **Database Connection Failed**
   ```bash
   # Check PostgreSQL status
   docker-compose logs postgres

   # Verify database exists
   docker exec -it blytz-postgres psql -U postgres -l
   ```

2. **Authentication Issues**
   ```bash
   # Check auth service logs
   docker-compose logs auth-service

   # Verify JWT secret
   echo $JWT_SECRET
   ```

3. **Auction Service Not Responding**
   ```bash
   # Check auction service logs
   docker-compose logs auction-service

   # Test database connection
   curl http://localhost:8083/health
   ```

4. **Real-time Updates Not Working**
   ```bash
   # Check Redis connection
   docker exec -it blytz-redis redis-cli ping

   # Verify bid placement
   curl -X POST http://localhost:8083/api/v1/auctions/test/bids
   ```

## 🎉 Success Metrics for Exhibition

### User Engagement
- Number of registrations
- Active auctions created
- Total bids placed
- Average auction duration
- User feedback/comments

### Technical Performance
- API response times (< 200ms target)
- Error rates (< 1% target)
- Database query performance
- Real-time update latency (< 1s target)

## 📞 Support During Exhibition

If issues arise during the exhibition:

1. **Check service health**: `docker-compose ps`
2. **Review logs**: `docker-compose logs [service-name]`
3. **Restart services**: `docker-compose restart [service-name]`
4. **Database issues**: Run database initialization script again
5. **Network issues**: Check firewall and port configurations

## 🎯 Next Steps After Exhibition

Based on user feedback, you can enhance with:
- Live streaming integration
- Advanced analytics
- Mobile app
- Payment gateway integration
- Social features
- Admin dashboard

**Good luck with your exhibition! 🎭 The Blytz Auction MVP is ready to impress your visitors!** 🚀

---

**Support**: For issues during deployment, check the logs and refer to the troubleshooting section above. The system is designed to be robust and handle multiple concurrent users. 🛡️

**Status**: ✅ **READY FOR PRODUCTION** 🎯
# 🌀 Blytz Live Auction MVP – Docker Deployment

## Overview

This repository hosts the **Blytz Live Auction MVP**, a real-time livestream commerce platform powered by **Go**, **Redis**, and **Nginx**.  
The architecture is designed for **speed**, **scalability**, and **microservice flexibility** using a containerized setup suitable for a **1vCore / 4GB Hostinger KVM**.

---

## 🧩 Architecture Summary

### Core Services
| Service | Stack | Role |
|----------|--------|------|
| **Auction Engine** | Go + Redis | Core bidding logic, anti-snipe, atomic Lua scripts |
| **Redis** | In-memory cache | Real-time bid state, session cache, product data |
| **Nginx API Gateway** | Nginx | Reverse proxy, load balancer, HTTPS routing |
| **Frontend (React Native)** | Expo | Live auction UI and stream display |
| **Firebase (External)** | Cloud Functions | Payments, authentication, notifications |

### Directory Layout
```

/srv/blytz/
├── docker-compose.yml
├── README.md
├── nginx/
│    └── default.conf
├── services/
│    ├── auction-engine/
│    │    ├── main.go
│    │    ├── go.mod
│    │    ├── internal/
│    │    └── Dockerfile
│    ├── redis/
│    └── gateway/ (optional future services)
└── logs/

````

---

## ⚙️ Quick Setup (Hostinger KVM / Ubuntu 22.04)

### 1️⃣ Install System Dependencies
```bash
sudo apt update && sudo apt install -y docker.io docker-compose git
````

### 2️⃣ Clone the Repository

```bash
sudo mkdir -p /srv/blytz
cd /srv/blytz
git clone https://github.com/gmsas95/blytz-redis.git .
```

### 3️⃣ Launch Containers

```bash
sudo docker-compose up -d --build
```

### 4️⃣ Verify Everything

```bash
sudo docker ps
curl http://localhost:8080/health
```

If you see `{ "status": "ok" }`, your auction engine is running ✅

---

## 🧱 Docker Compose Overview

### `docker-compose.yml`

```yaml
version: "3.8"

services:
  redis:
    image: redis:7-alpine
    container_name: blytz-redis
    restart: always
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data

  auction-engine:
    build: ./services/auction-engine
    container_name: blytz-auction
    restart: always
    environment:
      - REDIS_HOST=redis:6379
      - PORT=8080
    depends_on:
      - redis
    ports:
      - "8080:8080"

  nginx:
    image: nginx:alpine
    container_name: blytz-nginx
    restart: always
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx/default.conf:/etc/nginx/conf.d/default.conf:ro
    depends_on:
      - auction-engine

volumes:
  redis_data:
```

---

## 🌐 API Gateway (Nginx)

### `nginx/default.conf`

```nginx
server {
    listen 80;

    server_name _;

    location /api/auction/ {
        proxy_pass http://auction-engine:8080/;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }

    location / {
        return 200 'Blytz API Gateway Active';
        add_header Content-Type text/plain;
    }
}
```

This routes all requests to:

```
http://<your-server-ip>/api/auction/
```

which forwards to your Go Auction Engine microservice.

---

## 🧠 Useful Commands

| Action                 | Command                                      |
| ---------------------- | -------------------------------------------- |
| Rebuild all containers | `sudo docker-compose up -d --build`          |
| Stop all containers    | `sudo docker-compose down`                   |
| View logs              | `sudo docker-compose logs -f`                |
| Restart auction engine | `sudo docker restart blytz-auction`          |
| Connect to Redis shell | `sudo docker exec -it blytz-redis redis-cli` |

---

## 📡 Next Steps

* [ ] Add **API Gateway SSL** via Nginx + Let’s Encrypt
* [ ] Deploy Firebase Cloud Functions (for payments, user sync)
* [ ] Integrate LiveKit for livestreams
* [ ] Connect frontend mobile app to `/api/auction`
* [ ] Add Prometheus metrics endpoint for auction engine

---

## 🧾 License

© 2025 Blytz Ventures. All rights reserved.
Internal use only — not for public redistribution.

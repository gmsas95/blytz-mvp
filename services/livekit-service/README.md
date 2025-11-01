# LiveKit Service

## Overview

The LiveKit service is responsible for generating JWT tokens for LiveKit WebRTC connections. It provides secure token generation for both viewers and broadcasters in the auction platform.

## Features

- JWT token generation for LiveKit rooms
- Role-based permissions (viewer, host, broadcaster)
- Secure token signing with LiveKit API secret
- Health check endpoint
- CORS support for frontend applications

## API Endpoints

### Health Check
- **GET** `/api/v1/health`
- Returns service health status

### Token Generation
- **GET** `/api/v1/livekit/token`
- **GET** `/api/livekit/token` (legacy endpoint)
- Query Parameters:
  - `room` (required): Room name
  - `role` (optional): viewer, host, or broadcaster (default: viewer)
  - `name` (optional): Display name for participant

## Configuration

Environment Variables:
- `PORT`: Service port (default: 8089)
- `LIVEKIT_API_KEY`: LiveKit API key (required)
- `LIVEKIT_API_SECRET`: LiveKit API secret (required)
- `LIVEKIT_URL`: LiveKit server URL (default: https://livekit.blytz.app)
- `AUTH_SERVICE_URL`: Auth service URL (default: http://auth-service:8084)

## Token Permissions

### Viewer Role
- Can subscribe to audio/video
- Cannot publish audio/video
- Cannot publish data

### Host/Broadcaster Role
- Can subscribe to audio/video
- Can publish audio/video
- Can publish data
- Room admin privileges

## Usage Example

```bash
# Generate viewer token
curl "http://localhost:8089/api/livekit/token?room=demo-auction-123&role=viewer"

# Generate host token
curl "http://localhost:8089/api/livekit/token?room=demo-auction-123&role=host&name=AuctionHost"
```

## Integration

The service integrates with:
- LiveKit server for WebRTC streaming
- Auth service for user authentication
- Frontend applications for token retrieval
const express = require('express');
const jwt = require('jsonwebtoken');
const cors = require('cors');

const app = express();
const port = process.env.PORT || 8089;

// Middleware
app.use(cors());
app.use(express.json());

// LiveKit token endpoint
app.get('/api/livekit/token', (req, res) => {
  const { room = 'demo-room', name = 'demo-user' } = req.query;
  
  const apiKey = process.env.LIVEKIT_API_KEY;
  const apiSecret = process.env.LIVEKIT_API_SECRET;
  
  if (!apiKey || !apiSecret) {
    return res.status(500).json({ error: 'LiveKit credentials not configured' });
  }

  try {
    const payload = {
      iss: apiKey,
      nbf: Math.floor(Date.now() / 1000),
      exp: Math.floor(Date.now() / 1000) + (24 * 60 * 60), // 24 hours
      sub: name,
      video: {
        room: room,
        roomJoin: true
      }
    };

    const token = jwt.sign(payload, apiSecret);
    
    res.json({
      token: token,
      url: 'wss://livekit.blytz.app'
    });
  } catch (error) {
    console.error('Error generating token:', error);
    res.status(500).json({ error: 'Failed to generate token' });
  }
});

// Health endpoint
app.get('/health', (req, res) => {
  res.json({
    status: 'ok',
    service: 'livekit-token-service'
  });
});

app.listen(port, () => {
  console.log(`LiveKit token service running on port ${port}`);
});
# Blytz Seller Frontend

Broadcaster platform for the Blytz live auction platform. This frontend enables auction hosts to create and manage live auctions, stream video, and interact with bidders.

## Features

- **Live Broadcasting**: Real-time video streaming with LiveKit
- **Auction Management**: Create, start, and manage auctions
- **Real-time Analytics**: Track bidding activity and engagement
- **Chat Moderation**: Manage participant discussions
- **Dashboard**: Comprehensive auction analytics with Recharts
- **Responsive Design**: Mobile-first design with Tailwind CSS

## Development

### Prerequisites

- Node.js 18+ 
- npm or yarn
- Docker (for local development with services)

### Getting Started

1. **Install dependencies**:
   ```bash
   npm install
   ```

2. **Set up environment variables**:
   ```bash
   cp .env.local.example .env.local
   # Update .env.local with your configuration
   ```

3. **Start development server**:
   ```bash
   npm run dev
   ```

   The app will be available at `http://localhost:3002`

### Available Scripts

- `npm run dev` - Start development server on port 3002
- `npm run build` - Build for production
- `npm run start` - Start production server
- `npm run lint` - Run ESLint
- `npm run lint:fix` - Fix ESLint issues
- `npm run format` - Format code with Prettier
- `npm run format:check` - Check code formatting
- `npm run typecheck` - Run TypeScript type checking
- `npm run test` - Run Jest tests
- `npm run test:watch` - Run tests in watch mode
- `npm run test:coverage` - Run tests with coverage

### Technology Stack

- **Framework**: Next.js 14.2.13 with App Router
- **UI**: React 18 with TypeScript
- **Styling**: Tailwind CSS with Radix UI components
- **Real-time**: LiveKit for WebRTC streaming
- **State Management**: TanStack Query for server state
- **Analytics**: Recharts for data visualization
- **Testing**: Jest with React Testing Library
- **Code Quality**: ESLint, Prettier, Husky pre-commit hooks

### Project Structure

```
src/
├── app/                 # Next.js App Router pages
│   ├── globals.css     # Global styles
│   ├── layout.tsx      # Root layout
│   └── page.tsx        # Dashboard page
├── components/          # Reusable components
│   ├── ui/             # Base UI components
│   └── LiveKitBroadcaster.tsx # LiveKit streaming
└── lib/
    └── utils.ts        # Utility functions
```

### Environment Variables

Key environment variables in `.env.local`:

```bash
# API Configuration
NEXT_PUBLIC_API_URL=http://localhost:8080
NEXT_PUBLIC_WS_URL=ws://localhost:8080

# LiveKit Configuration
NEXT_PUBLIC_LIVEKIT_URL=ws://localhost:7880

# Seller Configuration
NEXT_PUBLIC_DEMO_MODE=broadcaster
NEXT_PUBLIC_ENABLE_ANALYTICS=false
NEXT_PUBLIC_ENABLE_DEBUG=true
```

### Development with Full Stack

To run the complete Blytz platform locally:

1. **Start backend services**:
   ```bash
   cd ../../
   docker-compose up -d
   ```

2. **Start this frontend**:
   ```bash
   npm run dev
   ```

3. **Access the application**:
   - Seller Dashboard: http://localhost:3002
   - Demo Viewer: http://localhost:3001
   - Main Frontend: http://localhost:3000
   - API Gateway: http://localhost:8080

### Broadcasting Features

The seller frontend provides:

- **Stream Setup**: Configure video and audio sources
- **Auction Controls**: Start, pause, and end auctions
- **Bid Management**: View and approve bids
- **Analytics Dashboard**: Real-time charts and metrics
- **Chat Management**: Moderate participant discussions

### Analytics Dashboard

Using Recharts for comprehensive data visualization:

- **Real-time Bidding**: Live bid activity charts
- **User Engagement**: Participant metrics
- **Revenue Tracking**: Auction performance data
- **Technical Metrics**: Stream quality and bandwidth

### Code Quality

This project maintains high code quality standards:

- **TypeScript**: Strict mode enabled for type safety
- **ESLint**: Next.js configuration with custom rules
- **Prettier**: Consistent code formatting
- **Husky**: Pre-commit hooks for code quality
- **Testing**: Jest with React Testing Library

### Accessibility

The application follows WCAG 2.1 AA guidelines:

- Semantic HTML5 elements
- ARIA labels and roles
- Keyboard navigation support
- Screen reader compatibility
- Focus management

### Performance

- **Bundle Optimization**: Next.js automatic code splitting
- **Image Optimization**: Next.js Image component
- **Caching**: TanStack Query for efficient data fetching
- **Lazy Loading**: Components and routes as needed
- **Stream Optimization**: Efficient WebRTC handling

## Deployment

### Docker Deployment

```bash
# Build image
docker build -t blytz-seller-frontend .

# Run container
docker run -p 3002:3002 blytz-seller-frontend
```

### Production Deployment

The application is configured for production deployment with:

- Standalone build output for containerization
- Environment-specific configuration
- Optimized bundles and assets
- Stream configuration for production LiveKit

## Broadcasting Setup

### Hardware Requirements

- **Camera**: HD webcam or professional camera
- **Microphone**: Clear audio input device
- **Lighting**: Proper lighting for video quality
- **Internet**: Stable upload connection (5+ Mbps)

### Software Configuration

1. **Browser Permissions**: Allow camera and microphone access
2. **LiveKit Room**: Configure room name and token
3. **Stream Quality**: Adjust resolution and bitrate
4. **Audio Settings**: Configure input source and levels

## Contributing

1. Follow the existing code style and patterns
2. Ensure all tests pass before submitting
3. Test streaming functionality thoroughly
4. Use semantic commit messages
5. Update documentation as needed

## Support

For issues and questions:

1. Check the [main documentation](../../README.md)
2. Review the [development guide](../../DEVELOPMENT.md)
3. Check LiveKit documentation for streaming issues
4. Create an issue in the project repository
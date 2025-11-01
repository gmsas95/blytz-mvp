# Blytz Demo Frontend

Viewer platform for the Blytz live auction platform. This frontend allows users to watch live auctions, place bids, and participate in real-time chat.

## Features

- **Live Auction Viewing**: Real-time streaming with LiveKit integration
- **Real-time Bidding**: Place bids with instant updates
- **Chat System**: Participate in auction discussions
- **Responsive Design**: Mobile-first design with Tailwind CSS
- **Accessibility**: WCAG 2.1 AA compliant components

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

   The app will be available at `http://localhost:3001`

### Available Scripts

- `npm run dev` - Start development server on port 3001
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
- **Testing**: Jest with React Testing Library
- **Code Quality**: ESLint, Prettier, Husky pre-commit hooks

### Project Structure

```
src/
├── app/                 # Next.js App Router pages
│   ├── globals.css     # Global styles
│   ├── layout.tsx      # Root layout
│   └── page.tsx        # Home page
├── components/          # Reusable components
│   ├── ui/             # Base UI components
│   └── LiveKitViewer.tsx # LiveKit integration
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

# Demo Configuration
NEXT_PUBLIC_DEMO_MODE=viewer
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
   - Demo Viewer: http://localhost:3001
   - Main Frontend: http://localhost:3000
   - API Gateway: http://localhost:8080

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

## Deployment

### Docker Deployment

```bash
# Build image
docker build -t blytz-demo-frontend .

# Run container
docker run -p 3001:3001 blytz-demo-frontend
```

### Production Deployment

The application is configured for production deployment with:

- Standalone build output for containerization
- Environment-specific configuration
- Optimized bundles and assets

## Contributing

1. Follow the existing code style and patterns
2. Ensure all tests pass before submitting
3. Use semantic commit messages
4. Update documentation as needed

## Support

For issues and questions:

1. Check the [main documentation](../../README.md)
2. Review the [development guide](../../DEVELOPMENT.md)
3. Create an issue in the project repository
# Blytz Frontend - Modern React Commerce Platform

A modern, mobile-first React + TypeScript frontend for the Blytz live auction
commerce platform. Built with Next.js 14, Tailwind CSS, shadcn/ui, and Framer
Motion for a polished user experience.

## 🚀 Features

- **Modern Architecture**: Next.js 14 with App Router, React 18, TypeScript
- **Mobile-First Design**: Responsive design with touch-friendly interactions
- **Live Auctions**: Real-time bidding interface with countdown timers
- **Live Streaming**: Video streaming integration for seller showcases
- **Shopping Cart**: Full cart functionality with auction integration
- **Seller Profiles**: Individual seller storefronts and ratings
- **Animations**: Smooth animations with Framer Motion
- **Accessibility**: WCAG compliant with keyboard navigation
- **Performance**: Optimized images, code splitting, and caching

## 🛠 Tech Stack

- **Framework**: Next.js 14 (React 18)
- **Language**: TypeScript
- **Styling**: Tailwind CSS with custom theme
- **UI Components**: shadcn/ui + Radix UI primitives
- **Animations**: Framer Motion
- **Icons**: Lucide React
- **Data Fetching**: TanStack Query (React Query)
- **Build Tool**: Next.js built-in (SWC)
- **Deployment**: Docker with Dokploy

## 📦 Installation

```bash
# Clone the repository
git clone <your-repo-url>
cd blytzmvp-clean/frontend

# Install dependencies
npm install

# Run development server
npm run dev
```

## 🎯 Environment Variables

Create a `.env.local` file:

```bash
# API Configuration
MODE=mock                    # 'mock' for mock data, 'remote' for real API
NEXT_PUBLIC_API_BASE=http://localhost:8080

# Optional: Real API base URL for production
REMOTE_API_BASE=https://api.blytz.app
```

## 🏗 Architecture

### Folder Structure

```
src/
├── app/                    # Next.js App Router pages
│   ├── api/               # API routes
│   ├── home/              # Home page components
│   ├── products/          # Product listing pages
│   ├── product/           # Individual product pages
│   ├── livestream/        # Live streaming pages
│   ├── cart/              # Shopping cart
│   ├── checkout/          # Checkout process
│   ├── auth/              # Authentication pages
│   └── profile/           # User profile pages
├── components/            # Reusable components
│   ├── ui/               # shadcn/ui components
│   ├── layout/           # Layout components (Header, Footer)
│   ├── home/             # Home page sections
│   ├── products/         # Product components
│   ├── livestream/       # Streaming components
│   └── cart/             # Cart components
├── lib/                  # Utility functions and configurations
│   ├── utils.ts          # Helper functions
│   ├── api-adapter.ts    # API adapter pattern
│   └── hooks/            # Custom React hooks
├── types/                # TypeScript type definitions
├── data/                 # Mock data and data layer
└── styles/               # Global styles and themes
```

### API Adapter Pattern

The frontend uses an adapter pattern to switch between mock and real APIs:

```typescript
// Mock mode (default)
MODE=mock

// Remote mode (production)
MODE=remote
NEXT_PUBLIC_API_BASE=https://api.blytz.app
```

### Key Components

#### UI Components (shadcn/ui)

- `Button`: Custom button variants with animations
- `Card`: Enhanced cards with hover effects
- `Input`: Styled form inputs
- Custom auction-specific components

#### Layout Components

- `Header`: Responsive navigation with mobile menu
- `Footer`: Comprehensive footer with links and newsletter
- `Hero`: Animated hero section with CTA buttons

#### Data Components

- `FeaturedProducts`: Product grid with loading states
- `ActiveAuctions`: Live auction cards with countdown timers
- `LiveStreams`: Live streaming preview cards

## 🎨 Design System

### Colors

```css
--background: 0 0% 99%; /* #F9FAFB */
--foreground: 222 47% 11%; /* #111827 */
--primary: 221 83% 53%; /* #2563EB */
--border-radius: 1rem; /* rounded-2xl */
```

### Typography

- **UI Font**: Geist Sans
- **Body Font**: Inter
- **Scale**: Mobile-first with responsive breakpoints

### Animations

- Page transitions with Framer Motion
- Micro-interactions on buttons and cards
- Smooth hover effects and loading states
- Real-time auction countdown animations

## 📱 Mobile-First Features

### Touch Interactions

- Optimized button sizes for touch
- Swipe gestures for image galleries
- Pull-to-refresh functionality
- Bottom navigation on mobile

### Responsive Design

- Breakpoints: `sm: 640px`, `md: 768px`, `lg: 1024px`, `xl: 1280px`
- Mobile-first CSS approach
- Optimized images with responsive sizing
- Touch-friendly form controls

## 🚀 Deployment

### Docker Deployment

```bash
# Build the Docker image
docker build -t blytz-frontend .

# Run with Docker Compose
docker-compose up -d
```

### Dokploy Deployment

```bash
# Deploy to Dokploy
dokploy deploy -f dokploy.json
```

### Environment Setup

The frontend supports two modes:

1. **Mock Mode** (Development/Default)
   - Uses mock data for all API calls
   - No backend dependencies
   - Perfect for development and demos

2. **Remote Mode** (Production)
   - Connects to real backend API
   - Requires backend services to be running
   - Full production functionality

## 🔧 Development

### Available Scripts

```bash
npm run dev          # Development server
npm run build        # Production build
npm run start        # Start production server
npm run lint         # ESLint
npm run type-check   # TypeScript checking
npm run format       # Prettier formatting
npm run test         # Run tests
```

### Code Quality

- ESLint configuration for React/Next.js
- Prettier for code formatting
- TypeScript strict mode
- Pre-commit hooks for quality checks

## 🧪 Testing

### Test Structure

```
tests/
├── unit/              # Unit tests
├── integration/       # Integration tests
└── e2e/              # End-to-end tests
```

### Test Commands

```bash
npm test              # Run all tests
npm run test:watch    # Watch mode
npm run test:e2e      # E2E tests
```

## 📊 Performance

### Optimization Features

- Image optimization with Next.js Image component
- Code splitting with dynamic imports
- Lazy loading for components
- Efficient re-renders with React.memo
- Optimized bundle size with tree shaking

### Core Web Vitals

- **LCP**: Optimized with image loading strategies
- **FID**: Minimal JavaScript execution
- **CLS**: Stable layouts with proper sizing

## 🔒 Security

### Security Features

- Content Security Policy headers
- XSS protection
- CSRF protection
- Input validation and sanitization
- Secure cookie handling

## 📚 Documentation

### Component Documentation

Each component includes:

- JSDoc comments for props and usage
- TypeScript interfaces for type safety
- Usage examples in Storybook (optional)

### API Documentation

- OpenAPI specification available at `/api/docs`
- Type-safe API client generation
- Mock data documentation

## 🤝 Contributing

### Development Workflow

1. Create feature branch from `main`
2. Implement changes with tests
3. Run quality checks
4. Submit pull request
5. Code review and merge

### Code Standards

- Follow TypeScript best practices
- Use semantic HTML
- Implement accessibility features
- Write comprehensive tests
- Document public APIs

## 🆘 Troubleshooting

### Common Issues

**Build fails with TypeScript errors**

```bash
npm run type-check
npm run lint
```

**Images not loading**

- Check image domains in `next.config.js`
- Verify image URLs are accessible

**API calls failing**

- Check environment variables
- Verify API adapter mode setting
- Check network connectivity

**Performance issues**

- Use React DevTools Profiler
- Check bundle size with `npm run build`
- Enable performance monitoring

## 📞 Support

For support and questions:

- Check the troubleshooting section
- Review component documentation
- Check the GitHub issues
- Contact the development team

## 📄 License

This project is part of the Blytz MVP and follows the same licensing terms. See
the main repository for license information."}

---

**Last Updated**: $(date) **Version**: 2.0.0 **Status**: Production Ready

# Demo Frontend - Viewer Platform with LiveKit
FROM node:18-alpine AS builder

WORKDIR /app

# Copy package files
COPY package*.json ./
RUN npm ci

# Copy source code
COPY . .

# Build the application
RUN npm run build

# Production stage
FROM node:18-alpine AS runner

WORKDIR /app

ENV NODE_ENV=production
ENV PORT=3001

# Create non-root user
RUN addgroup -g 1001 -S nodejs
RUN adduser -S nextjs -u 1001

# Copy built application
COPY --from=builder --chown=nextjs:nodejs /app/.next/standalone ./
COPY --from=builder --chown=nextjs:nodejs /app/.next/static ./.next/static
# Copy public directory if it exists (optional for Next.js) - create empty dir first
RUN mkdir -p public
COPY --from=builder --chown=nextjs:nodejs /app/public ./public || true

USER nextjs

EXPOSE 3001

ENV HOSTNAME="0.0.0.0"

CMD ["node", "server.js"]
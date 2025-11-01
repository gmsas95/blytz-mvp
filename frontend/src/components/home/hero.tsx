'use client';

import { ArrowRight, Play } from 'lucide-react';
import Link from 'next/link';

import { Button } from '@/components/ui/button';

export function Hero() {
  return (
    <section className="w-full bg-background" aria-labelledby="hero-heading">
      <div className="container mx-auto px-4 py-16 md:py-24">
        {/* Main content area */}
        <div className="flex flex-col lg:flex-row items-center gap-12 lg:gap-16">
          {/* Left column - Text content */}
          <div className="flex-1 space-y-8 text-center lg:text-left">
            {/* Status badge */}
            <div
              className="inline-flex items-center gap-2 bg-primary/10 text-primary px-4 py-2 rounded-full text-sm font-medium"
              role="status"
              aria-live="polite"
            >
              <div className="w-2 h-2 bg-primary rounded-full animate-pulse" aria-hidden="true" />
              Live Auctions Now Active
            </div>

            {/* Headline */}
            <div className="space-y-4">
              <h1
                id="hero-heading"
                className="text-4xl md:text-5xl lg:text-6xl font-bold tracking-tight text-foreground"
              >
                Shop Live.
                <span className="text-primary"> Bid Real.</span>
                Win Big.
              </h1>
              <p className="text-lg md:text-xl text-muted-foreground max-w-2xl">
                Join live streaming auctions where sellers showcase products in real-time. Bid,
                interact, and score amazing deals on unique items.
              </p>
            </div>

            {/* CTA buttons */}
            <div
              className="flex flex-col sm:flex-row gap-4 justify-center lg:justify-start"
              role="group"
              aria-label="Primary actions"
            >
              <Link href="/livestream">
                <Button size="lg" className="gap-2" aria-describedby="live-action-description">
                  Watch Live Now
                  <Play className="w-4 h-4" aria-hidden="true" />
                </Button>
              </Link>
              <span id="live-action-description" className="sr-only">
                Join current live streaming auctions and bid in real-time
              </span>

              <Link href="/products">
                <Button
                  variant="outline"
                  size="lg"
                  className="gap-2"
                  aria-describedby="browse-action-description"
                >
                  Browse Products
                  <ArrowRight className="w-4 h-4" aria-hidden="true" />
                </Button>
              </Link>
              <span id="browse-action-description" className="sr-only">
                Browse available products and upcoming auctions
              </span>
            </div>

            {/* Stats */}
            <div
              className="grid grid-cols-3 gap-8 pt-8"
              role="region"
              aria-label="Platform statistics"
            >
              <div className="text-center">
                <div
                  className="text-2xl md:text-3xl font-bold text-foreground"
                  aria-label="500 plus"
                >
                  500+
                </div>
                <div className="text-sm text-muted-foreground">Live Streams</div>
              </div>

              <div className="text-center">
                <div
                  className="text-2xl md:text-3xl font-bold text-foreground"
                  aria-label="10 thousand plus"
                >
                  10K+
                </div>
                <div className="text-sm text-muted-foreground">Active Bidders</div>
              </div>

              <div className="text-center">
                <div
                  className="text-2xl md:text-3xl font-bold text-foreground"
                  aria-label="50 thousand plus"
                >
                  50K+
                </div>
                <div className="text-sm text-muted-foreground">Products Sold</div>
              </div>
            </div>
          </div>

          {/* Right column - Visual element */}
          <div className="flex-1 w-full max-w-md lg:max-w-none">
            <div
              className="relative bg-secondary/30 rounded-3xl p-8 lg:p-12"
              role="img"
              aria-label="Live auction preview interface showing active bidding"
            >
              {/* Simple placeholder for product preview */}
              <div className="space-y-6">
                <div className="flex items-center gap-3">
                  <div
                    className="w-3 h-3 bg-red-500 rounded-full animate-pulse"
                    aria-hidden="true"
                  />
                  <span className="text-sm font-medium text-muted-foreground">Live Preview</span>
                </div>

                <div className="space-y-3" aria-hidden="true">
                  <div className="h-4 bg-muted rounded-full w-full" />
                  <div className="h-4 bg-muted rounded-full w-4/5" />
                  <div className="h-4 bg-muted rounded-full w-3/5" />
                </div>

                {/* Product grid placeholder */}
                <div className="grid grid-cols-2 gap-4 pt-4" aria-hidden="true">
                  <div className="aspect-square bg-muted rounded-2xl" />
                  <div className="aspect-square bg-muted rounded-2xl" />
                  <div className="aspect-square bg-muted rounded-2xl" />
                  <div className="aspect-square bg-muted rounded-2xl" />
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}

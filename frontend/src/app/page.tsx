'use client';

import { Loader2 } from 'lucide-react';
import dynamic from 'next/dynamic';
import { Suspense } from 'react';

// Dynamic imports with loading states
const Hero = dynamic(() => import('@/components/home/hero').then(mod => ({ default: mod.Hero })), {
  loading: () => <div className="h-96 flex items-center justify-center"><Loader2 className="w-8 h-8 animate-spin" /></div>,
  ssr: true
});

const FeaturedProducts = dynamic(() => import('@/components/home/featured-products').then(mod => ({ default: mod.FeaturedProducts })), {
  loading: () => <div className="h-64 flex items-center justify-center"><Loader2 className="w-6 h-6 animate-spin" /></div>,
  ssr: true
});

const ActiveAuctions = dynamic(() => import('@/components/home/active-auctions').then(mod => ({ default: mod.ActiveAuctions })), {
  loading: () => <div className="h-64 flex items-center justify-center"><Loader2 className="w-6 h-6 animate-spin" /></div>,
  ssr: true
});

const LiveStreams = dynamic(() => import('@/components/home/live-streams').then(mod => ({ default: mod.LiveStreams })), {
  loading: () => <div className="h-64 flex items-center justify-center"><Loader2 className="w-6 h-6 animate-spin" /></div>,
  ssr: true
});

export default function HomePage() {
  return (
    <div className="bg-background">
      <Suspense fallback={<div className="h-96 flex items-center justify-center"><Loader2 className="w-8 h-8 animate-spin" /></div>}>
        <Hero />
      </Suspense>
      <Suspense fallback={<div className="h-64 flex items-center justify-center"><Loader2 className="w-6 h-6 animate-spin" /></div>}>
        <FeaturedProducts />
      </Suspense>
      <Suspense fallback={<div className="h-64 flex items-center justify-center"><Loader2 className="w-6 h-6 animate-spin" /></div>}>
        <ActiveAuctions />
      </Suspense>
      <Suspense fallback={<div className="h-64 flex items-center justify-center"><Loader2 className="w-6 h-6 animate-spin" /></div>}>
        <LiveStreams />
      </Suspense>
    </div>
  );
}

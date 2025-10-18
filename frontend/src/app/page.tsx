'use client'

import { Hero } from '@/components/home/hero'
import { FeaturedProducts } from '@/components/home/featured-products'
import { ActiveAuctions } from '@/components/home/active-auctions'
import { LiveStreams } from '@/components/home/live-streams'
import { TestComponent } from '@/components/home/test-component'

export default function HomePage() {
  return (
    <div className="bg-background">
      <TestComponent />
      <Hero />
      <FeaturedProducts />
      <ActiveAuctions />
      <LiveStreams />
    </div>
  )
}

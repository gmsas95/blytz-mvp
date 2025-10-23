'use client'

import { Hero } from '@/components/home/hero'
import { FeaturedProducts } from '@/components/home/featured-products'
import { ActiveAuctions } from '@/components/home/active-auctions'
import { LiveStreams } from '@/components/home/live-streams'

export default function HomePage() {
  return (
    <div className="bg-background">
      <Hero />
      <FeaturedProducts />
      <ActiveAuctions />
      <LiveStreams />
    </div>
  )
}

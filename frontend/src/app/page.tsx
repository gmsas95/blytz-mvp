import { Header } from '@/components/layout/header'
import { Hero } from '@/components/home/hero'
import { FeaturedProducts } from '@/components/home/featured-products'
import { ActiveAuctions } from '@/components/home/active-auctions'
import { LiveStreams } from '@/components/home/live-streams'
import { Footer } from '@/components/layout/footer'

export default function HomePage() {
  return (
    <div className="min-h-screen bg-background">
      <Header />
      <main className="flex-1">
        <Hero />
        <FeaturedProducts />
        <ActiveAuctions />
        <LiveStreams />
      </main>
      <Footer />
    </div>
  )
}
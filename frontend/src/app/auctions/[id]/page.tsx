import { notFound } from 'next/navigation'
import { api } from '@/lib/api-adapter'
import { formatPrice, formatTimeRemaining } from '@/lib/utils'
import { Button } from '@/components/ui/button'

interface Props { params: { id: string } }

export default async function AuctionDetailPage({ params }: Props) {
  const res = await api.getAuction(params.id)
  if (!res.success || !res.data) return notFound()
  const auction = res.data

  const remaining = Math.max(0, Math.floor((new Date(auction.endTime).getTime() - Date.now())/1000))

  return (
    <section className="w-full py-16 md:py-24">
      <div className="container mx-auto px-4 grid gap-10 lg:grid-cols-2 items-start">
        <div className="w-full">
          <div className="aspect-square bg-muted rounded-2xl overflow-hidden">
            {/* eslint-disable-next-line @next/next/no-img-element */}
            <img src={auction.product.images[0]} alt={auction.product.title} className="object-cover w-full h-full" />
          </div>
        </div>
        <div className="space-y-6">
          <div className="flex items-center gap-2">
            <span className={`text-xs px-2 py-1 rounded-full ${auction.isLive ? 'bg-red-100 text-red-800' : 'bg-yellow-100 text-yellow-800'}`}>{auction.isLive ? 'LIVE' : auction.status.toUpperCase()}</span>
            <span className={`text-xs px-2 py-1 rounded-full ${remaining < 300 ? 'bg-orange-100 text-orange-800' : 'bg-gray-100 text-gray-800'}`}>{formatTimeRemaining(remaining)}</span>
          </div>
          <h1 className="text-3xl md:text-4xl font-bold tracking-tight">{auction.product.title}</h1>
          <p className="text-muted-foreground">{auction.product.description}</p>
          <div className="space-y-1">
            <div className="text-3xl font-bold text-primary">{formatPrice(auction.currentBid)}</div>
            <div className="text-sm text-muted-foreground">Min increment: +{formatPrice(auction.minBidIncrement)}</div>
          </div>
          <div className="flex gap-3">
            <Button size="lg">Place Bid</Button>
            <a href={`/product/${auction.productId}`}>
              <Button size="lg" variant="outline">View Product</Button>
            </a>
          </div>
        </div>
      </div>
    </section>
  )
}
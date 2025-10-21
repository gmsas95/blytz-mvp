import Link from 'next/link'
import { api } from '@/lib/api-adapter'
import { Card, CardContent, CardFooter, CardHeader } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { formatPrice, formatTimeRemaining } from '@/lib/utils'

export default async function AuctionsPage() {
  const res = await api.getAuctions()
  const items = res.success && res.data ? res.data.items : []

  return (
    <section className="w-full py-16 md:py-24">
      <div className="container mx-auto px-4">
        <div className="text-center space-y-4 mb-12">
          <h1 className="text-3xl md:text-4xl font-bold tracking-tight">Auctions</h1>
          <p className="text-muted-foreground">Join live and upcoming auctions</p>
        </div>
        {items.length === 0 ? (
          <div className="text-center text-muted-foreground">No auctions found.</div>
        ) : (
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
            {items.map((auction) => (
              <Link key={auction.id} href={`/auctions/${auction.id}`} className="group">
                <Card className="overflow-hidden h-full">
                  <CardHeader className="pb-3">
                    <div className="flex items-center justify-between">
                      <span className={`text-xs px-2 py-1 rounded-full ${auction.isLive ? 'bg-red-100 text-red-800' : 'bg-yellow-100 text-yellow-800'}`}>{auction.isLive ? 'LIVE' : auction.status.toUpperCase()}</span>
                      <span className="text-xs bg-gray-100 text-gray-800 px-2 py-1 rounded-full">{formatTimeRemaining(Math.max(0, Math.floor((new Date(auction.endTime).getTime() - Date.now())/1000)))}</span>
                    </div>
                  </CardHeader>
                  <CardContent className="pb-3">
                    <h3 className="font-semibold text-lg leading-tight">{auction.product.title}</h3>
                    <p className="text-sm text-muted-foreground line-clamp-2">{auction.product.description}</p>
                    <div className="mt-4 text-sm">
                      <div>Current Bid: <span className="font-semibold">{formatPrice(auction.currentBid)}</span></div>
                      <div className="text-muted-foreground">Min increment: +{formatPrice(auction.minBidIncrement)}</div>
                    </div>
                  </CardContent>
                  <CardFooter>
                    <Button className="w-full">View Auction</Button>
                  </CardFooter>
                </Card>
              </Link>
            ))}
          </div>
        )}
      </div>
    </section>
  )
}
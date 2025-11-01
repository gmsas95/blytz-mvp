'use client'

import { useState } from 'react'
import { Play, Eye, Clock, Users, Gavel } from 'lucide-react'
import LiveKitViewer from '../components/LiveKitViewer'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Badge } from '@/components/ui/badge'

export default function Home() {
  const [auctionId, setAuctionId] = useState('demo-auction-123')

  return (
    <div className="min-h-screen bg-background">
      {/* Modern Header */}
      <header className="border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
        <div className="container mx-auto px-4 py-6">
          <div className="flex items-center justify-between">
            <div>
              <h1 className="text-3xl font-bold tracking-tight">Live Auction Demo</h1>
              <p className="text-muted-foreground">Viewer Platform - Watch live auctions in real-time</p>
            </div>
            <Badge variant="secondary" className="gap-2">
              <div className="w-2 h-2 bg-green-500 rounded-full animate-pulse"/>
              Live Now
            </Badge>
          </div>
        </div>
      </header>

      <main className="container mx-auto px-4 py-8">
        {/* Auction Selection */}
        <Card className="mb-8">
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Gavel className="w-5 h-5" />
              Join Auction Room
            </CardTitle>
            <CardDescription>
              Enter the auction room ID to start watching the live stream
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="flex flex-col sm:flex-row gap-4">
              <div className="flex-1">
                <Label htmlFor="auction-id">Auction Room ID</Label>
                <Input 
                  id="auction-id"
                  type="text"
                  value={auctionId}
                  onChange={(e) => setAuctionId(e.target.value)}
                  placeholder="Enter auction ID"
                  className="max-w-md"
                />
              </div>
              <div className="flex items-end">
                <Button size="lg" className="gap-2">
                  <Play className="w-4 h-4" />
                  Join Stream
                </Button>
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Live Stream Viewer */}
        <Card className="mb-8 overflow-hidden">
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <div className="w-3 h-3 bg-red-500 rounded-full animate-pulse"/>
              Live Stream
            </CardTitle>
            <CardDescription>
              Real-time auction broadcast with interactive bidding
            </CardDescription>
          </CardHeader>
          <CardContent className="p-0">
            <LiveKitViewer 
              auctionId={auctionId}
              onConnected={() => console.log('Connected to live stream')}
              onDisconnected={() => console.log('Disconnected from live stream')}
              onError={(error) => console.error('Stream error:', error)}
            />
          </CardContent>
        </Card>

        {/* Auction Information */}
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Gavel className="w-5 h-5" />
                Current Auction
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                <div>
                  <p className="text-sm text-muted-foreground">Item</p>
                  <p className="font-semibold">Vintage Watch Collection</p>
                </div>
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <p className="text-sm text-muted-foreground">Current Bid</p>
                    <p className="text-2xl font-bold text-primary">$1,250</p>
                  </div>
                  <div>
                    <p className="text-sm text-muted-foreground">Bidders</p>
                    <p className="text-2xl font-bold text-primary">12</p>
                  </div>
                </div>
                <div>
                  <p className="text-sm text-muted-foreground">Time Remaining</p>
                  <p className="text-lg font-semibold flex items-center gap-2">
                    <Clock className="w-4 h-4" />
                    5m 32s
                  </p>
                </div>
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Eye className="w-5 h-5" />
                Quick Actions
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-3">
              <Button className="w-full" size="lg">
                Place Bid
              </Button>
              <Button variant="outline" className="w-full">
                View Item Details
              </Button>
              <Button variant="outline" className="w-full">
                Auction History
              </Button>
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Users className="w-5 h-5" />
                Live Activity
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                <div className="flex items-center justify-between">
                  <span className="text-sm text-muted-foreground">Viewers</span>
                  <Badge variant="secondary">247</Badge>
                </div>
                <div className="flex items-center justify-between">
                  <span className="text-sm text-muted-foreground">Total Bids</span>
                  <Badge variant="secondary">47</Badge>
                </div>
                <div className="flex items-center justify-between">
                  <span className="text-sm text-muted-foreground">Status</span>
                  <Badge className="gap-1">
                    <div className="w-2 h-2 bg-white rounded-full animate-pulse"/>
                    Active
                  </Badge>
                </div>
              </div>
            </CardContent>
          </Card>
        </div>
      </main>
    </div>
  )
}
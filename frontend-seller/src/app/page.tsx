'use client'

import { useState } from 'react'
import { Radio, Settings, Users, DollarSign, Clock, TrendingUp, Video, Mic, MicOff, VideoOff, Square } from 'lucide-react'
import LiveKitBroadcaster from '../components/LiveKitBroadcaster'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Badge } from '@/components/ui/badge'

export default function Home() {
  const [auctionId, setAuctionId] = useState('demo-auction-123')
  const [isBroadcasting, setIsBroadcasting] = useState(false)

  return (
    <div className="min-h-screen bg-background">
      {/* Modern Header */}
      <header className="border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
        <div className="container mx-auto px-4 py-6">
          <div className="flex items-center justify-between">
            <div>
              <h1 className="text-3xl font-bold tracking-tight">Seller Dashboard</h1>
              <p className="text-muted-foreground">Broadcast Studio - Host live auctions</p>
            </div>
            <Badge variant={isBroadcasting ? "destructive" : "secondary"} className="gap-2">
              <div className={`w-2 h-2 rounded-full ${isBroadcasting ? 'bg-white animate-pulse' : 'bg-current'}`}/>
              {isBroadcasting ? 'LIVE' : 'OFFLINE'}
            </Badge>
          </div>
        </div>
      </header>

      <main className="container mx-auto px-4 py-8">
        {/* Auction Setup */}
        <Card className="mb-8">
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Radio className="w-5 h-5" />
              Auction Room Setup
            </CardTitle>
            <CardDescription>
              Configure your live auction broadcast room
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
                  placeholder="Enter auction room ID"
                  className="max-w-md"
                />
              </div>
              <div className="flex items-end">
                <Button size="lg" className="gap-2" variant={isBroadcasting ? "destructive" : "default"}>
                  {isBroadcasting ? (
                    <>
                      <Square className="w-4 h-4" />
                      End Broadcast
                    </>
                  ) : (
                    <>
                      <Radio className="w-4 h-4" />
                      Start Broadcast
                    </>
                  )}
                </Button>
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Main Content Grid */}
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          {/* Broadcast Studio - Takes 2 columns */}
          <div className="lg:col-span-2">
            <Card className="overflow-hidden">
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <div className={`w-3 h-3 rounded-full ${isBroadcasting ? 'bg-red-500 animate-pulse' : 'bg-gray-500'}`}/>
                  Broadcast Studio
                </CardTitle>
                <CardDescription>
                  {isBroadcasting ? 'Live auction in progress' : 'Ready to start broadcasting'}
                </CardDescription>
              </CardHeader>
              <CardContent className="p-0">
                <LiveKitBroadcaster 
                  auctionId={auctionId}
                  onBroadcastStart={() => {
                    setIsBroadcasting(true)
                    console.log('Broadcast started')
                  }}
                  onBroadcastEnd={() => {
                    setIsBroadcasting(false)
                    console.log('Broadcast ended')
                  }}
                  onError={(error) => console.error('Broadcast error:', error)}
                  onViewerCount={(count) => console.log('Viewers:', count)}
                />
              </CardContent>
            </Card>
          </div>

          {/* Sidebar Controls */}
          <div className="space-y-6">
            {/* Auction Controls */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Settings className="w-5 h-5" />
                  Auction Controls
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-3">
                <Button 
                  className="w-full" 
                  variant={isBroadcasting ? "default" : "secondary"}
                  disabled={!isBroadcasting}
                >
                  Start Bidding
                </Button>
                <Button 
                  variant="outline" 
                  className="w-full"
                  disabled={!isBroadcasting}
                >
                  Pause Auction
                </Button>
                <Button 
                  variant="destructive" 
                  className="w-full"
                  disabled={!isBroadcasting}
                >
                  End Auction
                </Button>
              </CardContent>
            </Card>

            {/* Live Stats */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <TrendingUp className="w-5 h-5" />
                  Live Stats
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="grid grid-cols-2 gap-4">
                  <div className="text-center">
                    <div className="text-2xl font-bold text-primary">$1,250</div>
                    <div className="text-sm text-muted-foreground">Current Bid</div>
                  </div>
                  <div className="text-center">
                    <div className="text-2xl font-bold text-primary">12</div>
                    <div className="text-sm text-muted-foreground">Bidders</div>
                  </div>
                  <div className="text-center">
                    <div className="text-2xl font-bold text-primary">47</div>
                    <div className="text-sm text-muted-foreground">Total Bids</div>
                  </div>
                  <div className="text-center">
                    <div className="text-2xl font-bold text-primary">247</div>
                    <div className="text-sm text-muted-foreground">Viewers</div>
                  </div>
                </div>
              </CardContent>
            </Card>

            {/* Broadcast Status */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Users className="w-5 h-5" />
                  Status
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="flex items-center justify-between">
                  <span className="text-sm text-muted-foreground">Broadcast</span>
                  <Badge variant={isBroadcasting ? "destructive" : "secondary"} className="gap-1">
                    <div className={`w-2 h-2 rounded-full ${isBroadcasting ? 'bg-white animate-pulse' : 'bg-current'}`}/>
                    {isBroadcasting ? 'LIVE' : 'OFFLINE'}
                  </Badge>
                </div>
                <div className="flex items-center justify-between">
                  <span className="text-sm text-muted-foreground">Current Item</span>
                  <span className="text-sm font-medium">Vintage Watch</span>
                </div>
                <div className="flex items-center justify-between">
                  <span className="text-sm text-muted-foreground">Time Remaining</span>
                  <span className="text-sm font-medium flex items-center gap-1">
                    <Clock className="w-3 h-3" />
                    5:32
                  </span>
                </div>
              </CardContent>
            </Card>
          </div>
        </div>

        {/* Quick Stats Bar */}
        <div className="mt-8 grid grid-cols-2 md:grid-cols-4 gap-4">
          <Card>
            <CardContent className="p-4 text-center">
              <div className="flex items-center justify-center gap-2 mb-2">
                <DollarSign className="w-5 h-5 text-primary" />
                <div className="text-2xl font-bold text-primary">$1,250</div>
              </div>
              <div className="text-sm text-muted-foreground">Current Bid</div>
            </CardContent>
          </Card>
          <Card>
            <CardContent className="p-4 text-center">
              <div className="flex items-center justify-center gap-2 mb-2">
                <Users className="w-5 h-5 text-primary" />
                <div className="text-2xl font-bold text-primary">12</div>
              </div>
              <div className="text-sm text-muted-foreground">Active Bidders</div>
            </CardContent>
          </Card>
          <Card>
            <CardContent className="p-4 text-center">
              <div className="flex items-center justify-center gap-2 mb-2">
                <TrendingUp className="w-5 h-5 text-primary" />
                <div className="text-2xl font-bold text-primary">47</div>
              </div>
              <div className="text-sm text-muted-foreground">Total Bids</div>
            </CardContent>
          </Card>
          <Card>
            <CardContent className="p-4 text-center">
              <div className="flex items-center justify-center gap-2 mb-2">
                <Clock className="w-5 h-5 text-primary" />
                <div className="text-2xl font-bold text-primary">5:32</div>
              </div>
              <div className="text-sm text-muted-foreground">Time Remaining</div>
            </CardContent>
          </Card>
        </div>
      </main>
    </div>
  )
}
'use client'

import { useState } from 'react'
import LiveKitBroadcaster from '../components/LiveKitBroadcaster'

export default function Home() {
  const [auctionId, setAuctionId] = useState('demo-auction-123')
  const [isBroadcasting, setIsBroadcasting] = useState(false)

  return (
    <div className="min-h-screen bg-gray-950 text-white">
      <div className="container mx-auto px-4 py-8">
        <header className="mb-8">
          <h1 className="text-4xl font-bold mb-2">Blytz Seller Dashboard</h1>
          <p className="text-gray-400">Broadcaster Platform - Host live auctions</p>
        </header>

        <main>
          {/* Auction Setup */}
          <div className="mb-6">
            <label htmlFor="auction-id" className="block text-sm font-medium mb-2">
              Auction Room ID:
            </label>
            <input
              id="auction-id"
              type="text"
              value={auctionId}
              onChange={(e) => setAuctionId(e.target.value)}
              className="bg-gray-800 border border-gray-700 rounded-lg px-4 py-2 w-full max-w-md focus:outline-none focus:ring-2 focus:ring-red-500"
              placeholder="Enter auction room ID"
            />
          </div>

          {/* Live Stream Broadcaster */}
          <div className="mb-8">
            <h2 className="text-2xl font-semibold mb-4">Broadcast Studio</h2>
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
          </div>

          {/* Auction Controls */}
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div className="bg-gray-900 rounded-lg p-6">
              <h3 className="text-xl font-semibold mb-4">Auction Controls</h3>
              <div className="space-y-3">
                <button 
                  className={`w-full py-2 px-4 rounded-lg transition-colors ${
                    isBroadcasting 
                      ? 'bg-green-600 hover:bg-green-700' 
                      : 'bg-gray-700 hover:bg-gray-600'
                  } text-white`}
                  disabled={!isBroadcasting}
                >
                  Start Bidding
                </button>
                <button 
                  className={`w-full py-2 px-4 rounded-lg transition-colors ${
                    isBroadcasting 
                      ? 'bg-yellow-600 hover:bg-yellow-700' 
                      : 'bg-gray-700 hover:bg-gray-600'
                  } text-white`}
                  disabled={!isBroadcasting}
                >
                  Pause Auction
                </button>
                <button 
                  className={`w-full py-2 px-4 rounded-lg transition-colors ${
                    isBroadcasting 
                      ? 'bg-red-600 hover:bg-red-700' 
                      : 'bg-gray-700 hover:bg-gray-600'
                  } text-white`}
                  disabled={!isBroadcasting}
                >
                  End Auction
                </button>
              </div>
            </div>

            <div className="bg-gray-900 rounded-lg p-6">
              <h3 className="text-xl font-semibold mb-4">Auction Status</h3>
              <div className="space-y-2">
                <div className="flex justify-between">
                  <span className="text-gray-400">Status:</span>
                  <span className={isBroadcasting ? 'text-green-500' : 'text-gray-500'}>
                    {isBroadcasting ? 'LIVE' : 'OFFLINE'}
                  </span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-400">Current Item:</span>
                  <span>Vintage Watch Collection</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-400">Current Bid:</span>
                  <span className="text-green-500 font-semibold">$1,250</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-400">Active Bidders:</span>
                  <span>12</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-400">Total Bids:</span>
                  <span>47</span>
                </div>
              </div>
            </div>
          </div>

          {/* Quick Stats */}
          <div className="mt-8 grid grid-cols-1 md:grid-cols-4 gap-4">
            <div className="bg-gray-900 rounded-lg p-4 text-center">
              <div className="text-2xl font-bold text-blue-500">$1,250</div>
              <div className="text-sm text-gray-400">Current Bid</div>
            </div>
            <div className="bg-gray-900 rounded-lg p-4 text-center">
              <div className="text-2xl font-bold text-green-500">12</div>
              <div className="text-sm text-gray-400">Active Bidders</div>
            </div>
            <div className="bg-gray-900 rounded-lg p-4 text-center">
              <div className="text-2xl font-bold text-yellow-500">47</div>
              <div className="text-sm text-gray-400">Total Bids</div>
            </div>
            <div className="bg-gray-900 rounded-lg p-4 text-center">
              <div className="text-2xl font-bold text-purple-500">5:32</div>
              <div className="text-sm text-gray-400">Time Remaining</div>
            </div>
          </div>
        </main>
      </div>
    </div>
  )
}
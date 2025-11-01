'use client'

import { useState } from 'react'
import LiveKitViewer from '../components/LiveKitViewer'

export default function Home() {
  const [auctionId, setAuctionId] = useState('demo-auction-123')

  return (
    <div className="min-h-screen bg-gray-950 text-white">
      <div className="container mx-auto px-4 py-8">
        <header className="mb-8">
          <h1 className="text-4xl font-bold mb-2">Blytz Live Auction Demo</h1>
          <p className="text-gray-400">Viewer Platform - Watch live auctions in real-time</p>
        </header>

        <main>
          {/* Auction Selection */}
          <div className="mb-6">
            <label htmlFor="auction-id" className="block text-sm font-medium mb-2">
              Auction ID:
            </label>
            <input
              id="auction-id"
              type="text"
              value={auctionId}
              onChange={(e) => setAuctionId(e.target.value)}
              className="bg-gray-800 border border-gray-700 rounded-lg px-4 py-2 w-full max-w-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              placeholder="Enter auction ID"
            />
          </div>

          {/* Live Stream Viewer */}
          <div className="mb-8">
            <h2 className="text-2xl font-semibold mb-4">Live Stream</h2>
            <LiveKitViewer 
              auctionId={auctionId}
              onConnected={() => console.log('Connected to live stream')}
              onDisconnected={() => console.log('Disconnected from live stream')}
              onError={(error) => console.error('Stream error:', error)}
            />
          </div>

          {/* Auction Information */}
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div className="bg-gray-900 rounded-lg p-6">
              <h3 className="text-xl font-semibold mb-4">Current Auction</h3>
              <div className="space-y-2">
                <div className="flex justify-between">
                  <span className="text-gray-400">Item:</span>
                  <span>Vintage Watch Collection</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-400">Current Bid:</span>
                  <span className="text-green-500 font-semibold">$1,250</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-400">Time Remaining:</span>
                  <span>5m 32s</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-400">Bidders:</span>
                  <span>12</span>
                </div>
              </div>
            </div>

            <div className="bg-gray-900 rounded-lg p-6">
              <h3 className="text-xl font-semibold mb-4">Quick Actions</h3>
              <div className="space-y-3">
                <button className="w-full bg-blue-600 hover:bg-blue-700 text-white py-2 px-4 rounded-lg transition-colors">
                  Place Bid
                </button>
                <button className="w-full bg-gray-700 hover:bg-gray-600 text-white py-2 px-4 rounded-lg transition-colors">
                  View Item Details
                </button>
                <button className="w-full bg-gray-700 hover:bg-gray-600 text-white py-2 px-4 rounded-lg transition-colors">
                  Auction History
                </button>
              </div>
            </div>
          </div>
        </main>
      </div>
    </div>
  )
}
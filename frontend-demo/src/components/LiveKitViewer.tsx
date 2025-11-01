'use client'

import { useState, useEffect } from 'react'
import { 
  LiveKitRoom, 
  VideoConference, 
  formatChatMessageLinks,
  useToken
} from '@livekit/components-react'
import '@livekit/components-styles'
import { Room, RoomConnectOptions } from 'livekit-client'

interface LiveKitViewerProps {
  auctionId: string
  serverUrl?: string
  token?: string
  onConnected?: () => void
  onDisconnected?: () => void
  onError?: (error: Error) => void
}

export default function LiveKitViewer({ 
  auctionId, 
  serverUrl = process.env.NEXT_PUBLIC_LIVEKIT_URL || 'wss://livekit.blytz.app',
  token,
  onConnected,
  onDisconnected,
  onError
}: LiveKitViewerProps) {
  const [isConnecting, setIsConnecting] = useState(false)
  const [isConnected, setIsConnected] = useState(false)
  const [connectionError, setConnectionError] = useState<string | null>(null)
  const [viewerToken, setViewerToken] = useState<string | null>(token || null)

  // Fetch token if not provided
  useEffect(() => {
    if (!token && auctionId) {
      fetchViewerToken()
    }
  }, [auctionId, token])

  const fetchViewerToken = async () => {
    try {
      setIsConnecting(true)
      setConnectionError(null)

      const response = await fetch(`/api/livekit/token?room=${auctionId}&role=viewer`)
      
      if (!response.ok) {
        throw new Error(`Failed to fetch token: ${response.statusText}`)
      }

      const data = await response.json()
      setViewerToken(data.token)
    } catch (error) {
      console.error('Error fetching viewer token:', error)
      setConnectionError(error instanceof Error ? error.message : 'Failed to connect')
      onError?.(error instanceof Error ? error : new Error('Connection failed'))
    } finally {
      setIsConnecting(false)
    }
  }

  const handleConnected = () => {
    setIsConnected(true)
    onConnected?.()
  }

  const handleDisconnected = () => {
    setIsConnected(false)
    onDisconnected?.()
  }

  const handleConnectionError = (error: Error) => {
    console.error('LiveKit connection error:', error)
    setConnectionError(error.message)
    onError?.(error)
  }

  const connectionOptions: RoomConnectOptions = {
    autoSubscribe: true,
  }

  if (isConnecting) {
    return (
      <div className="flex items-center justify-center min-h-[400px] bg-gray-900 rounded-lg">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-white mx-auto mb-4"></div>
          <p className="text-white">Connecting to live stream...</p>
        </div>
      </div>
    )
  }

  if (connectionError) {
    return (
      <div className="flex items-center justify-center min-h-[400px] bg-gray-900 rounded-lg">
        <div className="text-center">
          <div className="text-red-500 mb-4">
            <svg className="w-16 h-16 mx-auto" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </div>
          <p className="text-white mb-2">Connection Error</p>
          <p className="text-gray-400 text-sm mb-4">{connectionError}</p>
          <button
            onClick={fetchViewerToken}
            className="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg transition-colors"
          >
            Retry Connection
          </button>
        </div>
      </div>
    )
  }

  if (!viewerToken) {
    return (
      <div className="flex items-center justify-center min-h-[400px] bg-gray-900 rounded-lg">
        <div className="text-center">
          <p className="text-white">Unable to connect to stream</p>
        </div>
      </div>
    )
  }

  return (
    <div className="bg-gray-900 rounded-lg overflow-hidden">
      <LiveKitRoom
        token={viewerToken}
        serverUrl={serverUrl}
        connectOptions={connectionOptions}
        onConnected={handleConnected}
        onDisconnected={handleDisconnected}
        onError={handleConnectionError}
        className="h-[600px]"
      >
        <VideoConference 
          chatMessageFormatter={formatChatMessageLinks}
        />
      </LiveKitRoom>
    </div>
  )
}
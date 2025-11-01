'use client'

import { useState, useEffect, useCallback } from 'react'
import { 
  LiveKitRoom, 
  VideoConference, 
  formatChatMessageLinks,
  useToken
} from '@livekit/components-react'
import '@livekit/components-styles'
import { Room, RoomConnectOptions } from 'livekit-client'
import { Button } from '@/components/ui/button'

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
  const fetchViewerToken = useCallback(async () => {
    try {
      setIsConnecting(true)
      setConnectionError(null)

      const response = await fetch(`/api/v1/livekit/token?room=${auctionId}&role=viewer`)
      
      if (!response.ok) {
        throw new Error(`Failed to fetch token: ${response.statusText}`)
      }

      const data = await response.json()
      setViewerToken(data.token)
    } catch (error) {
      setConnectionError(error instanceof Error ? error.message : 'Failed to fetch token')
    } finally {
      setIsConnecting(false)
    }
  }, [auctionId])

  useEffect(() => {
    if (!token && auctionId) {
      fetchViewerToken()
    }
  }, [auctionId, token, fetchViewerToken])

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
      <div className="flex items-center justify-center min-h-[400px] bg-card rounded-3xl">
        <div className="text-center space-y-4">
          <div className="animate-spin rounded-full h-12 w-12 border-2 border-primary border-t-transparent mx-auto"></div>
          <div>
            <p className="text-foreground font-medium">Connecting to live stream...</p>
            <p className="text-muted-foreground text-sm">Please wait while we establish connection</p>
          </div>
        </div>
      </div>
    )
  }

  if (connectionError) {
    return (
      <div className="flex items-center justify-center min-h-[400px] bg-card rounded-3xl">
        <div className="text-center space-y-4 max-w-md mx-auto p-6">
          <div className="text-destructive mx-auto">
            <svg className="w-16 h-16" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </div>
          <div>
            <h3 className="text-foreground font-semibold text-lg mb-2">Connection Error</h3>
            <p className="text-muted-foreground text-sm mb-4">{connectionError}</p>
          </div>
          <Button
            onClick={fetchViewerToken}
            className="gap-2"
          >
            Retry Connection
          </Button>
        </div>
      </div>
    )
  }

  if (!viewerToken) {
    return (
      <div className="flex items-center justify-center min-h-[400px] bg-card rounded-3xl">
        <div className="text-center space-y-4">
          <div className="text-muted-foreground">
            <svg className="w-16 h-16 mx-auto" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 10l4.553-2.276A1 1 0 0121 8.618v6.764a1 1 0 01-1.447.894L15 14M5 18h8a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z" />
            </svg>
          </div>
          <div>
            <p className="text-foreground font-medium">Stream Unavailable</p>
            <p className="text-muted-foreground text-sm">Unable to connect to the live stream</p>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="bg-card rounded-3xl overflow-hidden border">
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
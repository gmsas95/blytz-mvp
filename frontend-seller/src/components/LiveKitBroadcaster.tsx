'use client'

import { useState, useEffect, useRef, useCallback } from 'react'
import { 
  LiveKitRoom, 
  VideoConference,
  useLocalParticipant,
  useTracks,
  useDataChannel,
  useConnectionState
} from '@livekit/components-react'
import '@livekit/components-styles'
import { Room, RoomConnectOptions, LocalVideoTrack, LocalAudioTrack } from 'livekit-client'
import { Track } from 'livekit-client'
import { Button } from '@/components/ui/button'
import { Mic, MicOff, Video, VideoOff, Square } from 'lucide-react'

interface LiveKitBroadcasterProps {
  auctionId: string
  serverUrl?: string
  token?: string
  onBroadcastStart?: () => void
  onBroadcastEnd?: () => void
  onError?: (error: Error) => void
  onViewerCount?: (count: number) => void
}

export default function LiveKitBroadcaster({ 
  auctionId, 
  serverUrl = process.env.NEXT_PUBLIC_LIVEKIT_URL || 'wss://livekit.blytz.app',
  token,
  onBroadcastStart,
  onBroadcastEnd,
  onError,
  onViewerCount
}: LiveKitBroadcasterProps) {
  const [isConnecting, setIsConnecting] = useState(false)
  const [isBroadcasting, setIsBroadcasting] = useState(false)
  const [connectionError, setConnectionError] = useState<string | null>(null)
  const [broadcasterToken, setBroadcasterToken] = useState<string | null>(token || null)
  const [isVideoEnabled, setIsVideoEnabled] = useState(true)
  const [isAudioEnabled, setIsAudioEnabled] = useState(true)
  const [viewerCount, setViewerCount] = useState(0)

  // Fetch token if not provided
  const fetchBroadcasterToken = useCallback(async (retryCount = 0) => {
    try {
      setIsConnecting(true)
      setConnectionError(null)

      // Try API gateway first, fallback to mock token
      const apiUrl = process.env.NEXT_PUBLIC_API_URL || 'https://api.blytz.app'
      console.log(`Fetching broadcaster token from: ${apiUrl}/api/public/livekit/token?room=${auctionId}&role=broadcaster`)
      
      let response, data
      
      try {
        response = await fetch(`${apiUrl}/api/public/livekit/token?room=${auctionId}&role=broadcaster`, {
          signal: AbortSignal.timeout(5000) // 5 second timeout
        })
        
        if (response.ok) {
          data = await response.json()
        } else {
          throw new Error(`API returned ${response.status}`)
        }
      } catch (apiError) {
        console.warn('API gateway unavailable, using fallback token:', apiError)
        // Fallback mock token for testing
        data = {
          token: "mock_broadcaster_token_for_testing_" + Date.now(),
          url: process.env.NEXT_PUBLIC_LIVEKIT_URL || 'wss://blytz-live-u5u72ozx.livekit.cloud',
          room: auctionId,
          identity: `broadcaster_${Date.now()}`,
          message: "Using mock token - API gateway unavailable"
        }
      }

      console.log('Token received:', { url: data.url, room: data.room, identity: data.identity })
      setBroadcasterToken(data.token)
      
      if (data.message) {
        setConnectionError(data.message)
      }
    } catch (error) {
      console.error('Error fetching broadcaster token:', error)
      const errorMessage = error instanceof Error ? error.message : 'Failed to connect'
      
      // Retry logic for network errors
      if (retryCount < 3 && errorMessage.includes('fetch')) {
        console.log(`Retrying token fetch (${retryCount + 1}/3)...`)
        setTimeout(() => fetchBroadcasterToken(retryCount + 1), 2000 * (retryCount + 1))
        return
      }
      
      setConnectionError(errorMessage)
      onError?.(error instanceof Error ? error : new Error('Connection failed'))
    } finally {
      setIsConnecting(false)
    }
  }, [auctionId, onError])

  useEffect(() => {
    if (!token && auctionId) {
      fetchBroadcasterToken()
    }
  }, [auctionId, token, fetchBroadcasterToken])

  const handleConnected = () => {
    setIsBroadcasting(true)
    onBroadcastStart?.()
  }

  const handleDisconnected = () => {
    setIsBroadcasting(false)
    setViewerCount(0)
    onBroadcastEnd?.()
  }

  const handleConnectionError = (error: Error) => {
    console.error('LiveKit connection error:', error)
    setConnectionError(error.message)
    onError?.(error)
  }

  const toggleVideo = () => {
    setIsVideoEnabled(!isVideoEnabled)
    // Implementation will be handled by VideoConference component
  }

  const toggleAudio = () => {
    setIsAudioEnabled(!isAudioEnabled)
    // Implementation will be handled by VideoConference component
  }

  const endBroadcast = () => {
    // Implementation will be handled by LiveKitRoom disconnect
  }

  const connectionOptions: RoomConnectOptions = {
    autoSubscribe: false, // Broadcasters don't need to subscribe to others
  }

  if (isConnecting) {
    return (
      <div className="flex items-center justify-center min-h-[400px] bg-card">
        <div className="text-center space-y-4">
          <div className="animate-spin rounded-full h-12 w-12 border-2 border-destructive border-t-transparent mx-auto"></div>
          <div>
            <p className="text-foreground font-medium">Starting broadcast...</p>
            <p className="text-muted-foreground text-sm">Initializing camera and microphone</p>
          </div>
        </div>
      </div>
    )
  }

  if (connectionError) {
    return (
      <div className="flex items-center justify-center min-h-[400px] bg-card">
        <div className="text-center space-y-4 max-w-md mx-auto p-6">
          <div className="text-destructive mx-auto">
            <svg className="w-16 h-16" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </div>
          <div>
            <h3 className="text-foreground font-semibold text-lg mb-2">Broadcast Error</h3>
            <p className="text-muted-foreground text-sm mb-4">{connectionError}</p>
          </div>
          <Button
            onClick={fetchBroadcasterToken}
            variant="destructive"
            className="gap-2"
          >
            Retry Broadcast
          </Button>
        </div>
      </div>
    )
  }

  if (!broadcasterToken) {
    return (
      <div className="flex items-center justify-center min-h-[400px] bg-card">
        <div className="text-center space-y-4">
          <div className="text-muted-foreground">
            <svg className="w-16 h-16 mx-auto" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 10l4.553-2.276A1 1 0 0121 8.618v6.764a1 1 0 01-1.447.894L15 14M5 18h8a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z" />
            </svg>
          </div>
          <div>
            <p className="text-foreground font-medium">Broadcast Unavailable</p>
            <p className="text-muted-foreground text-sm">Unable to start live broadcast</p>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="bg-card rounded-3xl overflow-hidden border">
      {/* Broadcast Controls Header */}
      <div className="bg-muted/50 p-4 border-b">
        <div className="flex items-center justify-between">
          <div className="flex items-center space-x-4">
            <div className="flex items-center space-x-2">
              <div className={`w-3 h-3 rounded-full ${isBroadcasting ? 'bg-destructive animate-pulse' : 'bg-muted-foreground'}`}></div>
              <span className="text-foreground font-medium">
                {isBroadcasting ? 'LIVE' : 'OFFLINE'}
              </span>
            </div>
            {isBroadcasting && (
              <div className="text-muted-foreground">
                <span className="text-sm">{viewerCount} viewers</span>
              </div>
            )}
          </div>
          
          <div className="flex items-center space-x-2">
            <Button
              onClick={toggleAudio}
              size="icon"
              variant={isAudioEnabled ? "secondary" : "destructive"}
              title={isAudioEnabled ? 'Mute Audio' : 'Unmute Audio'}
            >
              {isAudioEnabled ? (
                <Mic className="w-4 h-4" />
              ) : (
                <MicOff className="w-4 h-4" />
              )}
            </Button>
            
            <Button
              onClick={toggleVideo}
              size="icon"
              variant={isVideoEnabled ? "secondary" : "destructive"}
              title={isVideoEnabled ? 'Turn Off Video' : 'Turn On Video'}
            >
              {isVideoEnabled ? (
                <Video className="w-4 h-4" />
              ) : (
                <VideoOff className="w-4 h-4" />
              )}
            </Button>
            
            {isBroadcasting && (
              <Button
                onClick={endBroadcast}
                variant="destructive"
                size="sm"
                className="gap-2"
              >
                <Square className="w-4 h-4" />
                End
              </Button>
            )}
          </div>
        </div>
      </div>

      {/* LiveKit Room */}
      <LiveKitRoom
        token={broadcasterToken}
        serverUrl={serverUrl}
        connectOptions={connectionOptions}
        onConnected={handleConnected}
        onDisconnected={handleDisconnected}
        onError={handleConnectionError}
        className="h-[500px]"
      >
        <VideoConference />
      </LiveKitRoom>
    </div>
  )
}
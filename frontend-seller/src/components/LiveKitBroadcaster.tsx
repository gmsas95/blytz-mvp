'use client'

import { useState, useEffect, useRef } from 'react'
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
  useEffect(() => {
    if (!token && auctionId) {
      fetchBroadcasterToken()
    }
  }, [auctionId, token])

  const fetchBroadcasterToken = async () => {
    try {
      setIsConnecting(true)
      setConnectionError(null)

      const response = await fetch(`/api/livekit/token?room=${auctionId}&role=broadcaster`)
      
      if (!response.ok) {
        throw new Error(`Failed to fetch token: ${response.statusText}`)
      }

      const data = await response.json()
      setBroadcasterToken(data.token)
    } catch (error) {
      console.error('Error fetching broadcaster token:', error)
      setConnectionError(error instanceof Error ? error.message : 'Failed to connect')
      onError?.(error instanceof Error ? error : new Error('Connection failed'))
    } finally {
      setIsConnecting(false)
    }
  }

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
      <div className="flex items-center justify-center min-h-[400px] bg-gray-900 rounded-lg">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-red-500 mx-auto mb-4"></div>
          <p className="text-white">Starting broadcast...</p>
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
          <p className="text-white mb-2">Broadcast Error</p>
          <p className="text-gray-400 text-sm mb-4">{connectionError}</p>
          <button
            onClick={fetchBroadcasterToken}
            className="bg-red-600 hover:bg-red-700 text-white px-4 py-2 rounded-lg transition-colors"
          >
            Retry Broadcast
          </button>
        </div>
      </div>
    )
  }

  if (!broadcasterToken) {
    return (
      <div className="flex items-center justify-center min-h-[400px] bg-gray-900 rounded-lg">
        <div className="text-center">
          <p className="text-white">Unable to start broadcast</p>
        </div>
      </div>
    )
  }

  return (
    <div className="bg-gray-900 rounded-lg overflow-hidden">
      {/* Broadcast Controls Header */}
      <div className="bg-gray-800 p-4 border-b border-gray-700">
        <div className="flex items-center justify-between">
          <div className="flex items-center space-x-4">
            <div className="flex items-center space-x-2">
              <div className={`w-3 h-3 rounded-full ${isBroadcasting ? 'bg-red-500 animate-pulse' : 'bg-gray-500'}`}></div>
              <span className="text-white font-medium">
                {isBroadcasting ? 'LIVE' : 'OFFLINE'}
              </span>
            </div>
            {isBroadcasting && (
              <div className="text-gray-400">
                <span className="text-sm">{viewerCount} viewers</span>
              </div>
            )}
          </div>
          
          <div className="flex items-center space-x-2">
            <button
              onClick={toggleAudio}
              className={`p-2 rounded-lg transition-colors ${
                isAudioEnabled 
                  ? 'bg-gray-700 text-white hover:bg-gray-600' 
                  : 'bg-red-600 text-white hover:bg-red-700'
              }`}
              title={isAudioEnabled ? 'Mute Audio' : 'Unmute Audio'}
            >
              {isAudioEnabled ? (
                <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 11a7 7 0 01-7 7m0 0a7 7 0 01-7-7m7 7v4m0 0H8m4 0h4m-4-8a3 3 0 01-3-3V5a3 3 0 116 0v6a3 3 0 01-3 3z" />
                </svg>
              ) : (
                <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5.586 15H4a1 1 0 01-1-1v-4a1 1 0 011-1h1.586l4.707-4.707C10.923 3.663 12 4.109 12 5v14c0 .891-1.077 1.337-1.707.707L5.586 15z" />
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2" />
                </svg>
              )}
            </button>
            
            <button
              onClick={toggleVideo}
              className={`p-2 rounded-lg transition-colors ${
                isVideoEnabled 
                  ? 'bg-gray-700 text-white hover:bg-gray-600' 
                  : 'bg-red-600 text-white hover:bg-red-700'
              }`}
              title={isVideoEnabled ? 'Turn Off Video' : 'Turn On Video'}
            >
              {isVideoEnabled ? (
                <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 10l4.553-2.276A1 1 0 0121 8.618v6.764a1 1 0 01-1.447.894L15 14M5 18h8a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z" />
                </svg>
              ) : (
                <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M18.364 18.364A9 9 0 005.636 5.636m12.728 12.728A9 9 0 015.636 5.636m12.728 12.728L5.636 5.636" />
                </svg>
              )}
            </button>
            
            {isBroadcasting && (
              <button
                onClick={endBroadcast}
                className="bg-red-600 hover:bg-red-700 text-white px-4 py-2 rounded-lg transition-colors"
              >
                End Broadcast
              </button>
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
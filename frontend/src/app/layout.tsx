import type { Metadata } from 'next'
import { Inter } from 'next/font/google'
import '@/styles/globals.css'
import { cn } from '@/lib/utils'

const inter = Inter({
  subsets: ['latin'],
  variable: '--font-inter',
})

export const metadata: Metadata = {
  title: 'Blytz - Live Auction Commerce',
  description: 'Discover amazing products through live auctions and streaming',
  keywords: 'auction, livestream, ecommerce, bidding, deals',
  authors: [{ name: 'Blytz' }],
}

export const viewport = {
  width: 'device-width',
  initialScale: 1,
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="en" suppressHydrationWarning>
      <body className={cn(
        'min-h-screen bg-background font-body antialiased',
        inter.variable
      )}>
        {children}
      </body>
    </html>
  )
}
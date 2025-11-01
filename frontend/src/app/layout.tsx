import type { Metadata } from 'next';
import { Inter } from 'next/font/google';

import '@/styles/globals.css';
import { AccessibilityInit } from '@/components/accessibility/accessibility-init';
import { Footer } from '@/components/layout/footer';
import { Header } from '@/components/layout/header';
import { AuthProvider } from '@/contexts/auth-context';
import { cn } from '@/lib/utils';

const inter = Inter({
  subsets: ['latin'],
  variable: '--font-inter',
});

export const metadata: Metadata = {
  title: 'Blytz - Live Auction Commerce',
  description: 'Discover amazing products through live auctions and streaming',
  keywords: 'auction, livestream, ecommerce, bidding, deals',
  authors: [{ name: 'Blytz' }],
};

export const viewport = {
  width: 'device-width',
  initialScale: 1,
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en" suppressHydrationWarning>
      <body className={cn('min-h-screen bg-background font-body antialiased', inter.variable)}>
        {/* Skip to main content link for screen readers */}
        <a
          href="#main-content"
          className="sr-only focus:not-sr-only focus:absolute focus:top-4 focus:left-4 bg-primary text-primary-foreground px-4 py-2 rounded-md z-50 focus:outline-none focus:ring-2 focus:ring-ring"
        >
          Skip to main content
        </a>

        <AuthProvider>
          <AccessibilityInit />
          <div className="flex min-h-screen flex-col">
            <Header />
            <main id="main-content" className="flex-1" role="main">
              {children}
            </main>
            <Footer />
          </div>
        </AuthProvider>
      </body>
    </html>
  );
}

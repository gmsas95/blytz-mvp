'use client';

import { Menu, ShoppingCart, User, Search, X, LogOut, Store } from 'lucide-react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { useState } from 'react';

import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Button } from '@/components/ui/button';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { Input } from '@/components/ui/input';
import { useAuth } from '@/contexts/auth-context';

export function Header() {
  const [isMenuOpen, setIsMenuOpen] = useState(false);
  const [searchQuery, setSearchQuery] = useState('');
  const { user, isAuthenticated, logout, isLoading } = useAuth();
  const router = useRouter();

  const navItems = [
    { href: '/products', label: 'Products' },
    { href: '/livestream', label: 'Live Streams' },
    { href: '/auctions', label: 'Auctions' },
    { href: '/sellers', label: 'Sellers' },
  ];

  const handleLogout = async () => {
    await logout();
    router.push('/');
    router.refresh();
  };

  const getUserInitials = (name?: string, email?: string) => {
    if (name) {
      return name
        .split(' ')
        .map((n) => n[0])
        .join('')
        .toUpperCase()
        .slice(0, 2);
    }
    return email?.slice(0, 2).toUpperCase() || 'U';
  };

  return (
    <header
      className="sticky top-0 z-50 w-full border-b border-border/40 bg-background/95 backdrop-blur"
      role="banner"
    >
      <div className="container mx-auto px-4">
        <div className="flex items-center justify-between h-16">
          {/* Logo */}
          <Link
            href="/"
            className="flex items-center space-x-2 focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2 rounded-lg p-1"
            aria-label="Blytz homepage"
          >
            <div
              className="w-8 h-8 rounded-xl bg-primary flex items-center justify-center"
              aria-hidden="true"
            >
              <span className="text-white font-bold text-sm">B</span>
            </div>
            <span className="font-bold text-xl font-sans">Blytz</span>
          </Link>

          {/* Desktop Navigation */}
          <nav className="hidden md:flex items-center space-x-6" aria-label="Main navigation">
            {navItems.map((item) => (
              <Link
                key={item.href}
                href={item.href}
                className="text-sm font-medium text-foreground hover:text-primary transition-colors focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2 rounded-lg px-2 py-1"
              >
                {item.label}
              </Link>
            ))}
          </nav>

          {/* Search Bar */}
          <div className="hidden lg:flex flex-1 max-w-md mx-8" role="search">
            <div className="relative w-full">
              <Search
                className="absolute left-3 top-1/2 w-4 h-4 -translate-y-1/2 text-muted-foreground"
                aria-hidden="true"
              />
              <Input
                type="search"
                placeholder="Search products, auctions, sellers..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="pl-10 w-full"
                aria-label="Search products, auctions, and sellers"
                aria-describedby="search-description"
              />
              <span id="search-description" className="sr-only">
                Use this search bar to find products, live auctions, and sellers
              </span>
            </div>
          </div>

          {/* Actions */}
          <div className="flex items-center space-x-2" role="toolbar" aria-label="User actions">
            <Button
              variant="ghost"
              size="icon"
              className="lg:hidden"
              aria-label="Toggle search"
              onClick={() => {
                // Focus mobile search input
                const mobileSearch = document.getElementById(
                  'mobile-search-input'
                ) as HTMLInputElement;
                mobileSearch?.focus();
              }}
            >
              <Search className="w-5 h-5" />
            </Button>

            <Link href="/cart" aria-label="Shopping cart">
              <Button variant="ghost" size="icon" className="relative">
                <ShoppingCart className="w-5 h-5" />
                {/* TODO: Add cart item count from state */}
              </Button>
            </Link>

            {/* User Actions */}
            {isLoading ? (
              <div className="w-9 h-9 rounded-full bg-muted animate-pulse" />
            ) : isAuthenticated && user ? (
              <DropdownMenu>
                <DropdownMenuTrigger asChild>
                  <Button
                    variant="ghost"
                    className="relative h-9 w-9 rounded-full"
                    aria-label="User menu"
                  >
                    <Avatar className="h-9 w-9">
                      <AvatarImage src={user.avatar} alt={user.name} />
                      <AvatarFallback>{getUserInitials(user.name, user.email)}</AvatarFallback>
                    </Avatar>
                  </Button>
                </DropdownMenuTrigger>
                <DropdownMenuContent className="w-56" align="end" forceMount>
                  <DropdownMenuLabel className="font-normal">
                    <div className="flex flex-col space-y-1">
                      <p className="text-sm font-medium leading-none">{user.name}</p>
                      <p className="text-xs leading-none text-muted-foreground">{user.email}</p>
                    </div>
                  </DropdownMenuLabel>
                  <DropdownMenuSeparator />
                  <DropdownMenuItem asChild>
                    <Link href="/profile" className="w-full cursor-pointer">
                      <User className="mr-2 h-4 w-4" />
                      <span>Profile</span>
                    </Link>
                  </DropdownMenuItem>
                  {user.isSeller && (
                    <DropdownMenuItem asChild>
                      <Link href="/seller" className="w-full cursor-pointer">
                        <Store className="mr-2 h-4 w-4" />
                        <span>My Store</span>
                      </Link>
                    </DropdownMenuItem>
                  )}
                  <DropdownMenuItem asChild>
                    <Link href="/orders" className="w-full cursor-pointer">
                      <ShoppingCart className="mr-2 h-4 w-4" />
                      <span>Orders</span>
                    </Link>
                  </DropdownMenuItem>
                  <DropdownMenuSeparator />
                  <DropdownMenuItem className="w-full cursor-pointer" onClick={handleLogout}>
                    <LogOut className="mr-2 h-4 w-4" />
                    <span>Log out</span>
                  </DropdownMenuItem>
                </DropdownMenuContent>
              </DropdownMenu>
            ) : (
              <Link href="/auth" aria-label="Sign in to your account">
                <Button variant="ghost" size="icon">
                  <User className="w-5 h-5" />
                </Button>
              </Link>
            )}

            <Button
              variant="ghost"
              size="icon"
              className="md:hidden"
              onClick={() => setIsMenuOpen(!isMenuOpen)}
              aria-label={isMenuOpen ? 'Close navigation menu' : 'Open navigation menu'}
              aria-expanded={isMenuOpen}
              aria-controls="mobile-menu"
            >
              {isMenuOpen ? <X className="w-5 h-5" /> : <Menu className="w-5 h-5" />}
            </Button>
          </div>
        </div>

        {/* Mobile Menu */}
        {isMenuOpen && (
          <div
            id="mobile-menu"
            className="md:hidden border-t border-border/40"
            role="navigation"
            aria-label="Mobile navigation"
          >
            <div className="container mx-auto px-4 py-4 space-y-4">
              <div className="relative" role="search">
                <Search
                  className="absolute left-3 top-1/2 w-4 h-4 -translate-y-1/2 text-muted-foreground"
                  aria-hidden="true"
                />
                <Input
                  id="mobile-search-input"
                  type="search"
                  placeholder="Search..."
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                  className="pl-10 w-full"
                  aria-label="Search products, auctions, and sellers"
                />
              </div>

              <nav className="flex flex-col space-y-2">
                {navItems.map((item) => (
                  <Link
                    key={item.href}
                    href={item.href}
                    className="text-sm font-medium py-2 text-foreground hover:text-primary transition-colors focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2 rounded-lg px-2 py-1"
                    onClick={() => setIsMenuOpen(false)}
                  >
                    {item.label}
                  </Link>
                ))}
              </nav>

              {/* Mobile User Actions */}
              {isAuthenticated && user && (
                <div className="border-t border-border/40 pt-4">
                  <div className="flex items-center space-x-3 mb-4">
                    <Avatar className="h-8 w-8">
                      <AvatarImage src={user.avatar} alt={user.name} />
                      <AvatarFallback>{getUserInitials(user.name, user.email)}</AvatarFallback>
                    </Avatar>
                    <div>
                      <p className="text-sm font-medium">{user.name}</p>
                      <p className="text-xs text-muted-foreground">{user.email}</p>
                    </div>
                  </div>
                  <nav className="flex flex-col space-y-2">
                    <Link
                      href="/profile"
                      className="text-sm py-2 text-foreground hover:text-primary transition-colors px-2"
                      onClick={() => setIsMenuOpen(false)}
                    >
                      Profile
                    </Link>
                    {user.isSeller && (
                      <Link
                        href="/seller"
                        className="text-sm py-2 text-foreground hover:text-primary transition-colors px-2"
                        onClick={() => setIsMenuOpen(false)}
                      >
                        My Store
                      </Link>
                    )}
                    <Link
                      href="/orders"
                      className="text-sm py-2 text-foreground hover:text-primary transition-colors px-2"
                      onClick={() => setIsMenuOpen(false)}
                    >
                      Orders
                    </Link>
                    <button
                      onClick={() => {
                        handleLogout();
                        setIsMenuOpen(false);
                      }}
                      className="text-sm py-2 text-left text-foreground hover:text-primary transition-colors px-2"
                    >
                      Log out
                    </button>
                  </nav>
                </div>
              )}
            </div>
          </div>
        )}
      </div>
    </header>
  );
}

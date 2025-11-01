'use client';

import { User, Mail, Calendar, DollarSign, Star } from 'lucide-react';
import { useEffect } from 'react';

import { ProtectedRoute } from '@/components/auth/protected-route';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Separator } from '@/components/ui/separator';
import { useAuth } from '@/contexts/auth-context';

export default function ProfilePage() {
  const { user, isLoading, refreshUser } = useAuth();

  useEffect(() => {
    // Refresh user data when component mounts
    if (user) {
      refreshUser();
    }
  }, [user, refreshUser]);

  if (isLoading) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="flex items-center justify-center min-h-[50vh]">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
        </div>
      </div>
    );
  }

  return (
    <ProtectedRoute>
      <div className="container mx-auto px-4 py-8">
        <div className="max-w-4xl mx-auto space-y-6">
          {/* Header */}
          <div className="text-center">
            <h1 className="text-3xl font-bold tracking-tight">My Profile</h1>
            <p className="text-muted-foreground mt-2">Manage your account and view your activity</p>
          </div>

          {/* Profile Overview Card */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <User className="h-5 w-5" />
                Profile Information
              </CardTitle>
              <CardDescription>Your account details and statistics</CardDescription>
            </CardHeader>
            <CardContent className="space-y-6">
              {/* User Info */}
              <div className="flex items-center gap-4">
                <div className="w-16 h-16 rounded-full bg-primary/10 flex items-center justify-center">
                  <span className="text-2xl font-bold text-primary">
                    {user?.name?.slice(0, 2).toUpperCase() ||
                      user?.email?.slice(0, 2).toUpperCase()}
                  </span>
                </div>
                <div className="space-y-1">
                  <h2 className="text-xl font-semibold">{user?.name}</h2>
                  <div className="flex items-center gap-2 text-sm text-muted-foreground">
                    <Mail className="h-4 w-4" />
                    {user?.email}
                  </div>
                  <div className="flex items-center gap-2">
                    {user?.isSeller && <Badge variant="secondary">Seller</Badge>}
                    <div className="flex items-center gap-1 text-sm text-muted-foreground">
                      <Calendar className="h-4 w-4" />
                      Joined {new Date(user?.createdAt || '').toLocaleDateString()}
                    </div>
                  </div>
                </div>
              </div>

              <Separator />

              {/* Stats */}
              <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                <div className="text-center">
                  <div className="text-2xl font-bold">{user?.rating.toFixed(1) || '0.0'}</div>
                  <div className="text-sm text-muted-foreground flex items-center justify-center gap-1">
                    <Star className="h-4 w-4 fill-yellow-400 text-yellow-400" />
                    Rating
                  </div>
                </div>
                <div className="text-center">
                  <div className="text-2xl font-bold">{user?.totalSales || 0}</div>
                  <div className="text-sm text-muted-foreground flex items-center justify-center gap-1">
                    <DollarSign className="h-4 w-4" />
                    Total Sales
                  </div>
                </div>
                <div className="text-center">
                  <div className="text-2xl font-bold">--</div>
                  <div className="text-sm text-muted-foreground">Items Listed</div>
                </div>
                <div className="text-center">
                  <div className="text-2xl font-bold">--</div>
                  <div className="text-sm text-muted-foreground">Bids Placed</div>
                </div>
              </div>

              <Separator />

              {/* Store Information (if seller) */}
              {user?.isSeller && (
                <div className="space-y-3">
                  <h3 className="text-lg font-semibold">Store Information</h3>
                  <div className="grid gap-4 md:grid-cols-2">
                    <div>
                      <label className="text-sm font-medium text-muted-foreground">
                        Store Name
                      </label>
                      <p className="font-medium">{user.storeName || 'Not set'}</p>
                    </div>
                    <div>
                      <label className="text-sm font-medium text-muted-foreground">
                        Store Description
                      </label>
                      <p className="text-sm">{user.storeDescription || 'Not set'}</p>
                    </div>
                  </div>
                  <Button variant="outline">Edit Store Settings</Button>
                </div>
              )}
            </CardContent>
          </Card>

          {/* Quick Actions */}
          <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
            <Card className="cursor-pointer hover:shadow-md transition-shadow">
              <CardHeader>
                <CardTitle className="text-lg">My Bids</CardTitle>
                <CardDescription>View your active and past bids</CardDescription>
              </CardHeader>
              <CardContent>
                <Button variant="outline" className="w-full">
                  View Bids
                </Button>
              </CardContent>
            </Card>

            <Card className="cursor-pointer hover:shadow-md transition-shadow">
              <CardHeader>
                <CardTitle className="text-lg">Order History</CardTitle>
                <CardDescription>Track your purchases and sales</CardDescription>
              </CardHeader>
              <CardContent>
                <Button variant="outline" className="w-full">
                  View Orders
                </Button>
              </CardContent>
            </Card>

            {user?.isSeller && (
              <Card className="cursor-pointer hover:shadow-md transition-shadow">
                <CardHeader>
                  <CardTitle className="text-lg">My Listings</CardTitle>
                  <CardDescription>Manage your products and auctions</CardDescription>
                </CardHeader>
                <CardContent>
                  <Button variant="outline" className="w-full">
                    Manage Listings
                  </Button>
                </CardContent>
              </Card>
            )}
          </div>
        </div>
      </div>
    </ProtectedRoute>
  );
}

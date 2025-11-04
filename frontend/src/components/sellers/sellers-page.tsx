'use client';

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';

export function SellersPage() {
  return (
    <div className="container mx-auto px-4 py-8">
      <div className="text-center mb-12">
        <h1 className="text-4xl font-bold mb-4">Become a Blytz Seller</h1>
        <p className="text-xl text-muted-foreground max-w-2xl mx-auto">
          Start your live auction business and reach thousands of eager buyers
        </p>
      </div>

      <div className="grid md:grid-cols-3 gap-8 mb-12">
        <Card>
          <CardHeader>
            <CardTitle>Start Selling</CardTitle>
            <CardDescription>
              Create your seller account and start listing products in minutes
            </CardDescription>
          </CardHeader>
          <CardContent>
            <Button className="w-full">Get Started</Button>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Go Live</CardTitle>
            <CardDescription>
              Host live auctions and interact with buyers in real-time
            </CardDescription>
          </CardHeader>
          <CardContent>
            <Button variant="outline" className="w-full">Learn More</Button>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Get Paid</CardTitle>
            <CardDescription>
              Secure payments and fast payouts to your bank account
            </CardDescription>
          </CardHeader>
          <CardContent>
            <Button variant="outline" className="w-full">View Pricing</Button>
          </CardContent>
        </Card>
      </div>

      <div className="text-center">
        <h2 className="text-2xl font-semibold mb-4">Ready to start selling?</h2>
        <Button size="lg">Create Seller Account</Button>
      </div>
    </div>
  );
}
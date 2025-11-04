'use client';

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';

export function HelpPage() {
  return (
    <div className="container mx-auto px-4 py-8">
      <div className="text-center mb-12">
        <h1 className="text-4xl font-bold mb-4">Help Center</h1>
        <p className="text-xl text-muted-foreground max-w-2xl mx-auto">
          Find answers to common questions and get support
        </p>
      </div>

      <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-6 mb-12">
        <Card>
          <CardHeader>
            <CardTitle>Getting Started</CardTitle>
            <CardDescription>
              Learn how to use Blytz and place your first bid
            </CardDescription>
          </CardHeader>
          <CardContent>
            <Button variant="outline" className="w-full">View Guide</Button>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Account & Billing</CardTitle>
            <CardDescription>
              Manage your account, payment methods, and billing
            </CardDescription>
          </CardHeader>
          <CardContent>
            <Button variant="outline" className="w-full">Account Help</Button>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Bidding & Buying</CardTitle>
            <CardDescription>
              Understand how bidding works and buyer protection
            </CardDescription>
          </CardHeader>
          <CardContent>
            <Button variant="outline" className="w-full">Bidding Guide</Button>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Selling</CardTitle>
            <CardDescription>
              Tips for successful auctions and seller tools
            </CardDescription>
          </CardHeader>
          <CardContent>
            <Button variant="outline" className="w-full">Seller Guide</Button>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Shipping</CardTitle>
            <CardDescription>
              Shipping options, tracking, and delivery info
            </CardDescription>
          </CardHeader>
          <CardContent>
            <Button variant="outline" className="w-full">Shipping Info</Button>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Contact Support</CardTitle>
            <CardDescription>
              Can't find what you're looking for? Get in touch
            </CardDescription>
          </CardHeader>
          <CardContent>
            <Button className="w-full">Contact Us</Button>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
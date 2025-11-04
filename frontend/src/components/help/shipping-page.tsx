'use client';

import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';

export function ShippingPage() {
  return (
    <div className="container mx-auto px-4 py-8 max-w-4xl">
      <div className="mb-8">
        <h1 className="text-4xl font-bold mb-4">Shipping Information</h1>
        <p className="text-lg text-muted-foreground">
          Shipping options, delivery times, and tracking information.
        </p>
      </div>

      <div className="space-y-6">
        <Card>
          <CardHeader>
            <CardTitle>Shipping Options</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="grid md:grid-cols-2 gap-4">
              <div>
                <h4 className="font-semibold">Standard Shipping</h4>
                <p>5-7 business days</p>
                <p>Free on orders over $50</p>
              </div>
              <div>
                <h4 className="font-semibold">Express Shipping</h4>
                <p>2-3 business days</p>
                <p>$15.99 flat rate</p>
              </div>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>International Shipping</CardTitle>
          </CardHeader>
          <CardContent>
            <p>We ship to most countries worldwide:</p>
            <ul className="list-disc list-inside mt-2 space-y-1">
              <li>Canada: 7-14 business days</li>
              <li>Europe: 10-20 business days</li>
              <li>Asia: 10-20 business days</li>
              <li>Australia: 10-20 business days</li>
            </ul>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Order Processing</CardTitle>
          </CardHeader>
          <CardContent>
            <p>
              Orders are typically processed within 1-2 business days. 
              You'll receive a shipping confirmation email with tracking information 
              once your order ships.
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Tracking Your Order</CardTitle>
          </CardHeader>
          <CardContent>
            <p>
              Track your order by:
            </p>
            <ul className="list-disc list-inside mt-2 space-y-1">
              <li>Checking your order history</li>
              <li>Using the tracking number in your email</li>
              <li>Contacting our support team</li>
            </ul>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
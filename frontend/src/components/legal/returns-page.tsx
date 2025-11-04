'use client';

import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';

export function ReturnsPage() {
  return (
    <div className="container mx-auto px-4 py-8 max-w-4xl">
      <div className="mb-8">
        <h1 className="text-4xl font-bold mb-4">Return Policy</h1>
        <p className="text-lg text-muted-foreground">
          Our return policy and instructions for returning items purchased on Blytz.
        </p>
      </div>

      <div className="space-y-6">
        <Card>
          <CardHeader>
            <CardTitle>30-Day Return Window</CardTitle>
          </CardHeader>
          <CardContent>
            <p>
              You can return most items within 30 days of delivery for a full refund. 
              Items must be in their original condition and packaging.
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>How to Return</CardTitle>
          </CardHeader>
          <CardContent>
            <ol className="list-decimal list-inside space-y-2">
              <li>Go to your order history and select the item to return</li>
              <li>Choose a reason for the return</li>
              <li>Print the return shipping label</li>
              <li>Package the item and drop it off at any shipping location</li>
              <li>Refund will be processed within 5-7 business days</li>
            </ol>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Non-Returnable Items</CardTitle>
          </CardHeader>
          <CardContent>
            <p>The following items cannot be returned:</p>
            <ul className="list-disc list-inside mt-2 space-y-1">
              <li>Perishable items</li>
              <li>Custom or personalized items</li>
              <li>Digital downloads</li>
              <li>Intimate items (for health/hygiene reasons)</li>
            </ul>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
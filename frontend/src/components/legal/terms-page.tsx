'use client';

import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';

export function TermsPage() {
  return (
    <div className="container mx-auto px-4 py-8 max-w-4xl">
      <div className="mb-8">
        <h1 className="text-4xl font-bold mb-4">Terms of Service</h1>
        <p className="text-lg text-muted-foreground">
          By using Blytz, you agree to these terms and conditions.
        </p>
      </div>

      <div className="space-y-6">
        <Card>
          <CardHeader>
            <CardTitle>Account Responsibilities</CardTitle>
          </CardHeader>
          <CardContent>
            <p>You are responsible for:</p>
            <ul className="list-disc list-inside mt-2 space-y-1">
              <li>Maintaining the confidentiality of your account</li>
              <li>All activities that occur under your account</li>
              <li>Notifying us of unauthorized use</li>
              <li>Providing accurate and current information</li>
            </ul>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Auction Rules</CardTitle>
          </CardHeader>
          <CardContent>
            <p>When participating in auctions:</p>
            <ul className="list-disc list-inside mt-2 space-y-1">
              <li>All bids are binding commitments</li>
              <li>Winning bidders must complete payment</li>
              <li>Sellers must deliver items as described</li>
              <li>No shill bidding or fake bids allowed</li>
            </ul>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Prohibited Activities</CardTitle>
          </CardHeader>
          <CardContent>
            <p>You may not:</p>
            <ul className="list-disc list-inside mt-2 space-y-1">
              <li>Sell illegal or prohibited items</li>
              <li>Use the service for fraudulent purposes</li>
              <li>Interfere with or disrupt the service</li>
              <li>Violate applicable laws or regulations</li>
            </ul>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Limitation of Liability</CardTitle>
          </CardHeader>
          <CardContent>
            <p>
              Blytz is provided "as is" without warranties of any kind. 
              We are not liable for damages arising from your use of the service.
            </p>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
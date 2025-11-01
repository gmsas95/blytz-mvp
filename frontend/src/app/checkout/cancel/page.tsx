'use client';

import { XCircle, ArrowLeft } from 'lucide-react';
import { useRouter } from 'next/navigation';

import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';

export default function CheckoutCancelPage() {
  const router = useRouter();

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="max-w-md mx-auto text-center">
        <XCircle className="w-16 h-16 text-red-500 mx-auto mb-4" />
        <h1 className="text-3xl font-bold mb-4">Payment Cancelled</h1>
        <p className="text-muted-foreground mb-6">
          Your payment has been cancelled. No charges were made to your account.
        </p>

        <Card className="mb-6">
          <CardHeader>
            <CardTitle>What would you like to do?</CardTitle>
          </CardHeader>
          <CardContent className="space-y-3">
            <Button onClick={() => router.push('/checkout')} className="w-full">
              <ArrowLeft className="w-4 h-4 mr-2" />
              Try Payment Again
            </Button>
            <Button variant="outline" onClick={() => router.push('/cart')} className="w-full">
              Review Cart
            </Button>
            <Button variant="ghost" onClick={() => router.push('/')} className="w-full">
              Continue Shopping
            </Button>
          </CardContent>
        </Card>

        <div className="text-sm text-muted-foreground">
          <p>If you encountered any issues during payment, please contact our support team.</p>
        </div>
      </div>
    </div>
  );
}

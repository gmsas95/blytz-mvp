'use client';

import { useRouter } from 'next/navigation';
import { CheckCircle } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { PaymentResponse } from '@/types/payment';

interface PaymentSuccessProps {
  paymentResponse: PaymentResponse;
}

export default function PaymentSuccess({ paymentResponse }: PaymentSuccessProps) {
  const router = useRouter();

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="max-w-md mx-auto text-center">
        <CheckCircle className="w-16 h-16 text-green-500 mx-auto mb-4" />
        <h1 className="text-3xl font-bold mb-4">Payment Successful!</h1>
        <p className="text-muted-foreground mb-6">
          Your order #{paymentResponse.orderNumber} has been confirmed.
        </p>
        <div className="space-y-3">
          <Button onClick={() => router.push('/orders')} className="w-full">
            View Order
          </Button>
          <Button variant="outline" onClick={() => router.push('/')} className="w-full">
            Continue Shopping
          </Button>
        </div>
      </div>
    </div>
  );
}
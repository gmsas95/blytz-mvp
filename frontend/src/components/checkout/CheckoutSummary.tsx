'use client';

import { Badge } from '@/components/ui/badge';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Separator } from '@/components/ui/separator';
import { Cart, PaymentMethodInfo } from '@/types/payment';

interface CheckoutSummaryProps {
  cart: Cart;
  selectedPaymentMethod: string;
  paymentMethods: PaymentMethodInfo[];
  total: number;
}

export default function CheckoutSummary({
  cart,
  selectedPaymentMethod,
  paymentMethods,
  total
}: CheckoutSummaryProps) {
  const selectedMethod = paymentMethods.find((m) => m.id === selectedPaymentMethod);
  const processingFee = selectedMethod?.fee || 0;

  return (
    <div>
      <h2 className="text-2xl font-bold mb-6">Order Summary</h2>
      <Card>
        <CardHeader>
          <CardTitle>Cart Items ({cart.itemCount})</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          {cart.items.map((item) => (
            <div key={item.id} className="flex items-center space-x-4">
              <div className="w-16 h-16 bg-gray-200 rounded-lg flex items-center justify-center">
                <span className="text-2xl">📦</span>
              </div>
              <div className="flex-1">
                <h3 className="font-medium">{item.product.title}</h3>
                <p className="text-sm text-muted-foreground">
                  Quantity: {item.quantity} × RM{item.product.price.toFixed(2)}
                </p>
                {item.selectedAuction && (
                  <Badge variant="secondary" className="mt-1">
                    Auction: {item.selectedAuction.product.title}
                  </Badge>
                )}
              </div>
              <div className="text-right">
                <p className="font-medium">
                  RM{(item.product.price * item.quantity).toFixed(2)}
                </p>
              </div>
            </div>
          ))}

          <Separator />

          <div className="space-y-2">
            <div className="flex justify-between">
              <span>Subtotal</span>
              <span>RM{cart.total.toFixed(2)}</span>
            </div>
            {processingFee > 0 && (
              <div className="flex justify-between">
                <span>Processing Fee</span>
                <span>RM{processingFee.toFixed(2)}</span>
              </div>
            )}
            <Separator />
            <div className="flex justify-between font-bold text-lg">
              <span>Total</span>
              <span>RM{total.toFixed(2)}</span>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
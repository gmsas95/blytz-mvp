'use client';

import { CreditCard, Smartphone, Building2 } from 'lucide-react';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Label } from '@/components/ui/label';
import { RadioGroup, RadioGroupItem } from '@/components/ui/radio-group';
import { Separator } from '@/components/ui/separator';
import { Cart, PaymentMethodInfo } from '@/types/payment';

interface PaymentMethodSelectorProps {
  paymentMethods: PaymentMethodInfo[];
  selectedMethod: string;
  onMethodChange: (method: string) => void;
  cart: Cart;
  processing: boolean;
  onPayment: () => void;
  error: string | null;
  total: number;
}

export default function PaymentMethodSelector({
  paymentMethods,
  selectedMethod,
  onMethodChange,
  cart,
  processing,
  onPayment,
  error,
  total
}: PaymentMethodSelectorProps) {
  const getPaymentIcon = (method: PaymentMethodInfo) => {
    switch (method.type) {
      case 'bank_transfer':
        return <Building2 className="w-5 h-5" />;
      case 'ewallet':
        return <Smartphone className="w-5 h-5" />;
      case 'card':
        return <CreditCard className="w-5 h-5" />;
      default:
        return <CreditCard className="w-5 h-5" />;
    }
  };

  const selectedPaymentMethod = paymentMethods.find((m) => m.id === selectedMethod);
  const processingFee = selectedPaymentMethod?.fee || 0;

  return (
    <div>
      <h2 className="text-2xl font-bold mb-6">Payment Method</h2>

      <Card>
        <CardHeader>
          <CardTitle>Select Payment Method</CardTitle>
        </CardHeader>
        <CardContent>
          <RadioGroup value={selectedMethod} onValueChange={onMethodChange}>
            {paymentMethods.map((method) => (
              <div
                key={method.id}
                className="flex items-center space-x-3 p-3 border rounded-lg"
              >
                <RadioGroupItem value={method.id} id={method.id} />
                <Label
                  htmlFor={method.id}
                  className="flex items-center space-x-3 flex-1 cursor-pointer"
                >
                  {getPaymentIcon(method)}
                  <div className="flex-1">
                    <div className="font-medium">{method.name}</div>
                    <div className="text-sm text-muted-foreground">{method.description}</div>
                  </div>
                  {method.fee > 0 && (
                    <Badge variant="secondary">+RM{method.fee.toFixed(2)}</Badge>
                  )}
                </Label>
              </div>
            ))}
          </RadioGroup>

          <Separator className="my-6" />

          <div className="space-y-4">
            <div className="bg-gray-50 p-4 rounded-lg">
              <h3 className="font-medium mb-2">Payment Details</h3>
              <div className="space-y-1 text-sm">
                <div className="flex justify-between">
                  <span>Items:</span>
                  <span>{cart.itemCount}</span>
                </div>
                <div className="flex justify-between">
                  <span>Subtotal:</span>
                  <span>RM{cart.total.toFixed(2)}</span>
                </div>
                {processingFee > 0 && (
                  <div className="flex justify-between">
                    <span>Processing Fee:</span>
                    <span>RM{processingFee.toFixed(2)}</span>
                  </div>
                )}
                <Separator className="my-2" />
                <div className="flex justify-between font-bold">
                  <span>Total Amount:</span>
                  <span>RM{total.toFixed(2)}</span>
                </div>
              </div>
            </div>

            <Button onClick={onPayment} disabled={processing} className="w-full" size="lg">
              {processing ? (
                <>
                  Processing Payment...
                </>
              ) : (
                `Pay RM${total.toFixed(2)}`
              )}
            </Button>

            <p className="text-xs text-muted-foreground text-center">
              By completing this payment, you agree to our terms of service and privacy policy.
              Your payment is secured by Fiuu payment gateway.
            </p>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
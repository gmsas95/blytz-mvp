'use client';

import { useRouter } from 'next/navigation';
import { useState, useEffect } from 'react';

import { Alert, AlertDescription } from '@/components/ui/alert';
import { Button } from '@/components/ui/button';
import { Loader2 } from 'lucide-react';
import { api } from '@/lib/api-adapter';
import { Cart, PaymentMethodInfo, PaymentResponse } from '@/types/payment';

import CheckoutSummary from '@/components/checkout/CheckoutSummary';
import PaymentMethodSelector from '@/components/checkout/PaymentMethodSelector';
import { usePaymentProcessor } from '@/components/checkout/PaymentProcessor';
import PaymentSuccess from '@/components/checkout/PaymentSuccess';

export default function CheckoutPage() {
  const router = useRouter();
  const [cart, setCart] = useState<Cart | null>(null);
  const [paymentMethods, setPaymentMethods] = useState<PaymentMethodInfo[]>([]);
  const [selectedMethod, setSelectedMethod] = useState<string>('fpx');
  const [loading, setLoading] = useState(true);
  const [processing, setProcessing] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [paymentComplete, setPaymentComplete] = useState(false);
  const [paymentResponse, setPaymentResponse] = useState<PaymentResponse | null>(null);

  useEffect(() => {
    loadCheckoutData();
  }, []);

  const loadCheckoutData = async () => {
    try {
      setLoading(true);
      const [cartResponse, methodsResponse] = await Promise.all([
        api.getCart(),
        api.getPaymentMethods(),
      ]);

      if (cartResponse.success && cartResponse.data) {
        setCart(cartResponse.data);
      }

      if (methodsResponse.success && methodsResponse.data) {
        setPaymentMethods(methodsResponse.data);
      }
    } catch (err) {
      setError('Failed to load checkout data');
    } finally {
      setLoading(false);
    }
  };

  const calculateTotal = () => {
    if (!cart) return 0;
    const selectedPaymentMethod = paymentMethods.find((m) => m.id === selectedMethod);
    const processingFee = selectedPaymentMethod?.fee || 0;
    return cart.total + processingFee;
  };

  const handlePaymentSuccess = (response: PaymentResponse) => {
    setPaymentResponse(response);
    setPaymentComplete(true);
  };

  const handlePaymentError = (errorMessage: string) => {
    setError(errorMessage);
  };

  const handleProcessingChange = (isProcessing: boolean) => {
    setProcessing(isProcessing);
  };

  const { processPayment } = usePaymentProcessor({
    cart: cart!,
    selectedPaymentMethod: selectedMethod,
    total: cart ? calculateTotal() : 0,
    onPaymentSuccess: handlePaymentSuccess,
    onPaymentError: handlePaymentError,
    onProcessingChange: handleProcessingChange,
  });

  if (loading) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="flex items-center justify-center min-h-[400px]">
          <Loader2 className="w-8 h-8 animate-spin" />
        </div>
      </div>
    );
  }

  if (!cart || cart.items.length === 0) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="text-center">
          <h1 className="text-3xl font-bold mb-4">Your cart is empty</h1>
          <Button onClick={() => router.push('/products')}>Continue Shopping</Button>
        </div>
      </div>
    );
  }

  if (paymentComplete && paymentResponse) {
    return <PaymentSuccess paymentResponse={paymentResponse} />;
  }

  if (!cart || cart.items.length === 0) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="text-center">
          <h1 className="text-3xl font-bold mb-4">Your cart is empty</h1>
          <Button onClick={() => router.push('/products')}>Continue Shopping</Button>
        </div>
      </div>
    );
  }

  const total = calculateTotal();

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="grid lg:grid-cols-2 gap-8">
        {/* Order Summary */}
        <CheckoutSummary
          cart={cart}
          selectedPaymentMethod={selectedMethod}
          paymentMethods={paymentMethods}
          total={total}
        />

        {/* Payment Method */}
        <div>
          {error && (
            <Alert variant="destructive" className="mb-6">
              <AlertDescription>{error}</AlertDescription>
            </Alert>
          )}

          <PaymentMethodSelector
            paymentMethods={paymentMethods}
            selectedMethod={selectedMethod}
            onMethodChange={setSelectedMethod}
            cart={cart}
            processing={processing}
            onPayment={processPayment}
            error={error}
            total={total}
          />
        </div>
      </div>
    </div>
  );
}

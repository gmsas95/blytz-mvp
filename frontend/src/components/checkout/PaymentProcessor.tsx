'use client';

import { useState, useCallback } from 'react';
import { paymentService } from '@/services/payment.service';
import { Cart, PaymentRequest, PaymentResponse, FiuuConfig } from '@/types/payment';

export function usePaymentProcessor({
  cart,
  selectedPaymentMethod,
  total,
  onPaymentSuccess,
  onPaymentError,
  onProcessingChange
}: {
  cart: Cart;
  selectedPaymentMethod: string;
  total: number;
  onPaymentSuccess: (response: PaymentResponse) => void;
  onPaymentError: (error: string) => void;
  onProcessingChange: (processing: boolean) => void;
}) {
  const [isScriptLoaded, setIsScriptLoaded] = useState(false);

  const loadFiuuScript = useCallback((scriptUrl: string): Promise<void> => {
    return new Promise((resolve, reject) => {
      console.log('Loading Fiuu script from:', scriptUrl);

      if (window.IPGSeamless) {
        console.log('Fiuu script already loaded');
        setIsScriptLoaded(true);
        resolve();
        return;
      }

      const script = document.createElement('script');
      script.src = scriptUrl;
      script.async = true;
      script.crossOrigin = 'anonymous';

      // Add timeout handling
      const timeout = setTimeout(() => {
        console.error('Fiuu script loading timeout');
        reject(new Error('Fiuu script loading timeout'));
      }, 10000); // 10 second timeout

      script.onload = () => {
        clearTimeout(timeout);
        console.log('Fiuu script loaded successfully');
        setIsScriptLoaded(true);

        // Verify the script is actually available
        if (window.IPGSeamless) {
          resolve();
        } else {
          reject(new Error('Fiuu script loaded but IPGSeamless not found on window'));
        }
      };

      script.onerror = (error) => {
        clearTimeout(timeout);
        console.error('Failed to load Fiuu script:', error);
        reject(new Error(`Failed to load Fiuu payment script from ${scriptUrl}`));
      };

      document.head.appendChild(script);
    });
  }, []);

  const processPayment = useCallback(async () => {
    try {
      onProcessingChange(true);

      // Generate order details
      const orderNumber = `BLYTZ_${Date.now()}`;
      const productDescription = `Payment for ${cart.itemCount} item(s)`;

      // Create payment request
      const paymentRequest: PaymentRequest = {
        amount: total,
        currency: 'MYR',
        paymentMethod: selectedPaymentMethod,
        orderNumber,
        description: productDescription,
        returnUrl: `${window.location.origin}/checkout/success`,
        cancelUrl: `${window.location.origin}/checkout/cancel`,
        webhookUrl: `${window.location.origin}/api/payments/webhook`,
      };

      // Get Fiuu configuration
      const fiuuConfig = await paymentService.getFiuuSeamlessConfig(
        orderNumber,
        total,
        'John Doe',
        'john@example.com',
        '+60123456789',
        productDescription
      );

      // Load Fiuu script if not already loaded
      await loadFiuuScript(fiuuConfig.scriptUrl);

      // Create payment
      const paymentResponse = await paymentService.createPayment(paymentRequest);

      // Initialize Fiuu seamless payment
      if (window.IPGSeamless) {
        // Update Fiuu config with current cart data
        const updatedConfig: FiuuConfig = {
          ...fiuuConfig,
          amount: total,
          orderNumber,
          productDescription,
          paymentMethod: selectedPaymentMethod === 'fpx' ? 'all' : selectedPaymentMethod,
        };

        const seamless = await paymentService.initializeFiuuPayment(updatedConfig);
        const result = await paymentService.processFiuuPayment(seamless);

        onPaymentSuccess(result);
      } else if (paymentResponse.redirectUrl) {
        // Fallback to redirect method
        window.location.href = paymentResponse.redirectUrl;
      } else {
        throw new Error('Payment initialization failed');
      }
    } catch (error) {
      console.error('Payment processing error:', error);
      const errorMessage = error instanceof Error ? error.message : 'Payment failed';
      onPaymentError(errorMessage);
    } finally {
      onProcessingChange(false);
    }
  }, [cart, selectedPaymentMethod, total, onPaymentSuccess, onPaymentError, onProcessingChange, loadFiuuScript]);

  return {
    processPayment,
    isScriptLoaded,
  };
}
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
      console.log('Loading jQuery and Fiuu scripts from:', scriptUrl);

      // Load jQuery first if not already loaded
      const loadjQuery = () => {
        return new Promise<void>((resolveJQuery, rejectJQuery) => {
          if (window.jQuery || window.$) {
            resolveJQuery();
            return;
          }

          const jqueryScript = document.createElement('script');
          jqueryScript.src = 'https://ajax.googleapis.com/ajax/libs/jquery/3.5.1/jquery.min.js';
          jqueryScript.async = true;

          jqueryScript.onload = () => {
            console.log('jQuery loaded successfully');
            resolveJQuery();
          };

          jqueryScript.onerror = () => {
            rejectJQuery(new Error('Failed to load jQuery'));
          };

          document.head.appendChild(jqueryScript);
        });
      };

      // Load Fiuu script
      const loadFiuuScript = () => {
        return new Promise<void>((resolveFiuu, rejectFiuu) => {
          if (window.MOLPaySeamless) {
            console.log('Fiuu script already loaded');
            setIsScriptLoaded(true);
            resolveFiuu();
            return;
          }

          const script = document.createElement('script');
          script.src = scriptUrl;
          script.async = true;

          script.onload = () => {
            console.log('Fiuu script loaded successfully');
            setIsScriptLoaded(true);

            if (window.MOLPaySeamless) {
              resolveFiuu();
            } else {
              rejectFiuu(new Error('Fiuu script loaded but MOLPaySeamless not found on window'));
            }
          };

          script.onerror = () => {
            rejectFiuu(new Error(`Failed to load Fiuu payment script from ${scriptUrl}`));
          };

          document.head.appendChild(script);
        });
      };

      // Load both scripts in sequence
      loadjQuery()
        .then(() => loadFiuuScript())
        .then(() => resolve())
        .catch((error) => reject(error));
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

      // Initialize Fiuu seamless payment using jQuery and MOLPaySeamless
      if (window.$ && window.MOLPaySeamless) {
        console.log('Initializing Fiuu payment with config:', fiuuConfig);

        // Create payment options for Fiuu
        const paymentOptions = {
          mpsmerchantid: fiuuConfig.mpsmerchantid,
          mpschannel: selectedPaymentMethod === 'fpx' ? 'fpx' : selectedPaymentMethod,
          mpsamount: fiuuConfig.mpsamount,
          mpsorderid: fiuuConfig.mpsorderid,
          mpsbill_name: fiuuConfig.mpsbill_name,
          mpsbill_email: fiuuConfig.mpsbill_email,
          mpsbill_mobile: fiuuConfig.mpsbill_mobile,
          mpsbill_desc: fiuuConfig.mpsbill_desc,
          mpscurrency: fiuuConfig.mpscurrency,
          mpslangcode: fiuuConfig.mpslangcode,
          mpscountry: fiuuConfig.mpscountry,
          vcode: fiuuConfig.vcode
        };

        console.log('Payment options:', paymentOptions);

        // Create a temporary button for MOLPaySeamless
        const tempButton = document.createElement('button');
        tempButton.id = 'temp-fiuu-button';
        tempButton.style.display = 'none';
        document.body.appendChild(tempButton);

        try {
          // Initialize MOLPaySeamless
          window.$('#temp-fiuu-button').MOLPaySeamless(paymentOptions);

          console.log('Fiuu payment initialized successfully');

          // Simulate payment success for now (in real implementation, this would be handled by callbacks)
          setTimeout(() => {
            const mockResponse: PaymentResponse = {
              id: `PAY_${Date.now()}`,
              orderNumber: orderNumber,
              amount: total,
              currency: 'MYR',
              status: 'completed',
              paymentMethod: selectedPaymentMethod,
              transactionId: `TXN_${Date.now()}`,
              createdAt: new Date().toISOString(),
              updatedAt: new Date().toISOString(),
            };

            // Clean up
            document.body.removeChild(tempButton);
            onPaymentSuccess(mockResponse);
          }, 2000);

        } catch (fiuuError) {
          console.error('Fiuu payment initialization error:', fiuuError);
          document.body.removeChild(tempButton);
          throw new Error('Failed to initialize Fiuu payment');
        }
      } else if (paymentResponse.redirectUrl) {
        // Fallback to redirect method
        window.location.href = paymentResponse.redirectUrl;
      } else {
        throw new Error('Payment initialization failed - jQuery or MOLPaySeamless not available');
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
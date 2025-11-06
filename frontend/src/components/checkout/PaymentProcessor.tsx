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
            console.log('jQuery already loaded');
            resolveJQuery();
            return;
          }

          console.log('Loading jQuery...');
          const jqueryScript = document.createElement('script');
          // Try multiple CDNs for reliability
          const cdnUrls = [
            'https://code.jquery.com/jquery-3.6.0.min.js',
            'https://ajax.googleapis.com/ajax/libs/jquery/3.6.0/jquery.min.js',
            'https://cdnjs.cloudflare.com/ajax/libs/jquery/3.6.0/jquery.min.js'
          ];

          let currentCdnIndex = 0;

          const tryLoadFromCdn = () => {
            if (currentCdnIndex >= cdnUrls.length) {
              rejectJQuery(new Error('Failed to load jQuery from all CDNs'));
              return;
            }

            const script = document.createElement('script');
            script.src = cdnUrls[currentCdnIndex];
            script.async = true;

            script.onload = () => {
              console.log(`jQuery loaded successfully from CDN ${currentCdnIndex + 1}`);
              resolveJQuery();
            };

            script.onerror = () => {
              console.warn(`Failed to load jQuery from CDN ${currentCdnIndex + 1}: ${cdnUrls[currentCdnIndex]}`);
              currentCdnIndex++;
              tryLoadFromCdn();
            };

            document.head.appendChild(script);
          };

          tryLoadFromCdn();
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
          // Initialize MOLPaySeamless with proper callbacks
          window.$('#temp-fiuu-button').MOLPaySeamless(paymentOptions);

          console.log('Fiuu payment initialized successfully');

          // Set up global callbacks for Fiuu responses
          window.molpay_seamless_acct_type = '0'; // Account type for callback handling

          // Set up success callback
          window.molpay_callback = (response: any) => {
            console.log('Fiuu payment success callback:', response);

            const paymentResponse: PaymentResponse = {
              id: response.txnID || `PAY_${Date.now()}`,
              orderNumber: response.orderNumber || orderNumber,
              amount: parseFloat(response.amount) || total,
              currency: response.currency || 'MYR',
              status: 'completed',
              paymentMethod: response.channel || selectedPaymentMethod,
              transactionId: response.txnID,
              createdAt: new Date().toISOString(),
              updatedAt: new Date().toISOString(),
            };

            // Clean up
            if (document.body.contains(tempButton)) {
              document.body.removeChild(tempButton);
            }
            onPaymentSuccess(paymentResponse);
          };

          // Set up error callback
          window.molpay_error = (error: any) => {
            console.error('Fiuu payment error callback:', error);

            const errorMessage = error.error_description || error.message || 'Payment failed';

            // Clean up
            if (document.body.contains(tempButton)) {
              document.body.removeChild(tempButton);
            }
            onPaymentError(errorMessage);
          };

          console.log('Fiuu callbacks set up successfully');

        } catch (fiuuError) {
          console.error('Fiuu payment initialization error:', fiuuError);
          if (document.body.contains(tempButton)) {
            document.body.removeChild(tempButton);
          }
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
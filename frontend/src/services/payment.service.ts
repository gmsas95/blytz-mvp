import { FiuuConfig, PaymentMethodInfo, PaymentRequest, PaymentResponse, IPGSeamlessInstance } from '@/types/payment';

export class PaymentService {
  private static instance: PaymentService;

  static getInstance(): PaymentService {
    if (!PaymentService.instance) {
      PaymentService.instance = new PaymentService();
    }
    return PaymentService.instance;
  }

  async getPaymentMethods(): Promise<PaymentMethodInfo[]> {
    try {
      const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/v1/public/payment-methods`);
      const data = await response.json();

      if (data.success) {
        return data.data;
      } else {
        // Fallback to mock payment methods
        return this.getMockPaymentMethods();
      }
    } catch (error) {
      console.warn('Failed to fetch payment methods, using mock data:', error);
      return this.getMockPaymentMethods();
    }
  }

  async getFiuuSeamlessConfig(orderNumber: string, amount: number, billName: string, billEmail: string, billMobile: string, billDesc: string): Promise<FiuuConfig> {
    try {
      const params = new URLSearchParams({
        order_id: orderNumber,
        amount: (amount * 100).toString(), // Convert to cents
        bill_name: billName,
        bill_email: billEmail,
        bill_mobile: billMobile,
        bill_desc: billDesc,
        channel: 'FPX'
      });

      const response = await fetch(
        `${process.env.NEXT_PUBLIC_API_URL}/api/v1/public/seamless/config?${params}`
      );
      const data = await response.json();

      if (data.success) {
        return data.data;
      } else {
        return this.getMockFiuuConfig(orderNumber, amount, billName, billEmail, billMobile, billDesc);
      }
    } catch (error) {
      console.warn('Failed to fetch Fiuu config, using mock data:', error);
      return this.getMockFiuuConfig(orderNumber, amount, billName, billEmail, billMobile, billDesc);
    }
  }

  async createPayment(paymentRequest: PaymentRequest): Promise<PaymentResponse> {
    try {
      const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/v1/payments`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(paymentRequest),
      });
      const data = await response.json();

      if (data.success) {
        return data.data;
      } else {
        throw new Error(data.error || 'Failed to create payment');
      }
    } catch (error) {
      console.warn('Failed to create payment, using mock response:', error);
      return this.getMockPaymentResponse(paymentRequest);
    }
  }

  initializeFiuuPayment(config: FiuuConfig): Promise<IPGSeamlessInstance> {
    return new Promise((resolve, reject) => {
      if (!window.IPGSeamless) {
        reject(new Error('Fiuu payment script not loaded'));
        return;
      }

      try {
        const seamless = new window.IPGSeamless(config);
        resolve(seamless);
      } catch (error) {
        reject(error);
      }
    });
  }

  processFiuuPayment(seamless: IPGSeamlessInstance): Promise<PaymentResponse> {
    return new Promise((resolve, reject) => {
      seamless.setCompleteCallback((response: any) => {
        // Process successful payment response
        const paymentResponse: PaymentResponse = {
          id: `PAY_${Date.now()}`,
          orderNumber: response.orderNumber || `ORDER_${Date.now()}`,
          amount: parseFloat(response.amount) || 0,
          currency: response.currency || 'MYR',
          status: 'completed',
          paymentMethod: response.paymentMethod || 'fpx',
          transactionId: response.transactionId,
          createdAt: new Date().toISOString(),
          updatedAt: new Date().toISOString(),
        };
        resolve(paymentResponse);
      });

      seamless.setErrorCallback((error: any) => {
        reject(new Error(error.message || 'Payment failed'));
      });

      // Start payment
      seamless.makePayment();
    });
  }

  private getMockPaymentMethods(): PaymentMethodInfo[] {
    return [
      {
        id: 'fpx',
        name: 'Online Banking (FPX)',
        type: 'bank_transfer',
        icon: '🏦',
        description: 'Pay with your Malaysian bank account',
        available: true,
        fee: 0.0,
      },
      {
        id: 'credit_card',
        name: 'Credit/Debit Card',
        type: 'card',
        icon: '💳',
        description: 'Visa, Mastercard, AMEX',
        available: true,
        fee: 1.5,
      },
      {
        id: 'shopeepay',
        name: 'ShopeePay',
        type: 'ewallet',
        icon: '🛒',
        description: 'Quick payment with ShopeePay',
        available: true,
        fee: 0.5,
      },
    ];
  }

  private getMockFiuuConfig(orderNumber: string, amount: number, billName: string, billEmail: string, billMobile: string, billDesc: string): FiuuConfig {
    const isSandbox = process.env.NEXT_PUBLIC_FIUU_SANDBOX === 'true';

    return {
      version: '7.5.0',
      actionType: 'Pay',
      merchantID: 'MERCHANT_ID',
      paymentMethod: 'all',
      orderNumber,
      amount,
      currency: 'MYR',
      productDescription: billDesc,
      userName: billName,
      userEmail: billEmail,
      userContact: billMobile,
      remark: 'Payment for auction items',
      lang: 'en',
      vcode: 'MOCK_VCODE_123456',
      callbackURL: `${process.env.NEXT_PUBLIC_API_URL}/api/v1/webhooks/fiuu`,
      returnURL: `${window.location.origin}/checkout/success`,
      backgroundUrl: `${process.env.NEXT_PUBLIC_API_URL}/api/v1/webhooks/fiuu`,
      sandbox: isSandbox,
      scriptUrl: isSandbox
        ? 'https://sandbox.merchant.fiuu.com.my/RMS2/IPGSeamless/IPGSeamless.js'
        : 'https://api.merchant.fiuu.com.my/RMS2/IPGSeamless/IPGSeamless.js',
    };
  }

  private getMockPaymentResponse(paymentRequest: PaymentRequest): PaymentResponse {
    return {
      id: `PAY_${Date.now()}`,
      orderNumber: paymentRequest.orderNumber,
      amount: paymentRequest.amount,
      currency: paymentRequest.currency,
      status: 'pending',
      paymentMethod: paymentRequest.paymentMethod,
      redirectUrl: `https://api.merchant.fiuu.com/payment/mock?order=${paymentRequest.orderNumber}`,
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
    };
  }
}

export const paymentService = PaymentService.getInstance();
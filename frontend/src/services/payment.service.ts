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

  // Fiuu payment processing is now handled directly in the PaymentProcessor component
  // using jQuery and MOLPaySeamless, so these methods are deprecated

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
      mpsmerchantid: 'blytzventures', // Your merchant ID
      mpschannel: 'fpx',
      mpsamount: amount.toString(),
      mpsorderid: orderNumber,
      mpsbill_name: billName,
      mpsbill_email: billEmail,
      mpsbill_mobile: billMobile,
      mpsbill_desc: billDesc,
      mpscurrency: 'MYR',
      mpslangcode: 'en',
      mpscountry: 'MY',
      vcode: 'MOCK_VCODE_123456',
      scriptUrl: 'https://pay.fiuu.com/RMS/API/seamless/3.28/js/MOLPay_seamless.deco.js',
      sandbox: isSandbox,
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
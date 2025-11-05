export interface FiuuConfig {
  version: string;
  actionType: string;
  merchantID: string;
  paymentMethod: string;
  orderNumber: string;
  amount: number;
  currency: string;
  productDescription: string;
  userName: string;
  userEmail: string;
  userContact: string;
  remark: string;
  lang: string;
  vcode: string;
  callbackURL: string;
  returnURL: string;
  backgroundUrl: string;
  sandbox: boolean;
  scriptUrl: string;
}

export interface PaymentMethodInfo {
  id: string;
  name: string;
  type: 'bank_transfer' | 'ewallet' | 'card';
  icon: string;
  description: string;
  available: boolean;
  fee: number;
}

export interface PaymentRequest {
  amount: number;
  currency: string;
  paymentMethod: string;
  orderNumber: string;
  description: string;
  returnUrl: string;
  cancelUrl: string;
  webhookUrl: string;
}

export interface PaymentResponse {
  id: string;
  orderNumber: string;
  amount: number;
  currency: string;
  status: 'pending' | 'processing' | 'completed' | 'failed' | 'cancelled';
  paymentMethod: string;
  redirectUrl?: string;
  transactionId?: string;
  createdAt: string;
  updatedAt: string;
}

export interface Cart {
  id: string;
  items: CartItem[];
  total: number;
  itemCount: number;
  createdAt: string;
  updatedAt: string;
}

export interface CartItem {
  id: string;
  product: {
    id: string;
    title: string;
    price: number;
    images: string[];
  };
  quantity: number;
  selectedAuction?: {
    id: string;
    product: {
      title: string;
    };
  };
}

export interface IPGSeamlessConfig {
  version: string;
  actionType: string;
  merchantID: string;
  paymentMethod: string;
  orderNumber: string;
  amount: number;
  currency: string;
  productDescription: string;
  userName: string;
  userEmail: string;
  userContact: string;
  remark: string;
  lang: string;
  vcode: string;
  callbackURL: string;
  returnURL: string;
  backgroundUrl: string;
  sandbox: boolean;
}

declare global {
  interface Window {
    IPGSeamless: IPGSeamlessConstructor;
  }
}

export interface IPGSeamlessConstructor {
  new(config: IPGSeamlessConfig): IPGSeamlessInstance;
}

export interface IPGSeamlessInstance {
  setCompleteCallback(callback: (response: any) => void): void;
  setErrorCallback(callback: (error: any) => void): void;
  makePayment(): void;
}

export type PaymentStatus = 'idle' | 'loading' | 'processing' | 'success' | 'error';

export interface PaymentState {
  status: PaymentStatus;
  error: string | null;
  response: PaymentResponse | null;
}
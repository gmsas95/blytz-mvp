export interface FiuuConfig {
  mpsmerchantid: string;
  mpschannel: string;
  mpsamount: string;
  mpsorderid: string;
  mpsbill_name: string;
  mpsbill_email: string;
  mpsbill_mobile: string;
  mpsbill_desc: string;
  mpscurrency: string;
  mpslangcode: string;
  mpscountry: string;
  vcode: string;
  scriptUrl: string;
  sandbox: boolean;
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
    jQuery: any;
    $: any;
    MOLPaySeamless: any;
  }
}

export type PaymentStatus = 'idle' | 'loading' | 'processing' | 'success' | 'error';

export interface PaymentState {
  status: PaymentStatus;
  error: string | null;
  response: PaymentResponse | null;
}
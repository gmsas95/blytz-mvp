import { Metadata } from 'next';
import { ShippingPage } from '@/components/help/shipping-page';

export const metadata: Metadata = {
  title: 'Shipping Info - Blytz',
  description: 'Blytz shipping information and delivery options',
};

export default function Shipping() {
  return <ShippingPage />;
}
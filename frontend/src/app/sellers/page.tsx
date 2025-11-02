import { Metadata } from 'next';
import { SellersPage } from '@/components/sellers/sellers-page';

export const metadata: Metadata = {
  title: 'Sellers - Blytz',
  description: 'Become a seller on Blytz and start your live auction business',
};

export default function Sellers() {
  return <SellersPage />;
}
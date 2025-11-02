import { Metadata } from 'next';
import { ReturnsPage } from '@/components/legal/returns-page';

export const metadata: Metadata = {
  title: 'Returns - Blytz',
  description: 'Blytz return policy and instructions',
};

export default function Returns() {
  return <ReturnsPage />;
}
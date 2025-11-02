import { Metadata } from 'next';
import { CookiesPage } from '@/components/legal/cookies-page';

export const metadata: Metadata = {
  title: 'Cookie Policy - Blytz',
  description: 'Blytz cookie policy and usage',
};

export default function Cookies() {
  return <CookiesPage />;
}
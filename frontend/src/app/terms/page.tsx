import { Metadata } from 'next';
import { TermsPage } from '@/components/legal/terms-page';

export const metadata: Metadata = {
  title: 'Terms of Service - Blytz',
  description: 'Blytz terms of service and user agreement',
};

export default function Terms() {
  return <TermsPage />;
}
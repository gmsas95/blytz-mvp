import { Metadata } from 'next';
import { PrivacyPage } from '@/components/legal/privacy-page';

export const metadata: Metadata = {
  title: 'Privacy Policy - Blytz',
  description: 'Blytz privacy policy and data protection',
};

export default function Privacy() {
  return <PrivacyPage />;
}
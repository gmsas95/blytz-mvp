import { Metadata } from 'next';
import { HelpPage } from '@/components/help/help-page';

export const metadata: Metadata = {
  title: 'Help Center - Blytz',
  description: 'Get help with your Blytz account, auctions, and more',
};

export default function Help() {
  return <HelpPage />;
}
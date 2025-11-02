import { Metadata } from 'next';
import { ContactPage } from '@/components/help/contact-page';

export const metadata: Metadata = {
  title: 'Contact Us - Blytz',
  description: 'Contact Blytz support and customer service',
};

export default function Contact() {
  return <ContactPage />;
}
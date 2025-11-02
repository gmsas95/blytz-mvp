import { Metadata } from 'next';
import { CategoriesPage } from '@/components/categories/categories-page';

export const metadata: Metadata = {
  title: 'Categories - Blytz',
  description: 'Browse products by category on Blytz',
};

export default function Categories() {
  return <CategoriesPage />;
}
import { notFound } from 'next/navigation';

import { Button } from '@/components/ui/button';
import { api } from '@/lib/api-adapter';
import { formatPrice } from '@/lib/utils';

interface Props {
  params: { id: string };
}

export default async function ProductDetailPage({ params }: Props) {
  const res = await api.getProduct(params.id);
  if (!res.success || !res.data) return notFound();
  const product = res.data;

  return (
    <section className="w-full py-16 md:py-24">
      <div className="container mx-auto px-4 grid gap-10 lg:grid-cols-2 items-start">
        <div className="w-full">
          <div className="aspect-square bg-muted rounded-2xl overflow-hidden">
            {/* eslint-disable-next-line @next/next/no-img-element */}
            <img
              src={product.images[0]}
              alt={product.title}
              className="object-cover w-full h-full"
            />
          </div>
        </div>
        <div className="space-y-6">
          <h1 className="text-3xl md:text-4xl font-bold tracking-tight">{product.title}</h1>
          <p className="text-muted-foreground">{product.description}</p>
          <div className="space-y-1">
            <div className="text-3xl font-bold text-primary">{formatPrice(product.price)}</div>
            {product.originalPrice ? (
              <div className="text-sm text-muted-foreground line-through">
                {formatPrice(product.originalPrice)}
              </div>
            ) : null}
          </div>
          <div className="space-y-2">
            <h2 className="font-semibold">Specifications</h2>
            <ul className="text-sm text-muted-foreground grid grid-cols-1 sm:grid-cols-2 gap-2">
              {Object.entries(product.specifications).map(([k, v]) => (
                <li key={k}>
                  <span className="font-medium text-foreground">{k}:</span> {v}
                </li>
              ))}
            </ul>
          </div>
          <div className="flex gap-3">
            <Button size="lg">Add to Cart</Button>
            {product.auction ? (
              <a href={`/auctions/${product.auction.id}`}>
                <Button size="lg" variant="outline">
                  Bid Now
                </Button>
              </a>
            ) : null}
          </div>
        </div>
      </div>
    </section>
  );
}

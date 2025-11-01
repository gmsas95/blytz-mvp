import Link from 'next/link';

import { Button } from '@/components/ui/button';
import { Card, CardContent, CardFooter, CardHeader } from '@/components/ui/card';
import { api } from '@/lib/api-adapter';
import { formatPrice } from '@/lib/utils';

export default async function ProductsPage() {
  const res = await api.getProducts();
  const items = res.success && res.data ? res.data.items : [];

  return (
    <section className="w-full py-16 md:py-24">
      <div className="container mx-auto px-4">
        <div className="text-center space-y-4 mb-12">
          <h1 className="text-3xl md:text-4xl font-bold tracking-tight">Products</h1>
          <p className="text-muted-foreground">Browse our latest products from top sellers</p>
        </div>

        {items.length === 0 ? (
          <div className="text-center text-muted-foreground">No products found.</div>
        ) : (
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
            {items.map((product) => (
              <Link key={product.id} href={`/product/${product.id}`} className="group">
                <Card className="overflow-hidden h-full border-border hover:shadow-lg transition-shadow duration-300">
                  <div className="relative aspect-square bg-muted">
                    {/* eslint-disable-next-line @next/next/no-img-element */}
                    <img
                      src={product.images[0]}
                      alt={product.title}
                      className="object-cover w-full h-full group-hover:scale-105 transition-transform duration-300"
                    />
                  </div>
                  <CardHeader className="pb-3">
                    <h3 className="font-semibold text-lg leading-tight group-hover:text-primary transition-colors">
                      {product.title}
                    </h3>
                    <p className="text-sm text-muted-foreground line-clamp-2">
                      {product.description}
                    </p>
                  </CardHeader>
                  <CardContent className="pb-3">
                    <div className="text-2xl font-bold">{formatPrice(product.price)}</div>
                  </CardContent>
                  <CardFooter>
                    <Button variant="outline" className="w-full">
                      View
                    </Button>
                  </CardFooter>
                </Card>
              </Link>
            ))}
          </div>
        )}
      </div>
    </section>
  );
}

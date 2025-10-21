'use client'

import { useEffect, useState } from 'react'
import Link from 'next/link'
import { Product } from '@/types'
import { Card, CardContent, CardFooter, CardHeader } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { formatPrice } from '@/lib/utils'
import { api } from '@/lib/api-adapter'
import { ShoppingBag } from 'lucide-react'

export function FeaturedProducts() {
  const [products, setProducts] = useState<Product[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    async function loadFeaturedProducts() {
      try {
        // Add a small delay to show loading state, then load real data
        await new Promise(resolve => setTimeout(resolve, 50))
        const response = await api.getFeaturedProducts()
        if (response.success && response.data) {
          setProducts(response.data)
        }
      } catch (error) {
        // If API fails, show empty state
        setProducts([])
      } finally {
        setLoading(false)
      }
    }

    loadFeaturedProducts()
  }, [])

  if (loading) {
    return (
      <section className="w-full bg-secondary/30 py-16 md:py-24">
        <div className="container mx-auto px-4">
          <div className="text-center space-y-4 mb-12">
            <h2 className="text-3xl md:text-4xl font-bold tracking-tight text-foreground">
              Featured Products
            </h2>
            <p className="text-lg text-muted-foreground max-w-2xl mx-auto">
              Discover amazing products from our top sellers
            </p>
          </div>

          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
            {[...Array(3)].map((_, i) => (
              <Card key={i} className="overflow-hidden border-border hover:shadow-lg transition-shadow duration-300">
                <div className="aspect-square bg-muted flex items-center justify-center">
                  <span className="text-muted-foreground text-sm">Loading...</span>
                </div>
                <CardHeader className="pb-3">
                  <div className="h-4 bg-muted rounded animate-pulse mb-2" />
                  <div className="h-3 bg-muted rounded w-2/3 animate-pulse" />
                </CardHeader>
                <CardContent className="pb-3">
                  <div className="h-6 bg-muted rounded w-1/3 animate-pulse" />
                </CardContent>
                <CardFooter>
                  <div className="h-10 bg-muted rounded-lg animate-pulse w-full" />
                </CardFooter>
              </Card>
            ))}
          </div>
        </div>
      </section>
    )
  }

  return (
    <section className="w-full bg-secondary/30 py-16 md:py-24">
      <div className="container mx-auto px-4">
        <div className="text-center space-y-4 mb-12">
          <h2 className="text-3xl md:text-4xl font-bold tracking-tight text-foreground">
            Featured Products
          </h2>
          <p className="text-lg text-muted-foreground max-w-2xl mx-auto">
            Discover amazing products from our top sellers, available for immediate purchase or live auction
          </p>
        </div>

        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
          {products.map((product, index) => (
            <Link key={product.id} href={`/product/${product.id}`} className="group">
              <Card className="overflow-hidden h-full border-border hover:shadow-lg transition-shadow duration-300">
                <div className="relative aspect-square bg-muted">
                  <img
                    src={product.images[0]}
                    alt={product.title}
                    className="object-cover w-full h-full group-hover:scale-105 transition-transform duration-300"
                  />

                  {product.originalPrice && (
                    <div className="absolute top-3 right-3 bg-destructive text-destructive-foreground px-2 py-1 rounded-full text-xs font-medium">
                      {Math.round(((product.originalPrice - product.price) / product.originalPrice) * 100)}% OFF
                    </div>
                  )}

                  <div className="absolute bottom-3 left-3 bg-background/90 backdrop-blur-sm px-3 py-1 rounded-full text-xs font-medium">
                    {product.category}
                  </div>
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
                  <div className="flex items-center justify-between">
                    <div className="space-y-1">
                      <div className="text-2xl font-bold text-foreground">
                        {formatPrice(product.price)}
                      </div>
                      {product.originalPrice && (
                        <div className="text-sm text-muted-foreground line-through">
                          {formatPrice(product.originalPrice)}
                        </div>
                      )}
                    </div>

                    <div className="text-right">
                      <div className="text-sm text-muted-foreground">
                        {product.seller.storeName}
                      </div>
                      <div className="flex items-center gap-1 text-xs text-muted-foreground">
                        <span>⭐ {product.seller.rating}</span>
                        <span>({product.seller.totalSales} sales)</span>
                      </div>
                    </div>
                  </div>
                </CardContent>

                <CardFooter className="gap-2">
                  <div className="flex-1">
                    <Button variant="outline" className="w-full gap-2">
                      <ShoppingBag className="w-4 h-4" />
                      View Details
                    </Button>
                  </div>

                  {product.auction && (
                    <div className="flex-1">
                      <Link href={`/auctions/${product.auction.id}`}>
                        <Button className="w-full gap-2">
                          Bid Now
                        </Button>
                      </Link>
                    </div>
                  )}
                </CardFooter>
              </Card>
            </Link>
          ))}
        </div>

        <div className="text-center mt-12">
          <Link href="/products">
            <Button size="lg" variant="outline">
              View All Products
            </Button>
          </Link>
        </div>
      </div>
    </section>
  )
}
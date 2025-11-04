'use client';

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';

export function CategoriesPage() {
  const categories = [
    { name: 'Electronics', count: 1234, icon: '📱' },
    { name: 'Fashion', count: 856, icon: '👗' },
    { name: 'Home & Garden', count: 642, icon: '🏠' },
    { name: 'Sports', count: 423, icon: '⚽' },
    { name: 'Books', count: 312, icon: '📚' },
    { name: 'Toys & Games', count: 289, icon: '🎮' },
    { name: 'Beauty', count: 276, icon: '💄' },
    { name: 'Automotive', count: 198, icon: '🚗' },
    { name: 'Collectibles', count: 187, icon: '🏆' },
  ];

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="text-center mb-12">
        <h1 className="text-4xl font-bold mb-4">Browse Categories</h1>
        <p className="text-xl text-muted-foreground max-w-2xl mx-auto">
          Find exactly what you're looking for in our organized categories
        </p>
      </div>

      <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
        {categories.map((category) => (
          <Card key={category.name} className="hover:shadow-lg transition-shadow cursor-pointer">
            <CardHeader>
              <div className="flex items-center justify-between">
                <div className="flex items-center space-x-3">
                  <span className="text-2xl">{category.icon}</span>
                  <CardTitle className="text-lg">{category.name}</CardTitle>
                </div>
                <span className="text-sm text-muted-foreground">
                  {category.count} items
                </span>
              </div>
            </CardHeader>
            <CardContent>
              <Button variant="outline" className="w-full">
                Browse {category.name}
              </Button>
            </CardContent>
          </Card>
        ))}
      </div>
    </div>
  );
}
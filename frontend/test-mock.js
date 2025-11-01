const mockProducts = [
  {
    id: '1',
    title: 'Test Product',
    description: 'Test description',
    price: 99.99,
    images: ['https://via.placeholder.com/300'],
    category: 'Electronics',
    seller: {
      id: '1',
      name: 'Test Seller',
      storeName: 'Test Store',
      rating: 4.5,
      totalSales: 100,
    },
  },
];

console.log('Mock data test:', mockProducts.length, 'products');
console.log('First product:', mockProducts[0].title);

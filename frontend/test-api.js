// Simple test to verify mock data works
const { createApiAdapter } = require('./src/lib/api-adapter');

async function testApi() {
  console.log('Testing API adapter...');

  try {
    const api = createApiAdapter();
    console.log('API adapter created successfully');

    console.log('Testing getFeaturedProducts...');
    const result = await api.getFeaturedProducts();
    console.log('API Result:', JSON.stringify(result, null, 2));

    if (result.success && result.data) {
      console.log(`✅ SUCCESS: Found ${result.data.length} products`);
      console.log('First product:', result.data[0]);
    } else {
      console.log('❌ FAILED:', result.error);
    }
  } catch (error) {
    console.log('❌ ERROR:', error.message);
    console.log('Stack:', error.stack);
  }
}

testApi();

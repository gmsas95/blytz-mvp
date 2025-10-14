import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate } from 'k6/metrics';

// Custom metrics
const bidSuccessRate = new Rate('bid_success_rate');
const bidFailureRate = new Rate('bid_failure_rate');

// Test configuration
export const options = {
  stages: [
    { duration: '30s', target: 10 },  // Ramp up to 10 users
    { duration: '1m', target: 50 },   // Stay at 50 users
    { duration: '30s', target: 0 },   // Ramp down to 0 users
  ],
  thresholds: {
    http_req_duration: ['p(95)<300'], // 95% of requests should be below 300ms
    http_req_failed: ['rate<0.1'],    // Error rate should be below 10%
    bid_success_rate: ['rate>0.8'],    // At least 80% of bids should succeed
  },
};

// Test data
const API_URL = __ENV.API_URL || 'http://localhost:8080';
const AUCTION_ID = 'test-auction-001';
const TEST_USERS = [
  'user-001', 'user-002', 'user-003', 'user-004', 'user-005',
  'user-006', 'user-007', 'user-008', 'user-009', 'user-010'
];

export function setup() {
  // Create a test auction
  const createPayload = JSON.stringify({
    auction_id: AUCTION_ID,
    start_time: Math.floor(Date.now() / 1000),
    end_time: Math.floor(Date.now() / 1000) + 3600, // 1 hour from now
    reserve_price: 1000
  });

  const createResponse = http.post(`${API_URL}/auction/auctions`, createPayload, {
    headers: { 'Content-Type': 'application/json' }
  });

  check(createResponse, {
    'auction created': (r) => r.status === 201
  });

  return { auctionId: AUCTION_ID };
}

export default function(data) {
  const auctionId = data.auctionId;
  const userId = TEST_USERS[Math.floor(Math.random() * TEST_USERS.length)];

  // Generate random bid amount (between 1000 and 50000 cents)
  const amount = Math.floor(Math.random() * 49000) + 1000;

  const payload = JSON.stringify({
    user_id: userId,
    amount: amount
  });

  const params = {
    headers: {
      'Content-Type': 'application/json',
      'User-Agent': 'k6-load-test'
    },
    timeout: '10s'
  };

  const response = http.post(`${API_URL}/auction/auctions/${auctionId}/bid`, payload, params);

  const success = check(response, {
    'bid request successful': (r) => r.status === 200,
    'response time acceptable': (r) => r.timings.duration < 300,
    'no server errors': (r) => r.status < 500
  });

  bidSuccessRate.add(success);
  bidFailureRate.add(!success);

  // Small random sleep between requests (100-500ms)
  sleep(Math.random() * 0.4 + 0.1);
}

export function teardown(data) {
  // End the test auction
  const endResponse = http.post(`${API_URL}/auction/auctions/${data.auctionId}/end`);

  check(endResponse, {
    'auction ended': (r) => r.status === 200
  });

  console.log('Load test completed. Auction ended.');
}

// Additional test scenarios
export function handleSummary(data) {
  return {
    'stdout': textSummary(data, { indent: ' ', enableColors: true }),
    'infra/load-test-results.json': JSON.stringify(data),
  };
}

function textSummary(data, options) {
  const { indent = '', enableColors = false } = options || {};
  const colors = enableColors ? {
    green: '\x1b[32m',
    red: '\x1b[31m',
    yellow: '\x1b[33m',
    reset: '\x1b[0m'
  } : { green: '', red: '', yellow: '', reset: '' };

  const successRate = (data.metrics.http_req_failed.values.rate * 100).toFixed(2);
  const avgResponseTime = data.metrics.http_req_duration.values.avg.toFixed(2);
  const p95ResponseTime = data.metrics.http_req_duration.values['p(95)'].toFixed(2);

  let status = colors.green;
  if (successRate > 10) status = colors.red;
  else if (successRate > 5) status = colors.yellow;

  return `
${indent}🚀 Blytz Auction Load Test Results
${indent}====================================
${indent}
${indent}📊 Summary:
${indent}  • Total Requests: ${data.metrics.http_reqs.values.count}
${indent}  • Success Rate: ${status}${(100 - successRate).toFixed(2)}%${colors.reset}
${indent}  • Avg Response Time: ${avgResponseTime}ms
${indent}  • P95 Response Time: ${p95ResponseTime}ms
${indent}  • Data Received: ${(data.metrics.data_received.values.rate / 1024).toFixed(2)} KB/s
${indent}
${indent}🎯 Performance Targets:
${indent}  • P95 < 300ms: ${p95ResponseTime < 300 ? colors.green + '✓ PASS' + colors.reset : colors.red + '✗ FAIL' + colors.reset}
${indent}  • Error Rate < 10%: ${successRate < 10 ? colors.green + '✓ PASS' + colors.reset : colors.red + '✗ FAIL' + colors.reset}
${indent}  • Bid Success Rate > 80%: ${(data.metrics.bid_success_rate?.values?.rate * 100 || 0).toFixed(2)}%
${indent}
${indent}⏱️  Duration: ${data.state.testRunDurationMs / 1000}s
${indent}
${indent}${status}${successRate < 10 ? '🎉 Load test passed!' : '⚠️  Load test issues detected'}${colors.reset}
`;
}
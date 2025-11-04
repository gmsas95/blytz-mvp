class ApiConstants {
  // Base URLs for each microservice
  static const String gatewayBaseUrl = 'http://localhost:8080';
  static const String authBaseUrl = 'http://localhost:8084';
  static const String auctionBaseUrl = 'http://localhost:8083';
  static const String productBaseUrl = 'http://localhost:8082';
  static const String orderBaseUrl = 'http://localhost:8085';
  static const String paymentBaseUrl = 'http://localhost:8086';
  static const String chatBaseUrl = 'http://localhost:8088';
  static const String logisticsBaseUrl = 'http://localhost:8087';

  // Default base URL for general API calls
  static const String baseUrl = gatewayBaseUrl;

  // Timeout constants
  static const Duration connectTimeout = Duration(seconds: 30);
  static const Duration receiveTimeout = Duration(seconds: 30);

  // WebSocket URLs
  static const String wsAuctionUpdates = 'ws://localhost:8083/ws/auctions';
  static const String wsLiveStream = 'wss://api.blytz.app/ws/live';
  static const String wsChat = 'ws://localhost:8088/chat';

  // Auth Service Endpoints (Port 8084)
  static const String authServiceLogin = '/api/v1/auth/login';
  static const String authServiceRegister = '/api/v1/auth/register';
  static const String authServiceRefresh = '/api/v1/auth/refresh';
  static const String authServiceLogout = '/api/v1/auth/logout';
  static const String authServiceProfile = '/api/v1/auth/profile';
  static const String authServiceMe = '/api/v1/auth/me';
  static const String authServiceVerify = '/api/v1/auth/verify';

  // Product Service Endpoints (Port 8082)
  static const String products = '/api/v1/products';
  static const String productsFeatured = '/api/v1/products/featured';
  static const String productsMe = '/api/v1/products/me';
  static const String productsInventory = '/api/v1/products/inventory';

  // Order Service Endpoints (Port 8083)
  static const String orders = '/api/v1/orders';
  static const String ordersStatus = '/status';
  static const String ordersCancel = '/cancel';

  // Auction Service Endpoints (Port 8081)
  static const String auctions = '/api/v1/auctions';
  static const String auctionsActive = '/api/v1/auctions/active';
  static const String auctionBids = '/bids';

  // Payment Service Endpoints (Port 8086)
  static const String payments = '/api/v1/payments';
  static const String paymentsProcess = '/process';
  static const String paymentsMethods = '/methods';
  static const String paymentsRefund = '/refund';

  // Chat Service Endpoints (Port 8088)
  static const String chatRooms = '/api/v1/chat/rooms';
  static const String chatMessages = '/api/v1/chat/messages';

  // Health check endpoints
  static const String healthCheck = '/health';
}
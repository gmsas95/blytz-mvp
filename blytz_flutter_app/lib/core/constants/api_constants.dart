class ApiConstants {
  static const String baseUrl = 'https://api.blytz.app';
  static const String wsUrl = 'wss://api.blytz.app';
  static const String livekitUrl = 'wss://livekit.blytz.app';
  
  // API Endpoints
  static const String auth = '/api/v1/auth';
  static const String login = '$auth/login';
  static const String register = '$auth/register';
  static const String refresh = '$auth/refresh';
  static const String logout = '$auth/logout';
  
  static const String auctions = '/api/v1/auctions';
  static const String products = '/api/v1/products';
  static const String orders = '/api/v1/orders';
  static const String payments = '/api/v1/payments';
  static const String chat = '/api/v1/chat';
  static const String logistics = '/api/v1/logistics';
  
  // WebSocket endpoints
  static const String chatWs = '$chat/ws';
  static const String bidsWs = '$auctions/ws';
  
  // Timeouts
  static const Duration connectTimeout = Duration(seconds: 30);
  static const Duration receiveTimeout = Duration(seconds: 30);
  static const Duration sendTimeout = Duration(seconds: 30);
}
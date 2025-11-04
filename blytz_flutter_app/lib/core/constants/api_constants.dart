class ApiConstants {
  static const String baseUrl = 'http://localhost:8080/api/v1';
  
  // Timeout constants
  static const Duration connectTimeout = Duration(seconds: 30);
  static const Duration receiveTimeout = Duration(seconds: 30);
  
  // Auth endpoints
  static const String login = '/auth/login';
  static const String register = '/auth/register';
  static const String refresh = '/auth/refresh';
  static const String logout = '/auth/logout';
  
  // Auction endpoints
  static const String auctions = '/auctions';
}
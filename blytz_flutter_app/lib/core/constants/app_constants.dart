class AppConstants {
  static const String appName = 'Blytz';
  static const String appVersion = '1.0.0';
  
  // Storage Keys
  static const String authTokenKey = 'auth_token';
  static const String refreshTokenKey = 'refresh_token';
  static const String userKey = 'user_data';
  static const String userIdKey = 'user_id';
  static const String themeKey = 'theme_mode';
  static const String languageKey = 'language_code';
  
  // Pagination
  static const int defaultPageSize = 20;
  static const int maxPageSize = 100;
  
  // Bid limits
  static const double minBidAmount = 1;
  static const double maxBidAmount = 10000;
  static const double minBidIncrement = 1;
  
  // Image constraints
  static const int maxImageSize = 5 * 1024 * 1024; // 5MB
  static const List<String> supportedImageFormats = ['jpg', 'jpeg', 'png', 'webp'];
  
  // Rate limiting
  static const int maxBidAttemptsPerMinute = 10;
  static const int maxChatMessagesPerMinute = 30;
  
  // Cache durations
  static const Duration auctionCacheDuration = Duration(minutes: 5);
  static const Duration userCacheDuration = Duration(hours: 1);
  static const Duration imageCacheDuration = Duration(days: 7);
}
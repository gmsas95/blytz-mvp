class RouteConstants {
  static const String splash = '/splash';
  static const String onboarding = '/onboarding';
  static const String auth = '/auth';
  static const String login = '/auth/login';
  static const String register = '/auth/register';
  static const String forgotPassword = '/auth/forgot-password';
  
  static const String home = '/home';
  static const String auctions = '/home/auctions';
  static const String auctionDetail = '/home/auctions/:id';
  static const String liveAuction = '/live/:auctionId';
  static const String search = '/home/search';
  static const String categories = '/home/categories';
  
  static const String profile = '/profile';
  static const String editProfile = '/profile/edit';
  static const String settings = '/profile/settings';
  static const String notifications = '/profile/notifications';
  
  static const String orders = '/orders';
  static const String orderDetail = '/orders/:id';
  
  static const String payments = '/payments';
  static const String paymentMethods = '/payments/methods';
  static const String paymentHistory = '/payments/history';
  
  static const String chat = '/chat';
  static const String chatRoom = '/chat/:roomId';
  
  static const String createAuction = '/auctions/create';
  static const String editAuction = '/auctions/:id/edit';
}
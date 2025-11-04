class RouteConstants {
  // Core
  static const String splash = '/splash';
  static const String onboarding = '/onboarding';

  // Authentication
  static const String auth = '/auth';
  static const String login = '/auth/login';
  static const String register = '/auth/register';
  static const String forgotPassword = '/auth/forgot-password';

  // Home & Discovery
  static const String home = '/home';
  static const String auctions = '/home/auctions';
  static const String auctionDetail = '/home/auctions/:id';
  static const String liveAuction = '/live/:auctionId';
  static const String liveStream = '/live/:streamId';
  static const String search = '/home/search';
  static const String categories = '/home/categories';
  static const String discovery = '/discovery';
  static const String categoryBrowse = '/categories/:category';

  // Community
  static const String community = '/community';
  static const String communityFeed = '/community/feed';
  static const String communitySellers = '/community/sellers';
  static const String communityForums = '/community/forums';
  static const String communityEvents = '/community/events';

  // Profile & User
  static const String profile = '/profile';
  static const String editProfile = '/profile/edit';
  static const String settings = '/profile/settings';
  static const String notifications = '/profile/notifications';
  static const String publicProfile = '/profile/:userId';

  // Orders & History
  static const String orders = '/orders';
  static const String orderDetail = '/orders/:id';
  static const String orderHistory = '/orders/history';
  static const String orderTracking = '/orders/:id/tracking';

  // Payments & Checkout
  static const String payments = '/payments';
  static const String paymentMethods = '/payments/methods';
  static const String paymentHistory = '/payments/history';
  static const String checkout = '/checkout';
  static const String checkoutSuccess = '/checkout/success';
  static const String checkoutFailed = '/checkout/failed';

  // Seller Pages
  static const String sellerDashboard = '/seller/dashboard';
  static const String sellerProfile = '/seller/:sellerId';
  static const String createStream = '/seller/stream/create';
  static const String sellerAnalytics = '/seller/analytics';
  static const String sellerOrders = '/seller/orders';
  static const String sellerProducts = '/seller/products';
  static const String createProduct = '/seller/products/create';
  static const String editProduct = '/seller/products/:id/edit';

  // Chat & Communication
  static const String chat = '/chat';
  static const String chatRoom = '/chat/:roomId';
  static const String chatList = '/chat/list';
  static const String sellerMessages = '/seller/messages';

  // Wishlist & Favorites
  static const String wishlist = '/wishlist';
  static const String favorites = '/favorites';
  static const String watchlist = '/watchlist';

  // Help & Support
  static const String help = '/help';
  static const String faq = '/help/faq';
  static const String contact = '/help/contact';
  static const String terms = '/help/terms';
  static const String privacy = '/help/privacy';

  // Miscellaneous
  static const String about = '/about';
  static const String careers = '/careers';
  static const String press = '/press';
}
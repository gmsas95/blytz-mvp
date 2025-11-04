import 'package:blytz_flutter_app/core/constants/route_constants.dart';
import 'package:blytz_flutter_app/features/auth/presentation/pages/login_page.dart';
import 'package:blytz_flutter_app/features/community/presentation/pages/community_page.dart';
import 'package:blytz_flutter_app/features/discovery/presentation/pages/categories_page.dart';
import 'package:blytz_flutter_app/features/discovery/presentation/pages/discovery_page.dart';
import 'package:blytz_flutter_app/features/home/presentation/pages/home_page.dart';
import 'package:blytz_flutter_app/features/livestream/presentation/pages/live_stream_page.dart';
import 'package:blytz_flutter_app/features/onboarding/presentation/pages/onboarding_page.dart';
import 'package:blytz_flutter_app/features/payments/presentation/pages/checkout_page.dart';
import 'package:blytz_flutter_app/features/profile/presentation/pages/enhanced_profile_page.dart';
import 'package:blytz_flutter_app/features/seller/presentation/pages/create_stream_page.dart';
import 'package:blytz_flutter_app/features/seller/presentation/pages/seller_dashboard_page.dart';
import 'package:blytz_flutter_app/shared/pages/not_found_page.dart';
import 'package:blytz_flutter_app/shared/pages/splash_page.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

// Router provider
final appRouterProvider = Provider<GoRouter>((ref) {
  return GoRouter(
    initialLocation: '/splash',
    debugLogDiagnostics: true,
    routes: [
      // Splash
      GoRoute(
        path: RouteConstants.splash,
        builder: (context, state) => const SplashPage(),
      ),

      // Auth routes
      GoRoute(
        path: RouteConstants.login,
        builder: (context, state) => const LoginPage(),
      ),

      // Main navigation
      GoRoute(
        path: RouteConstants.home,
        builder: (context, state) => const HomePage(),
      ),

      // Discovery
      GoRoute(
        path: RouteConstants.discovery,
        builder: (context, state) => const DiscoveryPage(),
      ),

      // Categories
      GoRoute(
        path: RouteConstants.categories,
        builder: (context, state) => const CategoriesPage(),
      ),

      // Community
      GoRoute(
        path: RouteConstants.community,
        builder: (context, state) => const CommunityPage(),
      ),

      // Profile
      GoRoute(
        path: RouteConstants.profile,
        builder: (context, state) => const EnhancedProfilePage(),
      ),

      // Live Stream
      GoRoute(
        path: RouteConstants.liveStream,
        builder: (context, state) {
          final streamId = state.pathParameters['streamId'] ?? '1';
          return LiveStreamPage(
            streamId: streamId,
            sellerName: 'TechCollector',
            productTitle: 'Vintage Electronics Auction',
          );
        },
      ),

      // Seller Dashboard
      GoRoute(
        path: RouteConstants.sellerDashboard,
        builder: (context, state) => const SellerDashboardPage(),
      ),

      // Create Stream
      GoRoute(
        path: RouteConstants.createStream,
        builder: (context, state) => const CreateStreamPage(),
      ),

      // Checkout
      GoRoute(
        path: RouteConstants.checkout,
        builder: (context, state) {
          // For demo purposes, use mock data
          return const CheckoutPage(
            productTitle: 'Vintage Watch',
            productImage: 'https://example.com/watch.jpg',
            price: 250,
            auctionId: 'auction_123',
          );
        },
      ),

      // Onboarding
      GoRoute(
        path: RouteConstants.onboarding,
        builder: (context, state) => const OnboardingPage(),
      ),
    ],
    errorBuilder: (context, state) => NotFoundPage(error: state.error.toString()),
  );
});
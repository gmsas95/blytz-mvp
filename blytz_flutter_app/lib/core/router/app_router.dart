import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../core/constants/route_constants.dart';
import '../core/providers/app_providers.dart';
import '../features/auth/presentation/pages/login_page.dart';
import '../features/auth/presentation/pages/register_page.dart';
import '../features/auth/presentation/pages/forgot_password_page.dart';
import '../features/onboarding/presentation/pages/onboarding_page.dart';
import '../features/home/presentation/pages/home_page.dart';
import '../features/auctions/presentation/pages/auctions_page.dart';
import '../features/auctions/presentation/pages/auction_detail_page.dart';
import '../features/auctions/presentation/pages/live_auction_page.dart';
import '../features/auctions/presentation/pages/create_auction_page.dart';
import '../features/profile/presentation/pages/profile_page.dart';
import '../features/profile/presentation/pages/edit_profile_page.dart';
import '../features/profile/presentation/pages/settings_page.dart';
import '../features/orders/presentation/pages/orders_page.dart';
import '../features/orders/presentation/pages/order_detail_page.dart';
import '../features/payments/presentation/pages/payments_page.dart';
import '../features/payments/presentation/pages/payment_methods_page.dart';
import '../features/chat/presentation/pages/chat_page.dart';
import '../features/chat/presentation/pages/chat_room_page.dart';
import '../shared/pages/splash_page.dart';
import '../shared/pages/not_found_page.dart';

final routerProvider = Provider<GoRouter>((ref) {
  final authState = ref.watch(authStateProvider);
  final isOnboardingCompleted = ref.watch(appPreferencesProvider
      .select((prefs) => prefs.isOnboardingCompleted()));

  return GoRouter(
    initialLocation: RouteConstants.splash,
    debugLogDiagnostics: true,
    redirect: (context, state) {
      final isAuthenticated = authState.isAuthenticated;
      final isAuthRoute = state.location.startsWith('/auth');
      final isOnboardingRoute = state.location == RouteConstants.onboarding;
      final isSplashRoute = state.location == RouteConstants.splash;

      // If not authenticated and not on auth/onboarding/splash route, redirect to login
      if (!isAuthenticated && !isAuthRoute && !isOnboardingRoute && !isSplashRoute) {
        return RouteConstants.login;
      }

      // If authenticated and on auth route, redirect to home
      if (isAuthenticated && isAuthRoute) {
        return RouteConstants.home;
      }

      // If authenticated and on onboarding route, redirect to home
      if (isAuthenticated && isOnboardingRoute) {
        return RouteConstants.home;
      }

      return null;
    },
    routes: [
      // Splash
      GoRoute(
        path: RouteConstants.splash,
        name: 'splash',
        builder: (context, state) => const SplashPage(),
      ),

      // Onboarding
      GoRoute(
        path: RouteConstants.onboarding,
        name: 'onboarding',
        builder: (context, state) => const OnboardingPage(),
      ),

      // Authentication
      GoRoute(
        path: RouteConstants.auth,
        name: 'auth',
        builder: (context, state) => const LoginPage(),
        routes: [
          GoRoute(
            path: RouteConstants.login.split('/').last,
            name: 'login',
            builder: (context, state) => const LoginPage(),
          ),
          GoRoute(
            path: RouteConstants.register.split('/').last,
            name: 'register',
            builder: (context, state) => const RegisterPage(),
          ),
          GoRoute(
            path: RouteConstants.forgotPassword.split('/').last,
            name: 'forgot-password',
            builder: (context, state) => const ForgotPasswordPage(),
          ),
        ],
      ),

      // Main App
      GoRoute(
        path: RouteConstants.home,
        name: 'home',
        builder: (context, state) => const HomePage(),
        routes: [
          // Auctions
          GoRoute(
            path: RouteConstants.auctions.split('/').last,
            name: 'auctions',
            builder: (context, state) => const AuctionsPage(),
            routes: [
              GoRoute(
                path: ':id',
                name: 'auction-detail',
                builder: (context, state) {
                  final auctionId = state.pathParameters['id']!;
                  return AuctionDetailPage(auctionId: auctionId);
                },
              ),
            ],
          ),

          // Search
          GoRoute(
            path: RouteConstants.search.split('/').last,
            name: 'search',
            builder: (context, state) {
              final query = state.uri.queryParameters['q'];
              return SearchPage(initialQuery: query);
            },
          ),

          // Categories
          GoRoute(
            path: RouteConstants.categories.split('/').last,
            name: 'categories',
            builder: (context, state) => const CategoriesPage(),
          ),
        ],
      ),

      // Live Auction
      GoRoute(
        path: RouteConstants.liveAuction,
        name: 'live-auction',
        builder: (context, state) {
          final auctionId = state.pathParameters['auctionId']!;
          return LiveAuctionPage(auctionId: auctionId);
        },
      ),

      // Profile
      GoRoute(
        path: RouteConstants.profile,
        name: 'profile',
        builder: (context, state) => const ProfilePage(),
        routes: [
          GoRoute(
            path: RouteConstants.editProfile.split('/').last,
            name: 'edit-profile',
            builder: (context, state) => const EditProfilePage(),
          ),
          GoRoute(
            path: RouteConstants.settings.split('/').last,
            name: 'settings',
            builder: (context, state) => const SettingsPage(),
          ),
          GoRoute(
            path: RouteConstants.notifications.split('/').last,
            name: 'notifications',
            builder: (context, state) => const NotificationsPage(),
          ),
        ],
      ),

      // Orders
      GoRoute(
        path: RouteConstants.orders,
        name: 'orders',
        builder: (context, state) => const OrdersPage(),
        routes: [
          GoRoute(
            path: ':id',
            name: 'order-detail',
            builder: (context, state) {
              final orderId = state.pathParameters['id']!;
              return OrderDetailPage(orderId: orderId);
            },
          ),
        ],
      ),

      // Payments
      GoRoute(
        path: RouteConstants.payments,
        name: 'payments',
        builder: (context, state) => const PaymentsPage(),
        routes: [
          GoRoute(
            path: RouteConstants.paymentMethods.split('/').last,
            name: 'payment-methods',
            builder: (context, state) => const PaymentMethodsPage(),
          ),
          GoRoute(
            path: RouteConstants.paymentHistory.split('/').last,
            name: 'payment-history',
            builder: (context, state) => const PaymentHistoryPage(),
          ),
        ],
      ),

      // Chat
      GoRoute(
        path: RouteConstants.chat,
        name: 'chat',
        builder: (context, state) => const ChatPage(),
        routes: [
          GoRoute(
            path: ':roomId',
            name: 'chat-room',
            builder: (context, state) {
              final roomId = state.pathParameters['roomId']!;
              return ChatRoomPage(roomId: roomId);
            },
          ),
        ],
      ),

      // Create Auction
      GoRoute(
        path: RouteConstants.createAuction,
        name: 'create-auction',
        builder: (context, state) => const CreateAuctionPage(),
      ),

      // Edit Auction
      GoRoute(
        path: RouteConstants.editAuction,
        name: 'edit-auction',
        builder: (context, state) {
          final auctionId = state.pathParameters['id']!;
          return EditAuctionPage(auctionId: auctionId);
        },
      ),
    ],

    errorBuilder: (context, state) => NotFoundPage(error: state.error),
  );
});

// Navigation Service
class NavigationService {
  static GlobalKey<NavigatorState> navigatorKey = GlobalKey<NavigatorState>();

  static BuildContext? get context => navigatorKey.currentContext;

  static void push(String route, {Object? extra}) {
    context?.push(route, extra: extra);
  }

  static void pushReplacement(String route, {Object? extra}) {
    context?.pushReplacement(route, extra: extra);
  }

  static void pop<T>([T? result]) {
    context?.pop<T>(result);
  }

  static void pushNamed(String name, {Map<String, String>? pathParameters, Object? extra}) {
    context?.pushNamed(name, pathParameters: pathParameters, extra: extra);
  }

  static void pushReplacementNamed(String name, {Map<String, String>? pathParameters, Object? extra}) {
    context?.pushReplacementNamed(name, pathParameters: pathParameters, extra: extra);
  }

  static bool canPop() {
    return context?.canPop() ?? false;
  }

  static void go(String route, {Object? extra}) {
    context?.go(route, extra: extra);
  }

  static void goNamed(String name, {Map<String, String>? pathParameters, Object? extra}) {
    context?.goNamed(name, pathParameters: pathParameters, extra: extra);
  }
}

// Route Guards
class RouteGuard {
  static bool isAuthenticated(WidgetRef ref) {
    return ref.read(authStateProvider).isAuthenticated;
  }

  static bool isSeller(WidgetRef ref) {
    final user = ref.read(currentUserProvider);
    return user?.isSeller ?? false;
  }

  static bool isOnboardingCompleted(WidgetRef ref) {
    // This would check from preferences
    return true; // Placeholder
  }

  static String? redirectIfNotAuthenticated(WidgetRef ref, GoRouterState state) {
    if (!isAuthenticated(ref)) {
      return RouteConstants.login;
    }
    return null;
  }

  static String? redirectIfNotSeller(WidgetRef ref, GoRouterState state) {
    if (!isSeller(ref)) {
      return RouteConstants.home;
    }
    return null;
  }
}

// Import placeholder pages (these would be created separately)
class SearchPage extends StatelessWidget {
  final String? initialQuery;
  
  const SearchPage({super.key, this.initialQuery});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Search')),
      body: Center(child: Text('Search Page: $initialQuery')),
    );
  }
}

class CategoriesPage extends StatelessWidget {
  const CategoriesPage({super.key});

  @override
  Widget build(BuildContext context) {
    return const Scaffold(
      appBar: AppBar(title: Text('Categories')),
      body: Center(child: Text('Categories Page')),
    );
  }
}

class NotificationsPage extends StatelessWidget {
  const NotificationsPage({super.key});

  @override
  Widget build(BuildContext context) {
    return const Scaffold(
      appBar: AppBar(title: Text('Notifications')),
      body: Center(child: Text('Notifications Page')),
    );
  }
}

class PaymentHistoryPage extends StatelessWidget {
  const PaymentHistoryPage({super.key});

  @override
  Widget build(BuildContext context) {
    return const Scaffold(
      appBar: AppBar(title: Text('Payment History')),
      body: Center(child: Text('Payment History Page')),
    );
  }
}

class EditAuctionPage extends StatelessWidget {
  final String auctionId;
  
  const EditAuctionPage({super.key, required this.auctionId});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: Text('Edit Auction: $auctionId')),
      body: Center(child: Text('Edit Auction Page: $auctionId')),
    );
  }
}
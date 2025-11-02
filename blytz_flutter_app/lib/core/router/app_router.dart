import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../features/auth/presentation/pages/login_page.dart';
import '../../features/home/presentation/pages/home_page.dart';
import '../../shared/pages/splash_page.dart';
import '../../shared/pages/not_found_page.dart';

// Router provider
final appRouterProvider = Provider<GoRouter>((ref) {
  return GoRouter(
    initialLocation: '/splash',
    debugLogDiagnostics: true,
    routes: [
      // Splash
      GoRoute(
        path: '/splash',
        builder: (context, state) => const SplashPage(),
      ),
      
      // Auth routes
      GoRoute(
        path: '/login',
        builder: (context, state) => const LoginPage(),
      ),
      
      // Home
      GoRoute(
        path: '/',
        builder: (context, state) => const HomePage(),
      ),
      
      // Onboarding
      GoRoute(
        path: '/onboarding',
        builder: (context, state) => const SplashPage(), // Reuse for now
      ),
    ],
    errorBuilder: (context, state) => NotFoundPage(error: state.error.toString()),
  );
});
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'core/router/app_router.dart';
import 'shared/themes/app_theme.dart';
import 'shared/widgets/loading/app_loading_indicator.dart';
import 'features/auth/presentation/pages/login_page.dart';
import 'features/auth/presentation/pages/register_page.dart';
import 'features/auth/presentation/pages/forgot_password_page.dart';
import 'features/onboarding/presentation/pages/onboarding_page.dart';
import 'features/home/presentation/pages/home_page.dart';
import 'features/auctions/presentation/pages/auctions_page.dart';
import 'features/profile/presentation/pages/profile_page.dart';
import 'features/orders/presentation/pages/orders_page.dart';
import 'features/payments/presentation/pages/payments_page.dart';
import 'features/chat/presentation/pages/chat_page.dart';
import 'shared/pages/splash_page.dart';
import 'shared/pages/not_found_page.dart';

void main() {
  runApp(const ProviderScope(child: BlytzApp()));
}

class BlytzApp extends ConsumerWidget {
  const BlytzApp({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final router = ref.watch(appRouterProvider);
    
    return MaterialApp.router(
      title: 'Blytz Live Auction',
      debugShowCheckedModeBanner: false,
      theme: AppTheme.lightTheme,
      darkTheme: AppTheme.darkTheme,
      themeMode: ThemeMode.system,
      routerConfig: router,
    );
  }
}
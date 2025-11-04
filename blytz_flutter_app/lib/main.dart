import 'package:blytz_flutter_app/core/router/app_router.dart';
import 'package:blytz_flutter_app/shared/themes/app_theme.dart';
import 'package:flutter/material.dart';
import 'package:flutter_platform_widgets/flutter_platform_widgets.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

void main() {
  runApp(const ProviderScope(child: BlytzApp()));
}

class BlytzApp extends ConsumerWidget {
  const BlytzApp({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final router = ref.watch(appRouterProvider);

    return PlatformApp.router(
      title: 'Blytz Live Auction',
      debugShowCheckedModeBanner: false,
      routerConfig: router,
      material: (_, __) => MaterialAppRouterData(
        theme: AppTheme.lightTheme,
        darkTheme: AppTheme.darkTheme,
      ),
      cupertino: (_, __) => CupertinoAppRouterData(),
    );
  }
}

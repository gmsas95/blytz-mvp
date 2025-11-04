import 'package:blytz_flutter_app/core/router/app_router.dart';
import 'package:blytz_flutter_app/core/di/service_locator.dart';
import 'package:blytz_flutter_app/shared/themes/app_theme.dart';
import 'package:flutter/material.dart';
import 'package:flutter_platform_widgets/flutter_platform_widgets.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

void main() async {
  WidgetsFlutterBinding.ensureInitialized();

  // Initialize services
  final container = ProviderContainer();
  await initializeServices(container);

  runApp(
    UncontrolledProviderScope(
      container: container,
      child: const BlytzApp(),
    ),
  );
}

class BlytzApp extends ConsumerWidget {
  const BlytzApp({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final router = ref.watch(appRouterProvider);

    // Listen to connectivity changes
    final connectivity = ref.watch(connectivityProvider);

    return PlatformApp.router(
      title: 'Blytz Live Auction',
      debugShowCheckedModeBanner: false,
      routerConfig: router,
      material: (_, __) => MaterialAppRouterData(
        theme: AppTheme.lightTheme,
        darkTheme: AppTheme.darkTheme,
        banner: connectivity.when(
          data: (isConnected) => !isConnected
              ? const MaterialBanner(
                  content: Text('No internet connection'),
                  actions: [SizedBox.shrink()],
                  backgroundColor: Colors.red,
                )
              : null,
          loading: () => null,
          error: (_, __) => null,
        ),
      ),
      cupertino: (_, __) => CupertinoAppRouterData(),
    );
  }
}

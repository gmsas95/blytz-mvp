import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'core/router/app_router.dart';
import 'core/providers/app_providers.dart';
import 'core/themes/app_theme.dart';
import 'shared/widgets/common/loading/app_loading_indicator.dart';

void main() async {
  WidgetsFlutterBinding.ensureInitialized();
  
  await SystemChrome.setPreferredOrientations([
    DeviceOrientation.portraitUp,
    DeviceOrientation.portraitDown,
  ]);

  runApp(
    const ProviderScope(
      child: BlytzApp(),
    ),
  );
}

class BlytzApp extends ConsumerWidget {
  const BlytzApp({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final appInitialization = ref.watch(appInitializationProvider);
    final themeState = ref.watch(themeProvider);
    final router = ref.watch(routerProvider);

    return appInitialization.when(
      data: (_) {
        return MaterialApp.router(
          title: 'Blytz - Live Auction Platform',
          debugShowCheckedModeBanner: false,
          theme: AppTheme.lightTheme,
          darkTheme: AppTheme.darkTheme,
          themeMode: themeState.theme == 'dark' ? ThemeMode.dark : ThemeMode.light,
          routerConfig: router,
          builder: (context, child) {
            return MediaQuery(
              data: MediaQuery.of(context).copyWith(
                textScaleFactor: 1.0, // Prevent text scaling
              ),
              child: child!,
            );
          },
        );
      },
      loading: () {
        return MaterialApp(
          title: 'Blytz - Live Auction Platform',
          debugShowCheckedModeBanner: false,
          home: const Scaffold(
            backgroundColor: Color(0xFF1E3A8A),
            body: Center(
              child: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  AppLoadingIndicator(),
                  SizedBox(height: 24),
                  Text(
                    'Blytz',
                    style: TextStyle(
                      color: Colors.white,
                      fontSize: 32,
                      fontWeight: FontWeight.bold,
                      letterSpacing: 2,
                    ),
                  ),
                  SizedBox(height: 8),
                  Text(
                    'Live Auction Platform',
                    style: TextStyle(
                      color: Colors.white70,
                      fontSize: 16,
                    ),
                  ),
                ],
              ),
            ),
          ),
        );
      },
      error: (error, stack) {
        return MaterialApp(
          title: 'Blytz - Error',
          debugShowCheckedModeBanner: false,
          home: Scaffold(
            backgroundColor: const Color(0xFF1E3A8A),
            body: Center(
              child: Padding(
                padding: const EdgeInsets.all(24.0),
                child: Column(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    const Icon(
                      Icons.error_outline,
                      color: Colors.red,
                      size: 64,
                    ),
                    const SizedBox(height: 16),
                    const Text(
                      'Initialization Error',
                      style: TextStyle(
                        color: Colors.white,
                        fontSize: 24,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                    const SizedBox(height: 8),
                    const Text(
                      'Failed to initialize the app. Please restart the application.',
                      style: TextStyle(
                        color: Colors.white70,
                        fontSize: 16,
                      ),
                      textAlign: TextAlign.center,
                    ),
                    const SizedBox(height: 24),
                    ElevatedButton(
                      onPressed: () {
                        // Restart the app
                        SystemNavigator.pop();
                      },
                      style: ElevatedButton.styleFrom(
                        backgroundColor: Colors.white,
                        foregroundColor: const Color(0xFF1E3A8A),
                        padding: const EdgeInsets.symmetric(
                          horizontal: 32,
                          vertical: 16,
                        ),
                      ),
                      child: const Text('Restart App'),
                    ),
                  ],
                ),
              ),
            ),
          ),
        );
      },
    );
  }
}
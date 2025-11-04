import 'package:blytz_flutter_app/core/network/api_client.dart';
import 'package:blytz_flutter_app/core/network/interceptors.dart';
import 'package:blytz_flutter_app/core/storage/secure_storage.dart';
import 'package:blytz_flutter_app/core/services/websocket_service.dart';
import 'package:blytz_flutter_app/features/auth/data/repositories/auth_repository_impl.dart';
import 'package:blytz_flutter_app/features/auction/data/repositories/auction_repository_impl.dart';
import 'package:blytz_flutter_app/core/network/network_info.dart';
import 'package:blytz_flutter_app/core/utils/logger.dart';
import 'package:dio/dio.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:connectivity_plus/connectivity_plus.dart';

// Core Services
final dioProvider = Provider<Dio>((ref) {
  final dio = Dio();

  // Configure base options
  dio.options = BaseOptions(
    connectTimeout: const Duration(seconds: 30),
    receiveTimeout: const Duration(seconds: 30),
    headers: {
      'Content-Type': 'application/json',
      'Accept': 'application/json',
    },
  );

  // Add interceptors
  dio.interceptors.add(LogInterceptor(
    requestBody: true,
    responseBody: true,
    logPrint: (obj) => AppLogger.debug(obj.toString()),
  ));

  dio.interceptors.add(AuthInterceptor());

  return dio;
});

final apiClientProvider = Provider<ApiClient>((ref) {
  return ApiClient(ref.watch(dioProvider));
});

final secureStorageProvider = Provider<SecureStorage>((ref) {
  return SecureStorage();
});

final networkInfoProvider = Provider<NetworkInfo>((ref) {
  return NetworkInfoImpl(Connectivity());
});

final loggerProvider = Provider<AppLogger>((ref) {
  return AppLogger();
});

final webSocketServiceProvider = Provider<WebSocketService>((ref) {
  return WebSocketService();
});

// Repository Providers
final authRepositoryProvider = Provider<AuthRepositoryImpl>((ref) {
  return AuthRepositoryImpl(
    apiClient: ref.watch(apiClientProvider),
    secureStorage: ref.watch(secureStorageProvider),
    dio: ref.watch(dioProvider),
  );
});

final auctionRepositoryProvider = Provider<AuctionRepositoryImpl>((ref) {
  return AuctionRepositoryImpl(
    apiClient: ref.watch(apiClientProvider),
  );
});

// Service Providers (to be added later)
// final productRepositoryProvider = Provider<ProductRepositoryImpl>((ref) {
//   return ProductRepositoryImpl(
//     apiClient: ref.watch(apiClientProvider),
//   );
// });

// final orderRepositoryProvider = Provider<OrderRepositoryImpl>((ref) {
//   return OrderRepositoryImpl(
//     apiClient: ref.watch(apiClientProvider),
//   );
// });

// final paymentRepositoryProvider = Provider<PaymentRepositoryImpl>((ref) {
//   return PaymentRepositoryImpl(
//     apiClient: ref.watch(apiClientProvider),
//   );
// });

// Utility Providers
final connectivityProvider = StreamProvider<bool>((ref) {
  return ref.watch(networkInfoProvider).connection;
});

// Current user provider (can be used to track authentication state)
final currentUserProvider = FutureProvider((ref) async {
  final authRepo = ref.watch(authRepositoryProvider);
  final tokenResult = await authRepo.getStoredToken();

  if (tokenResult.isLeft || tokenResult.right == null) {
    return null;
  }

  // Get current user if token exists
  final userResult = await authRepo.getCurrentUser();
  return userResult.fold(
    (failure) => null,
    (user) => user,
  );
});

// Initialize services
Future<void> initializeServices(ProviderContainer container) async {
  // Initialize any async services here
  final logger = container.read(loggerProvider);
  logger.log('Initializing services...');

  // Test connectivity
  final networkInfo = container.read(networkInfoProvider);
  final isConnected = await networkInfo.isConnected;
  logger.log('Network connected: $isConnected');

  // Check stored authentication
  final authRepo = container.read(authRepositoryProvider);
  final tokenValid = await authRepo.isTokenValid();
  logger.log('Token valid: ${tokenValid.fold((l) => false, (r) => r)}');

  logger.log('Services initialized successfully');
}
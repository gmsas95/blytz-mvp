import 'dart:convert';

import 'package:blytz_flutter_app/core/errors/exceptions.dart';
import 'package:blytz_flutter_app/core/network/api_client.dart';
import 'package:blytz_flutter_app/core/network/interceptors.dart';
import 'package:blytz_flutter_app/core/network/network_info.dart';
import 'package:blytz_flutter_app/core/storage/local_database.dart';
import 'package:blytz_flutter_app/core/storage/preferences.dart';
import 'package:blytz_flutter_app/core/storage/secure_storage.dart';
import 'package:blytz_flutter_app/core/utils/logger.dart';
import 'package:blytz_flutter_app/features/auth/data/models/auth_model.dart';
import 'package:connectivity_plus/connectivity_plus.dart';
import 'package:dio/dio.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

// API Client Provider
final apiClientProvider = Provider<ApiClient>((ref) {
  final dio = Dio();
  
  // Configure base options
  dio.options = BaseOptions(
    baseUrl: 'https://api.blytz.app',
    connectTimeout: const Duration(seconds: 30),
    receiveTimeout: const Duration(seconds: 30),
    sendTimeout: const Duration(seconds: 30),
  );

  // Add interceptors
  dio.interceptors.add(AuthInterceptor());
  
  if (true) { // kDebugMode
    dio.interceptors.add(LoggingInterceptor());
  }

  return ApiClient(dio);
});

// Network Info Provider
final networkInfoProvider = Provider<NetworkInfo>((ref) {
  return NetworkInfo(Connectivity());
});

// Secure Storage Provider
final secureStorageProvider = Provider<SecureStorage>((ref) {
  return SecureStorage();
});

// Local Database Provider
final localDatabaseProvider = Provider<LocalDatabase>((ref) {
  return LocalDatabase();
});

// App Preferences Provider
final appPreferencesProvider = Provider<AppPreferences>((ref) {
  return AppPreferences();
});

// Logger Provider
final loggerProvider = Provider<AppLogger>((ref) {
  return AppLogger();
});

// Authentication State Provider
final authStateProvider = StateNotifierProvider<AuthNotifier, AuthState>((ref) {
  return AuthNotifier(
    ref.read(apiClientProvider),
  );
});

// Current User Provider
final currentUserProvider = Provider<AuthModel?>((ref) {
  return ref.watch(authStateProvider).user;
});

// Is Authenticated Provider
final isAuthenticatedProvider = Provider<bool>((ref) {
  return ref.watch(authStateProvider).isAuthenticated;
});

// Theme Provider
final themeProvider = StateNotifierProvider<ThemeNotifier, ThemeState>((ref) {
  return ThemeNotifier();
});

// Language Provider
final languageProvider = StateNotifierProvider<LanguageNotifier, String>((ref) {
  return LanguageNotifier();
});

// Connectivity Provider
final connectivityProvider = StreamProvider<List<ConnectivityResult>>((ref) {
  return Connectivity().onConnectivityChanged;
});

// Is Online Provider
final isOnlineProvider = Provider<bool>((ref) {
  final connectivity = ref.watch(connectivityProvider);
  return connectivity.when(
    data: (results) => results.isNotEmpty && !results.contains(ConnectivityResult.none),
    loading: () => true,
    error: (_, __) => false,
  );
});

// Notification Settings Provider
final notificationSettingsProvider = StateNotifierProvider<NotificationSettingsNotifier, NotificationSettings>((ref) {
  return NotificationSettingsNotifier();
});

// App Initialization Provider
final appInitializationProvider = FutureProvider<void>((ref) async {
  try {
    // Initialize local database
    await LocalDatabase.init();
    
    // Check authentication status
    await ref.read(authStateProvider.notifier).checkAuthStatus();
    
    // Load preferences
    await ref.read(themeProvider.notifier).loadTheme();
    await ref.read(languageProvider.notifier).loadLanguage();
    await ref.read(notificationSettingsProvider.notifier).loadSettings();
    
    AppLogger.info('App initialized successfully');
  } catch (e) {
    AppLogger.error('Failed to initialize app', error: e);
    rethrow;
  }
});

// Error Handler Provider
final errorHandlerProvider = Provider<ErrorHandler>((ref) {
  return ErrorHandler(ref.read(loggerProvider));
});

class AuthNotifier extends StateNotifier<AuthState> {

  AuthNotifier(
    this._apiClient,
  ) : super(AuthState.initial());
  final ApiClient _apiClient;

  Future<void> login(String email, String password) async {
    state = state.copyWith(isLoading: true);
    
    try {
      final response = await _apiClient.login(
        LoginRequest(email: email, password: password),
      );

      await SecureStorage.storeToken(response.accessToken);
      await SecureStorage.storeRefreshToken(response.refreshToken);
      await SecureStorage.storeUser(jsonEncode(response.user.toJson()));
      await AppPreferences.setOnboardingCompleted(true);

      state = state.copyWith(
        isLoading: false,
        isAuthenticated: true,
        user: response.user,
      );

      AppLogger.info('User logged in successfully: ${response.user.email}');
    } catch (e) {
      final error = ErrorHandler(AppLogger()).handleException(e as Exception);
      state = state.copyWith(isLoading: false, error: error);
      AppLogger.error('Login failed', error: e);
    }
  }

  Future<void> register(String email, String password, String firstName, String lastName) async {
    state = state.copyWith(isLoading: true);
    
    try {
      final response = await _apiClient.register(
        RegisterRequest(
          email: email,
          password: password,
          firstName: firstName,
          lastName: lastName,
        ),
      );

      await SecureStorage.storeToken(response.accessToken);
      await SecureStorage.storeRefreshToken(response.refreshToken);
      await SecureStorage.storeUser(jsonEncode(response.user.toJson()));
      await AppPreferences.setOnboardingCompleted(true);

      state = state.copyWith(
        isLoading: false,
        isAuthenticated: true,
        user: response.user,
      );

      AppLogger.info('User registered successfully: ${response.user.email}');
    } catch (e) {
      final error = ErrorHandler(AppLogger()).handleException(e as Exception);
      state = state.copyWith(isLoading: false, error: error);
      AppLogger.error('Registration failed', error: e);
    }
  }

  Future<void> logout() async {
    try {
      await _apiClient.logout();
    } catch (e) {
      AppLogger.warning('Logout API call failed', error: e);
    }

    await SecureStorage.clearAll();
    state = AuthState.initial();
    
    AppLogger.info('User logged out');
  }

  Future<void> checkAuthStatus() async {
    final hasToken = await SecureStorage.hasToken();
    if (!hasToken) {
      state = AuthState.initial();
      return;
    }

    try {
      final userJson = await SecureStorage.getUser();
      if (userJson != null) {
        final user = AuthModel.fromJson(jsonDecode(userJson) as Map<String, dynamic>);
        state = state.copyWith(
          isAuthenticated: true,
          user: user,
        );
      }
    } catch (e) {
      AppLogger.error('Failed to check auth status', error: e);
      await logout();
    }
  }

  void clearError() {
    state = state.copyWith();
  }
}

class ThemeNotifier extends StateNotifier<ThemeState> {

  ThemeNotifier() : super(const ThemeState());

  Future<void> loadTheme() async {
    final theme = await AppPreferences.getTheme();
    state = ThemeState(theme: theme);
  }

  Future<void> setTheme(String theme) async {
    await AppPreferences.setTheme(theme);
    state = ThemeState(theme: theme);
  }
}

class LanguageNotifier extends StateNotifier<String> {

  LanguageNotifier() : super('en');

  Future<void> loadLanguage() async {
    final language = await AppPreferences.getLanguage();
    state = language;
  }

  Future<void> setLanguage(String language) async {
    await AppPreferences.setLanguage(language);
    state = language;
  }
}

class NotificationSettingsNotifier extends StateNotifier<NotificationSettings> {

  NotificationSettingsNotifier() : super(const NotificationSettings());

  Future<void> loadSettings() async {
    final enabled = await AppPreferences.isNotificationsEnabled();
    final autoBid = await AppPreferences.isAutoBidEnabled();
    final maxAutoBid = await AppPreferences.getMaxAutoBidAmount();

    state = NotificationSettings(
      enabled: enabled,
      autoBid: autoBid,
      maxAutoBidAmount: maxAutoBid,
    );
  }

  Future<void> setNotificationsEnabled(bool enabled) async {
    await AppPreferences.setNotificationsEnabled(enabled);
    state = state.copyWith(enabled: enabled);
  }

  Future<void> setAutoBidEnabled(bool enabled) async {
    await AppPreferences.setAutoBidEnabled(enabled);
    state = state.copyWith(autoBid: enabled);
  }

  Future<void> setMaxAutoBidAmount(double amount) async {
    await AppPreferences.setMaxAutoBidAmount(amount);
    state = state.copyWith(maxAutoBidAmount: amount);
  }
}

class ErrorHandler {

  ErrorHandler(this._logger);
  final AppLogger _logger;

  String handleException(Exception exception) {
    if (exception is ApiException) {
      return exception.message;
    } else if (exception is NetworkException) {
      return 'Network error. Please check your connection.';
    } else if (exception is AuthException) {
      return 'Authentication error. Please login again.';
    } else if (exception is ValidationException) {
      return exception.message;
    } else {
      AppLogger.error('Unknown error occurred', error: exception);
      return 'An unexpected error occurred. Please try again.';
    }
  }
}

// State Classes
class AuthState {

  const AuthState({
    this.isLoading = false,
    this.isAuthenticated = false,
    this.user,
    this.error,
  });
  final bool isLoading;
  final bool isAuthenticated;
  final AuthModel? user;
  final String? error;

  static AuthState initial() => const AuthState();

  AuthState copyWith({
    bool? isLoading,
    bool? isAuthenticated,
    AuthModel? user,
    String? error,
  }) {
    return AuthState(
      isLoading: isLoading ?? this.isLoading,
      isAuthenticated: isAuthenticated ?? this.isAuthenticated,
      user: user ?? this.user,
      error: error ?? this.error,
    );
  }
}

class ThemeState {

  const ThemeState({this.theme = 'light'});
  final String theme;
}

class NotificationSettings {

  const NotificationSettings({
    this.enabled = true,
    this.autoBid = false,
    this.maxAutoBidAmount = 1000.0,
  });
  final bool enabled;
  final bool autoBid;
  final double maxAutoBidAmount;

  NotificationSettings copyWith({
    bool? enabled,
    bool? autoBid,
    double? maxAutoBidAmount,
  }) {
    return NotificationSettings(
      enabled: enabled ?? this.enabled,
      autoBid: autoBid ?? this.autoBid,
      maxAutoBidAmount: maxAutoBidAmount ?? this.maxAutoBidAmount,
    );
  }
}


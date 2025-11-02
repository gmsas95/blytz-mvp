import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../network/api_client.dart';
import '../network/interceptors.dart';
import '../network/network_info.dart';
import '../storage/secure_storage.dart';
import '../storage/local_database.dart';
import '../storage/preferences.dart';
import '../utils/logger.dart';
import '../../features/auth/data/models/auth_model.dart';

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
    ref.read(secureStorageProvider),
    ref.read(appPreferencesProvider),
    ref.read(loggerProvider),
  );
});

// Current User Provider
final currentUserProvider = Provider<UserModel?>((ref) {
  return ref.watch(authStateProvider).user;
});

// Is Authenticated Provider
final isAuthenticatedProvider = Provider<bool>((ref) {
  return ref.watch(authStateProvider).isAuthenticated;
});

// Theme Provider
final themeProvider = StateNotifierProvider<ThemeNotifier, ThemeState>((ref) {
  return ThemeNotifier(ref.read(appPreferencesProvider));
});

// Language Provider
final languageProvider = StateNotifierProvider<LanguageNotifier, String>((ref) {
  return LanguageNotifier(ref.read(appPreferencesProvider));
});

// Connectivity Provider
final connectivityProvider = StreamProvider<ConnectivityResult>((ref) {
  return Connectivity().onConnectivityChanged;
});

// Is Online Provider
final isOnlineProvider = Provider<bool>((ref) {
  final connectivity = ref.watch(connectivityProvider);
  return connectivity.when(
    data: (result) => result != ConnectivityResult.none,
    loading: () => true,
    error: (_, __) => false,
  );
});

// Notification Settings Provider
final notificationSettingsProvider = StateNotifierProvider<NotificationSettingsNotifier, NotificationSettings>((ref) {
  return NotificationSettingsNotifier(ref.read(appPreferencesProvider));
});

// App Initialization Provider
final appInitializationProvider = FutureProvider<void>((ref) async {
  try {
    // Initialize local database
    await ref.read(localDatabaseProvider).init();
    
    // Check authentication status
    await ref.read(authStateProvider.notifier).checkAuthStatus();
    
    // Load preferences
    await ref.read(themeProvider.notifier).loadTheme();
    await ref.read(languageProvider.notifier).loadLanguage();
    await ref.read(notificationSettingsProvider.notifier).loadSettings();
    
    AppLogger.info('App initialized successfully');
  } catch (e) {
    AppLogger.error('Failed to initialize app', error: e);
    throw e;
  }
});

// Error Handler Provider
final errorHandlerProvider = Provider<ErrorHandler>((ref) {
  return ErrorHandler(ref.read(loggerProvider));
});

class AuthNotifier extends StateNotifier<AuthState> {
  final ApiClient _apiClient;
  final SecureStorage _secureStorage;
  final AppPreferences _appPreferences;
  final AppLogger _logger;

  AuthNotifier(
    this._apiClient,
    this._secureStorage,
    this._appPreferences,
    this._logger,
  ) : super(AuthState.initial());

  Future<void> login(String email, String password) async {
    state = state.copyWith(isLoading: true, error: null);
    
    try {
      final response = await _apiClient.login(
        LoginRequest(email: email, password: password),
      );

      await _secureStorage.storeToken(response.token);
      await _secureStorage.storeRefreshToken(response.refreshToken);
      await _secureStorage.storeUser(response.user.toJson());
      await _appPreferences.setOnboardingCompleted(true);

      state = state.copyWith(
        isLoading: false,
        isAuthenticated: true,
        user: response.user,
      );

      _logger.info('User logged in successfully', error: response.user.email);
    } catch (e) {
      final error = ErrorHandler(_logger).handleException(e as Exception);
      state = state.copyWith(isLoading: false, error: error.message);
      _logger.error('Login failed', error: e);
    }
  }

  Future<void> register(String email, String password, String firstName, String lastName, String phone) async {
    state = state.copyWith(isLoading: true, error: null);
    
    try {
      final response = await _apiClient.register(
        RegisterRequest(
          email: email,
          password: password,
          firstName: firstName,
          lastName: lastName,
          phone: phone,
        ),
      );

      await _secureStorage.storeToken(response.token);
      await _secureStorage.storeRefreshToken(response.refreshToken);
      await _secureStorage.storeUser(response.user.toJson());
      await _appPreferences.setOnboardingCompleted(true);

      state = state.copyWith(
        isLoading: false,
        isAuthenticated: true,
        user: response.user,
      );

      _logger.info('User registered successfully', error: response.user.email);
    } catch (e) {
      final error = ErrorHandler(_logger).handleException(e as Exception);
      state = state.copyWith(isLoading: false, error: error.message);
      _logger.error('Registration failed', error: e);
    }
  }

  Future<void> logout() async {
    try {
      await _apiClient.logout();
    } catch (e) {
      _logger.warning('Logout API call failed', error: e);
    }

    await _secureStorage.clearAll();
    state = AuthState.initial();
    
    _logger.info('User logged out');
  }

  Future<void> checkAuthStatus() async {
    final hasToken = await _secureStorage.hasToken();
    if (!hasToken) {
      state = AuthState.initial();
      return;
    }

    try {
      final userJson = await _secureStorage.getUser();
      if (userJson != null) {
        final user = UserModel.fromJson(jsonDecode(userJson));
        state = state.copyWith(
          isAuthenticated: true,
          user: user,
        );
      }
    } catch (e) {
      _logger.error('Failed to check auth status', error: e);
      await logout();
    }
  }

  void clearError() {
    state = state.copyWith(error: null);
  }
}

class ThemeNotifier extends StateNotifier<ThemeState> {
  final AppPreferences _appPreferences;

  ThemeNotifier(this._appPreferences) : super(const ThemeState());

  Future<void> loadTheme() async {
    final theme = await _appPreferences.getTheme();
    state = ThemeState(theme: theme);
  }

  Future<void> setTheme(String theme) async {
    await _appPreferences.setTheme(theme);
    state = ThemeState(theme: theme);
  }
}

class LanguageNotifier extends StateNotifier<String> {
  final AppPreferences _appPreferences;

  LanguageNotifier(this._appPreferences) : super('en');

  Future<void> loadLanguage() async {
    final language = await _appPreferences.getLanguage();
    state = language;
  }

  Future<void> setLanguage(String language) async {
    await _appPreferences.setLanguage(language);
    state = language;
  }
}

class NotificationSettingsNotifier extends StateNotifier<NotificationSettings> {
  final AppPreferences _appPreferences;

  NotificationSettingsNotifier(this._appPreferences) : super(const NotificationSettings());

  Future<void> loadSettings() async {
    final enabled = await _appPreferences.isNotificationsEnabled();
    final autoBid = await _appPreferences.isAutoBidEnabled();
    final maxAutoBid = await _appPreferences.getMaxAutoBidAmount();

    state = NotificationSettings(
      enabled: enabled,
      autoBid: autoBid,
      maxAutoBidAmount: maxAutoBid,
    );
  }

  Future<void> setNotificationsEnabled(bool enabled) async {
    await _appPreferences.setNotificationsEnabled(enabled);
    state = state.copyWith(enabled: enabled);
  }

  Future<void> setAutoBidEnabled(bool enabled) async {
    await _appPreferences.setAutoBidEnabled(enabled);
    state = state.copyWith(autoBid: enabled);
  }

  Future<void> setMaxAutoBidAmount(double amount) async {
    await _appPreferences.setMaxAutoBidAmount(amount);
    state = state.copyWith(maxAutoBidAmount: amount);
  }
}

class ErrorHandler {
  final AppLogger _logger;

  ErrorHandler(this._logger);

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
      _logger.error('Unknown error occurred', error: exception);
      return 'An unexpected error occurred. Please try again.';
    }
  }
}

// State Classes
class AuthState {
  final bool isLoading;
  final bool isAuthenticated;
  final UserModel? user;
  final String? error;

  const AuthState({
    this.isLoading = false,
    this.isAuthenticated = false,
    this.user,
    this.error,
  });

  static AuthState initial() => const AuthState();

  AuthState copyWith({
    bool? isLoading,
    bool? isAuthenticated,
    UserModel? user,
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
  final String theme;

  const ThemeState({this.theme = 'light'});
}

class NotificationSettings {
  final bool enabled;
  final bool autoBid;
  final double maxAutoBidAmount;

  const NotificationSettings({
    this.enabled = true,
    this.autoBid = false,
    this.maxAutoBidAmount = 1000.0,
  });

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

// Import required dependencies
import 'dart:convert';
import 'package:dio/dio.dart';
import 'package:connectivity_plus/connectivity_plus.dart';
import '../../features/auth/data/models/auth_model.dart';
import '../network/api_client.dart';
import '../network/interceptors.dart';
import '../network/network_info.dart';
import '../storage/secure_storage.dart';
import '../storage/local_database.dart';
import '../storage/preferences.dart';
import '../utils/logger.dart';
import '../errors/exceptions.dart';
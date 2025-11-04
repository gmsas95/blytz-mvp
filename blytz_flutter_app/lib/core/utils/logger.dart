import 'package:logger/logger.dart';

class AppLogger {
  static final Logger _logger = Logger(
    printer: PrettyPrinter(
      printTime: true,
    ),
  );

  static void debug(String message, {dynamic error, StackTrace? stackTrace}) {
    _logger.d(message, error: error, stackTrace: stackTrace);
  }

  static void info(String message, {dynamic error, StackTrace? stackTrace}) {
    _logger.i(message, error: error, stackTrace: stackTrace);
  }

  static void warning(String message, {dynamic error, StackTrace? stackTrace}) {
    _logger.w(message, error: error, stackTrace: stackTrace);
  }

  static void error(String message, {dynamic error, StackTrace? stackTrace}) {
    _logger.e(message, error: error, stackTrace: stackTrace);
  }

  static void logApiRequest(String method, String url, {dynamic data}) {
    _logger.d('API Request: $method $url', error: data);
  }

  static void logApiResponse(String method, String url, int statusCode, {dynamic data}) {
    _logger.i('API Response: $method $url - $statusCode', error: data);
  }

  static void logApiError(String method, String url, dynamic error) {
    _logger.e('API Error: $method $url', error: error);
  }

  static void logWebSocketEvent(String event, {dynamic data}) {
    _logger.d('WebSocket: $event', error: data);
  }

  static void logUserAction(String action, {Map<String, dynamic>? params}) {
    _logger.i('User Action: $action', error: params);
  }

  static void logPerformance(String operation, Duration duration, {Map<String, dynamic>? params}) {
    _logger.i('Performance: $operation took ${duration.inMilliseconds}ms', error: params);
  }
}
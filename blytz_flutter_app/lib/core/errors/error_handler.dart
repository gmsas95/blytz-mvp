import 'package:dio/dio.dart';
import '../errors/exceptions.dart';
import '../errors/failures.dart';

class ErrorHandler {
  static Failure handleException(Exception exception) {
    if (exception is ApiException) {
      return ServerFailure(
        exception.message,
        statusCode: exception.statusCode,
      );
    } else if (exception is NetworkException) {
      return NetworkFailure(exception.message);
    } else if (exception is AuthException) {
      return AuthFailure(exception.message);
    } else if (exception is ValidationException) {
      return ValidationFailure(exception.message, field: exception.field);
    } else if (exception is StorageException) {
      return StorageFailure(exception.message);
    } else if (exception is DioException) {
      return _handleDioException(exception);
    } else {
      return UnknownFailure(exception.toString());
    }
  }
  
  static Failure _handleDioException(DioException exception) {
    switch (exception.type) {
      case DioExceptionType.connectionTimeout:
      case DioExceptionType.sendTimeout:
      case DioExceptionType.receiveTimeout:
        return const NetworkFailure('Connection timeout');
      case DioExceptionType.badResponse:
        final statusCode = exception.response?.statusCode;
        final message = exception.response?.data?['message'] ?? 
                       exception.response?.statusMessage ?? 
                       'Server error';
        return ServerFailure(message, statusCode: statusCode);
      case DioExceptionType.cancel:
        return const NetworkFailure('Request cancelled');
      case DioExceptionType.connectionError:
        return const NetworkFailure('No internet connection');
      case DioExceptionType.unknown:
      default:
        return NetworkFailure(exception.message ?? 'Network error');
    }
  }
}
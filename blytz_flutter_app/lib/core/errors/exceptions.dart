import 'package:dio/dio.dart';

class ApiException implements Exception {

  const ApiException(this.message, {this.statusCode, this.data});

  factory ApiException.fromDioError(DioException error) {
    if (error.response?.data is Map<String, dynamic>) {
      final data = error.response?.data as Map<String, dynamic>;
      final message = data['message']?.toString() ?? 'Unknown error';
      final statusCode = error.response?.statusCode;

      return ApiException(
        message,
        statusCode: statusCode,
        data: data,
      );
    }

    final message = error.message ?? 'Network error';
    final statusCode = error.response?.statusCode;

    return ApiException(
      message,
      statusCode: statusCode,
    );
  }

  final String message;
  final int? statusCode;
  final dynamic data;

  @override
  String toString() => 'ApiException: $message';
}

class NetworkException implements Exception {
  
  const NetworkException(this.message);
  final String message;
  
  @override
  String toString() => 'NetworkException: $message';
}

class AuthException implements Exception {
  
  const AuthException(this.message);
  final String message;
  
  @override
  String toString() => 'AuthException: $message';
}

class ValidationException implements Exception {
  
  const ValidationException(this.field, this.message);
  final String field;
  final String message;
  
  @override
  String toString() => 'ValidationException: $field - $message';
}

class StorageException implements Exception {
  
  const StorageException(this.message);
  final String message;
  
  @override
  String toString() => 'StorageException: $message';
}
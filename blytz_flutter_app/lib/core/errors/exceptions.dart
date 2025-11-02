class ApiException implements Exception {
  final String message;
  final int? statusCode;
  final dynamic data;
  
  const ApiException(this.message, {this.statusCode, this.data});
  
  @override
  String toString() => 'ApiException: $message';
  
  factory ApiException.fromDioError(dynamic error) {
    if (error.response?.data is Map<String, dynamic>) {
      final data = error.response?.data as Map<String, dynamic>;
      return ApiException(
        data['message'] ?? 'Unknown error',
        statusCode: error.response?.statusCode,
        data: data,
      );
    }
    return ApiException(
      error.message ?? 'Network error',
      statusCode: error.response?.statusCode,
    );
  }
}

class NetworkException implements Exception {
  final String message;
  
  const NetworkException(this.message);
  
  @override
  String toString() => 'NetworkException: $message';
}

class AuthException implements Exception {
  final String message;
  
  const AuthException(this.message);
  
  @override
  String toString() => 'AuthException: $message';
}

class ValidationException implements Exception {
  final String field;
  final String message;
  
  const ValidationException(this.field, this.message);
  
  @override
  String toString() => 'ValidationException: $field - $message';
}

class StorageException implements Exception {
  final String message;
  
  const StorageException(this.message);
  
  @override
  String toString() => 'StorageException: $message';
}
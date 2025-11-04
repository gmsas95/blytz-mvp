abstract class Failure {
  const Failure(this.message);
  final String message;
}

class ServerFailure extends Failure {
  const ServerFailure(super.message, {this.statusCode});
  
  final int? statusCode;
  
  static ServerFailure fromDioException(DioException e) {
    return ServerFailure(
      e.message ?? 'Server error occurred',
      statusCode: e.response?.statusCode,
    );
  }
  
  @override
  String toString() => 'ServerFailure: $message';
}

class NetworkFailure extends Failure {
  const NetworkFailure(super.message);
  
  @override
  String toString() => 'NetworkFailure: $message';
}

class AuthFailure extends Failure {
  const AuthFailure(super.message);
  
  @override
  String toString() => 'AuthFailure: $message';
}

class ValidationFailure extends Failure {
  const ValidationFailure(super.message, {this.field});
  
  final String? field;
  
  @override
  String toString() => 'ValidationFailure: $field - $message';
}

class StorageFailure extends Failure {
  const StorageFailure(super.message);
  
  @override
  String toString() => 'StorageFailure: $message';
}

class CacheFailure extends Failure {
  const CacheFailure(super.message);
  
  @override
  String toString() => 'CacheFailure: $message';
}

class NotFoundFailure extends Failure {
  const NotFoundFailure(super.message);
  
  @override
  String toString() => 'NotFoundFailure: $message';
}

class UnknownFailure extends Failure {
  const UnknownFailure(super.message);
  
  @override
  String toString() => 'UnknownFailure: $message';
}
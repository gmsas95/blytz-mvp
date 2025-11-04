import 'package:blytz_flutter_app/core/constants/api_constants.dart';
import 'package:blytz_flutter_app/core/storage/secure_storage.dart';
import 'package:dio/dio.dart';

class AuthInterceptor extends Interceptor {
  @override
  Future<void> onRequest(RequestOptions options, RequestInterceptorHandler handler) async {
    final token = await SecureStorage.getToken();
    if (token != null) {
      options.headers['Authorization'] = 'Bearer $token';
    }
    handler.next(options);
  }

  @override
  Future<void> onError(DioException err, ErrorInterceptorHandler handler) async {
    if (err.response?.statusCode == 401) {
      // Token expired, try to refresh
      try {
        final refreshToken = await SecureStorage.getRefreshToken();
        if (refreshToken != null) {
          final dio = Dio(BaseOptions(
            baseUrl: ApiConstants.baseUrl,
            connectTimeout: ApiConstants.connectTimeout,
            receiveTimeout: ApiConstants.receiveTimeout,
          ),);

          final response = await dio.post(
            ApiConstants.refresh,
            data: {'refreshToken': refreshToken},
          );

          final newToken = response.data['token'] as String;
          await SecureStorage.storeToken(newToken);

          // Retry original request with new token
          final originalRequest = err.requestOptions;
          originalRequest.headers['Authorization'] = 'Bearer $newToken';
          
          final retryResponse = await dio.fetch(originalRequest);
          handler.resolve(retryResponse);
          return;
        }
      } catch (e) {
        // Refresh failed, clear tokens and redirect to login
        await SecureStorage.clearAll();
      }
    }
    handler.next(err);
  }
}

class LoggingInterceptor extends Interceptor {
  @override
  void onRequest(RequestOptions options, RequestInterceptorHandler handler) {
    print('REQUEST: ${options.method} ${options.uri}');
    if (options.data != null) {
      print('DATA: ${options.data}');
    }
    handler.next(options);
  }

  @override
  void onResponse(Response response, ResponseInterceptorHandler handler) {
    print('RESPONSE: ${response.statusCode} ${response.requestOptions.uri}');
    print('DATA: ${response.data}');
    handler.next(response);
  }

  @override
  void onError(DioException err, ErrorInterceptorHandler handler) {
    print('ERROR: ${err.message}');
    print('RESPONSE: ${err.response?.data}');
    handler.next(err);
  }
}
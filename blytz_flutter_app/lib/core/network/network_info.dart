import 'package:connectivity_plus/connectivity_plus.dart';
import 'package:dio/dio.dart';

class NetworkInfo {
  
  NetworkInfo(this._connectivity);
  final Connectivity _connectivity;

  Future<bool> get isConnected async {
    final result = await _connectivity.checkConnectivity();
    return result.isNotEmpty && !result.contains(ConnectivityResult.none);
  }

  Future<bool> get hasInternetAccess async {
    if (!await isConnected) return false;
    
    try {
      final dio = Dio();
      dio.options.connectTimeout = const Duration(seconds: 5);
      dio.options.receiveTimeout = const Duration(seconds: 5);
      final response = await dio.get<String>(
        'https://www.google.com',
      );
      return response.statusCode == 200;
    } catch (e) {
      return false;
    }
  }
}
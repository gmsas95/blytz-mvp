import 'package:dio/dio.dart';
import 'package:connectivity_plus/connectivity_plus.dart';
import '../constants/api_constants.dart';
import '../errors/exceptions.dart';

class NetworkInfo {
  final Connectivity _connectivity;
  
  NetworkInfo(this._connectivity);

  Future<bool> get isConnected async {
    final result = await _connectivity.checkConnectivity();
    return result != ConnectivityResult.none;
  }

  Future<bool> get hasInternetAccess async {
    if (!await isConnected) return false;
    
    try {
      final dio = Dio();
      final response = await dio.get(
        'https://www.google.com',
        options: Options(
          connectTimeout: const Duration(seconds: 5),
          receiveTimeout: const Duration(seconds: 5),
        ),
      );
      return response.statusCode == 200;
    } catch (e) {
      return false;
    }
  }
}
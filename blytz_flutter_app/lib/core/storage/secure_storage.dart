import 'package:blytz_flutter_app/core/constants/app_constants.dart';
import 'package:blytz_flutter_app/core/errors/exceptions.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';

class SecureStorage {
  final FlutterSecureStorage _storage = const FlutterSecureStorage(
    aOptions: AndroidOptions(
      encryptedSharedPreferences: true,
    ),
  );

  Future<void> storeToken(String token) async {
    try {
      await _storage.write(key: AppConstants.authTokenKey, value: token);
    } catch (e) {
      throw StorageException('Failed to store token: $e');
    }
  }

  Future<String?> getToken() async {
    try {
      return await _storage.read(key: AppConstants.authTokenKey);
    } catch (e) {
      throw StorageException('Failed to get token: $e');
    }
  }

  Future<void> storeRefreshToken(String refreshToken) async {
    try {
      await _storage.write(key: AppConstants.refreshTokenKey, value: refreshToken);
    } catch (e) {
      throw StorageException('Failed to store refresh token: $e');
    }
  }

  Future<String?> getRefreshToken() async {
    try {
      return await _storage.read(key: AppConstants.refreshTokenKey);
    } catch (e) {
      throw StorageException('Failed to get refresh token: $e');
    }
  }

  Future<void> storeUserId(String userId) async {
    try {
      await _storage.write(key: AppConstants.userIdKey, value: userId);
    } catch (e) {
      throw StorageException('Failed to store user ID: $e');
    }
  }

  Future<String?> getUserId() async {
    try {
      return await _storage.read(key: AppConstants.userIdKey);
    } catch (e) {
      throw StorageException('Failed to get user ID: $e');
    }
  }

  Future<void> storeUser(String userData) async {
    try {
      await _storage.write(key: AppConstants.userKey, value: userData);
    } catch (e) {
      throw StorageException('Failed to store user data: $e');
    }
  }

  Future<String?> getUser() async {
    try {
      return await _storage.read(key: AppConstants.userKey);
    } catch (e) {
      throw StorageException('Failed to get user data: $e');
    }
  }

  Future<void> clearToken() async {
    try {
      await _storage.delete(key: AppConstants.authTokenKey);
    } catch (e) {
      throw StorageException('Failed to clear token: $e');
    }
  }

  Future<void> clearRefreshToken() async {
    try {
      await _storage.delete(key: AppConstants.refreshTokenKey);
    } catch (e) {
      throw StorageException('Failed to clear refresh token: $e');
    }
  }

  Future<void> clearUser() async {
    try {
      await _storage.delete(key: AppConstants.userKey);
    } catch (e) {
      throw StorageException('Failed to clear user data: $e');
    }
  }

  Future<void> clearAll() async {
    try {
      await _storage.deleteAll();
    } catch (e) {
      throw StorageException('Failed to clear all data: $e');
    }
  }

  Future<bool> hasToken() async {
    final token = await getToken();
    return token != null && token.isNotEmpty;
  }

  Future<bool> hasRefreshToken() async {
    final refreshToken = await getRefreshToken();
    return refreshToken != null && refreshToken.isNotEmpty;
  }

  // Static methods for backward compatibility
  static Future<void> storeTokenStatic(String token) async {
    final storage = SecureStorage();
    await storage.storeToken(token);
  }

  static Future<String?> getTokenStatic() async {
    final storage = SecureStorage();
    return await storage.getToken();
  }

  static Future<void> clearAllStatic() async {
    final storage = SecureStorage();
    await storage.clearAll();
  }
}
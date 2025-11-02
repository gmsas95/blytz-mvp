import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import '../constants/app_constants.dart';
import '../errors/exceptions.dart';

class SecureStorage {
  static const FlutterSecureStorage _storage = FlutterSecureStorage(
    aOptions: AndroidOptions(
      encryptedSharedPreferences: true,
    ),
    iOptions: IOSOptions(
      accessibility: KeychainItemAccessibility.first_unlock_this_device,
    ),
  );

  static Future<void> storeToken(String token) async {
    try {
      await _storage.write(key: AppConstants.authTokenKey, value: token);
    } catch (e) {
      throw StorageException('Failed to store token: $e');
    }
  }

  static Future<String?> getToken() async {
    try {
      return await _storage.read(key: AppConstants.authTokenKey);
    } catch (e) {
      throw StorageException('Failed to get token: $e');
    }
  }

  static Future<void> storeRefreshToken(String refreshToken) async {
    try {
      await _storage.write(key: AppConstants.refreshTokenKey, value: refreshToken);
    } catch (e) {
      throw StorageException('Failed to store refresh token: $e');
    }
  }

  static Future<String?> getRefreshToken() async {
    try {
      return await _storage.read(key: AppConstants.refreshTokenKey);
    } catch (e) {
      throw StorageException('Failed to get refresh token: $e');
    }
  }

  static Future<void> storeUser(String userData) async {
    try {
      await _storage.write(key: AppConstants.userKey, value: userData);
    } catch (e) {
      throw StorageException('Failed to store user data: $e');
    }
  }

  static Future<String?> getUser() async {
    try {
      return await _storage.read(key: AppConstants.userKey);
    } catch (e) {
      throw StorageException('Failed to get user data: $e');
    }
  }

  static Future<void> clearToken() async {
    try {
      await _storage.delete(key: AppConstants.authTokenKey);
    } catch (e) {
      throw StorageException('Failed to clear token: $e');
    }
  }

  static Future<void> clearRefreshToken() async {
    try {
      await _storage.delete(key: AppConstants.refreshTokenKey);
    } catch (e) {
      throw StorageException('Failed to clear refresh token: $e');
    }
  }

  static Future<void> clearUser() async {
    try {
      await _storage.delete(key: AppConstants.userKey);
    } catch (e) {
      throw StorageException('Failed to clear user data: $e');
    }
  }

  static Future<void> clearAll() async {
    try {
      await _storage.deleteAll();
    } catch (e) {
      throw StorageException('Failed to clear all data: $e');
    }
  }

  static Future<bool> hasToken() async {
    final token = await getToken();
    return token != null && token.isNotEmpty;
  }

  static Future<bool> hasRefreshToken() async {
    final refreshToken = await getRefreshToken();
    return refreshToken != null && refreshToken.isNotEmpty;
  }
}
import 'package:blytz_flutter_app/core/constants/app_constants.dart';
import 'package:blytz_flutter_app/core/errors/exceptions.dart';
import 'package:shared_preferences/shared_preferences.dart';

class AppPreferences {
  static Future<SharedPreferences> get _prefs async =>
      SharedPreferences.getInstance();

  // Theme management
  static Future<void> setTheme(String theme) async {
    try {
      final prefs = await _prefs;
      await prefs.setString(AppConstants.themeKey, theme);
    } catch (e) {
      throw StorageException('Failed to set theme: $e');
    }
  }

  static Future<String> getTheme() async {
    try {
      final prefs = await _prefs;
      return prefs.getString(AppConstants.themeKey) ?? 'light';
    } catch (e) {
      throw StorageException('Failed to get theme: $e');
    }
  }

  // Language management
  static Future<void> setLanguage(String language) async {
    try {
      final prefs = await _prefs;
      await prefs.setString(AppConstants.languageKey, language);
    } catch (e) {
      throw StorageException('Failed to set language: $e');
    }
  }

  static Future<String> getLanguage() async {
    try {
      final prefs = await _prefs;
      return prefs.getString(AppConstants.languageKey) ?? 'en';
    } catch (e) {
      throw StorageException('Failed to get language: $e');
    }
  }

  // Onboarding status
  static Future<void> setOnboardingCompleted(bool completed) async {
    try {
      final prefs = await _prefs;
      await prefs.setBool('onboarding_completed', completed);
    } catch (e) {
      throw StorageException('Failed to set onboarding status: $e');
    }
  }

  static Future<bool> isOnboardingCompleted() async {
    try {
      final prefs = await _prefs;
      return prefs.getBool('onboarding_completed') ?? false;
    } catch (e) {
      throw StorageException('Failed to get onboarding status: $e');
    }
  }

  // Notification preferences
  static Future<void> setNotificationsEnabled(bool enabled) async {
    try {
      final prefs = await _prefs;
      await prefs.setBool('notifications_enabled', enabled);
    } catch (e) {
      throw StorageException('Failed to set notifications preference: $e');
    }
  }

  static Future<bool> isNotificationsEnabled() async {
    try {
      final prefs = await _prefs;
      return prefs.getBool('notifications_enabled') ?? true;
    } catch (e) {
      throw StorageException('Failed to get notifications preference: $e');
    }
  }

  // Bid preferences
  static Future<void> setAutoBidEnabled(bool enabled) async {
    try {
      final prefs = await _prefs;
      await prefs.setBool('auto_bid_enabled', enabled);
    } catch (e) {
      throw StorageException('Failed to set auto bid preference: $e');
    }
  }

  static Future<bool> isAutoBidEnabled() async {
    try {
      final prefs = await _prefs;
      return prefs.getBool('auto_bid_enabled') ?? false;
    } catch (e) {
      throw StorageException('Failed to get auto bid preference: $e');
    }
  }

  static Future<void> setMaxAutoBidAmount(double amount) async {
    try {
      final prefs = await _prefs;
      await prefs.setDouble('max_auto_bid_amount', amount);
    } catch (e) {
      throw StorageException('Failed to set max auto bid amount: $e');
    }
  }

  static Future<double> getMaxAutoBidAmount() async {
    try {
      final prefs = await _prefs;
      return prefs.getDouble('max_auto_bid_amount') ?? 1000.0;
    } catch (e) {
      throw StorageException('Failed to get max auto bid amount: $e');
    }
  }

  // Search preferences
  static Future<void> setRecentSearches(List<String> searches) async {
    try {
      final prefs = await _prefs;
      await prefs.setStringList('recent_searches', searches);
    } catch (e) {
      throw StorageException('Failed to set recent searches: $e');
    }
  }

  static Future<List<String>> getRecentSearches() async {
    try {
      final prefs = await _prefs;
      return prefs.getStringList('recent_searches') ?? [];
    } catch (e) {
      throw StorageException('Failed to get recent searches: $e');
    }
  }

  static Future<void> addRecentSearch(String search) async {
    try {
      final searches = await getRecentSearches();
      searches.remove(search); // Remove if exists
      searches.insert(0, search); // Add to beginning

      // Keep only last 10 searches
      if (searches.length > 10) {
        searches.removeRange(10, searches.length);
      }

      await setRecentSearches(searches);
    } catch (e) {
      throw StorageException('Failed to add recent search: $e');
    }
  }

  // Clear all preferences
  static Future<void> clearAll() async {
    try {
      final prefs = await _prefs;
      await prefs.clear();
    } catch (e) {
      throw StorageException('Failed to clear preferences: $e');
    }
  }
}

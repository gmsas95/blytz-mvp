import 'package:shared_preferences/shared_preferences.dart';
import 'package:hive_flutter/hive_flutter.dart';
import '../constants/app_constants.dart';
import '../errors/exceptions.dart';

class LocalDatabase {
  static const String _auctionBoxName = 'auctions';
  static const String _userBoxName = 'users';
  static const String _productBoxName = 'products';
  static const String _chatBoxName = 'chats';

  static Future<void> init() async {
    try {
      await Hive.initFlutter();
      
      // Register adapters here when models are ready
      // Hive.registerAdapter(AuctionModelAdapter());
      // Hive.registerAdapter(UserModelAdapter());
      // Hive.registerAdapter(ProductModelAdapter());
      // Hive.registerAdapter(MessageModelAdapter());

      // Open boxes
      await Hive.openBox(_auctionBoxName);
      await Hive.openBox(_userBoxName);
      await Hive.openBox(_productBoxName);
      await Hive.openBox(_chatBoxName);
    } catch (e) {
      throw StorageException('Failed to initialize local database: $e');
    }
  }

  // Auction caching
  static Future<void> cacheAuctions(List<dynamic> auctions) async {
    try {
      final box = await Hive.openBox(_auctionBoxName);
      await box.clear();
      for (final auction in auctions) {
        await box.put(auction.id, auction);
      }
    } catch (e) {
      throw StorageException('Failed to cache auctions: $e');
    }
  }

  static Future<List<dynamic>> getCachedAuctions() async {
    try {
      final box = await Hive.openBox(_auctionBoxName);
      return box.values.toList();
    } catch (e) {
      throw StorageException('Failed to get cached auctions: $e');
    }
  }

  static Future<void> cacheAuction(dynamic auction) async {
    try {
      final box = await Hive.openBox(_auctionBoxName);
      await box.put(auction.id, auction);
    } catch (e) {
      throw StorageException('Failed to cache auction: $e');
    }
  }

  static Future<dynamic?> getCachedAuction(String id) async {
    try {
      final box = await Hive.openBox(_auctionBoxName);
      return box.get(id);
    } catch (e) {
      throw StorageException('Failed to get cached auction: $e');
    }
  }

  // User caching
  static Future<void> cacheUser(dynamic user) async {
    try {
      final box = await Hive.openBox(_userBoxName);
      await box.put('current_user', user);
    } catch (e) {
      throw StorageException('Failed to cache user: $e');
    }
  }

  static Future<dynamic?> getCachedUser() async {
    try {
      final box = await Hive.openBox(_userBoxName);
      return box.get('current_user');
    } catch (e) {
      throw StorageException('Failed to get cached user: $e');
    }
  }

  // Product caching
  static Future<void> cacheProducts(List<dynamic> products) async {
    try {
      final box = await Hive.openBox(_productBoxName);
      await box.clear();
      for (final product in products) {
        await box.put(product.id, product);
      }
    } catch (e) {
      throw StorageException('Failed to cache products: $e');
    }
  }

  static Future<List<dynamic>> getCachedProducts() async {
    try {
      final box = await Hive.openBox(_productBoxName);
      return box.values.toList();
    } catch (e) {
      throw StorageException('Failed to get cached products: $e');
    }
  }

  // Chat caching
  static Future<void> cacheMessages(String roomId, List<dynamic> messages) async {
    try {
      final box = await Hive.openBox(_chatBoxName);
      await box.put('messages_$roomId', messages);
    } catch (e) {
      throw StorageException('Failed to cache messages: $e');
    }
  }

  static Future<List<dynamic>> getCachedMessages(String roomId) async {
    try {
      final box = await Hive.openBox(_chatBoxName);
      return box.get('messages_$roomId', defaultValue: <dynamic>[]);
    } catch (e) {
      throw StorageException('Failed to get cached messages: $e');
    }
  }

  static Future<void> clearAll() async {
    try {
      await Hive.deleteBoxFromDisk(_auctionBoxName);
      await Hive.deleteBoxFromDisk(_userBoxName);
      await Hive.deleteBoxFromDisk(_productBoxName);
      await Hive.deleteBoxFromDisk(_chatBoxName);
    } catch (e) {
      throw StorageException('Failed to clear local database: $e');
    }
  }
}
import 'package:blytz_flutter_app/core/di/service_locator.dart';
import 'package:blytz_flutter_app/core/network/api_client.dart';
import 'package:blytz_flutter_app/features/auth/data/models/auth_model.dart';
import 'package:blytz_flutter_app/features/auth/data/repositories/auth_repository_impl.dart';
import 'package:blytz_flutter_app/features/auction/data/models/auction_model.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:mockito/mockito.dart';
import 'package:mockito/annotations.dart';

import 'backend_integration_test.mocks.dart';

@GenerateMocks([ApiClient, AuthRepository])
void main() {
  group('Backend Integration Tests', () {
    late ProviderContainer container;
    late MockApiClient mockApiClient;
    late MockAuthRepository mockAuthRepository;

    setUp(() {
      mockApiClient = MockApiClient();
      container = ProviderContainer(
        overrides: [
          apiClientProvider.overrideWithValue(mockApiClient),
        ],
      );
    });

    tearDown(() {
      container.dispose();
    });

    test('should initialize services successfully', () async {
      // Test service initialization
      await expectLater(initializeServices(container), completes);
    });

    test('should create login request correctly', () {
      final loginRequest = LoginRequest(
        email: 'test@example.com',
        password: 'password123',
      );

      expect(loginRequest.email, equals('test@example.com'));
      expect(loginRequest.password, equals('password123'));

      final json = loginRequest.toJson();
      expect(json['email'], equals('test@example.com'));
      expect(json['password'], equals('password123'));
    });

    test('should create register request correctly', () {
      final registerRequest = RegisterRequest(
        email: 'test@example.com',
        password: 'password123',
        firstName: 'Test',
        lastName: 'User',
      );

      expect(registerRequest.email, equals('test@example.com'));
      expect(registerRequest.firstName, equals('Test'));
      expect(registerRequest.lastName, equals('User'));

      final json = registerRequest.toJson();
      expect(json['firstName'], equals('Test'));
      expect(json['lastName'], equals('User'));
    });

    test('should parse auction model correctly', () {
      final auctionJson = {
        'id': '1',
        'title': 'Test Auction',
        'description': 'Test Description',
        'sellerId': 'seller1',
        'sellerName': 'Test Seller',
        'startingPrice': 100.0,
        'currentBid': 150.0,
        'startTime': '2024-01-01T10:00:00Z',
        'endTime': '2024-01-01T11:00:00Z',
        'status': 'active',
        'totalBids': 5,
        'images': ['image1.jpg', 'image2.jpg'],
        'categories': ['electronics', 'phones'],
        'isActive': true,
        'isFeatured': false,
        'watchers': 10,
      };

      final auction = AuctionModel.fromJson(auctionJson);

      expect(auction.id, equals('1'));
      expect(auction.title, equals('Test Auction'));
      expect(auction.currentBid, equals(150.0));
      expect(auction.totalBids, equals(5));
      expect(auction.isActive, isTrue);
      expect(auction.categories, contains('electronics'));
    });

    test('should format auction time left correctly', () {
      final now = DateTime.now();
      final auction = AuctionModel(
        id: '1',
        title: 'Test Auction',
        description: 'Test',
        product: null,
        sellerId: 'seller1',
        sellerName: 'Test Seller',
        startingPrice: 100.0,
        currentBid: 150.0,
        startTime: now.subtract(const Duration(minutes: 30)),
        endTime: now.add(const Duration(minutes: 15)),
      );

      expect(auction.timeLeft, contains('m'));
      expect(auction.formattedCurrentBid, equals('\$150.00'));
    });

    test('should detect if auction is ending soon', () {
      final now = DateTime.now();
      final endingSoonAuction = AuctionModel(
        id: '1',
        title: 'Ending Soon Auction',
        description: 'Test',
        product: null,
        sellerId: 'seller1',
        sellerName: 'Test Seller',
        startingPrice: 100.0,
        currentBid: 150.0,
        startTime: now.subtract(const Duration(hours: 1)),
        endTime: now.add(const Duration(minutes: 10)),
      );

      expect(endingSoonAuction.isEndingSoon, isTrue);
    });

    test('should handle authentication repository', () {
      final authRepo = container.read(authRepositoryProvider);
      expect(authRepo, isNotNull);
    });

    test('should create mock auth repository', () {
      mockAuthRepository = MockAuthRepository();
      expect(mockAuthRepository, isNotNull);
    });

    test('should handle auction repository', () {
      final auctionRepo = container.read(auctionRepositoryProvider);
      expect(auctionRepo, isNotNull);
    });

    test('should parse bid model correctly', () {
      final bidJson = {
        'id': 'bid1',
        'auctionId': 'auction1',
        'bidderId': 'bidder1',
        'bidderName': 'Test Bidder',
        'amount': 200.0,
        'timestamp': '2024-01-01T10:30:00Z',
        'isWinning': true,
      };

      final bid = BidModel.fromJson(bidJson);

      expect(bid.id, equals('bid1'));
      expect(bid.amount, equals(200.0));
      expect(bid.isWinning, isTrue);
      expect(bid.formattedAmount, equals('\$200.00'));
    });

    test('should handle WebSocket service creation', () {
      final wsService = container.read(webSocketServiceProvider);
      expect(wsService, isNotNull);
      expect(wsService.isConnected, isFalse);
    });

    test('should provide dio client with interceptors', () {
      final dio = container.read(dioProvider);
      expect(dio, isNotNull);
      expect(dio.interceptors, isNotEmpty);
    });

    test('should provide secure storage', () {
      final secureStorage = container.read(secureStorageProvider);
      expect(secureStorage, isNotNull);
    });

    test('should provide network info', () {
      final networkInfo = container.read(networkInfoProvider);
      expect(networkInfo, isNotNull);
    });
  });

  group('API Constants Tests', () {
    test('should have correct base URLs', () {
      expect(ApiConstants.authBaseUrl, equals('http://localhost:8084'));
      expect(ApiConstants.auctionBaseUrl, equals('http://localhost:8081'));
      expect(ApiConstants.productBaseUrl, equals('http://localhost:8082'));
      expect(ApiConstants.orderBaseUrl, equals('http://localhost:8083'));
      expect(ApiConstants.paymentBaseUrl, equals('http://localhost:8086'));
      expect(ApiConstants.chatBaseUrl, equals('http://localhost:8088'));
      expect(ApiConstants.logisticsBaseUrl, equals('http://localhost:8087'));
    });

    test('should have correct endpoint paths', () {
      expect(ApiConstants.authServiceLogin, equals('/api/v1/auth/login'));
      expect(ApiConstants.authServiceRegister, equals('/api/v1/auth/register'));
      expect(ApiConstants.auctions, equals('/api/v1/auctions'));
      expect(ApiConstants.products, equals('/api/v1/products'));
      expect(ApiConstants.orders, equals('/api/v1/orders'));
    });

    test('should have WebSocket URLs', () {
      expect(ApiConstants.wsAuctionUpdates, equals('ws://localhost:8081/ws/auctions'));
      expect(ApiConstants.wsChat, equals('ws://localhost:8088/chat'));
    });
  });
}